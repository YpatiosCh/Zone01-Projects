# camera

Perspective camera and view frustum. The camera converts pixel coordinates into primary rays, and the frustum culls objects to those visible on screen before rendering begins.

---

## `Camera`

```rust
use camera::Camera;
use glam::Vec3;

let cam = Camera::new(
    Vec3::new(0.0, 3.5, 9.0), // eye position
    Vec3::new(0.0, 0.8, 0.0), // look-at target
    Vec3::Y,                    // world up vector
    60.0,                       // vertical FOV in degrees
    800, 600,                   // pixel dimensions
);
```

### Constructor

`Camera::new` builds an orthonormal view basis from the three directional inputs:

```
forward  = normalize(look_at - position)
right    = normalize(forward × up)
up_basis = normalize(right × forward)
```

`right` and `up_basis` span the image plane. `half_height = tan(fov/2)` and `aspect = width/height` are precomputed once so `ray_for_pixel` has no trigonometry at call time.

### `ray_for_pixel`

```rust
let ray = cam.ray_for_pixel(x, y);
```

Maps pixel `(x, y)` to a primary ray through the center of that pixel. The mapping is:

```
u  = (x + 0.5) / width          // [0, 1] normalized, pixel-centered
v  = (y + 0.5) / height
sx = (2u - 1) * aspect * half_height  // [-aspect*tan(fov/2), +...]
sy = (1 - 2v) * half_height           // [+tan(fov/2), -...] (y flipped)
dir = forward + right*sx + up_basis*sy
```

The `+0.5` centers the sample within the pixel. The `1 - 2v` flip makes row 0 correspond to the top of the image (standard raster order). The resulting direction is passed to `Ray::new`, which normalizes it.

### `set_position` / `set_look_at`

Dynamic updates for the liveview camera. Both recompute the view basis immediately so the next call to `ray_for_pixel` or `frustum()` reflects the new orientation.

```rust
cam.set_position(Vec3::new(1.0, 2.0, 5.0));
cam.set_look_at(cam.position + forward_from_yaw_pitch);
```

---

## `Frustum`

A four-plane half-space representation of the camera's view pyramid. Used to cull objects that are entirely outside the screen.

```rust
let frustum = cam.frustum();
let visible_indices = frustum.cull_and_sort(&world.objects);
let visible: Vec<_> = visible_indices.iter().map(|&i| &world.objects[i]).collect();
```

### Construction

`Camera::frustum()` computes four corner rays (top-left, top-right, bottom-left, bottom-right) and derives four plane equations from their pairwise cross products:

```
left plane   = normalize(bl × tl)
right plane  = normalize(tr × br)
top plane    = normalize(tl × tr)
bottom plane = normalize(br × bl)
```

Each `FrustumPlane` stores a normal and a plane constant `d = -normal · camera_position`, so the half-space test is `normal · point + d >= 0`.

### `Frustum::contains`

```rust
frustum.contains(center, radius) -> bool
```

Tests whether a sphere overlaps the frustum. A sphere is outside the frustum if it lies fully on the wrong side of any plane: `normal · center + d < -radius`. Testing all four planes (logical AND) gives the conservative inclusion test. Large spheres that straddle a plane edge are always included.

### `Frustum::depth`

```rust
frustum.depth(point) -> f32
```

Signed distance from the camera along the forward axis. Used as the sort key for near-to-far ordering, which improves early-exit behavior when `trace` scans the object list.

### `Frustum::cull_and_sort`

```rust
frustum.cull_and_sort(&objects) -> Vec<usize>
```

1. Filters the object slice to those whose bounding sphere passes `contains`.
2. Sorts the surviving indices by `depth(object.position())` in ascending order (near first).
3. Returns the sorted index list.

`World::render` uses the indices to build a `visible` slice that is passed to the raycaster and lighting crate. Objects entirely off screen are never tested for intersection.

---

## FOV and aspect ratio

The camera stores the **vertical** FOV. The horizontal FOV is derived from the aspect ratio:

```
horizontal_half_angle = atan(aspect * tan(vertical_fov / 2))
```

At `fov = 60°` and 800×600 (`aspect ≈ 1.33`), the horizontal FOV is about 74°. Wider FOV produces more pronounced perspective distortion (fisheye effect); narrower FOV flattens the scene (telephoto look).

---

## `tile_frustum`

```rust
cam.tile_frustum(x0, y0, x1, y1) -> Frustum
```

Computes a sub-frustum covering only the pixel rect `[x0,x1) × [y0,y1)`. Currently not implemented (`todo!()`), reserved for tile-based parallel rendering.
