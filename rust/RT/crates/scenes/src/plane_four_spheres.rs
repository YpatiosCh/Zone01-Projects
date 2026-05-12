use glam::Vec3;
use world::{Camera, Color, LightSource, Material, Plane, Sphere, World};

/// Scene: a grey ground plane with four separated spheres of different sizes.
pub fn build() -> World {
    let cam = Camera::new(
        Vec3::new(0.0, 3.0, 10.0),
        Vec3::new(0.0, 1.4, 0.0),
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
        Vec3::new(-3.2, 0.58, -1.8),
        0.6,
        Material::matte(Color::rgb(0.85, 0.25, 0.25)),
    ));
    w.add_object(Sphere::new(
        Vec3::new(-1.0, 0.88, 1.2),
        0.9,
        Material::matte(Color::rgb(0.95, 0.75, 0.25)),
    ));
    w.add_object(Sphere::new(
        Vec3::new(1.8, 1.18, -0.5),
        1.2,
        Material::matte(Color::rgb(0.2, 0.6, 0.9)),
    ));
    w.add_object(Sphere::new(
        Vec3::new(4.4, 0.43, 1.8),
        0.45,
        Material::matte(Color::rgb(0.35, 0.85, 0.45)),
    ));

    w.add_lights(LightSource::Point { 
        pos: Vec3::new(4.0, 7.0, 4.0), 
        intensity: 1.0, 
        color: Color::WHITE }, 
        10);
    w
}