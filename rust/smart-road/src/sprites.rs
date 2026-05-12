//! Car sprite loading and frame selection.
//!
//! Each car sprite sheet (`assets/{red,old,yellow}.png`) is a
//! 6×4 grid of 24 rotation frames in 15° increments. Frame 0 faces
//! North; frames 6, 12, 18 face East, South, West.
//!
//! The texture lifetime is tied to the [`TextureCreator`] that loaded
//! it, so [`CarTextures`] is parameterised by that lifetime and must
//! live no longer than the creator.

use sdl2::image::LoadTexture;
use sdl2::rect::Rect;
use sdl2::render::{Texture, TextureCreator};
use sdl2::video::WindowContext;
use uuid::Uuid;

use crate::config::{
    // CAR_FRAME_EAST, CAR_FRAME_NORTH, CAR_FRAME_SOUTH, CAR_FRAME_WEST,
    CAR_SPRITE_COLS,
    CAR_SPRITE_ROWS,
};

/// Three car variants and the dimensions of one frame within a sheet.
pub struct CarTextures<'a> {
    variants: [Texture<'a>; 3],
    frame_w: u32,
    frame_h: u32,
}

impl<'a> CarTextures<'a> {
    /// Load all three car sheets. Assumes every sheet shares identical
    /// frame dimensions, so the first sheet defines `frame_w`/`frame_h`.
    pub fn load(tc: &'a TextureCreator<WindowContext>) -> Result<Self, String> {
        let red = tc.load_texture("assets/red.png")?;
        let old = tc.load_texture("assets/old.png")?;
        let yellow = tc.load_texture("assets/yellow.png")?;

        let q = red.query();
        let frame_w = q.width / CAR_SPRITE_COLS;
        let frame_h = q.height / CAR_SPRITE_ROWS;

        Ok(CarTextures {
            variants: [red, old, yellow],
            frame_w,
            frame_h,
        })
    }

    /// Deterministically pick a car variant for `id` so a given vehicle
    /// always renders with the same sprite across frames.
    pub fn texture_for(&self, id: &Uuid) -> &Texture<'a> {
        let idx = id.as_bytes()[0] as usize % self.variants.len();
        &self.variants[idx]
    }

    /// Source rectangle inside a sheet for the given cardinal direction.
    pub fn src_rect(&self, dir: f32) -> Rect {
        let idx = angle_to_sprite(dir);
        let col = idx % CAR_SPRITE_COLS;
        let row = idx / CAR_SPRITE_COLS;
        Rect::new(
            (col * self.frame_w) as i32,
            (row * self.frame_h) as i32,
            self.frame_w,
            self.frame_h,
        )
    }
}

/// Map a continuous heading in degrees to a sprite-sheet frame index
/// in `0..24`.
///
/// The sheet is laid out as 24 frames in 15° increments, but two
/// quirks force the formula:
/// 1. The frames step *clockwise*, while our angle convention is
///    counter-clockwise — hence the `15 - step` flip.
/// 2. Frame 15 (not frame 0) is the one that visually points North
///    (0°), so we offset by 15 to line the sheet up with our
///    direction convention.
///
/// `rem_euclid` keeps the result non-negative after the flip.
fn angle_to_sprite(angle: f32) -> u32 {
    let normalized = angle.rem_euclid(360.0);
    let step = (normalized / 15.0).round() as i32 % 24;
    let sprite = (15 - step).rem_euclid(24);
    sprite as u32
}
