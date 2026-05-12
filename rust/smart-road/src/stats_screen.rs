//! Post-simulation statistics screen.
//!
//! Loads a TTF font, builds a two-column label/value layout, and
//! draws a semi-transparent panel over the baked intersection
//! background. Also owns the post-stats event loop that holds the
//! window open until the user quits.

use sdl2::pixels::Color;
use sdl2::rect::Rect;
use sdl2::render::{BlendMode, Canvas, Texture, TextureCreator, TextureQuery};
use sdl2::ttf::{Font, Sdl2TtfContext};
use sdl2::video::{Window, WindowContext};
use sdl2::EventPump;

use crate::config::{TICK_DURATION_MS, WINDOW_SIZE};
use crate::input::{self, Action};
use crate::stats::Stats;

/// Render the final statistics panel on top of the baked intersection
/// `bg` and present it. The panel auto-sizes to fit the widest label +
/// value pair so numbers never clip, and is drawn with alpha blending
/// so the road underneath remains visible.
///
/// Called once at the end of the run; `wait_for_quit` then holds the
/// window open until the user exits.
pub fn show(
    canvas: &mut Canvas<Window>,
    texture_creator: &TextureCreator<WindowContext>,
    ttf_context: &Sdl2TtfContext,
    bg: &Texture,
    stats: &Stats,
    total_ticks: u64,
) -> Result<(), String> {
    // Load the bundled monospace font. Reloaded each time show() runs;
    // cheap, avoids threading a Font handle through the caller.
    let font = ttf_context.load_font("assets/courier.ttf", 22)?;
    // Build the (label, value) rows we want to display.
    let lines = build_lines(stats, total_ticks);

    // Layout constants — tweak for spacing/feel.
    let inner_pad = 24i32;
    let line_h = 36i32;
    let col_gap = 20i32;

    // Size the panel to fit the widest label + widest value so numbers
    // never clip. Centered in the window.
    let (max_label_w, max_value_w) = measure_columns(&font, &lines);
    let panel_w = (max_label_w + col_gap as u32 + max_value_w + 2 * inner_pad as u32).max(300);
    let panel_h = lines.len() as u32 * line_h as u32 + 2 * inner_pad as u32;
    let panel_x = (WINDOW_SIZE as i32 - panel_w as i32) / 2;
    let panel_y = (WINDOW_SIZE as i32 - panel_h as i32) / 2;

    // Draw the frozen intersection underneath, then layer the panel on top.
    canvas.copy(bg, None, None).ok();

    // Semi-transparent panel body — lets the intersection show through.
    canvas.set_blend_mode(BlendMode::Blend);
    canvas.set_draw_color(Color::RGBA(15, 15, 15, 210));
    canvas
        .fill_rect(Rect::new(panel_x, panel_y, panel_w, panel_h))
        .ok();

    // Subtle border around the panel.
    canvas.set_draw_color(Color::RGB(100, 100, 100));
    canvas
        .draw_rect(Rect::new(panel_x, panel_y, panel_w, panel_h))
        .ok();

    // Render each row. An empty label is a vertical spacer.
    for (i, (label, value)) in lines.iter().enumerate() {
        let y = panel_y + inner_pad + i as i32 * line_h;
        if label.is_empty() {
            continue;
        }

        if value.is_empty() {
            // Standalone row (title or footer) — centered.
            let color = if *label == "Press Esc to close" {
                Color::RGB(140, 140, 140)
            } else {
                Color::RGB(255, 210, 80)
            };
            let (tex, w, h) = render_text(&font, texture_creator, label, color)?;
            let x = panel_x + (panel_w as i32 - w as i32) / 2;
            canvas.copy(&tex, None, Rect::new(x, y, w, h)).ok();
        } else {
            // Two-column row: label flush left, value flush right.
            let (ltex, lw, lh) =
                render_text(&font, texture_creator, label, Color::RGB(180, 180, 180))?;
            canvas
                .copy(&ltex, None, Rect::new(panel_x + inner_pad, y, lw, lh))
                .ok();

            let (vtex, vw, vh) =
                render_text(&font, texture_creator, value, Color::RGB(230, 230, 230))?;
            let vx = panel_x + panel_w as i32 - inner_pad - vw as i32;
            canvas.copy(&vtex, None, Rect::new(vx, y, vw, vh)).ok();
        }
    }

    // Flip the back buffer once; the panel is static thereafter.
    canvas.present();
    Ok(())
}

/// Block the main thread, polling the event pump at ~20 Hz, until the
/// user triggers [`Action::Quit`] (Esc or window close). Nothing is
/// drawn here — the stats panel rendered by [`show`] stays on screen
/// because we never repaint.
pub fn wait_for_quit(event_pump: &mut EventPump) {
    loop {
        for event in event_pump.poll_iter() {
            if let Action::Quit = input::map_event(&event) {
                return;
            }
        }
        // Idle sleep — no animation running, so we don't need 60 Hz.
        std::thread::sleep(std::time::Duration::from_millis(50));
    }
}

/// Assemble the rows shown in the panel as `(label, value)` pairs.
///
/// Conventions the renderer relies on:
/// - An empty `label` is a vertical spacer (skipped by the draw loop).
/// - A non-empty `label` with empty `value` is a standalone centered
///   row (used for the title and the footer hint).
/// - Otherwise it is a two-column row: label flush left, value flush
///   right.
///
/// Unset stats (e.g. no vehicle ever finished, so `min_time` is still
/// `u64::MAX`) render as `"—"` rather than a misleading number.
fn build_lines(s: &Stats, total_ticks: u64) -> Vec<(&'static str, String)> {
    let ticks_fmt = |t: u64| -> String {
        let secs = t as f64 * TICK_DURATION_MS as f64 / 1000.0;
        format!("{} ticks ({:.2}s)", t, secs)
    };

    vec![
        ("SIMULATION RESULTS", String::new()),
        ("", String::new()),
        (
            "Total runtime",
            if total_ticks > 0 {
                ticks_fmt(total_ticks)
            } else {
                "—".into()
            },
        ),
        (
            "Vehicles through intersection",
            format!("{}", s.vehicle_count),
        ),
        (
            "Vehicles per min",
            format!(
                "{:.2}",
                vehicles_per_minute(s.vehicle_count as u64, total_ticks)
            ),
        ),
        ("Close calls", format!("{}", s.close_calls)),
        ("Collisions", format!("{}", s.collisions)),
        (
            "Slowest cell crossing",
            if s.min_speed != f32::MAX {
                format!("{:.2} cells/kTicks", s.min_speed)
            } else {
                "—".into()
            },
        ),
        (
            "Fastest cell crossing",
            if s.max_speed > 0. {
                format!("{:.2} cells/kTicks", s.max_speed)
            } else {
                "—".into()
            },
        ),
        (
            "Shortest transit",
            if s.min_time != u64::MAX {
                ticks_fmt(s.min_time)
            } else {
                "—".into()
            },
        ),
        (
            "Longest transit",
            if s.max_time > 0 {
                ticks_fmt(s.max_time)
            } else {
                "—".into()
            },
        ),
        ("", String::new()),
        ("Press Esc to close", String::new()),
    ]
}

/// Ask the TTF font for the pixel width of each label and value, and
/// return `(widest_label, widest_value)`. The caller uses these to size
/// the panel so text never clips regardless of string length.
///
/// Spacer rows (empty label) are skipped. Centered rows (empty value)
/// only contribute to the label column.
fn measure_columns(font: &Font, lines: &[(&'static str, String)]) -> (u32, u32) {
    let mut max_label_w = 0u32;
    let mut max_value_w = 0u32;
    for (label, value) in lines {
        if label.is_empty() {
            continue;
        }
        if let Ok(sz) = font.size_of(label) {
            max_label_w = max_label_w.max(sz.0);
        }
        if !value.is_empty() {
            if let Ok(sz) = font.size_of(value) {
                max_value_w = max_value_w.max(sz.0);
            }
        }
    }
    (max_label_w, max_value_w)
}

/// Rasterize `text` into an alpha-blended surface, upload it to a GPU
/// texture, and return its natural (width, height). The returned
/// Texture borrows from `texture_creator`, so its lifetime is tied to
/// the caller's canvas setup.
fn render_text<'a>(
    font: &Font,
    texture_creator: &'a TextureCreator<WindowContext>,
    text: &str,
    color: Color,
) -> Result<(Texture<'a>, u32, u32), String> {
    let surf = font
        .render(text)
        .blended(color)
        .map_err(|e| e.to_string())?;
    let tex = texture_creator
        .create_texture_from_surface(&surf)
        .map_err(|e| e.to_string())?;
    let TextureQuery { width, height, .. } = tex.query();
    Ok((tex, width, height))
}

fn vehicles_per_minute(vehicle_count: u64, total_ticks: u64) -> f64 {
    if total_ticks == 0 {
        return 0.0;
    }

    let total_time_minutes = (total_ticks as f64 * TICK_DURATION_MS as f64) / 1000.0 / 60.0;

    vehicle_count as f64 / total_time_minutes
}
