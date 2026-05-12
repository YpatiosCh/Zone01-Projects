use std::cmp::Ordering;
use glam::Vec3;
use objects::{Object, Ray};

pub struct Camera {
    pub position: Vec3,
    pub look_at: Vec3,
    pub up: Vec3,
    pub fov: f32,
    pub width: u32,
    pub height: u32,
    // Precomputed from width/height/fov — avoids recomputing per pixel.
    aspect: f32,
    forward: Vec3,
    right: Vec3,
    up_basis: Vec3,
    half_height: f32,
}

impl Camera {
    /// Construct a camera from position, look-at target, up vector, vertical FOV (degrees),
    /// and pixel dimensions. Aspect ratio is derived as `width / height`.
    pub fn new(position: Vec3, look_at: Vec3, up: Vec3, fov: f32, width: u32, height: u32) -> Self {
        let forward = (look_at - position).normalize();
        let right = forward.cross(up).normalize();
        let up_basis = right.cross(forward).normalize();
        let half_height = (fov.to_radians() / 2.0).tan();
        let aspect = width as f32 / height as f32;
        Self { position, look_at, up, fov, width, height, aspect, forward, right, up_basis, half_height }
    }

    /// Move the camera to a new position. Recomputes the view basis and forward direction;
    /// the next `ray_for_pixel` and `frustum` calls will reflect the change immediately.
    pub fn set_position(&mut self, position: Vec3) {
        self.position = position;
        self.recompute_basis();
    }

    /// Point the camera at a new target. Recomputes the view basis and forward direction;
    /// the next `ray_for_pixel` and `frustum` calls will reflect the change immediately.
    pub fn set_look_at(&mut self, look_at: Vec3) {
        self.look_at = look_at;
        self.recompute_basis();
    }

    fn recompute_basis(&mut self) {
        self.forward = (self.look_at - self.position).normalize();
        self.right = self.forward.cross(self.up).normalize();
        self.up_basis = self.right.cross(self.forward).normalize();
    }

    /// Primary ray through pixel `(x, y)`.
    pub fn ray_for_pixel(&self, x: u32, y: u32) -> Ray {
        let u = (x as f32 + 0.5) / self.width as f32;
        let v = (y as f32 + 0.5) / self.height as f32;
        let sx = (2.0 * u - 1.0) * self.aspect * self.half_height;
        let sy = (1.0 - 2.0 * v) * self.half_height;
        let dir = self.forward + self.right * sx + self.up_basis * sy;
        Ray::new(self.position, dir)
    }

    /// Full-screen view frustum (4 side planes from the screen corners).
    pub fn frustum(&self) -> Frustum {
        let corner = |u: f32, v: f32| -> Vec3 {
            let sx = (2.0 * u - 1.0) * self.aspect * self.half_height;
            let sy = (1.0 - 2.0 * v) * self.half_height;
            (self.forward + self.right * sx + self.up_basis * sy).normalize()
        };
        let tl = corner(0.0, 0.0);
        let tr = corner(1.0, 0.0);
        let bl = corner(0.0, 1.0);
        let br = corner(1.0, 1.0);

        let make_plane = |n: Vec3| -> FrustumPlane {
            let normal = n.normalize();
            FrustumPlane { normal, d: -normal.dot(self.position) }
        };

        Frustum {
            planes: [
                make_plane(bl.cross(tl)), // left
                make_plane(tr.cross(br)), // right
                make_plane(tl.cross(tr)), // top
                make_plane(br.cross(bl)), // bottom
            ],
            origin: self.position,
            forward: self.forward,
        }
    }

    /// Sub-frustum covering the pixel rect `[x0,x1) × [y0,y1)`.
    pub fn tile_frustum(&self, _x0: u32, _y0: u32, _x1: u32, _y1: u32) -> Frustum {
        todo!()
    }
}

/// A single half-space bounding the view frustum. Renamed from `Plane` to avoid
/// colliding with `objects::Plane` (the geometric shape) when re-exported by `world`.
#[derive(Clone, Copy)]
pub struct FrustumPlane {
    pub normal: Vec3,
    /// Plane constant: normal · p + d >= 0 for points on the inside.
    pub d: f32,
}

pub struct Frustum {
    pub planes: [FrustumPlane; 4],
    pub origin: Vec3,
    pub forward: Vec3,
}

impl Frustum {
    /// True if a sphere is at least partially inside all 4 side planes.
    pub fn contains(&self, center: Vec3, radius: f32) -> bool {
        self.planes.iter().all(|p| p.normal.dot(center) + p.d >= -radius)
    }

    /// Signed distance from the camera origin to `p` along the forward axis (sort key).
    pub fn depth(&self, p: Vec3) -> f32 {
        (p - self.origin).dot(self.forward)
    }

    /// Cull `objects` to those whose bounding sphere overlaps this frustum, then sort
    /// near→far by `depth`. Returns indices into the original slice in render order.
    ///
    /// Call once per frame before the pixel loop:
    /// ```ignore
    /// let frustum = world.camera.frustum();
    /// let order   = frustum.cull_and_sort(&world.objects);
    /// let visible: Vec<_> = order.iter().map(|&i| &world.objects[i]).collect();
    /// ```
    pub fn cull_and_sort(&self, objects: &[Box<dyn Object>]) -> Vec<usize> {
        let mut indices: Vec<usize> = (0..objects.len())
            .filter(|&i| {
                let (center, radius) = objects[i].bounding_sphere();
                self.contains(center, radius)
            })
            .collect();
        indices.sort_by(|&a, &b| {
            let da = self.depth(objects[a].position());
            let db = self.depth(objects[b].position());
            da.partial_cmp(&db).unwrap_or(Ordering::Equal)
        });
        indices
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use glam::Vec3;
    use objects::{Color, Material, Object, Sphere};

    /// 90° vertical FOV, square aspect (100×100). Used for frustum tests where
    /// the specific resolution doesn't matter, only the 1:1 aspect.
    fn cam90() -> Camera {
        Camera::new(Vec3::new(0.0, 0.0, 5.0), Vec3::ZERO, Vec3::Y, 90.0, 100, 100)
    }

    fn box_sphere(center: Vec3, radius: f32, material: Material) -> Box<dyn Object> {
        Box::new(Sphere::new(center, radius, material))
    }

    // ── ray_for_pixel ─────────────────────────────────────────────────────────

    #[test]
    fn ray_origin_is_camera_position() {
        let cam = Camera::new(Vec3::new(0.0, 0.0, 5.0), Vec3::ZERO, Vec3::Y, 90.0, 800, 600);
        let ray = cam.ray_for_pixel(0, 0);
        assert!((ray.origin - cam.position).length() < 1e-5);
    }

    #[test]
    fn center_pixel_points_forward() {
        // 1×1 image → single pixel IS the center → direction must be (0,0,-1).
        let cam = Camera::new(Vec3::new(0.0, 0.0, 5.0), Vec3::ZERO, Vec3::Y, 90.0, 1, 1);
        let ray = cam.ray_for_pixel(0, 0);
        assert!((ray.dir - Vec3::new(0.0, 0.0, -1.0)).length() < 1e-5);
    }

    #[test]
    fn left_pixel_has_negative_x() {
        let cam = Camera::new(Vec3::new(0.0, 0.0, 5.0), Vec3::ZERO, Vec3::Y, 90.0, 2, 1);
        let left  = cam.ray_for_pixel(0, 0);
        let right = cam.ray_for_pixel(1, 0);
        assert!(left.dir.x < 0.0);
        assert!(right.dir.x > 0.0);
    }

    #[test]
    fn top_pixel_has_positive_y() {
        let cam = Camera::new(Vec3::new(0.0, 0.0, 5.0), Vec3::ZERO, Vec3::Y, 90.0, 1, 2);
        let top    = cam.ray_for_pixel(0, 0);
        let bottom = cam.ray_for_pixel(0, 1);
        assert!(top.dir.y > 0.0);
        assert!(bottom.dir.y < 0.0);
    }

    #[test]
    fn ray_direction_is_unit_length() {
        let cam = Camera::new(Vec3::new(0.0, 0.0, 5.0), Vec3::ZERO, Vec3::Y, 90.0, 800, 600);
        for (x, y) in [(0u32, 0u32), (799, 599), (400, 300)] {
            let ray = cam.ray_for_pixel(x, y);
            assert!((ray.dir.length() - 1.0).abs() < 1e-5, "dir not normalized at ({x},{y})");
        }
    }

    // ── Frustum::depth ────────────────────────────────────────────────────────

    #[test]
    fn depth_of_look_at_target_is_positive() {
        let frustum = cam90().frustum(); // camera at z=5, looking at origin
        // Origin is 5 units in front → depth = 5.
        assert!((frustum.depth(Vec3::ZERO) - 5.0).abs() < 1e-4);
    }

    #[test]
    fn depth_behind_camera_is_negative() {
        let frustum = cam90().frustum();
        let behind = Vec3::new(0.0, 0.0, 10.0); // same side as camera, past it
        assert!(frustum.depth(behind) < 0.0);
    }

    // ── Frustum::contains ─────────────────────────────────────────────────────

    #[test]
    fn contains_sphere_in_front() {
        let frustum = cam90().frustum();
        assert!(frustum.contains(Vec3::ZERO, 0.5));
    }

    #[test]
    fn contains_rejects_far_off_axis() {
        // 90° fov: at z=0 (5 units ahead) visible range is ±5. x=100 is way outside.
        let frustum = cam90().frustum();
        assert!(!frustum.contains(Vec3::new(100.0, 0.0, 0.0), 0.5));
    }

    #[test]
    fn contains_large_sphere_straddles_edge() {
        // A huge sphere that straddles the left frustum boundary should still be included.
        let frustum = cam90().frustum();
        // x=-50, but radius=60 — sphere overlaps the frustum.
        assert!(frustum.contains(Vec3::new(-50.0, 0.0, 0.0), 60.0));
    }

    // ── set_position / set_look_at ───────────────────────────────────────────

    #[test]
    fn set_position_updates_ray_direction() {
        // Camera at (0,0,5) → (5,0,0): forward flips from -Z to -X.
        let mut cam = Camera::new(Vec3::new(0.0, 0.0, 5.0), Vec3::ZERO, Vec3::Y, 90.0, 1, 1);
        cam.set_position(Vec3::new(5.0, 0.0, 0.0));
        let ray = cam.ray_for_pixel(0, 0); // center pixel → pure forward direction
        assert!((ray.dir - Vec3::new(-1.0, 0.0, 0.0)).length() < 1e-5);
    }

    #[test]
    fn set_look_at_updates_ray_direction() {
        // Camera at (0,0,5) originally looks toward -Z. Retarget to (0,0,10) → looks toward +Z.
        let mut cam = Camera::new(Vec3::new(0.0, 0.0, 5.0), Vec3::ZERO, Vec3::Y, 90.0, 1, 1);
        cam.set_look_at(Vec3::new(0.0, 0.0, 10.0));
        let ray = cam.ray_for_pixel(0, 0);
        assert!((ray.dir - Vec3::new(0.0, 0.0, 1.0)).length() < 1e-5);
    }

    #[test]
    fn set_position_updates_frustum() {
        // After moving, frustum origin must match new position.
        let mut cam = Camera::new(Vec3::new(0.0, 0.0, 5.0), Vec3::ZERO, Vec3::Y, 90.0, 100, 100);
        let new_pos = Vec3::new(3.0, 0.0, 0.0);
        cam.set_position(new_pos);
        let frustum = cam.frustum();
        assert!((frustum.origin - new_pos).length() < 1e-5);
    }

    // ── Frustum::cull_and_sort ────────────────────────────────────────────────

    #[test]
    fn cull_and_sort_empty() {
        let frustum = cam90().frustum();
        assert!(frustum.cull_and_sort(&[]).is_empty());
    }

    #[test]
    fn cull_and_sort_order_near_far() {
        let frustum = cam90().frustum(); // camera at z=5 looking at origin
        let near = box_sphere(Vec3::new(0.0, 0.0, 2.0), 0.3, Material::matte(Color::WHITE));
        let far  = box_sphere(Vec3::new(0.0, 0.0, -2.0), 0.3, Material::matte(Color::WHITE));
        let objs: Vec<Box<dyn Object>> = vec![far, near]; // intentionally wrong order
        let sorted = frustum.cull_and_sort(&objs);
        // Both visible; index 1 (near, z=2) should come before index 0 (far, z=-2).
        assert_eq!(sorted, vec![1, 0]);
    }

    #[test]
    fn cull_and_sort_removes_off_screen() {
        let frustum = cam90().frustum(); // fov=90°, aspect=1 — at z=5 looking at origin
        let visible    = box_sphere(Vec3::new(0.0, 0.0, 0.0), 0.5, Material::matte(Color::WHITE));
        let off_screen = box_sphere(Vec3::new(100.0, 0.0, 0.0), 0.5, Material::matte(Color::WHITE));
        let objs: Vec<Box<dyn Object>> = vec![visible, off_screen];
        let result = frustum.cull_and_sort(&objs);
        assert_eq!(result.len(), 1);
        assert_eq!(result[0], 0);
    }
}
