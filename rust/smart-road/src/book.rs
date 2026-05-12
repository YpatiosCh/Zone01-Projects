use std::collections::HashMap;

use crate::config::{
    APPROACH_RESERVATION_TICKS, FREE_CELLS_FOR_SPAWNING, GRID_SIZE, MAX_TICKS_PER_CELL,
    MIN_TICKS_PER_CELL, RESERVATION_LOOKAHEAD, ROAD_END, ROAD_START,
};
use crate::models::*;

// A time window during which a cell is blocked by a booked vehicle.
struct Reservation {
    start: Timestamp,
    end: Timestamp,
}

// Intersection booking table.
// Per-cell sorted reservation lists.  Gaps between bookings are available to
// cars that arrive in time, giving better throughput than the single-timestamp
// approach without re-introducing the original double-booking bugs.
pub struct Book {
    cells: HashMap<(i32, i32), Vec<Reservation>>,
}

// ── Book impl ────────────────────────────────────────────────────────────────

impl Book {
    pub fn new() -> Self {
        Book {
            cells: HashMap::new(),
        }
    }

    /// Returns `true` if the first  on-grid cells of `route` are all free
    /// at `now`.  Does not register anything.
    pub fn can_spawn(&self, route: &Routes, now: Timestamp) -> Option<Vec<CellPosition>> {
        let path = path_for_route(route);
        for cp in path.iter().take(FREE_CELLS_FOR_SPAWNING) {
            if !is_on_grid(cp.cell) {
                break;
            }
            let key = (cp.cell.x, cp.cell.y);
            if self.earliest_free(key, now, APPROACH_RESERVATION_TICKS) != now {
                return None;
            }
        }
        Some(path)
    }

    pub fn book(
        &mut self,
        path: Vec<CellPosition>,
        now: Timestamp,
    ) -> Vec<(Timestamp, CellPosition)> {
        // let path = path_for_route(route);
        let n = path.len();
        let mut times = vec![0u64; n];
        times[0] = now;

        self.prune(now);

        let core_start = path
            .iter()
            .position(|cp| is_on_grid(cp.cell) && is_core_cell(cp.cell));
        let core_end = path
            .iter()
            .rposition(|cp| is_on_grid(cp.cell) && is_core_cell(cp.cell));

        if let (Some(cs), Some(ce)) = (core_start, core_end) {
            let n_core = ce - cs + 1;

            // Phase 1 — approach.
            // Reservation duration = 2×MIN (no core padding), so use that for
            // the gap-check so we can slot into gaps between prior bookings.
            for i in 0..cs {
                if is_on_grid(path[i].cell) {
                    times[i] = self.earliest_free(
                        (path[i].cell.x, path[i].cell.y),
                        times[i],
                        APPROACH_RESERVATION_TICKS,
                    );
                }
                times[i + 1] = times[i] + MIN_TICKS_PER_CELL;
            }

            // Phase 2 — core traversal.
            //
            // Gap-check duration for core cells must include the MIN_TICKS_PER_CELL
            // padding that is added to their reservation end at registration, so
            // our new window [t, t + actual_end + MIN] never overlaps the next
            // car's [r.start, …].
            //
            // Variable speed: enter at T_approach, each step ≥ MIN_TICKS_PER_CELL.
            // If any step would exceed MAX_TICKS_PER_CELL (looks like stalling) →
            // fall back to a fixed-speed sprint.
            // Matches the worst-case registered window: core cell i is held until
            // times[i + RESERVATION_LOOKAHEAD] + CORE_RESERVATION_PADDING.
            let var_dur = RESERVATION_LOOKAHEAD as u64 * MAX_TICKS_PER_CELL;

            let mut var_core = vec![0u64; n_core];
            var_core[0] =
                self.earliest_free((path[cs].cell.x, path[cs].cell.y), times[cs], var_dur);
            let mut var_ok = var_core[0] - times[cs] <= MAX_TICKS_PER_CELL;
            if var_ok {
                for k in 0..n_core - 1 {
                    let key = (path[cs + k + 1].cell.x, path[cs + k + 1].cell.y);
                    let next_t = self.earliest_free(key, var_core[k] + MIN_TICKS_PER_CELL, var_dur);
                    if next_t - var_core[k] > MAX_TICKS_PER_CELL {
                        var_ok = false;
                        break;
                    }
                    var_core[k + 1] = next_t;
                }
            }

            if var_ok {
                for k in 0..n_core {
                    times[cs + k] = var_core[k];
                }
                println!("Variable speed (entered at t={})", var_core[0]);
            } else {
                // Sprint fallback: wait outside the core until all cells can be
                // crossed at MIN speed without any mid-intersection gaps.
                let entry = self.sprint_entry(cs, n_core, &path, times[cs], MIN_TICKS_PER_CELL);

                let wait = entry - times[cs];
                println!(
                    "Sprint (entered at t={}, waited {} ticks outside core)",
                    entry, wait
                );

                times[cs] = entry;
                for k in 1..n_core {
                    times[cs + k] = entry + k as u64 * MIN_TICKS_PER_CELL;
                }
            }

            // Phase 3 — exit at MAX speed.
            for i in (cs + n_core - 1)..n - 1 {
                times[i + 1] = times[i] + MIN_TICKS_PER_CELL;
            }

            // Phase 4 — smooth approach.
            smooth_times(&mut times[..=cs]);
        } else {
            for i in 0..n {
                if is_on_grid(path[i].cell) {
                    times[i] = self.earliest_free(
                        (path[i].cell.x, path[i].cell.y),
                        times[i],
                        APPROACH_RESERVATION_TICKS,
                    );
                }
                if i + 1 < n {
                    times[i + 1] = times[i] + MIN_TICKS_PER_CELL;
                }
            }
            smooth_times(&mut times);
        }

        // Register: cell i is held until times[i+2]; core cells get extra padding.
        for i in 0..n {
            if !is_on_grid(path[i].cell) {
                continue;
            }
            let key = (path[i].cell.x, path[i].cell.y);
            let end = times[i + RESERVATION_LOOKAHEAD];
            let bucket = self.cells.entry(key).or_default();
            let pos = bucket.partition_point(|r| r.start <= times[i]);
            bucket.insert(
                pos,
                Reservation {
                    start: times[i],
                    end,
                },
            );
        }

        let mut schedule: Vec<_> = times.into_iter().zip(path).collect();
        schedule.truncate(schedule.len() - (RESERVATION_LOOKAHEAD - 1));
        schedule
    }

    /// Earliest sprint entry for `n_core` consecutive cells starting at
    /// `path[cs]`, traversed at a uniform `speed` ticks/cell.
    ///
    /// Uses a fixpoint loop (capped at `n_core` iterations) because pushing
    /// the entry later for cell k may create a new conflict at an earlier cell.
    fn sprint_entry(
        &self,
        cs: usize,
        n_core: usize,
        path: &[CellPosition],
        t_approach: Timestamp,
        speed: u64,
    ) -> Timestamp {
        // Matches the worst-case registered window: core cell i is held until
        // times[i + RESERVATION_LOOKAHEAD] + CORE_RESERVATION_PADDING.
        let dur = RESERVATION_LOOKAHEAD as u64 * speed;
        let mut entry = t_approach;
        for _ in 0..n_core {
            let prev = entry;
            for k in 0..n_core {
                let cell = path[cs + k].cell;
                if !is_on_grid(cell) {
                    continue;
                }
                let t_at_k = entry + k as u64 * speed;
                let actual = self.earliest_free((cell.x, cell.y), t_at_k, dur);
                if actual > t_at_k {
                    let needed = actual - k as u64 * speed;
                    if needed > entry {
                        entry = needed;
                    }
                }
            }
            if entry == prev {
                break;
            }
        }
        entry
    }

    /// Earliest tick ≥ `min_t` at which a window of `duration` ticks fits at
    /// `key` without overlapping any existing reservation.
    fn earliest_free(&self, key: (i32, i32), min_t: Timestamp, duration: u64) -> Timestamp {
        let bucket = match self.cells.get(&key) {
            Some(b) => b,
            None => return min_t,
        };
        let mut candidate = min_t;
        for r in bucket {
            if candidate.saturating_add(duration) <= r.start {
                break;
            }
            if candidate >= r.end {
                continue;
            }
            candidate = r.end;
        }
        candidate
    }

    fn prune(&mut self, now: Timestamp) {
        for bucket in self.cells.values_mut() {
            bucket.retain(|r| r.end > now);
        }
    }
}

impl Default for Book {
    fn default() -> Self {
        Self::new()
    }
}

// ── Helpers ──────────────────────────────────────────────────────────────────

fn smooth_times(times: &mut [u64]) {
    let n = times.len();
    let mut seg_start = 0;
    for i in 1..n {
        let gap = times[i].saturating_sub(times[i - 1]);
        if gap > MIN_TICKS_PER_CELL {
            let steps = (i - seg_start) as u64;
            if steps > 1 {
                let start_t = times[seg_start];
                let end_t = times[i];
                for j in (seg_start + 1)..i {
                    let frac = (j - seg_start) as u64;
                    times[j] = start_t + frac * (end_t - start_t) / steps;
                }
            }
            seg_start = i;
        }
    }
}

fn is_on_grid(cell: Cell) -> bool {
    let gs = GRID_SIZE as i32;
    cell.x >= 0 && cell.x < gs && cell.y >= 0 && cell.y < gs
}

fn is_core_cell(cell: Cell) -> bool {
    let rs = ROAD_START as i32;
    let re = ROAD_END as i32;
    cell.x >= rs && cell.x <= re && cell.y >= rs && cell.y <= re
}

fn cp(x: i32, y: i32, dir: Direction) -> CellPosition {
    CellPosition {
        cell: Cell { x, y },
        dir,
    }
}

fn south_cells(x: i32, y_from: i32, y_to: i32) -> Vec<CellPosition> {
    (y_from..=y_to)
        .map(|y| cp(x, y, Direction::South))
        .collect()
}

fn north_cells(x: i32, y_from: i32, y_to: i32) -> Vec<CellPosition> {
    (y_to..=y_from)
        .rev()
        .map(|y| cp(x, y, Direction::North))
        .collect()
}

fn east_cells(y: i32, x_from: i32, x_to: i32) -> Vec<CellPosition> {
    (x_from..=x_to).map(|x| cp(x, y, Direction::East)).collect()
}

fn west_cells(y: i32, x_from: i32, x_to: i32) -> Vec<CellPosition> {
    (x_to..=x_from)
        .rev()
        .map(|x| cp(x, y, Direction::West))
        .collect()
}

fn exit_cells(x: i32, y: i32, dir: Direction) -> Vec<CellPosition> {
    let (dx, dy) = match dir {
        Direction::North => (0, -1),
        Direction::South => (0, 1),
        Direction::East => (1, 0),
        Direction::West => (-1, 0),
    };
    (0..RESERVATION_LOOKAHEAD)
        .map(|k| cp(x + k as i32 * dx, y + k as i32 * dy, dir))
        .collect()
}

// ── Route → cell path ────────────────────────────────────────────────────────

fn path_for_route(route: &Routes) -> Vec<CellPosition> {
    use Direction::*;
    use Routes::*;

    let gs = GRID_SIZE as i32;
    let rs = ROAD_START as i32;

    let nr = rs; // 6
    let ns = rs + 1; // 7
    let nl = rs + 2; // 8
    let sl = rs + 3; // 9
    let ss = rs + 4; // 10
    let sr = rs + 5; // 11

    match route {
        NorthS => {
            let mut p = south_cells(ns, 0, gs - 1);
            p.extend(exit_cells(ns, gs, South));
            p
        }
        NorthR => {
            let mut p = south_cells(nr, 0, nr - 1);
            p.push(cp(nr, nr, West));
            p.extend(west_cells(nr, nr - 1, 0));
            p.extend(exit_cells(-1, nr, West));
            p
        }
        NorthL => {
            let mut p = south_cells(nl, 0, nl);
            p.push(cp(nl, sl, East));
            p.extend(east_cells(sl, sl, gs - 1));
            p.extend(exit_cells(gs, sl, East));
            p
        }
        SouthS => {
            let mut p = north_cells(ss, gs - 1, 0);
            p.extend(exit_cells(ss, -1, North));
            p
        }
        SouthR => {
            let mut p = north_cells(sr, gs - 1, sr + 1);
            p.push(cp(sr, sr, East));
            p.extend(east_cells(sr, sr + 1, gs - 1));
            p.extend(exit_cells(gs, sr, East));
            p
        }
        SouthL => {
            let mut p = north_cells(sl, gs - 1, nl + 1);
            p.push(cp(sl, nl, West));
            p.extend(west_cells(nl, nl, 0));
            p.extend(exit_cells(-1, nl, West));
            p
        }
        EastS => {
            let mut p = west_cells(ns, gs - 1, 0);
            p.extend(exit_cells(-1, ns, West));
            p
        }
        EastR => {
            let mut p = west_cells(nr, gs - 1, sr + 1);
            p.push(cp(sr, nr, North));
            p.extend(north_cells(sr, nr - 1, 0));
            p.extend(exit_cells(sr, -1, North));
            p
        }
        EastL => {
            let mut p = west_cells(nl, gs - 1, nl + 1);
            p.push(cp(nl, nl, South));
            p.extend(south_cells(nl, nl + 1, gs - 1));
            p.extend(exit_cells(nl, gs, South));
            p
        }
        WestS => {
            let mut p = east_cells(ss, 0, gs - 1);
            p.extend(exit_cells(gs, ss, East));
            p
        }
        WestR => {
            let mut p = east_cells(sr, 0, nr - 1);
            p.push(cp(nr, sr, South));
            p.extend(south_cells(nr, sr + 1, gs - 1));
            p.extend(exit_cells(nr, gs, South));
            p
        }
        WestL => {
            let mut p = east_cells(sl, 0, nl);
            p.push(cp(sl, sl, North));
            p.extend(north_cells(sl, nl, 0));
            p.extend(exit_cells(sl, -1, North));
            p
        }
    }
}

// ── Tests ─────────────────────────────────────────────────────────────────────

#[cfg(test)]
mod book_test {
    use super::*;

    #[test]
    fn paths_have_exit_cells() {
        let routes = [
            Routes::NorthS,
            Routes::NorthR,
            Routes::NorthL,
            Routes::SouthS,
            Routes::SouthR,
            Routes::SouthL,
            Routes::EastS,
            Routes::EastR,
            Routes::EastL,
            Routes::WestS,
            Routes::WestR,
            Routes::WestL,
        ];
        for route in &routes {
            let path = path_for_route(route);
            assert!(path.len() >= RESERVATION_LOOKAHEAD);
            for k in 1..=RESERVATION_LOOKAHEAD {
                assert!(
                    !is_on_grid(path[path.len() - k].cell),
                    "cell -{} must be off-grid",
                    k
                );
            }
            for cp in &path[..path.len() - RESERVATION_LOOKAHEAD] {
                assert!(
                    is_on_grid(cp.cell),
                    "({},{}) unexpectedly off-grid",
                    cp.cell.x,
                    cp.cell.y
                );
            }
        }
    }

    #[test]
    fn no_collision_same_route() {
        let mut book = Book::new();
        let sched_a = book.book(path_for_route(&Routes::NorthS), 0);
        let sched_b = book.book(path_for_route(&Routes::NorthS), 0);
        for (ta, pa) in &sched_a {
            if !is_on_grid(pa.cell) {
                continue;
            }
            for (tb, pb) in &sched_b {
                if pa.cell == pb.cell {
                    assert_ne!(
                        ta, tb,
                        "collision at ({},{}) t={}",
                        pa.cell.x, pa.cell.y, ta
                    );
                }
            }
        }
    }

    #[test]
    fn no_collision_crossing_routes() {
        let mut book = Book::new();
        let sched_a = book.book(path_for_route(&Routes::NorthS), 0);
        let sched_b = book.book(path_for_route(&Routes::EastS), 0);
        for (ta, pa) in &sched_a {
            if !is_on_grid(pa.cell) {
                continue;
            }
            for (tb, pb) in &sched_b {
                if pa.cell == pb.cell {
                    assert_ne!(
                        ta, tb,
                        "crossing collision at ({},{}) t={}",
                        pa.cell.x, pa.cell.y, ta
                    );
                }
            }
        }
    }

    #[test]
    fn times_monotone() {
        let mut book = Book::new();
        let sched = book.book(path_for_route(&Routes::WestL), 0);
        for w in sched.windows(2) {
            assert!(w[1].0 > w[0].0, "timestamps not strictly increasing");
        }
    }

    // Every registered reservation must end no earlier than times[i+RESERVATION_LOOKAHEAD].
    #[test]
    fn reservations_cover_lookahead() {
        let mut book = Book::new();
        let sched = book.book(path_for_route(&Routes::NorthS), 0);
        for i in 0..sched.len() {
            let (t, cp) = &sched[i];
            if !is_on_grid(cp.cell) || i + RESERVATION_LOOKAHEAD >= sched.len() {
                continue;
            }
            let key = (cp.cell.x, cp.cell.y);
            let res = book.cells[&key]
                .iter()
                .find(|r| r.start == *t)
                .expect("reservation missing");
            assert!(
                res.end >= sched[i + RESERVATION_LOOKAHEAD].0,
                "cell ({},{}) res.end={} < times[i+{}]={}",
                cp.cell.x,
                cp.cell.y,
                res.end,
                RESERVATION_LOOKAHEAD,
                sched[i + RESERVATION_LOOKAHEAD].0
            );
        }
    }

    // A car that fits in a gap before an existing reservation should slot in
    // immediately rather than waiting until after it ends.
    #[test]
    fn gap_reuse() {
        let mut book = Book::new();
        let key = (7i32, 3i32); // non-core approach cell

        // Existing reservation [60, 100].
        book.cells.entry(key).or_default().push(Reservation {
            start: 60,
            end: 100,
        });

        // Duration 40 (= 2×MIN): 10 + 40 = 50 ≤ 60 → fits before the reservation.
        assert_eq!(
            book.earliest_free(key, 10, 40),
            10,
            "should slot into the gap before the reservation"
        );

        // Duration 60: 10 + 60 = 70 > 60 → overlaps, must wait until 100.
        assert_eq!(
            book.earliest_free(key, 10, 60),
            100,
            "should wait until after the reservation"
        );
    }
}
