use std::cmp::Ordering;
use objects::{Color, Object, Ray, Hit, SURFACE_EPSILON, reflect};

pub struct Raycaster;

impl Raycaster {
    fn trace(ray: Ray, objects: &[Box<dyn Object>]) -> Option<Hit> {
        objects
            .iter()
            .filter_map(|obj| obj.intersect(ray))
            .filter(|h| h.distance > 0.0)
            .min_by(|a, b| a.distance.partial_cmp(&b.distance).unwrap_or(Ordering::Equal))
    }

    pub fn nearest_hit_within(ray: Ray, max_distance: f32, objects: &[Box<dyn Object>]) -> Option<Hit> {
        Self::trace(ray, objects).filter(|h| h.distance < max_distance)
    }

    /// Return true if anything blocks `ray` within `max_distance`. Early-outs on first hit.
    ///
    /// Typical use in `lighting::shade`:
    /// ```ignore
    /// let shadow_origin = hit.point + hit.normal * SURFACE_EPSILON;
    /// let to_light      = light_pos - shadow_origin;
    /// let shadow_ray    = Ray::new(shadow_origin, to_light);
    /// if Raycaster::occluded(shadow_ray, to_light.length(), objects) {
    ///     continue; // in shadow — skip this light's contribution
    /// }
    /// ```
    pub fn occluded(ray: Ray, max_distance: f32, objects: &[Box<dyn Object>]) -> bool {
        Self::nearest_hit_within(ray, max_distance, objects).is_some()
    }

    /// Trace `ray` with mirror bounces up to `max_depth`.
    ///
    /// Returns `(tint, terminal_hit)`:
    /// - `tint` — product of every reflective surface's `material.color` along the bounce
    ///   chain. Start with `Color::WHITE`; each mirror multiplies by its own color.
    /// - `terminal_hit` — the first non-reflective `Hit`, or `None` if the ray escapes to
    ///   sky or `max_depth` bounces are exhausted.
    ///
    /// `World::render` uses this as:
    /// ```ignore
    /// let (tint, terminal) = Raycaster::trace_with_reflections(primary_ray, &visible, 8);
    /// let base_color = match terminal {
    ///     Some(hit) => lighting::shade(&hit, &world.lights, &visible),
    ///     None      => world.sky.color,
    /// };
    /// pixel_color = base_color * tint;
    /// ```
    ///
    /// **Bounce rules:**
    /// 1. Trace the current ray.
    /// 2. No hit → return `(tint, None)`.
    /// 3. Non-reflective hit → return `(tint, Some(hit))`.
    /// 4. Reflective hit → multiply `tint` by `hit.material.color`, spawn reflected ray,
    ///    repeat.
    /// 5. `max_depth` exhausted → return `(tint, None)` (treat as escaped to sky).
    pub fn trace_with_reflections(
        mut ray: Ray,
        objects: &[Box<dyn Object>],
        max_depth: u32,
    ) -> (Color, Option<Hit>) {
        let mut tint = Color::WHITE;
        for _ in 0..max_depth {
            match Self::trace(ray, objects) {
                None => return (tint, None),
                Some(hit) if !hit.material.reflective => return (tint, Some(hit)),
                Some(hit) => {
                    tint = tint * hit.material.color;
                    let origin = hit.point + hit.normal * SURFACE_EPSILON;
                    ray = Ray::new(origin, reflect(ray.dir, hit.normal));
                }
            }
        }
        (tint, None)
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use glam::Vec3;
    use objects::{Color, Material, Sphere};

    fn box_sphere(center: Vec3, radius: f32, material: Material) -> Box<dyn Object> {
        Box::new(Sphere::new(center, radius, material))
    }

    // ── trace ─────────────────────────────────────────────────────────────────

    #[test]
    fn trace_nearest_hit() {
        let near = box_sphere(Vec3::new(0.0, 0.0, 2.0), 0.5, Material::matte(Color::WHITE));
        let far  = box_sphere(Vec3::new(0.0, 0.0, -3.0), 0.5, Material::matte(Color::WHITE));
        let ray  = Ray::new(Vec3::new(0.0, 0.0, 5.0), Vec3::new(0.0, 0.0, -1.0));

        let hit = Raycaster::trace(ray, &[near, far]).expect("expected a hit");
        // Near sphere surface is at z ≈ 2.5; far at z ≈ -2.5. t_near < t_far.
        assert!(hit.distance < 5.0, "should hit the nearer sphere first");
    }

    #[test]
    fn trace_no_hit_returns_none() {
        let sphere = box_sphere(Vec3::new(10.0, 0.0, 0.0), 0.5, Material::matte(Color::WHITE));
        let ray    = Ray::new(Vec3::ZERO, Vec3::new(0.0, 0.0, -1.0));
        assert!(Raycaster::trace(ray, &[sphere]).is_none());
    }

    #[test]
    fn trace_ignores_hits_behind_origin() {
        let sphere = box_sphere(Vec3::new(0.0, 0.0, 5.0), 0.5, Material::matte(Color::WHITE));
        let ray    = Ray::new(Vec3::ZERO, Vec3::new(0.0, 0.0, -1.0)); // pointing away
        assert!(Raycaster::trace(ray, &[sphere]).is_none());
    }

    // ── occluded ──────────────────────────────────────────────────────────────

    #[test]
    fn occluded_blocked() {
        let blocker = box_sphere(Vec3::new(0.0, 0.0, 0.0), 0.5, Material::matte(Color::WHITE));
        let ray     = Ray::new(Vec3::new(0.0, 0.0, 3.0), Vec3::new(0.0, 0.0, -1.0));
        assert!(Raycaster::occluded(ray, 10.0, &[blocker]));
    }

    #[test]
    fn occluded_beyond_max_distance() {
        let sphere = box_sphere(Vec3::new(0.0, 0.0, 0.0), 0.5, Material::matte(Color::WHITE));
        let ray    = Ray::new(Vec3::new(0.0, 0.0, 3.0), Vec3::new(0.0, 0.0, -1.0));
        // max_distance of 1.0 — sphere is ~2.5 units away, so not occluded.
        assert!(!Raycaster::occluded(ray, 1.0, &[sphere]));
    }

    // ── trace_with_reflections ────────────────────────────────────────────────

    #[test]
    fn nearest_hit_within_returns_hit_distance() {
        let sphere = box_sphere(Vec3::new(0.0, 0.0, 0.0), 0.5, Material::matte(Color::WHITE));
        let ray    = Ray::new(Vec3::new(0.0, 0.0, 3.0), Vec3::new(0.0, 0.0, -1.0));
        let hit = Raycaster::nearest_hit_within(ray, 10.0, &[sphere]).expect("expected hit");
        assert!((hit.distance - 2.5).abs() < 1e-4, "nearest hit distance should be preserved");
    }

    #[test]
    fn nearest_hit_within_respects_max_distance() {
        let sphere = box_sphere(Vec3::new(0.0, 0.0, 0.0), 0.5, Material::matte(Color::WHITE));
        let ray    = Ray::new(Vec3::new(0.0, 0.0, 3.0), Vec3::new(0.0, 0.0, -1.0));
        assert!(Raycaster::nearest_hit_within(ray, 1.0, &[sphere]).is_none());
    }

    #[test]
    fn reflections_non_reflective_terminal() {
        let sphere = box_sphere(Vec3::new(0.0, 0.0, 0.0), 1.0, Material::matte(Color::WHITE));
        let ray    = Ray::new(Vec3::new(0.0, 0.0, 5.0), Vec3::new(0.0, 0.0, -1.0));
        let (tint, hit) = Raycaster::trace_with_reflections(ray, &[sphere], 4);
        assert!(hit.is_some());
        assert!((tint.r - 1.0).abs() < 1e-5); // no bounces → tint stays white
    }

    #[test]
    fn reflections_no_hit_returns_sky() {
        let ray = Ray::new(Vec3::ZERO, Vec3::new(0.0, 1.0, 0.0)); // aimed at sky
        let (tint, hit) = Raycaster::trace_with_reflections(ray, &[], 4);
        assert!(hit.is_none());
        assert!((tint.r - 1.0).abs() < 1e-5);
    }

    #[test]
    fn reflections_mirror_accumulates_tint() {
        // Red mirror in front, white matte behind it.
        let mirror = box_sphere(
            Vec3::new(0.0, 0.0, 1.0), 0.4,
            Material::mirror(Color::rgb(1.0, 0.0, 0.0)),
        );
        let back = box_sphere(
            Vec3::new(0.0, 0.0, -2.0), 1.0,
            Material::matte(Color::WHITE),
        );
        let ray = Ray::new(Vec3::new(0.0, 0.0, 5.0), Vec3::new(0.0, 0.0, -1.0));
        let (tint, _hit) = Raycaster::trace_with_reflections(ray, &[mirror, back], 4);
        assert!(tint.r > 0.5, "red channel should survive mirror tint");
        assert!(tint.g < 0.1, "green should be near 0 after red mirror");
    }

    #[test]
    fn reflections_max_depth_terminates() {
        // Two mirrors facing each other would loop forever without a depth cap.
        let m1 = box_sphere(Vec3::new(0.0, 0.0,  1.0), 0.4, Material::mirror(Color::WHITE));
        let m2 = box_sphere(Vec3::new(0.0, 0.0, -1.0), 0.4, Material::mirror(Color::WHITE));
        let ray = Ray::new(Vec3::new(0.0, 0.0, 3.0), Vec3::new(0.0, 0.0, -1.0));
        let (_tint, hit) = Raycaster::trace_with_reflections(ray, &[m1, m2], 4);
        assert!(hit.is_none()); // exceeded depth → sky
    }
}
