use glam::Vec3;
use world::{Camera, Color, LightSource, Material, Plane, World};

/// Scene: one sphere resting above a grey ground plane.
pub fn build() -> World {
    let cam = Camera::new(
        Vec3::new(0.0, 2.0, 6.0),
        Vec3::new(0.0, 0.5, 0.0),
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
    // w.add_object(Sphere::new(
    //     Vec3::new(0.0, 0.98, 0.0),
    //     1.0,
    //     Material::matte(Color::rgb(0.8, 0.3, 0.3)),
    // ));
    w.add_lights(
        LightSource::Directional {
            dir: Vec3::new(-1.0, -2.0, -1.0),
            intensity: 1.0,
            color: Color::WHITE,
        },
        10,
    );
    w
}
