use crate::{Color, Hit, Material, Object, Ray};
use glam::Vec3;
use rand::{thread_rng, Rng};

pub struct Sphere {
    pub center: Vec3,
    pub radius: f32,
    pub material: Material,
}

impl Sphere {
    /// Creates a new [Sphere] with given center, radius and [Material] properties.
    /// # Arguments
    /// * `center` - The position of the sphere in 3D space
    /// * `radius` - The size of the sphere
    /// * [Material] - Defines its color and how the sphere interacts with light
    pub fn new(center: Vec3, radius: f32, material: Material) -> Self {
        Self {
            center,
            radius,
            material,
        }
    }

    /// Random [Sphere] with center in `[-5, 5]³`, radius in `[0.2, 1.5]`,
    /// random RGB color, ~20% chance of being a mirror. Uses a thread-local RNG.
    pub fn random() -> Self {
        Self::random_with(&mut thread_rng())
    }

    /// Same distribution as [`Sphere::random`] but draws from `rng`. The e2e
    /// fuzz harness passes a seeded `StdRng` here so failing scenes can be
    /// reproduced from the seed alone.
    pub fn random_with<R: Rng + ?Sized>(rng: &mut R) -> Self {
        let center = Vec3::new(
            rng.gen_range(-5.0..=5.0),
            rng.gen_range(-5.0..=5.0),
            rng.gen_range(-5.0..=5.0),
        );
        let radius = rng.gen_range(0.2..=1.5);
        let color = Color::rgb(rng.gen(), rng.gen(), rng.gen());
        let material = if rng.gen_bool(0.2) {
            Material::mirror(color)
        } else {
            Material::matte(color)
        };
        Self {
            center,
            radius,
            material,
        }
    }
}

impl Object for Sphere {
    fn intersect(&self, ray: Ray) -> Option<Hit> {
        // Vector from sphere center to ray origin.
        // This shifts the problem so the sphere is at the origin, simplifying the math.
        let oc = ray.origin - self.center;

        // Coefficients of the quadratic equation (t² + 2bt + c = 0)
        // derived from substituting the ray equation into the sphere equation.
        // Note: the leading `a` coefficient is 1 because ray.dir is a unit vector.
        let b = oc.dot(ray.dir);
        let c = oc.dot(oc) - self.radius * self.radius;

        // The discriminant tells us how many solutions exist:
        //   < 0 → no intersection (ray misses the sphere entirely)
        //   = 0 → one intersection (ray is tangent, grazes the surface)
        //   > 0 → two intersections (ray enters and exits the sphere)
        let disc = b * b - c;
        if disc < 0.0 {
            return None;
        }

        let sqrt_disc = disc.sqrt();

        // The two intersection distances along the ray.
        // t_near is always the smaller value (entry point).
        // t_far is always the larger value (exit point).
        let t_near = -b - sqrt_disc;
        let t_far = -b + sqrt_disc;

        // Pick the closest intersection that is in front of the camera (t > 0).
        // Negative t means the hit is behind the ray origin — we ignore those.
        let t = if t_near > 0.0 {
            t_near // Normal case: ray hits the front face of the sphere
        } else if t_far > 0.0 {
            t_far // Edge case: ray origin is inside the sphere, use exit point
        } else {
            return None; // Sphere is entirely behind the ray origin
        };

        // The actual 3D point on the sphere surface where the hit occurred.
        let point = ray.at(t);

        // Outward normal at the hit point: direction from center to surface point.
        // Dividing by radius normalizes it to unit length.
        let normal = (point - self.center) / self.radius;

        Some(Hit {
            distance: t,
            point,
            normal,
            material: self.material,
        })
    }
    fn position(&self) -> Vec3 {
        self.center
    }
    fn set_position(&mut self, p: Vec3) {
        self.center = p;
    }
    fn bounding_sphere(&self) -> (Vec3, f32) {
        (self.center, self.radius)
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::Color;

    fn unit_sphere() -> Sphere {
        Sphere {
            center: Vec3::ZERO,
            radius: 1.0,
            material: Material::matte(Color::WHITE),
        }
    }

    #[test]
    fn head_on_hit_returns_near_face() {
        let sphere = unit_sphere();
        let ray = Ray::new(Vec3::new(0.0, 0.0, 5.0), Vec3::new(0.0, 0.0, -1.0));
        let hit = sphere.intersect(ray).expect("expected a hit");
        // Sphere at origin radius 1, ray from z=5 toward -z → near face at z=1, t=4.
        assert!((hit.distance - 4.0).abs() < 1e-5);
        assert!((hit.point - Vec3::new(0.0, 0.0, 1.0)).length() < 1e-5);
    }

    #[test]
    fn miss_returns_none() {
        let sphere = unit_sphere();
        // Ray parallel to -z but offset far in +x → never crosses the unit sphere.
        let ray = Ray::new(Vec3::new(5.0, 0.0, 5.0), Vec3::new(0.0, 0.0, -1.0));
        assert!(sphere.intersect(ray).is_none());
    }

    #[test]
    fn ray_pointing_away_returns_none() {
        let sphere = unit_sphere();
        // Origin already past the sphere along its direction.
        let ray = Ray::new(Vec3::new(0.0, 0.0, 5.0), Vec3::new(0.0, 0.0, 1.0));
        assert!(sphere.intersect(ray).is_none());
    }

    #[test]
    fn tangent_ray_hits() {
        let sphere = unit_sphere();
        // Ray along -z at x=1 grazes the sphere — disc == 0, single root.
        let ray = Ray::new(Vec3::new(1.0, 0.0, 5.0), Vec3::new(0.0, 0.0, -1.0));
        let hit = sphere
            .intersect(ray)
            .expect("tangent should still register a hit");
        assert!((hit.point - Vec3::new(1.0, 0.0, 0.0)).length() < 1e-4);
    }

    #[test]
    fn ray_from_inside_returns_far_face() {
        let sphere = unit_sphere();
        // Origin at center, pointing along -z → exits at z=-1, t=1.
        let ray = Ray::new(Vec3::ZERO, Vec3::new(0.0, 0.0, -1.0));
        let hit = sphere
            .intersect(ray)
            .expect("ray from inside should hit far face");
        assert!((hit.distance - 1.0).abs() < 1e-5);
        assert!((hit.point - Vec3::new(0.0, 0.0, -1.0)).length() < 1e-5);
    }

    #[test]
    fn normal_is_outward_and_unit_length() {
        let sphere = Sphere {
            center: Vec3::new(2.0, 0.0, 0.0),
            radius: 1.0,
            material: Material::matte(Color::WHITE),
        };
        let ray = Ray::new(Vec3::new(2.0, 0.0, 5.0), Vec3::new(0.0, 0.0, -1.0));
        let hit = sphere.intersect(ray).expect("expected a hit");
        // Hits near face at (2,0,1); outward normal should be +z.
        assert!((hit.normal.length() - 1.0).abs() < 1e-5);
        assert!((hit.normal - Vec3::new(0.0, 0.0, 1.0)).length() < 1e-5);
    }

    #[test]
    fn random_is_in_bounds() {
        for _ in 0..50 {
            let s = Sphere::random();
            assert!(s.center.x >= -5.0 && s.center.x <= 5.0);
            assert!(s.center.y >= -5.0 && s.center.y <= 5.0);
            assert!(s.center.z >= -5.0 && s.center.z <= 5.0);
            assert!(s.radius >= 0.2 && s.radius <= 1.5);
            for ch in [s.material.color.r, s.material.color.g, s.material.color.b] {
                assert!((0.0..=1.0).contains(&ch));
            }
        }
    }

    #[test]
    fn random_with_same_seed_is_reproducible() {
        use rand::SeedableRng;
        let mut a = rand::rngs::StdRng::seed_from_u64(42);
        let mut b = rand::rngs::StdRng::seed_from_u64(42);
        let sa = Sphere::random_with(&mut a);
        let sb = Sphere::random_with(&mut b);
        assert_eq!(sa.center, sb.center);
        assert_eq!(sa.radius, sb.radius);
        assert_eq!(sa.material.color.r, sb.material.color.r);
        assert_eq!(sa.material.reflective, sb.material.reflective);
    }

    #[test]
    fn material_is_stamped_in() {
        let sphere = Sphere {
            center: Vec3::ZERO,
            radius: 1.0,
            material: Material::mirror(Color::rgb(0.2, 0.4, 0.6)),
        };
        let ray = Ray::new(Vec3::new(0.0, 0.0, 5.0), Vec3::new(0.0, 0.0, -1.0));
        let hit = sphere.intersect(ray).expect("expected a hit");
        assert!(hit.material.reflective);
        assert!((hit.material.color.r - 0.2).abs() < 1e-5);
    }
}
