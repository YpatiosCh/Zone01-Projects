# Ray Tracer — Documentation

## How it works

### Overview

Ray tracing renders a 3D scene by simulating light in reverse: instead of following photons from a light source, we shoot a ray from the camera through each pixel and ask "where does this ray go, and what does it see?"

```
Camera → ray → hits sphere → shadow ray to light → lit/shadowed? → pixel color
```

### 1. Ray generation

The camera precomputes an orthonormal basis (forward, right, up) from its position, look-at target, and up vector. For each pixel at screen coordinate `(x, y)`, a ray direction is computed by mapping the pixel to normalised device coordinates and interpolating within the view frustum:

```
dir = forward + right * sx + up * sy
```

The ray origin is always the camera position. All directions are unit-length.

### 2. Ray–object intersection

Each object implements an `intersect(ray)` function that returns the nearest hit point along the ray (distance > 0). The engine tests the ray against every object, collects all hits, and keeps the closest one. A tiny epsilon offset (`1e-4`) is applied when spawning shadow rays to prevent a surface from shadowing itself.

### 3. Shading

Once the closest hit is found, the surface color is computed in three parts:

1. **Ambient** — a constant base `(surface_color × 0.5)` so unlit faces are never fully black.
2. **Diffuse (Lambert)** — for each light source, the dot product of the surface normal and the direction to the light (`n·l`, clamped to ≥ 0) scales the light's contribution. The result is multiplied by both the surface color and the light color.
3. **Shadows** — before adding a light's diffuse contribution, a shadow ray is cast from the hit point toward the light. If another object blocks the path, the light contributes nothing.

### 4. Reflection

Mirror materials spawn a secondary ray reflected about the surface normal:

```
reflected = incoming − normal × 2(incoming · normal)
```

The engine follows up to **8 bounces**. At each reflective hit, a tint color accumulates by multiplying the mirror's color. When a non-reflective surface is reached (or the bounce limit is hit), shading is computed normally and the tint is applied.

### 5. Soft shadows via area lights

There is no per-pixel anti-aliasing (one ray per pixel). Soft shadows are approximated by `add_lights()`, which spreads one logical light across multiple point samples (on a disk for point lights, across a cone for directional lights). Total intensity is split equally across samples, so scene energy is preserved. More samples → wider, softer penumbra.

---

## Running the program

### Build

```bash
cargo build --release
```

### CLI usage

```
rt [--menu] [--live] [--render] [--scene <name>] [--out <file>] [--width N] [--height N] [--list-scenes]
```

| Flag | Description | Default |
|---|---|---|
| `--menu` | Open the interactive startup menu (default mode) | — |
| `--live` | Open the liveview window directly | — |
| `--render` | Render to file(s) without opening any window | — |
| `--scene <name>` | Scene to load | `sphere` |
| `--out <path>` | Output `.ppm` file path | `scenes_out/<scene>.ppm` |
| `--width N` | Image width in pixels | `800` |
| `--height N` | Image height in pixels | `600` |
| `--list-scenes` | Print all available scene names and exit | — |
| `-h`, `--help` | Show help text and exit | — |

### Common commands

```bash
# Open the interactive menu
cargo run --release

# Render the sphere scene to scenes_out/sphere.ppm (800×600)
cargo run -p rt --release -- --render --scene sphere

# Render all-objects scene at a lower resolution for a quick preview
cargo run -p rt --release -- --render --scene all_objects --width 400 --height 300

# Render to a custom output path
cargo run -p rt --release -- --render --scene all_objects --out my_render.ppm

# Open the liveview window for the cylinder scene
cargo run -p rt --release -- --live --scene cylinder

# List all available scenes
cargo run -p rt --release -- --list-scenes
```

### Available scenes

| Name | Description |
|---|---|
| `sphere` | Single red sphere |
| `cube` | Single cube |
| `cylinder` | Single cylinder |
| `plane` | Infinite flat plane |
| `plane_cube` | Flat plane + cube, reduced brightness |
| `plane_four_spheres` | Flat plane with four spheres |
| `plane_reflective` | Reflective floor plane |
| `all_objects` | One of every object (sphere, cube, cylinder, plane) |
| `all_objects_alt_cam` | Same as `all_objects` from a different camera position |

### Liveview controls

| Key | Action |
|---|---|
| `W / A / S / D` | Move forward / left / back / right |
| `Space` | Move up |
| `Ctrl` | Move down |
| Mouse | Rotate view |
| Arrow keys | Move / rotate lights |
| `F` | Save screenshot to `screenshots/` |
| `Esc` | Return to scene menu |

---

## Adding a new scene

Adding a scene touches **3 places**:

### 1. Create the scene file

Add `crates/scenes/src/<name>.rs` with a `pub fn build() -> World` function:

```rust
use glam::Vec3;
use world::{Camera, Color, LightSource, Material, Sphere, World};

pub fn build() -> World {
    let cam = Camera::new(Vec3::new(0.0, 2.0, 8.0), Vec3::ZERO, Vec3::Y, 60.0, 800, 600);
    let mut w = World::new(cam);
    w.add_object(Sphere::new(Vec3::ZERO, 1.0, Material::matte(Color::rgb(0.8, 0.3, 0.3))));
    w.add_light(LightSource::Point { pos: Vec3::new(3.0, 5.0, 3.0), intensity: 1.0, color: Color::WHITE });
    w
}
```

### 2. Register it in `crates/scenes/src/lib.rs`

Add a `pub mod` declaration and one entry to `REGISTRY`:

```rust
pub mod my_scene;   // add this line with the other pub mods

static REGISTRY: &[(&str, fn() -> World)] = &[
    // ... existing entries ...
    ("my_scene", my_scene::build),  // add this line
];
```

### 3. (Optional) Add it to the interactive menu — `crates/menu/src/lib.rs`

Add one line to the `SCENES` array if you want the scene selectable from the GUI:

```rust
static SCENES: &[(&str, &str, u32, u32)] = &[
    // ... existing entries ...
    ("my_scene", "My Scene", 0x6CB4D9, 0x0A1720),  // (cli_name, label, fill_color, text_color)
];
```

If you skip step 3, the scene is still reachable via CLI:

```bash
cargo run --release -- --render --scene my_scene
cargo run --release -- --live   --scene my_scene
```

---

## Quick start

```rust
use glam::Vec3;
use world::{Camera, Color, LightSource, Material, Sphere, World};

fn main() {
    let camera = Camera::new(
        Vec3::new(0.0, 2.0, 8.0), // position
        Vec3::new(0.0, 0.0, 0.0), // look-at target
        Vec3::Y,                   // up vector
        60.0,                      // vertical FOV (degrees)
        800, 600,                  // pixel dimensions
    );

    let mut world = World::new(camera);

    world.add_object(Sphere::new(
        Vec3::ZERO,
        1.0,
        Material::matte(Color::rgb(0.8, 0.3, 0.3)),
    ));

    world.add_light(LightSource::Point {
        pos: Vec3::new(4.0, 6.0, 4.0),
        intensity: 1.0,
        color: Color::WHITE,
    });

    let fb = world.render(800, 600);
    fb.save_ppm("out.ppm").unwrap();
    fb.save_png("out.png").unwrap();
}
```

---

## Camera

### Creating a camera

```rust
let camera = Camera::new(
    Vec3::new(0.0, 3.5, 9.0), // eye position
    Vec3::new(0.0, 0.8, 0.0), // point the camera looks at
    Vec3::Y,                   // world up — almost always Vec3::Y
    60.0,                      // vertical FOV in degrees (45–90 is typical)
    800, 600,                  // render resolution
);
```

- **position** — where the eye sits in world space.
- **look_at** — the point the camera is aimed at; the forward direction is `(look_at − position).normalize()`.
- **up** — used to orient the camera roll. `Vec3::Y` keeps the horizon level.
- **fov** — vertical field of view in degrees. Smaller values zoom in; larger values give a wide angle.

### Moving the camera at runtime

```rust
// Teleport the eye to a new position (look-at stays the same).
world.camera.set_position(Vec3::new(5.0, 2.0, 5.0));

// Retarget what the camera is pointing at.
world.camera.set_look_at(Vec3::new(0.0, 1.0, 0.0));

// Both can be combined to fully reposition the view.
world.camera.set_position(Vec3::new(-3.0, 4.0, 8.0));
world.camera.set_look_at(Vec3::ZERO);
```

Both methods recompute the view basis immediately, so the next `world.render()` reflects the change.

---

## Shapes

All shapes implement the `Object` trait. Add them to the world with `world.add_object(shape)`.

### Color and Material

Every shape takes a `Material`, which bundles its surface color and whether it acts as a mirror.

```rust
// Matte (diffuse) surface
let red   = Material::matte(Color::rgb(0.9, 0.2, 0.2));
let grey  = Material::matte(Color::rgb(0.5, 0.5, 0.5));
let white = Material::matte(Color::WHITE);

// Mirror surface — tinted reflections
let mirror      = Material::mirror(Color::WHITE);
let blue_mirror = Material::mirror(Color::rgb(0.6, 0.7, 1.0));
```

`Color::rgb(r, g, b)` takes values in `[0.0, 1.0]`.

---

### Sphere

```rust
use world::{Color, Material, Sphere};
use glam::Vec3;

// Sphere::new(center, radius, material)
world.add_object(Sphere::new(
    Vec3::new(-2.5, 0.6, 0.0), // center position
    0.6,                        // radius
    Material::matte(Color::rgb(0.95, 0.75, 0.2)),
));

// Mirror sphere
world.add_object(Sphere::new(
    Vec3::new(1.5, 1.0, 0.0),
    1.0,
    Material::mirror(Color::WHITE),
));

// Random sphere (useful for testing / fuzz scenes)
world.add_object(Sphere::random());
```

---

### Cube

```rust
use world::{Color, Cube, Material};
use glam::{Quat, Vec3};

// Cube::new(center, side_length, rotation, material)
world.add_object(Cube::new(
    Vec3::new(0.0, 0.5, 0.0),  // center position
    1.0,                         // side length
    Quat::IDENTITY,              // no rotation
    Material::matte(Color::rgb(0.85, 0.25, 0.25)),
));

// Rotated cube — 45° around the Y axis
use std::f32::consts::FRAC_PI_4;
world.add_object(Cube::new(
    Vec3::new(0.0, 0.5, 0.0),
    1.0,
    Quat::from_rotation_y(FRAC_PI_4),
    Material::matte(Color::rgb(0.3, 0.6, 0.9)),
));
```

---

### Plane

A `Plane` is infinite. Define it with any point on the surface and its outward normal.

```rust
use world::{Color, Material, Plane};
use glam::Vec3;

// Horizontal floor at y = 0, normal pointing up
world.add_object(Plane::new(
    Vec3::ZERO,
    Vec3::Y,
    Material::matte(Color::rgb(0.5, 0.5, 0.5)),
));

// Vertical back wall at z = -5, facing +Z
world.add_object(Plane::new(
    Vec3::new(0.0, 0.0, -5.0),
    Vec3::Z,
    Material::matte(Color::rgb(0.8, 0.8, 0.9)),
));

// Reflective floor
world.add_object(Plane::new(
    Vec3::ZERO,
    Vec3::Y,
    Material::mirror(Color::rgb(0.9, 0.9, 0.9)),
));
```

---

### Cylinder

A `Cylinder` has a base center, an axis direction, a radius, and a finite height. Both end caps are solid.

```rust
use world::{Color, Cylinder, Material};
use glam::Vec3;

// Cylinder::new(base_center, axis_direction, radius, height, material)
world.add_object(Cylinder::new(
    Vec3::new(2.5, 0.0, 0.0), // center of the bottom cap
    Vec3::Y,                   // axis direction (normalized internally)
    0.5,                       // radius
    1.5,                       // height along the axis
    Material::matte(Color::rgb(0.2, 0.7, 0.75)),
));

// Horizontal cylinder lying along the X axis
world.add_object(Cylinder::new(
    Vec3::new(-1.0, 1.0, 0.0),
    Vec3::X,
    0.4,
    2.0,
    Material::matte(Color::rgb(0.6, 0.3, 0.8)),
));
```

---

## Lights

### Point light

Emits light in all directions from a single position. Intensity falls off with distance.

```rust
world.add_light(LightSource::Point {
    pos:       Vec3::new(4.0, 7.0, 4.0),
    intensity: 1.0,
    color:     Color::WHITE,
});
```

### Directional light

Simulates a distant source (like the sun). `dir` is the direction the light *travels* — e.g. `Vec3::NEG_Y` shines straight down.

```rust
world.add_light(LightSource::Directional {
    dir:       Vec3::new(-1.0, -2.0, -1.0), // light travels in this direction
    intensity: 0.8,
    color:     Color::rgb(1.0, 0.95, 0.8),  // warm sunlight tint
});
```

### Configuring brightness

`intensity` is a linear multiplier on the light's contribution to the shading equation.

| `intensity` | Effect |
|---|---|
| `0.3` | Dim fill light |
| `0.8` | Moderate key light |
| `1.0` | Full white illumination |
| `> 1.0` | Oversaturates; output is clamped to 1.0 |

Coloured lights tint the lit surface by multiplying the light color with the surface color:

```rust
// Cool blue rim light
world.add_light(LightSource::Point {
    pos:       Vec3::new(-5.0, 3.0, -2.0),
    intensity: 0.4,
    color:     Color::rgb(0.5, 0.6, 1.0),
});
```

### Soft (area) light via `add_lights`

`World::add_lights` spreads one logical light across multiple samples to approximate a soft area source. Total intensity is divided equally across all samples so scene energy is preserved.

```rust
// 10-sample soft point light — wider penumbra
world.add_lights(
    LightSource::Point {
        pos:       Vec3::new(4.0, 7.0, 4.0),
        intensity: 1.0,
        color:     Color::WHITE,
    },
    10, // sample count: 1 = sharp, 8–16 = soft
);

// Soft directional light (spreads across a narrow elevation cone)
world.add_lights(
    LightSource::Directional {
        dir:       Vec3::new(-0.5, -1.0, -0.3),
        intensity: 0.9,
        color:     Color::WHITE,
    },
    8,
);
```

---

## Rendering

```rust
// Render to a Framebuffer and save as PPM and/or PNG.
let fb = world.render(800, 600); // width, height
fb.save_ppm("output.ppm").unwrap();
fb.save_png("output.png").unwrap();
```

The resolution passed to `render` may differ from the camera's construction resolution — the camera is reconstructed internally at the requested size.

---

## Liveview

`liveview::run` opens an interactive window. The camera flies with WASD and mouse look.

```rust
use liveview::{run, LiveviewAction};

match run(world, 1280, 720) {
    LiveviewAction::OpenSceneMenu => { /* return to menu */ }
    LiveviewAction::ExitApp       => { /* quit */ }
}
```

**Keyboard controls:**

| Key | Action |
|---|---|
| `W / A / S / D` | Move forward / left / back / right |
| `Space` | Move up |
| `Ctrl` | Move down |
| Mouse | Rotate view |
| Arrow keys | Move / rotate lights |
| `F` | Save screenshot to `screenshots/` |
| `Esc` | Return to scene menu |

---

## Complete example — all four shapes

```rust
use glam::{Quat, Vec3};
use world::{Camera, Color, Cube, Cylinder, LightSource, Material, Plane, Sphere, World};

fn build_scene() -> World {
    let camera = Camera::new(
        Vec3::new(0.0, 3.5, 9.0),
        Vec3::new(0.0, 0.8, 0.0),
        Vec3::Y,
        60.0,
        800, 600,
    );
    let mut world = World::new(camera);

    // Ground plane
    world.add_object(Plane::new(
        Vec3::ZERO,
        Vec3::Y,
        Material::matte(Color::rgb(0.5, 0.5, 0.5)),
    ));

    // Yellow sphere
    world.add_object(Sphere::new(
        Vec3::new(-2.5, 0.6, 0.0),
        0.6,
        Material::matte(Color::rgb(0.95, 0.75, 0.2)),
    ));

    // Red cube (no rotation)
    world.add_object(Cube::new(
        Vec3::new(0.0, 0.5, 0.0),
        1.0,
        Quat::IDENTITY,
        Material::matte(Color::rgb(0.85, 0.25, 0.25)),
    ));

    // Teal cylinder
    world.add_object(Cylinder::new(
        Vec3::new(2.5, 0.0, 0.0),
        Vec3::Y,
        0.5,
        1.5,
        Material::matte(Color::rgb(0.2, 0.7, 0.75)),
    ));

    // Soft area point light (10 samples)
    world.add_lights(
        LightSource::Point {
            pos:       Vec3::new(4.0, 7.0, 4.0),
            intensity: 1.0,
            color:     Color::WHITE,
        },
        10,
    );

    world
}

fn main() {
    let world = build_scene();
    let fb = world.render(800, 600);
    fb.save_png("scene.png").unwrap();
}
```
