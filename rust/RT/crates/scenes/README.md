# scenes

Pre-built scene definitions. Each scene is a `fn() -> World` builder that constructs a complete `World` — camera, objects, lights, and sky — ready to render or explore in liveview.

---

## Listing available scenes

```rust
let names = scenes::names(); // ["sphere", "cylinder", "cube", ...]
```

---

## Loading a scene by name

```rust
let builder: fn() -> World = scenes::get("plane_four_spheres")
    .expect("unknown scene name");
let world = builder();
let fb = world.render(800, 600);
```

`get` returns `None` for unknown names. The CLI (`rt`) uses `scenes::names()` to generate its `--list-scenes` output and its error messages.

---

## Available scenes

| Name | Description |
|---|---|
| `sphere` | Single matte sphere at the origin, default camera |
| `cylinder` | Single cylinder, default camera |
| `cube` | Single cube, default camera |
| `plane` | Infinite ground plane |
| `plane_four_spheres` | Ground plane with four colored matte spheres arranged on it |
| `plane_reflective` | Ground plane plus a mirror sphere; shows reflection |
| `plane_cube` | Ground plane with a cube sitting on it |
| `all_objects` | All four shapes together (sphere, cube, cylinder, plane) |
| `all_objects_alt_cam` | Same as `all_objects` from a different camera angle |

---

## Scene anatomy

Every scene builder follows the same pattern:

```rust
pub fn build() -> World {
    // 1. Create a camera
    let cam = Camera::new(position, look_at, up, fov_degrees, width, height);
    let mut w = World::new(cam);

    // 2. Add objects
    w.add_object(Plane::new(...));
    w.add_object(Sphere::new(...));

    // 3. Add lights
    w.add_lights(LightSource::Point { ... }, 10);

    w
}
```

The `width` and `height` passed to `Camera::new` in a scene builder are only defaults. `World::render(w, h)` creates its own `Camera` from the stored position/look_at/fov with the requested pixel dimensions, so scenes render correctly at any resolution.

---

## `all_objects` — example

```
Camera: position (0, 3.5, 9), looking at (0, 0.8, 0), 60° FOV
Objects:
  - Plane    at origin,         Y normal,  gray  (0.5, 0.5, 0.5)
  - Sphere   at (-2.5, 0.6, 0), r=0.6,     gold  (0.95, 0.75, 0.2)
  - Cube     at (0, 0.5, 0),    side=1,    red   (0.85, 0.25, 0.25)
  - Cylinder at (2.5, 0, 0),    r=0.5 h=1.5, teal (0.2, 0.7, 0.75)
Lights: 10-sample soft point light at (4, 7, 4)
```

---

## Adding a new scene

1. Create `crates/scenes/src/my_scene.rs` with a `pub fn build() -> World` function.
2. Declare `pub mod my_scene;` at the top of `crates/scenes/src/lib.rs`.
3. Add `("my_scene", my_scene::build)` to the `REGISTRY` array in `lib.rs`.
4. Add the scene to the `SCENES` table in `crates/menu/src/lib.rs` to make it appear in the GUI.

The registry is a static slice of `(&str, fn() -> World)` tuples. `get` does a linear scan, which is fine for the small number of scenes in this project.
