//! Application-wide constants for the Smart Road simulation.
//!
//! All layout dimensions derive from [`CELL_SIZE`] and [`GRID_SIZE`],
//! ensuring that a single change propagates consistently across the
//! entire rendering pipeline.

/// Side length of one grid cell in pixels.
pub const CELL_SIZE: u32 = 50;

/// Number of cells along each axis of the square grid.
pub const GRID_SIZE: u32 = 18;

/// Window width and height in pixels (derived: `CELL_SIZE * GRID_SIZE`).
pub const WINDOW_SIZE: u32 = CELL_SIZE * GRID_SIZE;

/// Minimum ticks to traverse one cell — the fastest a vehicle can move.
/// Book uses this to propagate earliest-possible arrival times. The actual
/// per-cell speed of a vehicle is implicit in its schedule:
/// times[i+1] - times[i] >= MIN_TICKS_PER_CELL.
pub const MIN_TICKS_PER_CELL: u64 = 20;

/// Maximum ticks to traverse one cell — the slowest a vehicle may move
/// through the intersection before it is considered stalling.
pub const MAX_TICKS_PER_CELL: u64 = 3 * MIN_TICKS_PER_CELL;

/// Duration in ticks for which an approach (non-core) cell reservation is held.
/// A vehicle occupies a cell for at least MIN_TICKS_PER_CELL; the second
/// MIN_TICKS_PER_CELL provides a safety buffer so the gap-check sees the full
/// occupied window.
pub const APPROACH_RESERVATION_TICKS: u64 = 2 * MIN_TICKS_PER_CELL;

/// How many steps ahead in the schedule a cell reservation must extend.
/// Cell i is held until times[i + RESERVATION_LOOKAHEAD], ensuring a vehicle's
/// booking covers not just its arrival tick but also its next departure tick.
pub const RESERVATION_LOOKAHEAD: usize = 4;

/// The number of cells at the start of the route that need to be unoccupied
/// at the time of a spawning for can_spawn to return true. Ensures
/// a healthy safety distance from the car in front.
pub const FREE_CELLS_FOR_SPAWNING: usize = 1;

/// Pixel distance between two car centers below which a collision is recorded.
pub const COLLISION_THRESHOLD_PX: f32 = CELL_SIZE as f32 * 2.0 / 3.0;

/// Milliseconds the main loop sleeps between ticks (~60 FPS at 16 ms).
pub const TICK_DURATION_MS: u64 = 16;

/// First column/row index occupied by the road (inclusive).
pub const ROAD_START: u32 = 6;

/// Last column/row index occupied by the road (inclusive).
pub const ROAD_END: u32 = 11;

/// Inner boundary of the central 4×4 collision zone (inclusive).
pub const ZONE_MIN: i32 = ROAD_START as i32 + 1; // 7

/// Outer boundary of the central 4×4 collision zone (inclusive).
pub const ZONE_MAX: i32 = ROAD_END as i32 - 1; // 10

// ── Tree sprite sheet layout ────────────────────────────────────────

/// Number of columns in `assets/trees.png`.
pub const TREE_COLS: u32 = 3;

/// Number of rows in `assets/trees.png`.
pub const TREE_ROWS: u32 = 2;

// ── Car sprite sheet layout ─────────────────────────────────────────

/// Columns in each car sprite sheet (one frame per 15° of rotation).
pub const CAR_SPRITE_COLS: u32 = 6;

/// Rows in each car sprite sheet (6 × 4 = 24 rotation frames total).
pub const CAR_SPRITE_ROWS: u32 = 4;

// Frame index (0..24, row-major) to use for each cardinal direction.
// The sheet starts at some unknown base angle and steps 15° per frame,
// so adjust these visually until each car faces the right way.
// Frames are numbered left-to-right, top-to-bottom.
pub const CAR_FRAME_NORTH: u32 = 15;
pub const CAR_FRAME_EAST: u32 = 9;
pub const CAR_FRAME_SOUTH: u32 = 3;
pub const CAR_FRAME_WEST: u32 = 21;

// ── Color palette (R, G, B) ─────────────────────────────────────────

/// Non-road background fill.
pub const COLOR_BACKGROUND: (u8, u8, u8) = (34, 34, 34);

/// Road surface.
pub const COLOR_ROAD: (u8, u8, u8) = (80, 80, 80);

/// Intersection zone (6×6 center block).
pub const COLOR_INTERSECTION: (u8, u8, u8) = (90, 90, 90);

/// Dashed lane divider lines between same-direction lanes.
pub const COLOR_LANE_DIVIDER: (u8, u8, u8) = (200, 200, 200);

/// Solid center line separating opposing traffic directions.
pub const COLOR_CENTER_LINE: (u8, u8, u8) = (220, 200, 50);

/// Solid lines at road boundaries.
pub const COLOR_ROAD_EDGE: (u8, u8, u8) = (255, 255, 255);
