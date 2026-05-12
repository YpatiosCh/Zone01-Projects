use crate::{Hit, Material, Object, Ray};
use glam::{Quat, Vec3};

pub struct Cube {
    pub center: Vec3,
    pub size: f32,
    pub rotation: Quat,
    pub material: Material,
}

impl Cube {
    /// Creates a new [Cube] with given center, side length, rotation and [Material].
    /// # Arguments
    /// * `center` - The position of the cube's center in 3D space
    /// * `size` - The side length of the cube
    /// * `rotation` - Rotation applied to the cube about its center
    /// * `material` - Surface color and how the cube interacts with light
    pub fn new(center: Vec3, size: f32, rotation: Quat, material: Material) -> Self {
        Self { center, size, rotation, material }
    }
}

impl Object for Cube {
    fn intersect(&self, ray: Ray) -> Option<Hit> {
        // Transform the ray into the cube's local frame, where the cube is the
        // axis-aligned box [-h, h]³. Rotation is rigid, so the ray parameter t
        // is preserved between local and world space.
        let h = self.size * 0.5;
        let inv_rot = self.rotation.inverse();
        let o = inv_rot * (ray.origin - self.center);
        let d = inv_rot * ray.dir;

        let o = [o.x, o.y, o.z];
        let d = [d.x, d.y, d.z];

        let mut t_min = f32::NEG_INFINITY;
        let mut t_max = f32::INFINITY;
        let mut entry_axis = 0usize;
        let mut entry_sign = 0.0f32;
        let mut exit_axis = 0usize;
        let mut exit_sign = 0.0f32;

        for i in 0..3 {
            if d[i].abs() < 1e-8 {
                // Ray parallel to this slab — no intersection unless origin is between the planes.
                if o[i] < -h || o[i] > h {
                    return None;
                }
                continue;
            }
            let inv_d = 1.0 / d[i];
            let t1 = (-h - o[i]) * inv_d;
            let t2 = (h - o[i]) * inv_d;
            let (t_near, t_far, near_sign, far_sign) = if t1 < t2 {
                (t1, t2, -1.0, 1.0)
            } else {
                (t2, t1, 1.0, -1.0)
            };
            if t_near > t_min {
                t_min = t_near;
                entry_axis = i;
                entry_sign = near_sign;
            }
            if t_far < t_max {
                t_max = t_far;
                exit_axis = i;
                exit_sign = far_sign;
            }
            if t_min > t_max {
                return None;
            }
        }

        let (t, axis, sign) = if t_min > 0.0 {
            (t_min, entry_axis, entry_sign)
        } else if t_max > 0.0 {
            // Ray origin is inside the cube — return the exit face.
            (t_max, exit_axis, exit_sign)
        } else {
            return None;
        };

        let mut n = [0.0; 3];
        n[axis] = sign;
        let normal_local = Vec3::new(n[0], n[1], n[2]);
        let normal = self.rotation * normal_local;
        let point = ray.at(t);

        Some(Hit { distance: t, point, normal, material: self.material })
    }

    fn position(&self) -> Vec3 { self.center }
    fn set_position(&mut self, p: Vec3) { self.center = p; }
    fn bounding_sphere(&self) -> (Vec3, f32) {
        // Half-diagonal of a cube with side `size`, rotation-invariant.
        (self.center, self.size * 0.5 * 3f32.sqrt())
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::Color;

    fn unit_cube() -> Cube {
        Cube::new(Vec3::ZERO, 2.0, Quat::IDENTITY, Material::matte(Color::WHITE))
    }

    #[test]
    fn head_on_hit_returns_near_face() {
        // Cube spans [-1,1]³. Ray from z=5 along -z hits +Z face at z=1, t=4.
        let cube = unit_cube();
        let ray = Ray::new(Vec3::new(0.0, 0.0, 5.0), Vec3::new(0.0, 0.0, -1.0));
        let hit = cube.intersect(ray).expect("expected a hit");
        assert!((hit.distance - 4.0).abs() < 1e-5);
        assert!((hit.point - Vec3::new(0.0, 0.0, 1.0)).length() < 1e-5);
        // Outward normal on +Z face is +z.
        assert!((hit.normal - Vec3::new(0.0, 0.0, 1.0)).length() < 1e-5);
        assert!((hit.normal.length() - 1.0).abs() < 1e-5);
    }

    #[test]
    fn miss_returns_none() {
        let cube = unit_cube();
        // Ray parallel to -z but offset far in +x — never enters the cube.
        let ray = Ray::new(Vec3::new(5.0, 0.0, 5.0), Vec3::new(0.0, 0.0, -1.0));
        assert!(cube.intersect(ray).is_none());
    }

    #[test]
    fn ray_pointing_away_returns_none() {
        let cube = unit_cube();
        let ray = Ray::new(Vec3::new(0.0, 0.0, 5.0), Vec3::new(0.0, 0.0, 1.0));
        assert!(cube.intersect(ray).is_none());
    }

    #[test]
    fn ray_from_inside_returns_exit_face() {
        // Origin at center, pointing -z → exits at z=-1 (the -Z face), t=1, normal -z.
        let cube = unit_cube();
        let ray = Ray::new(Vec3::ZERO, Vec3::new(0.0, 0.0, -1.0));
        let hit = cube.intersect(ray).expect("ray from inside should hit exit face");
        assert!((hit.distance - 1.0).abs() < 1e-5);
        assert!((hit.point - Vec3::new(0.0, 0.0, -1.0)).length() < 1e-5);
        assert!((hit.normal - Vec3::new(0.0, 0.0, -1.0)).length() < 1e-5);
    }

    #[test]
    fn rotated_cube_normal_rotates_with_it() {
        // Rotate 30° around Y so the camera (along -z) clearly hits the rotated
        // former-+Z face (45° would aim straight at an edge — ambiguous).
        // Whatever chirality glam uses, the normal must lie in the XZ plane,
        // be unit length, and project onto +Z by cos(30°).
        let theta = std::f32::consts::FRAC_PI_6;
        let cube = Cube::new(
            Vec3::ZERO,
            2.0,
            Quat::from_rotation_y(theta),
            Material::matte(Color::WHITE),
        );
        let ray = Ray::new(Vec3::new(0.0, 0.0, 5.0), Vec3::new(0.0, 0.0, -1.0));
        let hit = cube.intersect(ray).expect("expected a hit");
        assert!((hit.normal.length() - 1.0).abs() < 1e-5);
        assert!(hit.normal.y.abs() < 1e-5, "normal should stay in XZ plane");
        assert!((hit.normal.dot(Vec3::Z) - theta.cos()).abs() < 1e-4);
        assert!((hit.normal.x.abs() - theta.sin()).abs() < 1e-4);
    }

    #[test]
    fn off_center_cube_translates_correctly() {
        // Move the cube to (10, 0, 0). Ray from (10, 0, 5) along -z hits +Z face at z=1.
        let cube = Cube::new(
            Vec3::new(10.0, 0.0, 0.0),
            2.0,
            Quat::IDENTITY,
            Material::matte(Color::WHITE),
        );
        let ray = Ray::new(Vec3::new(10.0, 0.0, 5.0), Vec3::new(0.0, 0.0, -1.0));
        let hit = cube.intersect(ray).expect("expected a hit");
        assert!((hit.distance - 4.0).abs() < 1e-5);
        assert!((hit.point - Vec3::new(10.0, 0.0, 1.0)).length() < 1e-5);
    }

    #[test]
    fn bounding_sphere_is_half_diagonal() {
        let cube = Cube::new(
            Vec3::new(1.0, 2.0, 3.0),
            2.0,
            Quat::IDENTITY,
            Material::matte(Color::WHITE),
        );
        let (center, radius) = cube.bounding_sphere();
        assert_eq!(center, Vec3::new(1.0, 2.0, 3.0));
        // Side 2 → half-diagonal sqrt(3).
        assert!((radius - 3f32.sqrt()).abs() < 1e-5);
    }

    #[test]
    fn material_is_stamped_in() {
        let cube = Cube::new(
            Vec3::ZERO,
            2.0,
            Quat::IDENTITY,
            Material::mirror(Color::rgb(0.1, 0.2, 0.3)),
        );
        let ray = Ray::new(Vec3::new(0.0, 0.0, 5.0), Vec3::new(0.0, 0.0, -1.0));
        let hit = cube.intersect(ray).expect("expected a hit");
        assert!(hit.material.reflective);
        assert!((hit.material.color.b - 0.3).abs() < 1e-5);
    }
}
