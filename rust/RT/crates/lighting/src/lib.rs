use glam::Vec3;
use objects::{Object, Hit, Color, Ray, SURFACE_EPSILON};
use raycaster::Raycaster;

pub enum LightSource {
    Point { pos: Vec3, intensity: f32, color: Color },
    /// `dir` is the direction light travels (toward the scene); negated internally to get the
    /// surface-to-light direction.
    Directional { dir: Vec3, intensity: f32, color: Color },
}

pub struct SkyProperties { pub color: Color }

impl Default for SkyProperties {
    fn default() -> Self { Self { color: Color::rgb(0.5, 0.7, 1.0) } }
}

const AMBIENT: f32 = 0.5;
const CONTACT_SHADOW_MAX_HEIGHT: f32 = 0.35;
const CONTACT_SHADOW_STRENGTH: f32 = 0.8;

/// Shaded color at `hit`: ambient base + per-light Lambert diffuse with distance falloff.
/// Shadow rays are cast with `Raycaster::occluded`; occluded lights contribute nothing.
pub fn shade(hit: &Hit, lights: &[LightSource], objects: &[Box<dyn Object>]) -> Color {
    let shadow_origin = hit.point + hit.normal * SURFACE_EPSILON;
    let mut color = hit.material.color * AMBIENT;

    for light in lights {
        let (to_light_dir, dist, contribution) = match light {
            LightSource::Point { pos, intensity, color: lc } => {
                let to_light = *pos - shadow_origin;
                let dist = to_light.length();
                let dir = to_light / dist;
                (dir, dist, *lc * *intensity)
            }
            LightSource::Directional { dir, intensity, color: lc } => {
                let to_light_dir = (-*dir).normalize();
                (to_light_dir, f32::INFINITY, *lc * *intensity)
            }
        };

        let shadow_ray = Ray::new(shadow_origin, to_light_dir);
        if Raycaster::occluded(shadow_ray, dist, objects) {
            continue;
        }

        let n_dot_l = hit.normal.dot(to_light_dir).max(0.0);
        color = color + hit.material.color * contribution * n_dot_l;
    }

    color = color * contact_shadow_factor(hit, objects);

    Color::rgb(color.r.min(1.0), color.g.min(1.0), color.b.min(1.0))
}

fn contact_shadow_factor(hit: &Hit, objects: &[Box<dyn Object>]) -> f32 {
    if hit.normal.y < 0.95 {
        return 1.0;
    }

    let up_ray = Ray::new(hit.point + hit.normal * SURFACE_EPSILON, hit.normal);
    let Some(blocker) = Raycaster::nearest_hit_within(up_ray, CONTACT_SHADOW_MAX_HEIGHT, objects) else {
        return 1.0;
    };

    let proximity = 1.0 - (blocker.distance / CONTACT_SHADOW_MAX_HEIGHT);
    1.0 - CONTACT_SHADOW_STRENGTH * proximity * proximity
}

#[cfg(test)]
mod tests {
    use super::*;
    use glam::Vec3;
    use objects::{Material, Sphere};

    fn white_hit_up() -> Hit {
        Hit {
            distance: 1.0,
            point: Vec3::ZERO,
            normal: Vec3::Y,
            material: Material::matte(Color::WHITE),
        }
    }

    fn box_sphere(center: Vec3, radius: f32) -> Box<dyn Object> {
        Box::new(Sphere::new(center, radius, Material::matte(Color::WHITE)))
    }

    // ── ambient ───────────────────────────────────────────────────────────────

    #[test]
    fn ambient_only_when_no_lights() {
        let color = shade(&white_hit_up(), &[], &[]);
        assert!((color.r - AMBIENT).abs() < 1e-5, "expected ambient-only red channel");
        assert!((color.g - AMBIENT).abs() < 1e-5);
        assert!((color.b - AMBIENT).abs() < 1e-5);
    }

    // ── point light ───────────────────────────────────────────────────────────

    #[test]
    fn point_light_facing_surface_brightens_it() {
        let light = LightSource::Point {
            pos: Vec3::new(0.0, 10.0, 0.0),
            intensity: 1.0,
            color: Color::WHITE,
        };
        let color = shade(&white_hit_up(), &[light], &[]);
        assert!(color.r > AMBIENT, "facing point light should brighten the surface");
    }

    #[test]
    fn back_facing_point_light_adds_nothing() {
        let light = LightSource::Point {
            pos: Vec3::new(0.0, -10.0, 0.0),
            intensity: 1.0,
            color: Color::WHITE,
        };
        let color = shade(&white_hit_up(), &[light], &[]);
        assert!((color.r - AMBIENT).abs() < 1e-5, "back-facing light should not contribute");
    }

    #[test]
    fn occluder_blocks_point_light() {
        let light = LightSource::Point {
            pos: Vec3::new(0.0, 5.0, 0.0),
            intensity: 1.0,
            color: Color::WHITE,
        };
        let blocker = box_sphere(Vec3::new(0.0, 2.0, 0.0), 0.5);
        let color = shade(&white_hit_up(), &[light], &[blocker]);
        assert!((color.r - AMBIENT).abs() < 1e-4, "occluded light should leave only ambient");
    }

    #[test]
    fn point_light_color_tints_result() {
        let light = LightSource::Point {
            pos: Vec3::new(0.0, 10.0, 0.0),
            intensity: 1.0,
            color: Color::rgb(1.0, 0.0, 0.0),
        };
        let color = shade(&white_hit_up(), &[light], &[]);
        assert!(color.r > color.g, "red light should tint the red channel higher than green");
        assert!(color.r > color.b, "red light should tint the red channel higher than blue");
    }

    // ── directional light ─────────────────────────────────────────────────────

    #[test]
    fn directional_light_facing_surface_brightens_it() {
        let light = LightSource::Directional {
            dir: Vec3::NEG_Y,     // travels downward → surface-to-light is +Y
            intensity: 1.0,
            color: Color::WHITE,
        };
        let color = shade(&white_hit_up(), &[light], &[]);
        assert!(color.r > AMBIENT, "directional light should brighten the surface");
    }

    #[test]
    fn directional_light_dir_is_negated_correctly() {
        // dir = +Y means light travels upward away from the scene — back-facing, no contribution.
        let away = LightSource::Directional {
            dir: Vec3::Y,
            intensity: 1.0,
            color: Color::WHITE,
        };
        let color = shade(&white_hit_up(), &[away], &[]);
        assert!((color.r - AMBIENT).abs() < 1e-5, "upward-pointing directional should not hit a +Y normal");
    }

    // ── multi-light & clamping ────────────────────────────────────────────────

    #[test]
    fn two_lights_brighter_than_one() {
        let one_light = shade(
            &white_hit_up(),
            &[LightSource::Point { pos: Vec3::new(0.0, 10.0, 0.0), intensity: 0.4, color: Color::WHITE }],
            &[],
        );
        let two_lights = shade(
            &white_hit_up(),
            &[
                LightSource::Point { pos: Vec3::new(0.0, 10.0, 0.0), intensity: 0.4, color: Color::WHITE },
                LightSource::Point { pos: Vec3::new(1.0, 10.0, 0.0), intensity: 0.4, color: Color::WHITE },
            ],
            &[],
        );
        assert!(two_lights.r > one_light.r, "two lights should produce more brightness than one");
    }

    #[test]
    fn output_clamped_to_one() {
        let light = LightSource::Point {
            pos: Vec3::new(0.0, 1.0, 0.0),
            intensity: 10.0,
            color: Color::WHITE,
        };
        let color = shade(&white_hit_up(), &[light], &[]);
        assert!((color.r - 1.0).abs() < 1e-5, "output should be clamped to 1.0");
        assert!((color.g - 1.0).abs() < 1e-5);
        assert!((color.b - 1.0).abs() < 1e-5);
    }

    #[test]
    fn contact_shadow_darkens_upward_facing_receiver() {
        let blocker = box_sphere(Vec3::new(0.0, 0.2, 0.0), 0.1);
        let color = shade(&white_hit_up(), &[], &[blocker]);
        assert!(color.r < AMBIENT, "nearby geometry above the plane should darken contact area");
    }

    #[test]
    fn contact_shadow_does_not_affect_non_ground_surfaces() {
        let wall_hit = Hit {
            distance: 1.0,
            point: Vec3::ZERO,
            normal: Vec3::X,
            material: Material::matte(Color::WHITE),
        };
        let blocker = box_sphere(Vec3::new(0.2, 0.0, 0.0), 0.1);
        let color = shade(&wall_hit, &[], &[blocker]);
        assert!((color.r - AMBIENT).abs() < 1e-5, "contact shadow should stay limited to upward-facing receivers");
    }
}
