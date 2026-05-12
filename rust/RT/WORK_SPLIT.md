# Work Split

Pick one bundle from each section. Every person owns **one shape** plus one larger task. Shape work is small and self-contained — pair it with whatever bundle interests you.

## Shapes (pick one each)
- ✅ **Sphere** — `crates/objects/src/sphere.rs`
- ✅ **Cube** — `crates/objects/src/cube.rs`
- ✅ **Plane** — `crates/objects/src/plane.rs`
- ✅ **Cylinder** — `crates/objects/src/cylinder.rs`

For each shape:
- Implement `intersect` — return nearest forward `Hit`, stamp in own `material`, normal must be outward-facing and unit length
- Verify `bounding_sphere` is tight enough for culling
- Add unit tests in the same file

(No `bounce_ray` per shape — reflection is a single free function `reflect(incoming, normal)` in `objects`, called by the raycaster.)

## Bundles (pick one each)

### ✅ A — Shared core types  *(already largely done — see `objects/src/lib.rs`)*
- ✅ `Ray { origin, dir }` with `Ray::new` (normalizes dir) and `ray.at(t)`
- ✅ `Hit { distance, point, normal, material }` — outward unit normal
- ✅ `Material { color, reflective: bool }` + `matte` / `mirror` constructors
- ✅ `Color` with `+`, `*f32`, `*Color`, plus `BLACK`, `WHITE`, `rgb()`
- ✅ Free function `reflect(incoming, normal) -> Vec3`
- ✅ `Object` trait (no `bounce_ray` — kept off the trait on purpose)
- **Blocker for everyone — must land before B/C/shapes start**

### ✅ B — Camera & raycasting
- ✅ `Camera::new`, `ray_for_pixel`
- ✅ `Frustum`, `contains`, `depth`
- ✅ `Raycaster::cull_and_sort`, `trace`
- ✅ Reflection recursion: when `hit.material.reflective`, spawn next ray from `hit.point + EPS*normal` using `reflect(ray.dir, hit.normal)`, recurse up to `MAX_DEPTH` and combine colors along the path

### ✅ C — Lighting
- ✅ `LightSource` variants (point, directional)
- ✅ Shadow rays
- ✅ `shade()` — Lambert + falloff + occlusion

### ✅ D — Image output
- ✅ `Framebuffer::new`, `set`
- ✅ `write_ppm` (P3 format)
- ✅ Handle gamma / clamping

### ✅ E — World glue & binary
- ✅ `World::new`, `add_object`, `add_light`
- ✅ `World::render` orchestration
- ✅ `rt` CLI parsing
- ✅ Wire `--scene`, `--out`, `--width`, `--height`, `--list-scenes`

### F — Scenes (pure data)
- ✅ sphere scene
- ✅ plane scene
- ✅ cylinder scene
- ✅ `plane_cube` scene builder
- ✅ `all_objects` scene builder
- ✅ `all_objects_alt_cam` scene builder
- Generate the 4 final 800×600 `.ppm` images
- Add a couple extra fun scenes if time permits

### ✅ G — Liveview
- ✅ Pick a windowing crate (`minifb` or `pixels`)
- ✅ Implement `run(world, w, h)` — blit framebuffer per frame
- ✅ Keyboard input → mutate `world.camera`
- ✅ Wire `--live` flag in `rt`

### H — E2E fuzz tests
- Random world generator (seeded)
- Render at 100×100, assert no panic / no NaN
- Hook into `cargo test` at workspace level

### I — Documentation
- ✅ `DOCUMENTATION.md` — how to use the API, create objects, configure light brightness, move camera
- ✅ Code examples for each shape
- ✅ README quickstart

## Order of operations
1. **A first** — unblocks everything else
2. Shapes + B/C/D/E/F/G/H/I in parallel
3. F (final scenes) and I (docs) wrap up at the end once everything works
