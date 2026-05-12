# book.rs — Line-by-Line Explanation

This file implements the **intersection booking system**: the scheduler that decides
when each car is allowed to enter each cell, preventing collisions without traffic lights.

---

## Imports (lines 1–4)

```rust
use std::collections::HashMap;
```
Brings `HashMap` into scope. A `HashMap<K, V>` is a dictionary: given a key of type `K`
it gives you a value of type `V` in O(1) average time.

```rust
use crate::config::{MIN_TICKS_PER_CELL, ...};
```
`crate` means "this project". This line imports named constants from `config.rs`.

```rust
use crate::models::*;
```
The `*` imports everything that is `pub` in `models.rs` — structs like `Cell`,
`CellPosition`, `Direction`, and type aliases like `Timestamp`.

---

## `Reservation` struct (lines 6–10)

```rust
struct Reservation {
    start: Timestamp,
    end:   Timestamp,
}
```
A `struct` is a named bundle of data, like a record or class with only fields and no
methods (yet). `Timestamp` is a type alias for `u64` (an unsigned 64-bit integer).

A `Reservation` records a time interval `[start, end)` during which one car has
claimed a cell. Any other car wanting that cell must wait until `end`.

This struct has no `pub` — it is private to this file.

---

## `Book` struct (lines 12–18)

```rust
pub struct Book {
    cells: HashMap<(i32, i32), Vec<Reservation>>,
}
```
`pub` makes this type visible to other files. `Book` owns a single field `cells`.

The key `(i32, i32)` is a tuple of two signed 32-bit integers representing a cell's
`(x, y)` grid coordinates. The value `Vec<Reservation>` is a growable list of
reservations for that cell, kept in chronological order by `start`.

The booking table therefore looks like:
```
(7, 3)  → [ Reservation{0..40}, Reservation{80..120}, ... ]
(8, 3)  → [ Reservation{20..60}, ... ]
...
```

---

## `impl Book` block (lines 22–220)

`impl Book` is where methods that belong to `Book` are written.

### `new` (line 23)

```rust
pub fn new() -> Self { Book { cells: HashMap::new() } }
```
A constructor. `Self` is shorthand for `Book`. Returns a `Book` with an empty map.

---

### `can_spawn` (lines 27–37)

```rust
pub fn can_spawn(&self, route: &Routes, now: Timestamp) -> bool {
```
`&self` means this method borrows `Book` immutably — it reads but does not change it.
`&Routes` borrows the route enum without taking ownership.

```rust
    let path = path_for_route(route);
```
Computes the full ordered list of `CellPosition`s for this route (defined later in the file).

```rust
    for cp in path.iter().take(FREE_CELLS_FOR_SPAWNING) {
```
Iterates over only the first `FREE_CELLS_FOR_SPAWNING` (= 2) cells of the path.
`iter()` creates a borrowing iterator; `.take(n)` stops it after `n` items.

```rust
        if !is_on_grid(cp.cell) { break; }
```
If the cell is off the 18×18 grid (it is a virtual off-screen start cell), stop checking.

```rust
        let key = (cp.cell.x, cp.cell.y);
        if self.earliest_free(key, now, APPROACH_RESERVATION_TICKS) != now {
            return false;
        }
```
`earliest_free` returns the earliest tick at which a window of duration
`APPROACH_RESERVATION_TICKS` fits without overlapping existing reservations.
If that tick is not `now`, the cell is occupied right now, so spawning is blocked.

```rust
    true
```
All checked cells are free: a vehicle can be spawned.

---

### `book` (lines 39–159)

This is the main scheduling method. It takes a route and the current tick, and returns a
`Vec<(Timestamp, CellPosition)>` — the schedule: a list of (tick, cell) pairs telling
the vehicle exactly when to be in each cell.

```rust
    let path = path_for_route(route);
    let n = path.len();
    let mut times = vec![0u64; n];
    times[0] = now;
```
`vec![0u64; n]` creates a vector of `n` zeros of type `u64`. `times[i]` will hold the
tick at which the vehicle arrives at `path[i]`. The vehicle starts at `path[0]` at `now`.

```rust
    self.prune(now);
```
Removes all reservations that ended before `now`, keeping the map tidy.

```rust
    let core_start = path.iter().position(|cp| is_on_grid(cp.cell) && is_core_cell(cp.cell));
    let core_end   = path.iter().rposition(|cp| is_on_grid(cp.cell) && is_core_cell(cp.cell));
```
`position` finds the index of the **first** element matching the closure (a function written
inline with `|args| body` syntax). `rposition` finds the **last**. Together they give the
range of the "core" — the 6×6 road block where cars from different directions can cross.

Both return `Option<usize>`: either `Some(index)` if found, or `None` if not.

```rust
    if let (Some(cs), Some(ce)) = (core_start, core_end) {
```
`if let` destructures the tuple: if both are `Some`, bind `cs` ("core start index") and
`ce` ("core end index") and enter the block. If either is `None` (route does not cross the
core — currently unused but safe to handle), fall through to the `else` branch.

```rust
        let n_core = ce - cs + 1;
```
`n_core` = number of cells in the core segment. The three speed tiers come directly from
config: `MIN_TICKS_PER_CELL` (fastest), `MID_TICKS_PER_CELL`, `MAX_TICKS_PER_CELL` (slowest).

---

#### Phase 1 — Approach (lines 61–70)

```rust
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
```
For every cell before the core, push the vehicle's arrival as early as possible — but no
earlier than the earliest free window. Then set the next cell's arrival to be exactly
`MIN_TICKS_PER_CELL` (the fastest traversal) after that.

---

#### Phase 2 — Core traversal (lines 72–119)

```rust
        let var_dur = RESERVATION_LOOKAHEAD as u64 * MAX_TICKS_PER_CELL;
```
The gap-check window must match the worst-case registered window for a core cell.
Cell `i` is held until `times[i + RESERVATION_LOOKAHEAD]`, which at maximum slow speed is
at most `t + RESERVATION_LOOKAHEAD * MAX_TICKS_PER_CELL`. Using the same formula here
ensures the gap-check never passes for a slot too narrow to hold the actual registration.

```rust
        let mut var_core = vec![0u64; n_core];
        var_core[0] = self.earliest_free(
            (path[cs].cell.x, path[cs].cell.y),
            times[cs],
            var_dur,
        );
        let mut var_ok = var_core[0] - times[cs] <= MAX_TICKS_PER_CELL;
```
`times[cs]` is the earliest tick the car could physically reach the core entry (set during
Phase 1). `var_core[0]` is what `earliest_free` returned — the actual earliest tick the
car is *allowed* to enter given existing reservations. Their difference is the forced wait
at the core boundary.

`var_ok` is `true` when that wait is at most `MAX_TICKS_PER_CELL` ticks. The same threshold is
used for inter-cell gaps inside the core (line below), so this check keeps the entry
consistent with the rest of the variable-speed attempt: if the car would have to stall
for longer than `MAX_TICKS_PER_CELL` just to get in, the whole variable-speed strategy is already
over budget and the sprint fallback is used instead.

`var_core[0] >= times[cs]` is always guaranteed because `earliest_free` never returns
a tick earlier than its `min_t` argument, so the subtraction cannot underflow on `u64`.

```rust
        if var_ok {
            for k in 0..n_core - 1 {
                let key    = (path[cs + k + 1].cell.x, path[cs + k + 1].cell.y);
                let next_t = self.earliest_free(key, var_core[k] + MIN_TICKS_PER_CELL, var_dur);
                if next_t - var_core[k] > MAX_TICKS_PER_CELL { var_ok = false; break; }
                var_core[k + 1] = next_t;
            }
        }
```
For each subsequent core cell, find the earliest free slot reachable from the previous
cell. If any gap forces the vehicle to wait longer than `MAX_TICKS_PER_CELL`, the variable-speed
attempt fails and `var_ok` becomes `false`.

```rust
        if var_ok {
            for k in 0..n_core { times[cs + k] = var_core[k]; }
        } else {
```
If variable speed worked, copy its computed times into the main `times` array.

```rust
            let entry = self.sprint_entry(cs, n_core, &path, times[cs], MIN_TICKS_PER_CELL);
```

**Sprint fallback**: variable speed failed, so the car waits outside the core until all
cells can be crossed at `MIN_TICKS_PER_CELL` with no mid-intersection gaps. `sprint_entry`
returns the earliest entry tick that satisfies this constraint (see `sprint_entry` below).

```rust
            times[cs] = entry;
            for k in 1..n_core {
                times[cs + k] = entry + k as u64 * MIN_TICKS_PER_CELL;
            }
```
Fill in the core times: each step is exactly `MIN_TICKS_PER_CELL` ticks apart.

---

#### Phase 3 — Exit (lines 122–124)

```rust
        for i in (cs + n_core - 1)..n - 1 {
            times[i + 1] = times[i] + MIN_TICKS_PER_CELL;
        }
```
After leaving the core the vehicle drives at maximum speed (one cell per `MIN_TICKS_PER_CELL`).

---

#### Phase 4 — Smooth approach (line 127)

```rust
        smooth_times(&mut times[..=cs]);
```
`&mut times[..=cs]` is a mutable slice of `times` from index 0 through `cs` (inclusive).
`smooth_times` redistributes approach timestamps so a car that had to wait for the core
decelerates gradually rather than pausing abruptly at one cell (see `smooth_times` below).

---

#### No-core fallback (lines 128–140)

```rust
    } else {
        for i in 0..n {
            if is_on_grid(path[i].cell) {
                times[i] = self.earliest_free(
                    (path[i].cell.x, path[i].cell.y),
                    times[i],
                    APPROACH_RESERVATION_TICKS,
                );
            }
            if i + 1 < n { times[i + 1] = times[i] + MIN_TICKS_PER_CELL; }
        }
        smooth_times(&mut times);
    }
```
Routes that do not cross the core (not currently used) are scheduled with simple earliest-
free logic and then smoothed.

The three arguments to `earliest_free` are:

- `(path[i].cell.x, path[i].cell.y)` — the key: the `(x, y)` coordinate of the current
  cell, used to look up its reservation list in the booking table.
- `times[i]` — the minimum start tick: the earliest the vehicle could arrive at this cell
  given the times already assigned to previous cells. `earliest_free` will return this
  value unchanged if the cell is free, or a later tick if it is blocked.
- `APPROACH_RESERVATION_TICKS` — the window duration: how wide a free gap must be for
  the car to fit. This matches the size of the reservation that will be registered for
  this cell (`times[i]` to `times[i + RESERVATION_LOOKAHEAD]` at minimum speed), so the
  gap check and the registration are always consistent.

---

#### Registration (lines 142–154)

```rust
    for i in 0..n {
        if !is_on_grid(path[i].cell) { continue; }
```
Off-grid cells (the two phantom exit cells) are not registered in the booking table.

```rust
        let key = (path[i].cell.x, path[i].cell.y);
        let end = times[i + RESERVATION_LOOKAHEAD];
```
Cell `i` is held until `times[i + RESERVATION_LOOKAHEAD]` — the tick when the vehicle
will be `RESERVATION_LOOKAHEAD` cells ahead. All cells use the same formula; there is no
distinction between core and approach cells at registration time.

```rust
        let bucket = self.cells.entry(key).or_default();
```
`entry(key).or_default()` returns a mutable reference to the `Vec<Reservation>` for this
key, creating an empty one if it does not exist yet.

```rust
        let pos = bucket.partition_point(|r| r.start <= times[i]);
        bucket.insert(pos, Reservation { start: times[i], end });
```
`partition_point` is a binary search that returns the index where the predicate stops being
true — i.e., the insertion point that keeps the list sorted by `start`.
`bucket.insert(pos, ...)` shifts everything from `pos` onward right by one and places the
new reservation at `pos`.

---

#### Return value (lines 156–158)

```rust
    let mut schedule: Vec<_> = times.into_iter().zip(path).collect();
    schedule.truncate(schedule.len() - (RESERVATION_LOOKAHEAD - 1));
    schedule
```
`times.into_iter().zip(path)` pairs each tick with its corresponding `CellPosition`.
`collect()` assembles the pairs into a `Vec`. `Vec<_>` lets Rust infer the element type.

`.truncate(...)` removes `RESERVATION_LOOKAHEAD - 1` trailing elements, keeping exactly
one off-grid exit cell regardless of the lookahead value. That cell's tick is when the
vehicle has fully left the canvas.

The last expression without a semicolon is the return value of the function.

---

### `sprint_entry` (lines 166–193)

```rust
fn sprint_entry(&self, cs, n_core, path, t_approach, speed) -> Timestamp {
    let dur = RESERVATION_LOOKAHEAD as u64 * speed;
    let mut entry = t_approach;
    for _ in 0..n_core {
```
`speed` is always `MIN_TICKS_PER_CELL`. `dur` matches the registration window so gap
checks and registrations stay in sync. `for _ in 0..n_core` — `_` means we don't
care about the loop counter, only the number of iterations.

```rust
        let prev = entry;
        for k in 0..n_core {
            let t_at_k = entry + k as u64 * speed;
            let actual = self.earliest_free((cell.x, cell.y), t_at_k, dur);
            if actual > t_at_k {
                let needed = actual - k as u64 * speed;
                if needed > entry { entry = needed; }
            }
        }
        if entry == prev { break; }
    }
    entry
}
```
This is a **fixpoint loop**: keep adjusting `entry` until no core cell forces a later entry.

For each core cell at index `k` (0 = first core cell, 1 = second, …), the vehicle would
arrive at `entry + k * speed`. If the earliest free slot at that cell is later than that
(`actual > t_at_k`), back-calculate how much earlier the vehicle would need to enter
(`needed = actual - k * speed`) and update `entry`.

We repeat because bumping entry for cell `k` might create a new conflict at cell `k-1`.
Once a full pass produces no change (`entry == prev`), we stop early.

---

### `earliest_free` (lines 197–209)

```rust
fn earliest_free(&self, key: (i32, i32), min_t: Timestamp, duration: u64) -> Timestamp {
    let bucket = match self.cells.get(&key) {
        Some(b) => b,
        None    => return min_t,
    };
```
`match` is Rust's pattern-matching switch. If no reservations exist for this cell, return
`min_t` immediately (the cell is free from the start).

```rust
    let mut candidate = min_t;
    for r in bucket {
        if candidate.saturating_add(duration) <= r.start { break; }
        if candidate >= r.end { continue; }
        candidate = r.end;
    }
    candidate
}
```
Walk reservations in order (they are sorted by `start`). At each step:
- `saturating_add` is like `+` but clamps to `u64::MAX` instead of overflowing.
- If `[candidate, candidate + duration)` ends before `r.start`, the window fits in the
  gap before this reservation — we can stop.
- If `candidate` is already past `r.end`, this reservation doesn't block us — skip it.
- Otherwise the window overlaps the reservation: push `candidate` to `r.end` and try again.

---

### `prune` (lines 211–215)

```rust
fn prune(&mut self, now: Timestamp) {
    for bucket in self.cells.values_mut() {
        bucket.retain(|r| r.end > now);
    }
}
```
`values_mut()` iterates over all `Vec<Reservation>` values mutably.
`retain` keeps only elements where the closure returns `true`, removing the rest in-place.
This discards expired reservations so the map does not grow without bound.

---

### `Default` impl (lines 218–220)

```rust
impl Default for Book {
    fn default() -> Self { Self::new() }
}
```
`Default` is a standard Rust trait (interface). Implementing it lets code call
`Book::default()` or use `Book` in contexts that require a default value (e.g., `#[derive(Default)]`
on structs that contain a `Book`).

---

## Helper functions (lines 222–283)

### `smooth_times` (lines 224–242)

```rust
fn smooth_times(times: &mut [u64]) {
```
Takes a mutable slice (`[u64]`) — a view into part of a `Vec` (or the whole thing).

The algorithm walks forward through the approach timestamps looking for "wait points" — places where a car had to slow down because the core or a cell ahead was occupied. When it finds one, it linearly redistributes all the timestamps before it so the car decelerates smoothly from the last anchor rather than arriving at full speed and stopping abruptly.

```rust
    let mut seg_start = 0;
    for i in 1..n {
        let gap = times[i].saturating_sub(times[i - 1]);
        if gap > MIN_TICKS_PER_CELL {
```
`seg_start` is the index of the last anchor — either the very start (index 0) or the last cell where a wait was detected. It marks the beginning of the current segment being processed.

`saturating_sub` subtracts but clamps at 0 (avoids underflow on unsigned integers). A gap larger than `MIN_TICKS_PER_CELL` means the car had to wait at cell `i` — something was blocking it. Cell `i` becomes the end of the current segment.

```rust
            let steps = (i - seg_start) as u64;
            if steps > 1 {
                let start_t = times[seg_start];
                let end_t   = times[i];
                for j in (seg_start + 1)..i {
                    let frac = (j - seg_start) as u64;
                    times[j] = start_t + frac * (end_t - start_t) / steps;
                }
            }
            seg_start = i;
```
`steps` is the number of cells in the segment from `seg_start` to `i` (exclusive). `frac` is how far cell `j` is from `seg_start` within that segment — it goes from 1 up to `steps - 1`. The formula `start_t + frac * (end_t - start_t) / steps` is standard linear interpolation: it spaces the intermediate timestamps evenly between `start_t` and `end_t`.

After redistributing, `seg_start` advances to `i` — the wait point becomes the new anchor for the next segment.

Example: approach times before smoothing are `[0, 20, 40, 60, 120]` (the car rushed to the core entry and waited at index 4). One segment: `seg_start=0`, `i=4`, `steps=4`. Interpolation sets indices 1, 2, 3 to `30, 60, 90`. Result: `[0, 30, 60, 90, 120]` — the car decelerates uniformly all the way to the core.

---

### `is_on_grid` / `is_core_cell` (lines 244–253)

```rust
fn is_on_grid(cell: Cell) -> bool {
    let gs = GRID_SIZE as i32;
    cell.x >= 0 && cell.x < gs && cell.y >= 0 && cell.y < gs
}
```
`as i32` casts the `u32` constant to `i32` so it can be compared to the signed coordinates.

```rust
fn is_core_cell(cell: Cell) -> bool {
    let rs = ROAD_START as i32;
    let re = ROAD_END as i32;
    cell.x >= rs && cell.x <= re && cell.y >= rs && cell.y <= re
}
```
"Core" = the 6×6 block where all four roads overlap. Both x and y must be within
`[ROAD_START, ROAD_END]`.

---

### Cell-building helpers (lines 255–283)

```rust
fn cp(x: i32, y: i32, dir: Direction) -> CellPosition {
    CellPosition { cell: Cell { x, y }, dir }
}
```
Shorthand constructor. `Cell { x, y }` uses Rust's field-init shorthand: when the local
variable name matches the field name, you can write just `x` instead of `x: x`.

```rust
fn south_cells(x: i32, y_from: i32, y_to: i32) -> Vec<CellPosition> {
    (y_from..=y_to).map(|y| cp(x, y, Direction::South)).collect()
}
```
`(y_from..=y_to)` is an inclusive range. `.map(|y| ...)` transforms each `y` into a
`CellPosition`. `.collect()` consumes the iterator and assembles a `Vec`.

`north_cells`, `east_cells`, `west_cells` follow the same pattern, adjusting the axis and
direction. `north_cells` adds `.rev()` to iterate the range in reverse (decreasing y).

```rust
fn exit_cells(x: i32, y: i32, dir: Direction) -> Vec<CellPosition> {
    let (dx, dy) = match dir { ... };
    (0..RESERVATION_LOOKAHEAD).map(|k| cp(x + k*dx, y + k*dy, dir)).collect()
}
```
Returns `RESERVATION_LOOKAHEAD` off-grid cells stepping away from the grid edge. These
are never booked — they exist only so `times[i + RESERVATION_LOOKAHEAD]` stays in bounds
for the last on-grid cell.

---

## `path_for_route` (lines 287–355)

```rust
fn path_for_route(route: &Routes) -> Vec<CellPosition> {
    use Direction::*;
    use Routes::*;
```
`use X::*` inside a function brings the variants into scope locally, so you can write
`South` instead of `Direction::South`.

```rust
    let gs = GRID_SIZE as i32;   // 18
    let rs = ROAD_START as i32;  // 6

    let nr = rs;       // 6  — north-entry / right-turn column
    let ns = rs + 1;   // 7  — straight-through column
    let nl = rs + 2;   // 8  — left-turn column
    let sl = rs + 3;   // 9  — south left-turn column
    let ss = rs + 4;   // 10 — straight-through column
    let sr = rs + 5;   // 11 — south-entry / right-turn column
```
These are the six road column/row indices. The intersection is 6 cells wide; each direction
gets a "right", "straight", and "left" lane.

```rust
    match route {
        NorthS => { let mut p = south_cells(ns, 0, gs-1); p.extend(exit_cells(ns, gs, South)); p }
```
A `match` on an enum must handle every variant (Rust enforces exhaustiveness).

`NorthS` (North → Straight → South): start at the top of column `ns` (y=0), go south
to the bottom (y=17), then add 2 exit cells off the bottom of the grid.

`p.extend(...)` appends all items from the iterator/array to `p` in-place.

The more complex routes (Right-turn, Left-turn) follow the same pattern: drive straight
into the intersection, push a corner cell with the new direction, then drive straight out.

---

## Tests (lines 359–457)

Each test function is marked `#[test]`. Running `cargo test` finds and executes them.

- **`paths_have_exit_cells`** — verifies every route ends with exactly `RESERVATION_LOOKAHEAD` off-grid cells and that all cells before them are on-grid.
- **`no_collision_same_route`** — books two cars on the same route at the same start tick
  and checks they are never in the same cell at the same tick.
- **`no_collision_crossing_routes`** — same check for two cars on perpendicular routes.
- **`times_monotone`** — checks that a booked schedule has strictly increasing timestamps.
- **`reservations_cover_lookahead`** — verifies the registration invariant: cell `i`'s reservation must extend at least until `times[i + RESERVATION_LOOKAHEAD]`.
- **`gap_reuse`** — unit test for `earliest_free`: a short window fits before an existing
  reservation; a longer window is pushed past it.
