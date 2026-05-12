/// Scene: flat plane + cube with reduced brightness.
use glam::{Quat, Vec3};
use world::{Camera, Color, Cube, LightSource, Material, Plane, World};

/// Scene: a single teal cylinder lit by one point light.
pub fn build() -> World {
    let cam = Camera::new(
        Vec3::new(0.0, 1.5, 5.5),
        Vec3::ZERO,
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
    w.add_object(Cube::new(
        Vec3::new(0.0, 0.5, 0.0),
        1.,
        Quat::IDENTITY,
        Material::matte(Color::rgb(0.85, 0.25, 0.25)),
    ));
    w.add_lights(
        LightSource::Point {
            pos: Vec3::new(4.0, 7.0, 4.0),
            intensity: 0.4,
            color: Color::WHITE,
        },
        10,
    );
    w
}
