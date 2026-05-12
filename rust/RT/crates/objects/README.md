# objects

Core geometry primitives for the ray tracer. Every other crate depends on this one. It defines the shared vocabulary — `Color`, `Ray`, `Hit`, `Material`, and the `Object` trait — along with four concrete shapes.

---

## Types

### `Color`

An RGB triplet where each channel is a `f32` in `[0.0, 1.0]`.

```rust
let red   = Color::rgb(1.0, 0.0, 0.0);
let white = Color::WHITE;
let black = Color::BLACK;
```

Colors support component-wise addition (`Color + Color`) and scalar multiplication (`Color * f32`) as well as color multiplication (`Color * Color`), which is used to apply a light tint to a surface color.

---

### `Ray`

A half-line in 3D space: an origin point and a direction.

```rust
let ray = Ray::new(
    Vec3::new(0.0, 2.0, 5.0), // origin
    Vec3::new(0.0, 0.0, -1.0), // direction (normalized automatically)
);
let point_at_t3 = ray.at(3.0); // origin + dir * 3.0
```

`Ray::new` always normalizes the direction, so downstream code can rely on `ray.dir` being unit length. `ray.at(t)` evaluates the parametric equation `origin + dir * t`.

---

### `Hit`

The result returned by a shape's `intersect` method when a ray strikes a surface.

```rust
pub struct Hit {
    pub distance: f32,  // t value — how far along the ray the hit occurred
    pub point:    Vec3, // world-space position of the hit
    pub normal:   Vec3, // outward surface normal (unit length)
    pub material: Material,
}
```

`normal` is always the **outward** surface normal regardless of which side the ray came from. The raycaster uses it for shading and for orienting reflected rays.

---

### `Material`

Surface appearance properties: base color and whether the surface reflects like a mirror.

```rust
// Diffuse matte surface
let clay   = Material::matte(Color::rgb(0.8, 0.5, 0.3));

// Perfect mirror (retains its own color as a tint on reflections)
let mirror = Material::mirror(Color::WHITE);

// Colored mirror — reflections are tinted red
let red_mirror = Material::mirror(Color::rgb(1.0, 0.0, 0.0));
```

`reflective: false` means `shade()` is called directly. `reflective: true` means the raycaster spawns a bounce ray and accumulates a tint before shading the terminal surface.

---

### `Object` trait

The interface every shape must implement.

```rust
pub trait Object {
    fn intersect(&self, ray: Ray) -> Option<Hit>;
    fn position(&self) -> Vec3;
    fn set_position(&mut self, p: Vec3);
    fn bounding_sphere(&self) -> (Vec3, f32);
    fn floor_height_at(&self, position: Vec3) -> Option<f32> { None }
}
```

- `intersect` — ray vs. shape test. Returns the nearest forward hit (`distance > 0`), or `None` for a miss or a hit behind the origin.
- `position` / `set_position` — world-space anchor point used by the camera frustum for depth sorting.
- `bounding_sphere` — returns `(center, radius)` for frustum culling. This is a fast conservative test; if the sphere is outside the frustum the shape is skipped entirely.
- `floor_height_at` — optional; implemented only by `Plane`. The liveview system calls this to prevent the camera from clipping through floors.

---

### `SURFACE_EPSILON`

```rust
pub const SURFACE_EPSILON: f32 = 1e-4;
```

A small offset added when spawning shadow rays and reflected rays so they do not immediately re-intersect the surface they were spawned from (self-intersection artifact). Every shape's `intersect` returns `distance > 0` for forward hits only, but floating-point rounding can put a spawned ray origin exactly on the surface; `SURFACE_EPSILON` pushes it safely off.

---

### `reflect`

```rust
pub fn reflect(incoming: Vec3, normal: Vec3) -> Vec3 {
    incoming - normal * (2.0 * incoming.dot(normal))
}
```

Standard mirror-reflection formula. The raycaster calls this when `hit.material.reflective` is true.

---

## Shapes

### `Sphere`

```rust
use objects::{Sphere, Material, Color};
use glam::Vec3;

let sphere = Sphere::new(
    Vec3::new(0.0, 1.0, 0.0), // center
    1.0,                        // radius
    Material::matte(Color::rgb(0.9, 0.3, 0.3)),
);
```

**Intersection math** — the ray-sphere test solves the quadratic `|ray.origin + t*ray.dir - center|² = radius²`. Expanding and substituting `oc = ray.origin - center`:

```
t² + 2·(oc·dir)·t + (oc·oc - r²) = 0
discriminant = b² - c
```

where `b = oc·dir` and `c = oc·oc - r²`. A negative discriminant means no intersection. Two roots (`-b ± √disc`) are the entry and exit points; the code picks the smallest positive root (nearest forward hit). If the origin is inside the sphere the entry root is negative, so the exit root is used instead.

The bounding sphere is the sphere itself.

**Random generation** — `Sphere::random()` and `Sphere::random_with(rng)` generate random spheres for fuzz testing. Center is uniform in `[-5, 5]³`, radius in `[0.2, 1.5]`, 20% chance of being a mirror. The `_with` variant accepts a seeded RNG so failing fuzz seeds can be reproduced exactly.

---

### `Cube`

```rust
use objects::{Cube, Material, Color};
use glam::{Quat, Vec3};

let cube = Cube::new(
    Vec3::new(0.0, 0.5, 0.0), // center
    1.0,                        // side length
    Quat::from_rotation_y(0.4), // rotation
    Material::matte(Color::rgb(0.85, 0.25, 0.25)),
);
```

**Intersection math** — uses the **slab method**. The cube is axis-aligned in its own local frame `[-h, h]³` (where `h = size * 0.5`). The ray is transformed into local space by applying the inverse rotation: `o_local = inv_rot * (ray.origin - center)` and `d_local = inv_rot * ray.dir`. The rotation is rigid (no scale), so the parameter `t` is identical in both frames.

For each of the three axis pairs, the ray is tested against the two parallel planes `x = ±h`, `y = ±h`, `z = ±h`. The slab algorithm finds the range of `t` values where the ray is inside all three slabs simultaneously by tracking the running maximum of near-plane hits (`t_min`) and minimum of far-plane hits (`t_max`). When `t_min > t_max` the slabs do not overlap and the ray misses.

The face normal is reconstructed from `axis` and `sign` (which slab was entered, which side was hit) and rotated back to world space with `self.rotation`.

The bounding sphere is the circumscribed sphere with radius `size/2 * √3` (half the space diagonal), which is rotation-invariant.

---

### `Plane`

```rust
use objects::{Plane, Material, Color};
use glam::Vec3;

let floor = Plane::new(
    Vec3::ZERO,  // a point on the plane
    Vec3::Y,     // surface normal (normalized automatically)
    Material::matte(Color::rgb(0.5, 0.5, 0.5)),
);
```

**Intersection math** — a ray hits a plane when `(point - ray.origin) · normal = 0` along the parametric ray equation. The distance is:

```
distance = ((plane_point - ray.origin) · normal) / (ray.dir · normal)
```

If `ray.dir · normal` is near zero the ray is parallel and misses. If `distance ≤ 0` the hit is behind the origin. The normal is flipped to face the incoming ray (`denom < 0` means the ray is hitting the front face, so the plane's own normal is used; `denom > 0` means the ray comes from behind, so the normal is negated).

The bounding sphere has infinite radius, ensuring the plane is never culled.

`floor_height_at(position)` returns the world-space Y coordinate of the plane surface directly below an XZ position. The liveview system calls this to keep the camera above the floor:

```
y = plane_point.y - (normal.x * dx + normal.z * dz) / normal.y
```

Only planes with a meaningful upward component (`normal.y > PLANE_EPSILON`) implement this.

---

### `Cylinder`

```rust
use objects::{Cylinder, Material, Color};
use glam::Vec3;

let pillar = Cylinder::new(
    Vec3::new(2.5, 0.0, 0.0), // base center
    Vec3::Y,                    // axis direction (normalized automatically)
    0.5,                        // radius
    1.5,                        // height
    Material::matte(Color::rgb(0.2, 0.7, 0.75)),
);
```

**Intersection math** — the cylinder is a closed surface (side wall + two end caps). The test decomposes the ray into components parallel and perpendicular to the axis:

```
dir_perp   = ray.dir   - axis * (ray.dir   · axis)
offset_perp = origin_off - axis * (origin_off · axis)
```

where `origin_off = ray.origin - base`. The perpendicular components satisfy the 2D circle equation `|offset_perp + t * dir_perp|² = r²`, giving the quadratic `a·t² + b·t + c = 0`. Valid side hits must land within `[0, height]` along the axis.

Cap intersections are tested separately: for each cap (bottom at `axis_offset=0`, top at `axis_offset=height`) the ray is intersected with the cap plane and the radial distance is checked against the radius.

All candidates (side and caps) are collected and the nearest forward hit is returned.

A bounding sphere centered at the midpoint of the axis with radius `√((height/2)² + radius²)` is used for early rejection.
