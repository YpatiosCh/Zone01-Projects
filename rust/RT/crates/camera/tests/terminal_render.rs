// Smoke-renders a single sphere at 30×30 and prints ASCII art to the terminal.
// Run: cargo test -p camera --test terminal_render -- --nocapture

use camera::Camera;
use glam::Vec3;
use objects::{Color, Material, Object, Sphere};
use raycaster::Raycaster;

#[test]
fn terminal_render() {
    let cam = Camera::new(
        Vec3::new(0.0, 0.0, 5.0),
        Vec3::ZERO,
        Vec3::Y,
        40.0,
        30, 30,
    );

    let objects: Vec<Box<dyn Object>> = vec![
        Box::new(Sphere::new(Vec3::ZERO, 1.0, Material::matte(Color::rgb(0.8, 0.3, 0.3)))),
    ];

    // Lighting is not implemented yet, so we compute a one-liner diffuse term
    // directly from the hit normal: brightness = max(normal · to_light, 0).
    // This is Lambert shading without shadows, attenuation, or ambient — just
    // enough to show the sphere's 3D shape.
    let to_light = Vec3::new(1.0, 2.0, 1.0).normalize();

    // Index 0 (' ') is reserved for background misses. Hit pixels start at
    // index 1 ('.') even when completely unlit, so the sphere silhouette is
    // always visible against the empty background.
    let shades = " .:-=+*#%@";

    println!();
    for y in 0..cam.height {
        for x in 0..cam.width {
            let ray = cam.ray_for_pixel(x, y);
            let (_, hit) = Raycaster::trace_with_reflections(ray, &objects, 1);
            let ch = match hit {
                None => ' ',
                Some(h) => {
                    let brightness = h.normal.dot(to_light).max(0.0);
                    let idx = 1 + (brightness * 8.0) as usize;
                    shades.chars().nth(idx).unwrap()
                }
            };
            // Each character is printed twice. Terminal glyphs are roughly
            // 2× taller than they are wide, so doubling horizontally makes the
            // sphere appear round rather than vertically squashed.
            print!("{ch}{ch}");
        }
        println!();
    }
}
