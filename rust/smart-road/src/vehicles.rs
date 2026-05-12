use crate::book::*;
use crate::config::{CELL_SIZE, COLLISION_THRESHOLD_PX};
use crate::models::*;
use crate::stats::*;
use rand::seq::SliceRandom;
use uuid::Uuid;

pub struct Vehicles {
    pub veh: Vec<Vehicle>,
    pub stats: Stats,
}

#[derive(Clone, PartialEq)]
pub struct Vehicle {
    pub id: Uuid,
    pub positions: Vec<(Timestamp, CellPosition)>,
    pub highlight: Highlight,
}

impl Vehicles {
    pub fn new() -> Self {
        Vehicles {
            veh: Vec::new(),
            stats: Stats::new(),
        }
    }

    pub fn spawn_vehicle(&mut self, book: &mut Book, origin: Origin, now: Timestamp) {
        for route in shuffled_routes_array(origin) {
            if let Some(path) = book.can_spawn(&route, now) {
                self.veh.push(Vehicle {
                    id: Uuid::new_v4(),
                    positions: book.book(path, now),
                    highlight: Highlight::None,
                });
                break;
            }
        }
    }

    pub fn get_positions_and_update_stats(&mut self, now: Timestamp) -> ToRender {
        let mut result = Vec::new();
        let mut cell_data: Vec<(Uuid, Cell, Cell)> = Vec::new();

        self.veh.retain(|vehicle: &Vehicle| {
            if vehicle.positions.is_empty() {
                self.stats.clean_up_seen_pairs(vehicle.id);
                return false;
            }

            let (prev, next) = find_segment(&vehicle.positions, now);

            let render_pos = match (prev, next) {
                (Some((t0, p0)), Some((t1, p1))) => {
                    // update speed after each cell
                    self.stats.update_max_min_speed(t0, t1);
                    cell_data.push((vehicle.id, p0.cell, p1.cell));
                    interpolate_position(t0, p0, t1, p1, now)
                    // crate::new_interpolate::interpolate_position(t0, p0, p1, t1, now)
                }
                _ => {
                    // update total max min time,
                    // add to passed vehicles upon exit and
                    // clean up seen pairs
                    self.stats.update_max_min_time(&vehicle.positions);
                    self.stats.vehicle_finished();
                    self.stats.clean_up_seen_pairs(vehicle.id);
                    return false;
                }
            };

            result.push((vehicle.id, render_pos));
            true
        });

        let collisions = self
            .stats
            .detect_collisions(&result, COLLISION_THRESHOLD_PX);
        let close_calls = self.stats.detect_close_calls(&cell_data);

        for v in self.veh.iter_mut() {
            if close_calls.contains(&v.id) {
                v.highlight = Highlight::CloseCall;
            }
        }

        for (id, pos) in result.iter_mut() {
            if collisions.contains(id) {
                pos.highlight = Highlight::Collision;
            } else if close_calls.contains(id)
                || self
                    .veh
                    .iter()
                    .any(|v| v.id == *id && v.highlight == Highlight::CloseCall)
            {
                pos.highlight = Highlight::CloseCall;
            }
        }

        result
    }
}

// Returns a vehicle's current and next position in cell spacetime (cell and time of entry on cell).
fn find_segment(
    positions: &[(Timestamp, CellPosition)],
    now: Timestamp,
) -> (
    Option<(Timestamp, &CellPosition)>,
    Option<(Timestamp, &CellPosition)>,
) {
    match positions.binary_search_by_key(&now, |(t, _)| *t) {
        Ok(i) => {
            let curr = Some((positions[i].0, &positions[i].1));
            let next = positions.get(i + 1).map(|(t, p)| (*t, p));
            (curr, next)
        }
        Err(i) => {
            let prev = if i > 0 {
                Some((positions[i - 1].0, &positions[i - 1].1))
            } else {
                None
            };

            let next = positions.get(i).map(|(t, p)| (*t, p));

            (prev, next)
        }
    }
}

// Calculates the pixel position and angle in degrees of a vehicle
// based on the current and next cell (cell and time on cell) and current time (now).
fn interpolate_position(
    t0: Timestamp,
    p0: &CellPosition,
    t1: Timestamp,
    p1: &CellPosition,
    now: Timestamp,
) -> RenderPosition {
    if t1 == t0 {
        return cell_to_render(p0);
    }

    let total = (t1 - t0) as f32;
    let elapsed = (now - t0) as f32;

    let mut progress = elapsed / total;

    // clamp to [0,1]
    if progress < 0.0 {
        progress = 0.0;
    } else if progress > 1.0 {
        progress = 1.0;
    }

    let turn_delta = calculate_turn_delta(p0.dir, p1.dir, progress);

    let base_x = p0.cell.x as f32 * CELL_SIZE as f32;
    let base_y = p0.cell.y as f32 * CELL_SIZE as f32;

    let delta = CELL_SIZE as f32 * progress;

    let (x, y) = match p0.dir {
        Direction::North => (base_x, base_y - delta),
        Direction::South => (base_x, base_y + delta),
        Direction::East => (base_x + delta, base_y),
        Direction::West => (base_x - delta, base_y),
    };

    RenderPosition {
        angle: dir_to_angle(p0.dir) + turn_delta,
        x,
        y,
        highlight: Highlight::None,
    }
}

// Creates an array of three routes related to Origin in random order.
pub fn shuffled_routes_array(origin: Origin) -> [Routes; 3] {
    use crate::models::Routes::*;
    let mut routes = match origin {
        Origin::North => [NorthL, NorthR, NorthS],
        Origin::South => [SouthL, SouthR, SouthS],
        Origin::East => [EastL, EastR, EastS],
        Origin::West => [WestL, WestR, WestS],
    };

    routes.shuffle(&mut rand::thread_rng());
    routes
}

// Converts a CellPosition to a RenderPosition.
pub fn cell_to_render(p: &CellPosition) -> RenderPosition {
    RenderPosition {
        angle: dir_to_angle(p.dir),
        x: p.cell.x as f32 * CELL_SIZE as f32,
        y: p.cell.y as f32 * CELL_SIZE as f32,
        highlight: Highlight::None,
    }
}

// Converts Direction to angle degrees.
pub fn dir_to_angle(dir: Direction) -> f32 {
    match dir {
        Direction::North => 0.0,
        Direction::East => 90.0,
        Direction::South => 180.0,
        Direction::West => 270.0,
    }
}

// Sprite rotation angle in degrees
const TURN_ANGLE_STEP: f32 = 15.;
const TURN_STEPS: f32 = 90. / TURN_ANGLE_STEP - 1.0;

// Calculates the angle based on prev and next Direction and the progress in Cell rounded to
// the sprite TURN_ANGLE_STEP. The steps are
pub fn calculate_turn_delta(prev: Direction, next: Direction, progress: f32) -> f32 {
    if prev == next {
        return 0.;
    }

    let diff = (next.index() - prev.index()).rem_euclid(4);

    let d = match diff {
        0 => 0.,               // straight
        1 => TURN_ANGLE_STEP,  // right turn
        3 => -TURN_ANGLE_STEP, // left turn
        2 => 0.,               // opposite (U-turn) → treat as straight (or adjust if needed)
        _ => unreachable!(),
    };

    let stage = if progress <= 0.5 {
        0.0
    } else {
        let normalized = (progress - 0.5) / 0.5; // maps 0.5..1.0 → 0.0..1.0
        (normalized * TURN_STEPS).ceil()
    };

    (d * stage / TURN_ANGLE_STEP).round() * TURN_ANGLE_STEP
}

#[cfg(test)]
mod veh_test {
    use super::*;

    #[test]
    fn test_turn_delta() {
        assert_eq!(
            calculate_turn_delta(Direction::East, Direction::North, 0.3),
            0.
        ); // stage 0 — no turn yet
        assert_eq!(
            calculate_turn_delta(Direction::East, Direction::North, 0.55),
            -15.
        ); // stage 1
        assert_eq!(
            calculate_turn_delta(Direction::East, Direction::North, 0.65),
            -30.
        ); // stage 2
        assert_eq!(
            calculate_turn_delta(Direction::East, Direction::North, 0.75),
            -45.
        ); // stage 3
        assert_eq!(
            calculate_turn_delta(Direction::North, Direction::East, 0.55),
            15.
        ); // stage 1, right turn
        assert_eq!(
            calculate_turn_delta(Direction::South, Direction::West, 0.55),
            15.
        ); // stage 1, right turn
        assert_eq!(
            calculate_turn_delta(Direction::South, Direction::South, 0.55),
            0.
        ); // straight — always 0
    }

    #[test]
    fn test_straight() {
        let mut avs = Vehicles::new();
        let uuid = Uuid::new_v4();
        avs.veh.push(Vehicle {
            id: uuid,
            positions: vec![
                (
                    0,
                    CellPosition {
                        cell: Cell { x: 0, y: 0 },
                        dir: Direction::South,
                    },
                ),
                (
                    6,
                    CellPosition {
                        cell: Cell { x: 0, y: 1 },
                        dir: Direction::South,
                    },
                ),
                (
                    7,
                    CellPosition {
                        cell: Cell { x: 0, y: 2 },
                        dir: Direction::South,
                    },
                ),
            ],
            highlight: Highlight::None,
        });

        let res = avs.get_positions_and_update_stats(3);
        let expected: ToRender = vec![(
            uuid,
            RenderPosition {
                angle: dir_to_angle(Direction::South),
                x: 0.,
                y: 25.,
                highlight: Highlight::None,
            },
        )];
        assert_eq!(res[0], expected[0]);

        let res = avs.get_positions_and_update_stats(6);
        let expected: ToRender = vec![(
            uuid,
            RenderPosition {
                angle: dir_to_angle(Direction::South),
                x: 0.,
                y: 50.,
                highlight: Highlight::None,
            },
        )];
        assert_eq!(res, expected);

        let res = avs.get_positions_and_update_stats(7);
        assert!(res == Vec::new());
    }
}
