// E2E fuzz tests — generate random worlds with a fixed seed, render each at
// 100×100, assert no panic and no NaN/Inf pixels.  Goal: crash / edge-case
// coverage, not visual correctness.
//
// Also exercises every pre-built scene to make sure none of them panic on a
// full render at a tiny resolution.
//
// Run with: cargo test --test e2e

use glam::{Quat, Vec3};
use rand::rngs::StdRng;
use rand::{Rng, SeedableRng};
use world::{Camera, Color, Cube, Cylinder, LightSource, Material, Plane, Sphere, World};

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

/// Build a world seeded from `rng`.  Generates:
///   - 1–5 spheres (some possibly mirrors)
///   - 0–2 cubes
///   - 0–2 planes
///   - 0–2 cylinders
///   - 1–3 point lights and/or directional lights
///   - Camera at a random position looking at the origin
fn random_world(rng: &mut StdRng) -> World {
    // --- camera -----------------------------------------------------------
    let cam_pos = Vec3::new(
        rng.gen_range(-8.0..=8.0),
        rng.gen_range(0.5..=6.0),
        rng.gen_range(3.0..=10.0),
    );
    let look_at = Vec3::new(
        rng.gen_range(-1.0..=1.0),
        rng.gen_range(-0.5..=1.5),
        0.0,
    );
    let fov = rng.gen_range(30.0..=100.0_f32);
    let cam = Camera::new(cam_pos, look_at, Vec3::Y, fov, 100, 100);
    let mut world = World::new(cam);

    // --- spheres ----------------------------------------------------------
    let sphere_count: u32 = rng.gen_range(1..=5);
    for _ in 0..sphere_count {
        world.add_object(Sphere::random_with(rng));
    }

    // --- cubes ------------------------------------------------------------
    let cube_count: u32 = rng.gen_range(0..=2);
    for _ in 0..cube_count {
        let center = Vec3::new(
            rng.gen_range(-4.0..=4.0),
            rng.gen_range(-1.0..=3.0),
            rng.gen_range(-4.0..=4.0),
        );
        let size = rng.gen_range(0.3..=2.0_f32);
        let angle: f32 = rng.gen_range(0.0..std::f32::consts::TAU);
        let rotation = Quat::from_rotation_y(angle);
        let color = Color::rgb(rng.gen(), rng.gen(), rng.gen());
        let material = if rng.gen_bool(0.2) {
            Material::mirror(color)
        } else {
            Material::matte(color)
        };
        world.add_object(Cube::new(center, size, rotation, material));
    }

    // --- planes -----------------------------------------------------------
    let plane_count: u32 = rng.gen_range(0..=2);
    for _ in 0..plane_count {
        let point = Vec3::new(0.0, rng.gen_range(-2.0..=0.0), 0.0);
        let color = Color::rgb(rng.gen(), rng.gen(), rng.gen());
        let material = if rng.gen_bool(0.15) {
            Material::mirror(color)
        } else {
            Material::matte(color)
        };
        world.add_object(Plane::new(point, Vec3::Y, material));
    }

    // --- cylinders --------------------------------------------------------
    let cyl_count: u32 = rng.gen_range(0..=2);
    for _ in 0..cyl_count {
        let base = Vec3::new(
            rng.gen_range(-4.0..=4.0),
            rng.gen_range(-1.0..=0.0),
            rng.gen_range(-4.0..=4.0),
        );
        let radius = rng.gen_range(0.2..=1.0_f32);
        let height = rng.gen_range(0.5..=3.0_f32);
        let color = Color::rgb(rng.gen(), rng.gen(), rng.gen());
        let material = Material::matte(color);
        world.add_object(Cylinder::new(base, Vec3::Y, radius, height, material));
    }

    // --- lights -----------------------------------------------------------
    let light_count: u32 = rng.gen_range(1..=3);
    for _ in 0..light_count {
        if rng.gen_bool(0.5) {
            let pos = Vec3::new(
                rng.gen_range(-6.0..=6.0),
                rng.gen_range(2.0..=8.0),
                rng.gen_range(-6.0..=6.0),
            );
            world.add_light(LightSource::Point {
                pos,
                intensity: rng.gen_range(0.3..=1.5),
                color: Color::WHITE,
            });
        } else {
            let dir = Vec3::new(
                rng.gen_range(-1.0..=1.0),
                rng.gen_range(-1.0..=-0.1),
                rng.gen_range(-1.0..=1.0),
            );
            world.add_light(LightSource::Directional {
                dir,
                intensity: rng.gen_range(0.3..=1.0),
                color: Color::WHITE,
            });
        }
    }

    world
}

/// Assert every pixel in the framebuffer is finite (no NaN, no Inf).
fn assert_no_nan_pixels(fb: &world::Framebuffer, label: &str) {
    for (i, pixel) in fb.pixels.iter().enumerate() {
        assert!(
            pixel.r.is_finite() && pixel.g.is_finite() && pixel.b.is_finite(),
            "{label}: pixel {i} contains NaN/Inf: ({}, {}, {})",
            pixel.r,
            pixel.g,
            pixel.b
        );
    }
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

/// Core fuzz: 50 random worlds at 100×100, fixed seed for reproducibility.
/// Pass = no panic and no NaN pixels.
#[test]
fn random_scenes_render_without_panic() {
    const SEED: u64 = 0xDEAD_BEEF_1234_5678;
    const SCENE_COUNT: u32 = 50;
    const W: u32 = 100;
    const H: u32 = 100;

    let mut rng = StdRng::seed_from_u64(SEED);

    for i in 0..SCENE_COUNT {
        let world = random_world(&mut rng);
        let fb = world.render(W, H);

        assert_eq!(fb.w, W, "scene {i}: framebuffer width mismatch");
        assert_eq!(fb.h, H, "scene {i}: framebuffer height mismatch");
        assert_eq!(
            fb.pixels.len(),
            (W * H) as usize,
            "scene {i}: pixel count mismatch"
        );
        assert_no_nan_pixels(&fb, &format!("random scene {i}"));
    }
}

/// Verify the framebuffer dimensions match the requested render size.
#[test]
fn render_returns_correct_dimensions() {
    let cam = Camera::new(Vec3::new(0.0, 1.0, 5.0), Vec3::ZERO, Vec3::Y, 60.0, 800, 600);
    let mut world = World::new(cam);
    world.add_object(Sphere::new(Vec3::ZERO, 1.0, Material::matte(Color::WHITE)));
    world.add_light(LightSource::Point {
        pos: Vec3::new(2.0, 4.0, 2.0),
        intensity: 1.0,
        color: Color::WHITE,
    });

    let fb = world.render(64, 48);
    assert_eq!(fb.w, 64);
    assert_eq!(fb.h, 48);
    assert_eq!(fb.pixels.len(), 64 * 48);
}

/// A world with no objects should render without panicking (sky fallback).
#[test]
fn empty_world_renders_sky() {
    let cam = Camera::new(Vec3::new(0.0, 1.0, 5.0), Vec3::ZERO, Vec3::Y, 60.0, 10, 10);
    let world = World::new(cam);
    let fb = world.render(10, 10);
    assert_no_nan_pixels(&fb, "empty world");
    assert_eq!(fb.pixels.len(), 100);
}

/// A world with no lights should render without panicking (ambient only).
#[test]
fn world_with_no_lights_renders_without_panic() {
    let cam = Camera::new(Vec3::new(0.0, 1.0, 5.0), Vec3::ZERO, Vec3::Y, 60.0, 20, 20);
    let mut world = World::new(cam);
    world.add_object(Sphere::new(
        Vec3::ZERO,
        1.0,
        Material::matte(Color::rgb(0.8, 0.3, 0.3)),
    ));
    let fb = world.render(20, 20);
    assert_no_nan_pixels(&fb, "no lights");
}

/// A fully mirror scene must not recurse infinitely (bounded by MAX_DEPTH).
#[test]
fn all_mirror_scene_does_not_hang() {
    let cam = Camera::new(Vec3::new(0.0, 0.0, 5.0), Vec3::ZERO, Vec3::Y, 60.0, 20, 20);
    let mut world = World::new(cam);
    // Two mirror spheres facing each other — without depth-limit this would loop forever.
    world.add_object(Sphere::new(
        Vec3::new(0.0, 0.0, 0.5),
        0.5,
        Material::mirror(Color::WHITE),
    ));
    world.add_object(Sphere::new(
        Vec3::new(0.0, 0.0, -0.5),
        0.5,
        Material::mirror(Color::WHITE),
    ));
    world.add_light(LightSource::Point {
        pos: Vec3::new(0.0, 5.0, 5.0),
        intensity: 1.0,
        color: Color::WHITE,
    });
    let fb = world.render(20, 20);
    assert_no_nan_pixels(&fb, "all mirror");
}

/// PPM round-trip: render → write_ppm → re-parse pixel count.
#[test]
fn ppm_round_trip() {
    let cam = Camera::new(Vec3::new(0.0, 1.0, 5.0), Vec3::ZERO, Vec3::Y, 60.0, 8, 8);
    let mut world = World::new(cam);
    world.add_object(Sphere::new(
        Vec3::ZERO,
        1.0,
        Material::matte(Color::rgb(0.5, 0.5, 0.5)),
    ));
    world.add_light(LightSource::Point {
        pos: Vec3::new(2.0, 4.0, 2.0),
        intensity: 1.0,
        color: Color::WHITE,
    });

    let fb = world.render(8, 8);
    let mut buf = Vec::new();
    fb.write_ppm(&mut buf).expect("write_ppm should not fail");

    let text = String::from_utf8(buf).expect("PPM should be valid UTF-8");
    let mut lines = text.lines();
    assert_eq!(lines.next(), Some("P3"), "PPM magic number");
    assert_eq!(lines.next(), Some("8 8"), "PPM dimensions");
    assert_eq!(lines.next(), Some("255"), "PPM max value");
    let data_lines: Vec<_> = lines.collect();
    assert_eq!(data_lines.len(), 64, "expected 64 pixel lines in PPM");
}

/// Every registered named scene must render at 100×100 without panicking.
#[test]
fn all_named_scenes_render_without_panic() {
    for name in scenes::names() {
        let build = scenes::get(name).unwrap_or_else(|| panic!("scene '{name}' not found"));
        let world = build();
        let fb = world.render(100, 100);
        assert_eq!(fb.w, 100, "scene '{name}': width mismatch");
        assert_eq!(fb.h, 100, "scene '{name}': height mismatch");
        assert_no_nan_pixels(&fb, &format!("named scene '{name}'"));
    }
}

/// Render must be pixel-for-pixel identical given the same world (determinism).
#[test]
fn render_is_deterministic() {
    let build_world = || {
        let cam = Camera::new(Vec3::new(0.0, 2.0, 6.0), Vec3::ZERO, Vec3::Y, 60.0, 32, 32);
        let mut w = World::new(cam);
        w.add_object(Sphere::new(
            Vec3::ZERO,
            1.0,
            Material::matte(Color::rgb(0.7, 0.2, 0.5)),
        ));
        w.add_object(Plane::new(
            Vec3::new(0.0, -1.0, 0.0),
            Vec3::Y,
            Material::matte(Color::rgb(0.5, 0.5, 0.5)),
        ));
        w.add_light(LightSource::Point {
            pos: Vec3::new(3.0, 5.0, 3.0),
            intensity: 1.0,
            color: Color::WHITE,
        });
        w
    };

    let fb1 = build_world().render(32, 32);
    let fb2 = build_world().render(32, 32);

    for (i, (p1, p2)) in fb1.pixels.iter().zip(fb2.pixels.iter()).enumerate() {
        assert_eq!(p1.r, p2.r, "pixel {i} r differs between renders");
        assert_eq!(p1.g, p2.g, "pixel {i} g differs between renders");
        assert_eq!(p1.b, p2.b, "pixel {i} b differs between renders");
    }
}

/// Camera position changes are reflected in the next render (different pixels).
#[test]
fn camera_position_change_affects_render() {
    let make_cam =
        |z: f32| Camera::new(Vec3::new(0.0, 1.0, z), Vec3::ZERO, Vec3::Y, 60.0, 20, 20);

    let mut world_a = World::new(make_cam(5.0));
    world_a.add_object(Sphere::new(
        Vec3::ZERO,
        1.0,
        Material::matte(Color::rgb(0.8, 0.3, 0.3)),
    ));
    world_a.add_light(LightSource::Point {
        pos: Vec3::new(2.0, 4.0, 2.0),
        intensity: 1.0,
        color: Color::WHITE,
    });

    let mut world_b = World::new(make_cam(20.0));
    world_b.add_object(Sphere::new(
        Vec3::ZERO,
        1.0,
        Material::matte(Color::rgb(0.8, 0.3, 0.3)),
    ));
    world_b.add_light(LightSource::Point {
        pos: Vec3::new(2.0, 4.0, 2.0),
        intensity: 1.0,
        color: Color::WHITE,
    });

    let fb_a = world_a.render(20, 20);
    let fb_b = world_b.render(20, 20);

    let different = fb_a
        .pixels
        .iter()
        .zip(fb_b.pixels.iter())
        .any(|(a, b)| (a.r - b.r).abs() > 1e-4 || (a.g - b.g).abs() > 1e-4);
    assert!(different, "renders from different camera positions should differ");
}

/// A scene rendered at 1×1 (minimum valid resolution) must not panic.
#[test]
fn one_pixel_render_does_not_panic() {
    let cam = Camera::new(Vec3::new(0.0, 0.0, 5.0), Vec3::ZERO, Vec3::Y, 60.0, 1, 1);
    let mut world = World::new(cam);
    world.add_object(Sphere::new(Vec3::ZERO, 1.0, Material::matte(Color::WHITE)));
    world.add_light(LightSource::Point {
        pos: Vec3::new(0.0, 5.0, 5.0),
        intensity: 1.0,
        color: Color::WHITE,
    });
    let fb = world.render(1, 1);
    assert_eq!(fb.pixels.len(), 1);
    assert_no_nan_pixels(&fb, "1x1 render");
}

/// Extreme camera FOV values (very narrow / very wide) must not produce NaN.
#[test]
fn extreme_fov_renders_without_nan() {
    for &fov in &[1.0_f32, 179.0_f32] {
        let cam = Camera::new(Vec3::new(0.0, 1.0, 5.0), Vec3::ZERO, Vec3::Y, fov, 20, 20);
        let mut world = World::new(cam);
        world.add_object(Sphere::new(Vec3::ZERO, 1.0, Material::matte(Color::WHITE)));
        world.add_light(LightSource::Point {
            pos: Vec3::new(2.0, 4.0, 2.0),
            intensity: 1.0,
            color: Color::WHITE,
        });
        let fb = world.render(20, 20);
        assert_no_nan_pixels(&fb, &format!("fov={fov}"));
    }
}

/// A point light directly at the object surface must not produce NaN.
#[test]
fn light_at_surface_does_not_produce_nan() {
    let cam = Camera::new(Vec3::new(0.0, 0.0, 5.0), Vec3::ZERO, Vec3::Y, 60.0, 20, 20);
    let mut world = World::new(cam);
    world.add_object(Sphere::new(
        Vec3::ZERO,
        1.0,
        Material::matte(Color::rgb(0.5, 0.5, 0.5)),
    ));
    world.add_light(LightSource::Point {
        pos: Vec3::new(1.0, 0.0, 0.0),
        intensity: 1.0,
        color: Color::WHITE,
    });
    let fb = world.render(20, 20);
    assert_no_nan_pixels(&fb, "light at surface");
}

/// add_lights with count=0 should be a no-op (no lights added, no panic).
#[test]
fn add_lights_zero_count_is_noop() {
    let cam = Camera::new(Vec3::new(0.0, 1.0, 5.0), Vec3::ZERO, Vec3::Y, 60.0, 10, 10);
    let mut world = World::new(cam);
    world.add_object(Sphere::new(Vec3::ZERO, 1.0, Material::matte(Color::WHITE)));
    world.add_lights(
        LightSource::Point {
            pos: Vec3::new(2.0, 4.0, 2.0),
            intensity: 1.0,
            color: Color::WHITE,
        },
        0,
    );
    assert!(world.lights.is_empty(), "no lights should have been added");
    let fb = world.render(10, 10);
    assert_no_nan_pixels(&fb, "add_lights(0)");
}

/// Soft area light (add_lights with count>1) must render without NaN.
#[test]
fn soft_area_light_renders_without_nan() {
    let cam = Camera::new(Vec3::new(0.0, 2.0, 6.0), Vec3::ZERO, Vec3::Y, 60.0, 30, 30);
    let mut world = World::new(cam);
    world.add_object(Sphere::new(Vec3::ZERO, 1.0, Material::matte(Color::WHITE)));
    world.add_lights(
        LightSource::Point {
            pos: Vec3::new(4.0, 7.0, 4.0),
            intensity: 1.0,
            color: Color::WHITE,
        },
        10,
    );
    let fb = world.render(30, 30);
    assert_no_nan_pixels(&fb, "soft area light (10 samples)");
}

/// Directional soft light must render without NaN.
#[test]
fn soft_directional_light_renders_without_nan() {
    let cam = Camera::new(Vec3::new(0.0, 2.0, 6.0), Vec3::ZERO, Vec3::Y, 60.0, 30, 30);
    let mut world = World::new(cam);
    world.add_object(Plane::new(
        Vec3::ZERO,
        Vec3::Y,
        Material::matte(Color::rgb(0.5, 0.5, 0.5)),
    ));
    world.add_lights(
        LightSource::Directional {
            dir: Vec3::new(-0.5, -1.0, -0.3),
            intensity: 0.9,
            color: Color::WHITE,
        },
        8,
    );
    let fb = world.render(30, 30);
    assert_no_nan_pixels(&fb, "soft directional light (8 samples)");
}
