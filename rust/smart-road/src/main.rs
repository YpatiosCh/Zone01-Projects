//! Binary entry point. Boots SDL, wires up the rendering resources,
//! hands control to the simulation, then displays the stats screen.

use smart_road::config::WINDOW_SIZE;
use smart_road::intersection;
use smart_road::simulation::{self, Args};
use smart_road::sprites::CarTextures;
use smart_road::stats_screen;

fn main() -> Result<(), String> {
    // Parse CLI flags (--auto, --highlight).
    let args = Args::parse();

    // Initialize SDL and the subsystems we need. `sdl_context` owns the
    // library state; dropping it shuts SDL down, so it must outlive
    // every handle derived from it.
    let sdl_context = sdl2::init()?;
    // Video subsystem is required before we can open a window.
    let video_subsystem = sdl_context.video()?;
    // TTF is a separate library (SDL2_ttf) with its own init/quit
    // lifecycle. Kept alive for the stats screen's font loading.
    let ttf_context = sdl2::ttf::init().map_err(|e| e.to_string())?;

    // Create the OS window. `position_centered` lets the WM place it.
    let window = video_subsystem
        .window("Smart Road", WINDOW_SIZE, WINDOW_SIZE)
        .position_centered()
        .build()
        .map_err(|e| e.to_string())?;

    // Turn the window into a Canvas — the 2D drawing surface. `accelerated`
    // asks for a GPU-backed renderer (falls back to software if unavailable).
    let mut canvas = window
        .into_canvas()
        .accelerated()
        .build()
        .map_err(|e| e.to_string())?;
    // Enable alpha blending so semi-transparent draws (collision overlay,
    // stats panel background) composite correctly instead of overwriting.
    canvas.set_blend_mode(sdl2::render::BlendMode::Blend);

    // A TextureCreator is the factory for GPU textures. It's borrowed
    // from the canvas, so its lifetime bounds every texture we create.
    let texture_creator = canvas.texture_creator();

    // Bake the intersection (roads, markings, trees) once into an
    // off-screen texture. We blit this every frame instead of redrawing
    // all the static geometry from scratch.
    let mut bg_texture = texture_creator
        .create_texture_target(None, WINDOW_SIZE, WINDOW_SIZE)
        .map_err(|e| e.to_string())?;
    // `with_texture_canvas` temporarily retargets the canvas at our
    // texture so the draw calls land there instead of the window.
    canvas
        .with_texture_canvas(&mut bg_texture, |c| {
            intersection::draw_intersection(c, &texture_creator);
        })
        .map_err(|e| e.to_string())?;

    // Load the car sprite sheets into GPU textures.
    let sprites = CarTextures::load(&texture_creator)?;
    // The event pump drains OS input events (keyboard, window close).
    let mut event_pump = sdl_context.event_pump()?;

    // Run the simulation until the user quits.
    let (stats, total_ticks) =
        simulation::run(&mut canvas, &bg_texture, &sprites, &mut event_pump, &args);

    // Show the final results panel and hold the window open until quit.
    stats_screen::show(
        &mut canvas,
        &texture_creator,
        &ttf_context,
        &bg_texture,
        &stats,
        total_ticks,
    )?;
    stats_screen::wait_for_quit(&mut event_pump);

    Ok(())
}
