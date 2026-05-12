pub mod sphere;
pub mod cube;
pub mod plane;
pub mod cylinder;

pub use sphere::Sphere;
pub use cube::Cube;
pub use plane::Plane;
pub use cylinder::Cylinder;

use glam::Vec3;
use std::ops::{Add, Mul};

#[derive(Clone, Copy, Debug, Default)]
pub struct Color { pub r: f32, pub g: f32, pub b: f32 }

impl Color {
    pub const BLACK: Color = Color { r: 0.0, g: 0.0, b: 0.0 };
    pub const WHITE: Color = Color { r: 1.0, g: 1.0, b: 1.0 };
    pub fn rgb(r: f32, g: f32, b: f32) -> Self { Self { r, g, b } }
}

impl Add for Color {
    type Output = Color;
    fn add(self, o: Color) -> Color { Color::rgb(self.r + o.r, self.g + o.g, self.b + o.b) }
}
impl Mul<f32> for Color {
    type Output = Color;
    fn mul(self, s: f32) -> Color { Color::rgb(self.r * s, self.g * s, self.b * s) }
}
impl Mul<Color> for Color {
    type Output = Color;
    fn mul(self, o: Color) -> Color { Color::rgb(self.r * o.r, self.g * o.g, self.b * o.b) }
}

#[derive(Clone, Copy, Debug)]
pub struct Ray { pub origin: Vec3, pub dir: Vec3 }

impl Ray {
    /// `dir` is normalized on construction; downstream code may assume unit length.
    pub fn new(origin: Vec3, dir: Vec3) -> Self { Self { origin, dir: dir.normalize() } }
    /// Point along the ray at parameter `t`.
    pub fn at(&self, t: f32) -> Vec3 { self.origin + self.dir * t }
}

/// Surface intersection record. Created by the shape's `intersect`; consumed by the raycaster.
/// `normal` is always the outward-facing surface normal (length 1).
#[derive(Clone, Copy, Debug)]
pub struct Hit {
    pub distance: f32,
    pub point: Vec3,
    pub normal: Vec3,
    pub material: Material,
}

#[derive(Clone, Copy, Debug, Default)]
pub struct Material {
    pub color: Color,
    pub reflective: bool,
}

impl Material {
    pub fn matte(color: Color) -> Self { Self { color, reflective: false } }
    pub fn mirror(color: Color) -> Self { Self { color, reflective: true } }
}

/// Small offset used to push spawned rays off a surface to avoid self-intersection.
pub const SURFACE_EPSILON: f32 = 1e-4;

/// Mirror reflection of `incoming` about `normal`. Same math for every shape, so it lives
/// here as a free function instead of on the `Object` trait. The raycaster calls this when
/// `hit.material.reflective` is true.
pub fn reflect(incoming: Vec3, normal: Vec3) -> Vec3 {
    incoming - normal * (2.0 * incoming.dot(normal))
}

pub trait Object {
    /// Ray-vs-shape intersection. Returns the nearest forward hit (`distance > 0`).
    /// The shape fills in its own `material` on the returned `Hit`.
    fn intersect(&self, ray: Ray) -> Option<Hit>;
    /// World-space anchor point.
    fn position(&self) -> Vec3;
    /// Move the object before rendering.
    fn set_position(&mut self, p: Vec3);
    /// World-space bounding sphere `(center, radius)` for frustum culling.
    fn bounding_sphere(&self) -> (Vec3, f32);
    /// Optional floor height at a world-space position, used by liveview collision.
    fn floor_height_at(&self, _position: Vec3) -> Option<f32> { None }
}
