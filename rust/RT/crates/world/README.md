# world

Scene container and top-level renderer. `World` holds the camera, object list, light list, and sky color, and exposes a single `render` method that produces a complete `Framebuffer`.

This crate also re-exports everything from `camera`, `image`, `lighting`, `objects`, and `raycaster` so scene builders only need to import from `world`.

---

## `World`

```rust
use world::{World, Camera, Sphere, Material, Color, LightSource};
use glam::Vec3;

let cam = Camera::new(
    Vec3::new(0.0, 3.0, 8.0),
    Vec3::ZERO,
    Vec3::Y,
    60.0,
    800, 600,
);
let mut world = World::new(cam);
```

### Adding objects

```rust
world.add_object(Sphere::new(
    Vec3::ZERO,
    1.0,
    Material::matte(Color::rgb(0.8, 0.3, 0.3)),
));
```

`add_object` accepts any type that implements `Object + 'static`. It boxes the value for dynamic dispatch so the object list can hold mixed shape types.

### Adding a single light

```rust
world.add_light(LightSource::Point {
    pos: Vec3::new(4.0, 8.0, 4.0),
    intensity: 1.0,
    color: Color::WHITE,
});
```

### Soft shadows with `add_lights`

```rust
world.add_lights(
    LightSource::Point {
        pos: Vec3::new(4.0, 7.0, 4.0),
        intensity: 1.0,
        color: Color::WHITE,
    },
    10, // number of shadow samples
);
```

`add_lights` approximates a soft area light by spreading one logical light into `n` slightly varied copies, each with `intensity / n`. The total scene energy is preserved.

- **Point lights** — sample positions are placed evenly on a horizontal disk of radius `AREA_LIGHT_RADIUS = 1.0` centered on the original position. Soft penumbra width scales with the disk radius.
- **Directional lights** — sample directions are spread linearly across a narrow elevation cone of half-angle `DIRECTIONAL_SPREAD = 0.2 rad ≈ 11°`. The spread creates a soft terminator on shaded surfaces.

For `n = 1` the light is added unchanged (no disk or cone sampling).

The number of samples trades rendering speed for shadow softness. 10 samples is a reasonable default; 1 gives hard shadows; 20+ gives smooth penumbras at higher cost.

---

## `World::render`

```rust
let framebuffer = world.render(800, 600);
```

Renders the scene at `width × height` pixels. A fresh `Camera` is constructed from the world camera's fields so the aspect ratio matches the output dimensions.

**Per-pixel pipeline:**

1. `camera.ray_for_pixel(x, y)` — generate the primary ray.
2. `Raycaster::trace_with_reflections(ray, &objects, 8)` — trace with up to 8 mirror bounces; returns `(tint, terminal_hit)`.
3. If `terminal_hit` is `Some(hit)`, call `lighting::shade(&hit, &lights, &objects)` to get the diffuse color.  
   If `terminal_hit` is `None` (ray escaped to sky or bounce limit reached), use `sky.color`.
4. Multiply the color by `tint` to apply mirror color accumulation.
5. Write the result to the framebuffer with `fb.set(x, y, color)`.

The render is single-threaded; each row is computed in sequence.

---

## Sky

```rust
use world::SkyProperties;

world.sky = SkyProperties { color: Color::rgb(0.5, 0.7, 1.0) }; // default blue sky
world.sky = SkyProperties { color: Color::BLACK };                // studio black
```

The sky color is used whenever a ray exits the scene without hitting any object, including rays that reach the bounce depth limit. The default is a light blue (`0.5, 0.7, 1.0`).

---

## Re-exports

`world` re-exports the full public API of every foundational crate:

```rust
pub use camera::*;   // Camera, Frustum, FrustumPlane
pub use image::*;    // Framebuffer, png_path_for
pub use lighting::*; // LightSource, SkyProperties, shade
pub use objects::*;  // Color, Ray, Hit, Material, Object, Sphere, Cube, Plane, Cylinder, reflect, SURFACE_EPSILON
pub use raycaster::*;// Raycaster
```

Scene builders import only `world` and get access to everything they need.

---

## Constants

| Constant | Value | Purpose |
|---|---|---|
| `AREA_LIGHT_RADIUS` | `1.0` | Disk radius for point light soft-shadow sampling |
| `DIRECTIONAL_SPREAD` | `0.2` | Half-arc of elevation cone for directional light sampling |
