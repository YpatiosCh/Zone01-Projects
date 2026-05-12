use crate::models::*;
use std::collections::{HashMap, HashSet};
use uuid::Uuid;

pub struct Stats {
    pub vehicle_count: u32,
    pub max_speed: f32,
    pub min_speed: f32,
    pub max_time: u64,
    pub min_time: u64,
    pub close_calls: u32,
    pub collisions: u32,
    seen_pairs: HashSet<(Uuid, Uuid)>,
    seen_collision_pairs: HashSet<(Uuid, Uuid)>,
}

impl Stats {
    pub fn new() -> Self {
        Stats {
            vehicle_count: 0,
            max_speed: 0.,
            min_speed: std::f32::MAX,
            max_time: 0,
            min_time: std::u64::MAX,
            close_calls: 0,
            collisions: 0,
            seen_pairs: HashSet::new(),
            seen_collision_pairs: HashSet::new(),
        }
    }

    pub fn update_max_min_time(&mut self, positions: &[(u64, CellPosition)]) {
        if let (Some(first), Some(last)) = (positions.first(), positions.last()) {
            let total = last.0 - first.0;
            self.max_time = self.max_time.max(total);
            self.min_time = self.min_time.min(total);
        }
    }

    pub fn update_max_min_speed(&mut self, t0: u64, t1: u64) {
        let vel = 1000. / (t1 - t0) as f32;

        if vel > self.max_speed {
            self.max_speed = vel;
        }
        if vel < self.min_speed {
            self.min_speed = vel;
        }
    }

    pub fn vehicle_finished(&mut self) {
        self.vehicle_count += 1;
    }

    pub fn clean_up_seen_pairs(&mut self, removed: Uuid) {
        self.seen_pairs
            .retain(|(a, b)| *a != removed && *b != removed);
        self.seen_collision_pairs
            .retain(|(a, b)| *a != removed && *b != removed);
    }

    /// A collision is detected when two vehicles' rendered pixel positions are
    /// within `threshold` pixels of each other.
    pub fn detect_collisions(&mut self, render_data: &[(Uuid, RenderPosition)], threshold: f32) -> HashSet<Uuid> {
        let mut collisions = HashSet::new();
        for i in 0..render_data.len() {
            for j in (i + 1)..render_data.len() {
                let (id_a, pos_a) = &render_data[i];
                let (id_b, pos_b) = &render_data[j];
                let dx = pos_a.x - pos_b.x;
                let dy = pos_a.y - pos_b.y;
                if (dx * dx + dy * dy).sqrt() < threshold {
                    let pair = if id_a < id_b { (*id_a, *id_b) } else { (*id_b, *id_a) };
                    if !self.seen_collision_pairs.contains(&pair) {
                        self.seen_collision_pairs.insert(pair);
                        self.collisions += 1;
                        println!("Collision #{} at ({:.0},{:.0})", self.collisions, pos_a.x, pos_a.y);
                    }
                    collisions.insert(*id_a);
                    collisions.insert(*id_b);
                }
            }
        }
        collisions
    }

    /// A close call is detected when a vehicle's next scheduled cell is
    /// currently occupied by another vehicle.
    pub fn detect_close_calls(&mut self, cell_data: &[(Uuid, Cell, Cell)]) -> HashSet<Uuid> {
        let mut close_calls = HashSet::new();
        let mut occupancy: HashMap<Cell, Uuid> = HashMap::new();
        for (id, current, _) in cell_data {
            occupancy.insert(*current, *id);
        }

        for (id, current, next) in cell_data {
            if let Some(&occupant_id) = occupancy.get(next) {
                if occupant_id == *id {
                    continue;
                }
                let pair = if id < &occupant_id {
                    (*id, occupant_id)
                } else {
                    (occupant_id, *id)
                };
                if !self.seen_pairs.contains(&pair) {
                    self.seen_pairs.insert(pair);
                    self.close_calls += 1;
                    println!(
                        "Close call #{} — at ({},{}) → ({},{})",
                        self.close_calls, current.x, current.y, next.x, next.y
                    );
                    close_calls.insert(pair.0);
                    close_calls.insert(pair.1);
                }
            }
        }
        close_calls
    }
}
