use glam::Vec3;
use world::{World, Camera, Sphere, Material, Color, LightSource};

/// Scene: a single red sphere lit by one point light.
pub fn build() -> World {
    let cam = Camera::new(Vec3::new(0.0, 1.0, 5.0), Vec3::ZERO, Vec3::Y, 60.0, 800, 600);
    let mut w = World::new(cam);
    w.add_object(Sphere {
        center: Vec3::ZERO,
        radius: 1.0,
        material: Material::matte(Color::rgb(0.8, 0.3, 0.3)),
    });
    w.add_light(LightSource::Point {
        pos: Vec3::new(2.0, 4.0, 2.0),
        intensity: 1.0,
        color: Color::WHITE,
    });
    w
}
