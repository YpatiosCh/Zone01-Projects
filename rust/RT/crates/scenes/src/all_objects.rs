use glam::{Quat, Vec3};
use world::{Camera, Color, Cube, Cylinder, LightSource, Material, Plane, Sphere, World};

/// Scene: one of each shape (sphere, cube, plane, cylinder).
pub fn build() -> World {
    let cam = Camera::new(
        Vec3::new(0.0, 3.5, 9.0),
        Vec3::new(0.0, 0.8, 0.0),
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
        Vec3::new(-2.5, 0.6, 0.0),
        0.6,
        Material::matte(Color::rgb(0.95, 0.75, 0.2)),
    ));
    w.add_object(Cube::new(
        Vec3::new(0.0, 0.5, 0.0),
        1.0,
        Quat::IDENTITY,
        Material::matte(Color::rgb(0.85, 0.25, 0.25)),
    ));
    w.add_object(Cylinder::new(
        Vec3::new(2.5, 0.0, 0.0),
        Vec3::Y,
        0.5,
        1.5,
        Material::matte(Color::rgb(0.2, 0.7, 0.75)),
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
