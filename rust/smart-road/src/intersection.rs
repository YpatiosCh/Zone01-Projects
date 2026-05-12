//! Renders the static intersection background: roads, lane markings,
//! and decorative trees in the four corner quadrants.
//!
//! The intersection follows an 18×18 cell grid where columns/rows 6–11
//! form a six-lane cross. Each direction has three incoming lanes
//! (Right, Straight, Left) and three outgoing lanes.
//!
//! The result is drawn once and captured into a texture that is blitted
//! every frame — no per-frame recalculation required.

use sdl2::image::LoadTexture;
use sdl2::pixels::Color;
use sdl2::rect::Rect;
use sdl2::render::{Canvas, TextureCreator};
use sdl2::video::{Window, WindowContext};

use crate::config::*;

/// Draws the complete intersection scene onto `canvas`.
///
/// Rendering order (back-to-front):
/// 1. Green background
/// 2. Gray road surfaces (vertical + horizontal strips)
/// 3. Slightly lighter intersection zone
/// 4. Decorative trees in the four corner quadrants
/// 5. Road edge lines, center lines, and dashed lane dividers
pub fn draw_intersection(
    canvas: &mut Canvas<Window>,
    texture_creator: &TextureCreator<WindowContext>,
) {
    let cs = CELL_SIZE;
    let gs = GRID_SIZE;

    // 1. Green background for non-road areas
    canvas.set_draw_color(Color::RGB(34, 80, 34));
    canvas.clear();

    // 2. Road surfaces
    let (rr, rg, rb) = COLOR_ROAD;
    canvas.set_draw_color(Color::RGB(rr, rg, rb));

    // Vertical road (columns ROAD_START..=ROAD_END, full height)
    canvas
        .fill_rect(Rect::new(
            (ROAD_START * cs) as i32,
            0,
            (ROAD_END - ROAD_START + 1) * cs,
            gs * cs,
        ))
        .ok();

    // Horizontal road (rows ROAD_START..=ROAD_END, full width)
    canvas
        .fill_rect(Rect::new(
            0,
            (ROAD_START * cs) as i32,
            gs * cs,
            (ROAD_END - ROAD_START + 1) * cs,
        ))
        .ok();

    // 3. Intersection zone overlay
    let (ir, ig, ib) = COLOR_INTERSECTION;
    canvas.set_draw_color(Color::RGB(ir, ig, ib));
    canvas
        .fill_rect(Rect::new(
            (ROAD_START * cs) as i32,
            (ROAD_START * cs) as i32,
            (ROAD_END - ROAD_START + 1) * cs,
            (ROAD_END - ROAD_START + 1) * cs,
        ))
        .ok();

    // 4. Trees
    draw_trees(canvas, texture_creator);

    // 5. Road markings
    draw_road_markings(canvas);
}

// ── Road markings ───────────────────────────────────────────────────

/// Draws edge lines, center lines, and dashed lane dividers on the
/// approach roads (outside the intersection zone).
fn draw_road_markings(canvas: &mut Canvas<Window>) {
    let cs = CELL_SIZE;
    let gs = GRID_SIZE;
    let window = (gs * cs) as i32;

    let road_left = (ROAD_START * cs) as i32;
    let road_right = ((ROAD_END + 1) * cs) as i32;
    let road_top = (ROAD_START * cs) as i32;
    let road_bottom = ((ROAD_END + 1) * cs) as i32;

    // ── Solid white road-edge lines ──
    let (er, eg, eb) = COLOR_ROAD_EDGE;
    canvas.set_draw_color(Color::RGB(er, eg, eb));

    // Vertical road edges (top and bottom sections)
    canvas.draw_line((road_left, 0), (road_left, road_top)).ok();
    canvas.draw_line((road_right, 0), (road_right, road_top)).ok();
    canvas.draw_line((road_left, road_bottom), (road_left, window)).ok();
    canvas.draw_line((road_right, road_bottom), (road_right, window)).ok();

    // Horizontal road edges (left and right sections)
    canvas.draw_line((0, road_top), (road_left, road_top)).ok();
    canvas.draw_line((0, road_bottom), (road_left, road_bottom)).ok();
    canvas.draw_line((road_right, road_top), (window, road_top)).ok();
    canvas.draw_line((road_right, road_bottom), (window, road_bottom)).ok();

    // ── Solid yellow center line (between lane groups 3 and 4) ──
    let (cr, cg, cb) = COLOR_CENTER_LINE;
    canvas.set_draw_color(Color::RGB(cr, cg, cb));

    let center_v = ((ROAD_START + 3) * cs) as i32;
    canvas.draw_line((center_v, 0), (center_v, road_top)).ok();
    canvas.draw_line((center_v, road_bottom), (center_v, window)).ok();

    let center_h = ((ROAD_START + 3) * cs) as i32;
    canvas.draw_line((0, center_h), (road_left, center_h)).ok();
    canvas.draw_line((road_right, center_h), (window, center_h)).ok();

    // ── Dashed white lane dividers ──
    let (lr, lg, lb) = COLOR_LANE_DIVIDER;
    canvas.set_draw_color(Color::RGB(lr, lg, lb));

    let dash = 12_i32;
    let gap = 10_i32;

    // Offsets 1, 2 = between lanes in the incoming group
    // Offsets 4, 5 = between lanes in the outgoing group
    // (offset 3 is the center line drawn above)
    for offset in [1u32, 2, 4, 5] {
        let x = ((ROAD_START + offset) * cs) as i32;
        draw_dashed_line_v(canvas, x, 0, road_top, dash, gap);
        draw_dashed_line_v(canvas, x, road_bottom, window, dash, gap);

        let y = ((ROAD_START + offset) * cs) as i32;
        draw_dashed_line_h(canvas, 0, road_left, y, dash, gap);
        draw_dashed_line_h(canvas, road_right, window, y, dash, gap);
    }
}

// ── Trees ───────────────────────────────────────────────────────────

/// Fills the four corner quadrants with densely overlapping trees
/// loaded from `assets/trees.png`.
///
/// Trees are placed on a half-overlap grid with deterministic jitter
/// so the pattern looks organic yet is reproducible across frames.
fn draw_trees(canvas: &mut Canvas<Window>, texture_creator: &TextureCreator<WindowContext>) {
    let tree_texture = match texture_creator.load_texture("assets/trees.png") {
        Ok(t) => t,
        Err(_) => return,
    };

    let query = tree_texture.query();
    let tree_cell_w = query.width / TREE_COLS;
    let tree_cell_h = query.height / TREE_ROWS;

    let cs = CELL_SIZE;
    let tree_size = cs + cs / 2; // ~1.5 cells — large enough for overlap

    // Inset boundaries so trees don't spill onto the road.
    // Top quadrants pull y_max back by a full tree height because
    // the sprite's visual root extends downward from its origin.
    let margin = tree_size as i32 / 2;
    let corners: [(i32, i32, i32, i32); 4] = [
        (
            0,
            0,
            (ROAD_START * cs) as i32 - margin,
            (ROAD_START * cs) as i32 - tree_size as i32,
        ),
        (
            ((ROAD_END + 1) * cs) as i32 + margin,
            0,
            (GRID_SIZE * cs) as i32,
            (ROAD_START * cs) as i32 - tree_size as i32,
        ),
        (
            0,
            ((ROAD_END + 1) * cs) as i32 + margin,
            (ROAD_START * cs) as i32 - margin,
            (GRID_SIZE * cs) as i32,
        ),
        (
            ((ROAD_END + 1) * cs) as i32 + margin,
            ((ROAD_END + 1) * cs) as i32 + margin,
            (GRID_SIZE * cs) as i32,
            (GRID_SIZE * cs) as i32,
        ),
    ];

    // Deterministic LCG so the forest looks identical every frame.
    let mut seed: u32 = 42;

    for (x_min, y_min, x_max, y_max) in &corners {
        let mut y = *y_min - (tree_size as i32 / 3);
        while y < *y_max {
            let mut x = *x_min - (tree_size as i32 / 3);
            while x < *x_max {
                seed = seed.wrapping_mul(1103515245).wrapping_add(12345);
                let tree_idx = (seed >> 16) % (TREE_COLS * TREE_ROWS);
                let col = tree_idx % TREE_COLS;
                let row = tree_idx / TREE_COLS;

                let src = Rect::new(
                    (col * tree_cell_w) as i32,
                    (row * tree_cell_h) as i32,
                    tree_cell_w,
                    tree_cell_h,
                );

                let jitter_x = ((seed >> 8) % 15) as i32 - 7;
                let jitter_y = ((seed >> 12) % 15) as i32 - 7;
                let dst = Rect::new(x + jitter_x, y + jitter_y, tree_size, tree_size);

                canvas.copy(&tree_texture, src, dst).ok();

                x += tree_size as i32 / 2;
            }
            y += tree_size as i32 / 2;
        }
    }
}

// ── Dashed-line helpers ─────────────────────────────────────────────

/// Draws a vertical dashed line from `y_start` to `y_end` at column `x`.
fn draw_dashed_line_v(
    canvas: &mut Canvas<Window>,
    x: i32,
    y_start: i32,
    y_end: i32,
    dash: i32,
    gap: i32,
) {
    let mut y = y_start;
    while y < y_end {
        let end = (y + dash).min(y_end);
        canvas.draw_line((x, y), (x, end)).ok();
        y += dash + gap;
    }
}

/// Draws a horizontal dashed line from `x_start` to `x_end` at row `y`.
fn draw_dashed_line_h(
    canvas: &mut Canvas<Window>,
    x_start: i32,
    x_end: i32,
    y: i32,
    dash: i32,
    gap: i32,
) {
    let mut x = x_start;
    while x < x_end {
        let end = (x + dash).min(x_end);
        canvas.draw_line((x, y), (end, y)).ok();
        x += dash + gap;
    }
}
