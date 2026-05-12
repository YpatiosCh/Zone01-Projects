use glam::Vec3;
use world::{Camera, Color, LightSource, Material, Plane, Sphere, World};

/// Scene: a grey ground plane with one reflective sphere and one matte sphere.
pub fn build() -> World {
    let cam = Camera::new(
        Vec3::new(0.0, 2.5, 8.0),
        Vec3::new(0.0, 1.2, 0.0),
        Vec3::Y,
        60.0,
        800,
        600,
    );
    let mut w = World::new(cam);
    w.add_object(Plane::new(
        Vec3::ZERO,
        Vec3::Y,
        Material::matte(Color::rgb(0.5, 0.5, 0.5)),
    ));
    w.add_object(Sphere::new(
        Vec3::new(-1.7, 0.98, 0.0),
        1.0,
        Material::mirror(Color::WHITE),
    ));
    w.add_object(Sphere::new(
        Vec3::new(1.8, 0.73, -0.2),
        0.75,
        Material::matte(Color::rgb(0.85, 0.35, 0.3)),
    ));
    w.add_lights(
        LightSource::Point {
            pos: Vec3::new(4.0, 7.0, 4.0),
            intensity: 1.0,
            color: Color::WHITE,
        },
        10,
    );
    w
}
