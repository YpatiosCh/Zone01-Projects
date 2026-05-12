# lighting

Shading and light sources for the ray tracer. The public surface is two types (`LightSource`, `SkyProperties`) and one function (`shade`).

---

## `LightSource`

Two variants are supported. Both carry an **intensity** (a plain brightness multiplier — `1.0` is full brightness, `0.5` is half, `0.0` is off) and a **color** tint.

### Point light

Emits light from a single position in world space in every direction.

```rust
use world::{LightSource, Color};
use glam::Vec3;

LightSource::Point {
    pos:       Vec3::new(3.0, 8.0, 2.0), // world-space position
    intensity: 1.0,                       // brightness multiplier
    color:     Color::WHITE,              // tint (use Color::rgb for colored light)
}
```

**Brightness** scales linearly: `intensity: 2.0` is twice as bright as `intensity: 1.0`. There is no distance falloff — a point light at distance 2 and at distance 20 contribute the same diffuse energy given the same angle. Adjust `intensity` or move the light to control how bright areas are.

### Directional light

Simulates a distant light source (sun, sky). Every point in the scene receives parallel rays from the same direction. The `dir` field is the direction the light *travels* (toward the scene); the shader negates it internally to get the surface-to-light vector.

```rust
// Sun from above-right
LightSource::Directional {
    dir:       Vec3::new(-0.5, -1.0, -0.3).normalize(), // light travels downward toward scene
    intensity: 0.8,
    color:     Color::rgb(1.0, 0.95, 0.85),             // warm sunlight
}
```

---

## Changing brightness

Brightness is controlled by the `intensity` field on each light. A useful pattern is to keep one main light at `intensity: 1.0` and add fill lights at lower intensities:

```rust
// High-brightness scene
world.add_light(LightSource::Point {
    pos: Vec3::new(2.0, 6.0, 4.0),
    intensity: 1.0,
    color: Color::WHITE,
});

// Lower-brightness scene — same light, halved intensity
world.add_light(LightSource::Point {
    pos: Vec3::new(2.0, 6.0, 4.0),
    intensity: 0.4,
    color: Color::WHITE,
});
```

Multiple lights stack: two lights each at `0.5` produce roughly the same brightness as one at `1.0` (subject to surface angle). Output is clamped to `[0, 1]` per channel, so pushing intensity above `1.0` only brightens the shadowed ambient term's proportion.

---

## `SkyProperties`

The background color rendered when a ray hits nothing. Set it on the `World`:

```rust
use world::SkyProperties;

world.sky = SkyProperties { color: Color::rgb(0.5, 0.7, 1.0) }; // default blue sky
world.sky = SkyProperties { color: Color::BLACK };                // night / studio black
```

---

## `shade`

```rust
pub fn shade(hit: &Hit, lights: &[LightSource], objects: &[Box<dyn Object>]) -> Color
```

Called automatically by `World::render` for every surface hit. You do not need to call it directly. What it does per frame per pixel:

1. Starts with a small ambient term (`material.color × 0.05`) so surfaces are never fully black.
2. For each `LightSource`, computes the Lambert diffuse term (`NdotL = max(0, normal · to_light)`).
3. Casts a shadow ray via `Raycaster::occluded`; skips the light's contribution if anything blocks it.
4. Accumulates all contributions and clamps the result to `[0, 1]`.

---

## Shadows

Shadow testing is automatic. Any object placed between a surface and a light will cast a shadow on that surface. No extra configuration is needed.

---

## Full example

```rust
use glam::Vec3;
use world::{World, Camera, Sphere, Material, Color, LightSource, SkyProperties};

let cam = Camera::new(
    Vec3::new(0.0, 2.0, 6.0), // position
    Vec3::ZERO,                // look-at target
    Vec3::Y,                   // up vector
    60.0,                      // vertical FOV in degrees
    800, 600,
);
let mut world = World::new(cam);

// Objects
world.add_object(Sphere {
    center:   Vec3::ZERO,
    radius:   1.0,
    material: Material::matte(Color::rgb(0.8, 0.3, 0.3)),
});

// Key light — warm white from upper right
world.add_light(LightSource::Point {
    pos:       Vec3::new(4.0, 8.0, 4.0),
    intensity: 1.0,
    color:     Color::WHITE,
});

// Fill light — cool blue from the left, dimmer
world.add_light(LightSource::Directional {
    dir:       Vec3::new(1.0, -0.3, 0.0).normalize(),
    intensity: 0.2,
    color:     Color::rgb(0.6, 0.7, 1.0),
});

// Sky
world.sky = SkyProperties { color: Color::rgb(0.4, 0.6, 0.9) };

let framebuffer = world.render(800, 600);
```
