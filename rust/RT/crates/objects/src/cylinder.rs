use glam::Vec3;
use crate::{Hit, Material, Object, Ray};

const CYLINDER_EPSILON: f32 = 1e-5;

pub struct Cylinder {
    pub base: Vec3,
    pub axis: Vec3,
    pub radius: f32,
    pub height: f32,
    pub material: Material,
}

impl Cylinder {
    pub fn new(base: Vec3, axis: Vec3, radius: f32, height: f32, material: Material) -> Self {
        Self {
            base,
            axis: axis.normalize(),
            radius,
            height,
            material,
        }
    }

    fn ray_hits_bounding_sphere(&self, ray: Ray) -> bool {
        let (center, radius) = self.bounding_sphere();
        let to_center = center - ray.origin;
        let projection = to_center.dot(ray.dir);
        let center_distance_sq = to_center.dot(to_center);
        let radius_sq = radius * radius;

        if projection < 0.0 {
            return center_distance_sq <= radius_sq;
        }

        let closest_distance_sq = center_distance_sq - projection * projection;
        closest_distance_sq <= radius_sq
    }

}

impl Object for Cylinder {
    fn intersect(&self, ray: Ray) -> Option<Hit> {
        if !self.ray_hits_bounding_sphere(ray) {
            return None;
        }

        let axis = self.axis;
        let origin_offset = ray.origin - self.base;

        let dir_axis = ray.dir.dot(axis);
        let offset_axis = origin_offset.dot(axis);

        let dir_perp = ray.dir - axis * dir_axis;
        let offset_perp = origin_offset - axis * offset_axis;

        let mut best_hit: Option<Hit> = None;

        let a = dir_perp.dot(dir_perp);
        let b = 2.0 * offset_perp.dot(dir_perp);
        let c = offset_perp.dot(offset_perp) - self.radius * self.radius;

        if a > CYLINDER_EPSILON {
            let discriminant = b * b - 4.0 * a * c;
            if discriminant >= 0.0 {
                let sqrt_disc = discriminant.sqrt();
                let inv_denom = 0.5 / a;
                for t in [(-b - sqrt_disc) * inv_denom, (-b + sqrt_disc) * inv_denom] {
                    if t <= CYLINDER_EPSILON {
                        continue;
                    }

                    let axis_distance = offset_axis + t * dir_axis;
                    if !(0.0..=self.height).contains(&axis_distance) {
                        continue;
                    }

                    let point = ray.at(t);
                    let radial = point - self.base - axis * axis_distance;
                    let normal = radial.normalize();
                    let hit = Hit {
                        distance: t,
                        point,
                        normal,
                        material: self.material,
                    };

                    if best_hit
                        .as_ref()
                        .is_none_or(|current| hit.distance < current.distance)
                    {
                        best_hit = Some(hit);
                    }
                }
            }
        }

        if dir_axis.abs() > CYLINDER_EPSILON {
            for (cap_offset, normal) in [(0.0, -axis), (self.height, axis)] {
                let t = (cap_offset - offset_axis) / dir_axis;
                if t <= CYLINDER_EPSILON {
                    continue;
                }

                let point = ray.at(t);
                let cap_center = self.base + axis * cap_offset;
                let radial = point - cap_center;
                let radial_sq = radial.dot(radial);
                if radial_sq > self.radius * self.radius + CYLINDER_EPSILON {
                    continue;
                }

                let hit = Hit {
                    distance: t,
                    point,
                    normal,
                    material: self.material,
                };

                if best_hit
                    .as_ref()
                    .is_none_or(|current| hit.distance < current.distance)
                {
                    best_hit = Some(hit);
                }
            }
        }

        best_hit
    }
    fn position(&self) -> Vec3 { self.base }
    fn set_position(&mut self, p: Vec3) { self.base = p; }
    fn bounding_sphere(&self) -> (Vec3, f32) {
        let center = self.base + self.axis * (self.height * 0.5);
        let radius = ((self.height * 0.5).powi(2) + self.radius * self.radius).sqrt();
        (center, radius)
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::Color;

    fn unit_cylinder() -> Cylinder {
        Cylinder::new(
            Vec3::new(0.0, -1.0, 0.0),
            Vec3::Y,
            1.0,
            2.0,
            Material::matte(Color::WHITE),
        )
    }

    #[test]
    fn side_hit_returns_expected_point() {
        let cylinder = unit_cylinder();
        let ray = Ray::new(Vec3::new(3.0, 0.0, 0.0), Vec3::new(-1.0, 0.0, 0.0));
        let hit = cylinder.intersect(ray).expect("expected side hit");
        assert!((hit.distance - 2.0).abs() < 1e-5);
        assert!((hit.point - Vec3::new(1.0, 0.0, 0.0)).length() < 1e-5);
        assert!((hit.normal - Vec3::new(1.0, 0.0, 0.0)).length() < 1e-5);
    }

    #[test]
    fn top_cap_hit_uses_cap_normal() {
        let cylinder = unit_cylinder();
        let ray = Ray::new(Vec3::new(0.0, 3.0, 0.0), Vec3::new(0.0, -1.0, 0.0));
        let hit = cylinder.intersect(ray).expect("expected top cap hit");
        assert!((hit.point - Vec3::new(0.0, 1.0, 0.0)).length() < 1e-5);
        assert!((hit.normal - Vec3::Y).length() < 1e-5);
    }

    #[test]
    fn miss_returns_none() {
        let cylinder = unit_cylinder();
        let ray = Ray::new(Vec3::new(3.0, 3.0, 0.0), Vec3::new(1.0, 0.0, 0.0));
        assert!(cylinder.intersect(ray).is_none());
    }

    #[test]
    fn miss_outside_bounding_sphere_returns_none() {
        let cylinder = unit_cylinder();
        let ray = Ray::new(Vec3::new(5.0, 5.0, 0.0), Vec3::new(1.0, 0.0, 0.0));
        assert!(!cylinder.ray_hits_bounding_sphere(ray));
        assert!(cylinder.intersect(ray).is_none());
    }

    #[test]
    fn bounding_sphere_covers_midpoint() {
        let cylinder = unit_cylinder();
        let (center, radius) = cylinder.bounding_sphere();
        assert!((center - Vec3::ZERO).length() < 1e-5);
        assert!((radius - 2.0_f32.sqrt()).abs() < 1e-5);
    }
}
