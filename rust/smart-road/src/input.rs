//! Keyboard input handling for the Smart Road simulation.
//!
//! Maps SDL2 key events to high-level [`Action`] values that the main
//! loop dispatches. Arrow keys spawn a vehicle approaching from the
//! opposite side of the intersection (Up = car coming up from the
//! south, Down = coming down from the north, etc.). `R` toggles
//! continuous random spawning; `Esc` or window close quits.

use crate::models::Origin;
use rand::Rng;
use sdl2::event::Event;
use sdl2::keyboard::Keycode;

/// A high-level action derived from a single SDL event.
pub enum Action {
    Quit,
    Spawn(Origin),
    None,
}

/// Translate a single SDL event into an [`Action`].
///
/// Unhandled events (mouse motion, other keys, key repeats we ignore)
/// collapse to [`Action::None`] so the caller can uniformly skip them.
pub fn map_event(event: &Event) -> Action {
    match event {
        Event::Quit { .. } => Action::Quit,
        Event::KeyDown {
            keycode: Some(code),
            repeat: false,
            ..
        } => match *code {
            Keycode::Escape => Action::Quit,
            Keycode::Up => Action::Spawn(Origin::South),
            Keycode::Down => Action::Spawn(Origin::North),
            Keycode::Right => Action::Spawn(Origin::West),
            Keycode::Left => Action::Spawn(Origin::East),
            Keycode::R => Action::Spawn(random_origin()),
            _ => Action::None,
        },
        _ => Action::None,
    }
}

/// Pick a uniformly random route across all 12 variants. Fires once per
/// `R` key press.
pub fn random_origin() -> Origin {
    let origin = match rand::thread_rng().gen_range(0..4) {
        0 => Origin::North,
        1 => Origin::South,
        2 => Origin::East,
        _ => Origin::West,
    };
    origin
}
