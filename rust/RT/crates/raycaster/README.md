# raycaster

Ray-object intersection and mirror reflection traversal. This crate sits between the geometry layer (`objects`) and the shading layer (`lighting`). Its three public functions cover the three distinct query types the renderer needs.

---

## Functions

### `Raycaster::occluded`

```rust
pub fn occluded(ray: Ray, max_distance: f32, objects: &[Box<dyn Object>]) -> bool
```

Returns `true` if any object blocks `ray` within `max_distance`. Early-exits on the first hit found — it does not find the nearest hit, so it is faster than `nearest_hit_within` when you only care about presence.

Used by `lighting::shade` to test shadow rays:

```rust
let shadow_origin = hit.point + hit.normal * SURFACE_EPSILON;
let to_light      = light_pos - shadow_origin;
let shadow_ray    = Ray::new(shadow_origin, to_light);

if Raycaster::occluded(shadow_ray, to_light.length(), objects) {
    continue; // in shadow — skip this light's contribution
}
```

The `max_distance` cap is important: without it, objects on the far side of the light would incorrectly cast shadows.

---

### `Raycaster::nearest_hit_within`

```rust
pub fn nearest_hit_within(ray: Ray, max_distance: f32, objects: &[Box<dyn Object>]) -> Option<Hit>
```

Finds and returns the nearest `Hit` along `ray` with `distance < max_distance`. Returns `None` if no object is hit within that range. Used internally by `occluded` and by `contact_shadow_factor` in the lighting crate.

---

### `Raycaster::trace_with_reflections`

```rust
pub fn trace_with_reflections(
    ray: Ray,
    objects: &[Box<dyn Object>],
    max_depth: u32,
) -> (Color, Option<Hit>)
```

The main render-path function. Traces `ray` through the scene, bouncing off mirror surfaces up to `max_depth` times.

**Return value** — `(tint, terminal_hit)`:

- `tint` — the accumulated product of every mirror surface's `material.color` along the bounce chain. Starts as `Color::WHITE`. After one red-mirror bounce it becomes `Color::rgb(1,0,0)`. After two bounces through a 50%-gray mirror it becomes `Color::rgb(0.25, 0.25, 0.25)`. `World::render` multiplies the final shaded color by `tint`.
- `terminal_hit` — the first non-reflective `Hit` encountered, or `None` if the ray escaped to sky or the depth limit was reached.

**Bounce rules:**

1. Trace the current ray against all objects.
2. No hit → return `(tint, None)` — the ray escaped to sky.
3. Non-reflective hit → return `(tint, Some(hit))` — shade this surface.
4. Reflective hit → multiply `tint` by `hit.material.color`, push the origin off the surface by `SURFACE_EPSILON`, spawn a new ray in the reflected direction, and loop.
5. `max_depth` exhausted → return `(tint, None)` — treat as sky (e.g. two mirrors facing each other would loop forever without this cap).

**Usage in `World::render`:**

```rust
let (tint, terminal_hit) = Raycaster::trace_with_reflections(primary_ray, &objects, 8);
let base_color = match terminal_hit {
    Some(hit) => lighting::shade(&hit, &lights, &objects),
    None      => sky.color,
};
let pixel_color = base_color * tint;
```

---

## Internal function

### `Raycaster::trace` (private)

Iterates all objects, calls `intersect` on each, filters out negative-distance hits (behind the ray origin), and returns the nearest one. Used by both `nearest_hit_within` and `trace_with_reflections`.

```
Complexity: O(n) per call, n = number of objects.
```

`World::render` pre-filters objects with `Frustum::cull_and_sort` before passing them here, so the slice is already culled to the visible set.

---

## Reflection math

The reflected ray direction is computed by `objects::reflect`:

```
reflected = incoming - normal * 2 * (incoming · normal)
```

This is the standard specular reflection formula. The reflected ray origin is offset from the hit point along the surface normal by `SURFACE_EPSILON` to prevent floating-point self-intersection.

---

## Depth limit

The default bounce depth used by `World::render` is **8**. Scenes with many mirrors may appear slightly darker than reality at that depth, but infinite loops are not possible. Increasing the depth produces diminishing returns — each bounce multiplies tint by the mirror's color, so the contribution of very deep bounces is nearly zero.
