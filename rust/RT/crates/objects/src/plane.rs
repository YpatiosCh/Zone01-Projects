use glam::Vec3;
use crate::{Object, Ray, Hit, Material};

const PLANE_EPSILON: f32 = 1e-5;

pub struct Plane {
    pub point: Vec3,
    pub normal: Vec3,
    pub material: Material,
}

impl Plane {
    pub fn new(point: Vec3, normal: Vec3, material: Material) -> Self {
        Self {
            point,
            normal: normal.normalize(),
            material,
        }
    }
}

impl Object for Plane {
    fn intersect(&self, ray: Ray) -> Option<Hit> {
        let denom = self.normal.dot(ray.dir);
        if denom.abs() <= PLANE_EPSILON {
            return None;
        }

        let distance = (self.point - ray.origin).dot(self.normal) / denom;
        if distance <= PLANE_EPSILON {
            return None;
        }

        let outward_normal = if denom < 0.0 {
            self.normal
        } else {
            -self.normal
        };

        Some(Hit {
            distance,
            point: ray.at(distance),
            normal: outward_normal,
            material: self.material,
        })
    }
    fn position(&self) -> Vec3 { self.point }
    fn set_position(&mut self, p: Vec3) { self.point = p; }
    fn bounding_sphere(&self) -> (Vec3, f32) { (self.point, f32::INFINITY) }
    fn floor_height_at(&self, position: Vec3) -> Option<f32> {
        if self.normal.y <= PLANE_EPSILON {
            return None;
        }

        let xz_offset = Vec3::new(position.x - self.point.x, 0.0, position.z - self.point.z);
        let y = self.point.y - (self.normal.x * xz_offset.x + self.normal.z * xz_offset.z) / self.normal.y;
        Some(y)
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::Color;

    fn ground_plane() -> Plane {
        Plane::new(
            Vec3::ZERO,
            Vec3::Y,
            Material::matte(Color::rgb(0.5, 0.5, 0.5)),
        )
    }

    #[test]
    fn downward_ray_hits_plane() {
        let plane = ground_plane();
        let ray = Ray::new(Vec3::new(0.0, 2.0, 0.0), Vec3::new(0.0, -1.0, 0.0));
        let hit = plane.intersect(ray).expect("expected plane hit");
        assert!((hit.distance - 2.0).abs() < 1e-5);
        assert!((hit.point - Vec3::ZERO).length() < 1e-5);
        assert!((hit.normal - Vec3::Y).length() < 1e-5);
    }

    #[test]
    fn parallel_ray_misses_plane() {
        let plane = ground_plane();
        let ray = Ray::new(Vec3::new(0.0, 2.0, 0.0), Vec3::new(1.0, 0.0, 0.0));
        assert!(plane.intersect(ray).is_none());
    }

    #[test]
    fn upward_ray_from_below_flips_normal() {
        let plane = ground_plane();
        let ray = Ray::new(Vec3::new(0.0, -2.0, 0.0), Vec3::new(0.0, 1.0, 0.0));
        let hit = plane.intersect(ray).expect("expected plane hit from below");
        assert!((hit.normal + Vec3::Y).length() < 1e-5);
    }
}
