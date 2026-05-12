use glam::Vec3;
use world::{Camera, Color, Cylinder, LightSource, Material, World};

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
    w.add_object(Cylinder::new(
        Vec3::new(0.0, -1.0, 0.0),
        Vec3::Y,
        1.0,
        2.0,
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
