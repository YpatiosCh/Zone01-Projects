//! Real-time simulation loop.
//!
//! Owns the [`Book`] and [`Vehicles`] for the run, drains SDL events
//! through [`input::map_event`], calls into the backend each tick, and
//! blits the result on top of the baked intersection background.
//!
//! Returns the final [`Stats`] and total tick count so the caller can
//! hand them to the stats screen.

use sdl2::pixels::Color;
use sdl2::rect::Rect;
use sdl2::render::{Canvas, Texture};
use sdl2::video::Window;
use sdl2::EventPump;

use crate::book::Book;
use crate::config::{CELL_SIZE, TICK_DURATION_MS};
use crate::input::{self, random_origin, Action};
use crate::models::Highlight;
use crate::sprites::CarTextures;
use crate::stats::Stats;
use crate::vehicles::Vehicles;

/// Command-line flags parsed from `std::env::args()`.
pub struct Args {
    /// Spawn a random-origin vehicle every tick.
    pub auto: bool,
    /// Draw a rect around each vehicle (red while in a close-call pair).
    pub highlight: bool,
}

impl Args {
    pub fn parse() -> Self {
        let args: Vec<String> = std::env::args().collect();
        Self {
            auto: args.iter().any(|a| a == "--auto"),
            highlight: args.iter().any(|a| a == "--highlight"),
        }
    }
}

/// Run the simulation until the user quits. Returns the final
/// [`Stats`] and the total tick count reached before exit.
pub fn run(
    canvas: &mut Canvas<Window>,
    bg: &Texture,
    sprites: &CarTextures,
    event_pump: &mut EventPump,
    args: &Args,
) -> (Stats, u64) {
    // Book owns all cell reservations; Vehicles owns the spawned cars
    // and their per-car schedules. Both live for the whole run.
    let mut book = Book::new();
    let mut vehicles = Vehicles::new();
    // Tick counter drives all scheduling and interpolation.
    let mut tick: u64 = 0;

    'running: loop {
        // --auto mode: attempt to spawn one random-origin vehicle every
        // tick. Book::can_spawn rejects it if the lane is not clear.
        if args.auto {
            vehicles.spawn_vehicle(&mut book, random_origin(), tick);
        }
        // Drain every queued OS event this tick. `map_event` collapses
        // raw SDL events into our small Action enum.
        for event in event_pump.poll_iter() {
            match input::map_event(&event) {
                Action::Quit => break 'running,
                Action::Spawn(origin) => {
                    vehicles.spawn_vehicle(&mut book, origin, tick);
                }
                Action::None => {}
            }
        }

        // Advance every live vehicle, interpolate its render position,
        // and update stats (close calls, collisions, speed, transit).
        // Finished vehicles are dropped from the list here.
        let to_render = vehicles.get_positions_and_update_stats(tick);

        // Blit the baked intersection as the base frame, then draw
        // sprites on top. No need to redraw static geometry per tick.
        canvas.copy(bg, None, None).ok();
        for (id, pos) in &to_render {
            // Pick this car's sprite sheet (hash of id → variant) and
            // the frame matching its continuous angle.
            let texture = sprites.texture_for(id);
            let src = sprites.src_rect(pos.angle);
            let dst = Rect::new(pos.x as i32, pos.y as i32, CELL_SIZE, CELL_SIZE);
            canvas.copy(texture, src, dst).ok();

            // Overlay: collisions always flash red-filled; close calls
            // and the debug rect only render under --highlight.
            match pos.highlight {
                Highlight::Collision => {
                    canvas.set_draw_color(Color::RGBA(255, 0, 0, 160));
                    canvas.fill_rect(dst).ok();
                }
                Highlight::CloseCall if args.highlight => {
                    canvas.set_draw_color(Color::RGB(255, 0, 0));
                    canvas.draw_rect(dst).ok();
                }
                _ => {
                    if args.highlight {
                        canvas.set_draw_color(Color::RGB(0, 255, 0));
                        canvas.draw_rect(dst).ok();
                    }
                }
            }
        }
        // Flip the back buffer to the window.
        canvas.present();

        tick += 1;
        // Fixed-step pacing. Crude but adequate at 60 Hz — no catch-up
        // logic because the simulation is tick-indexed, not wall-clock.
        std::thread::sleep(std::time::Duration::from_millis(TICK_DURATION_MS));
    }

    // Partial move: pull stats out of `vehicles` as the run ends.
    (vehicles.stats, tick)
}
