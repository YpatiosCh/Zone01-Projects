# RT — Project Plan

## Workspace Layout

Cargo workspace at repo root. Each library crate is self-contained and owned by a different person; `world` glues them together; `rt` is the binary entrypoint.

```
RT/
├── Cargo.toml              # [workspace] members = [...]
├── README.md
├── DOCUMENTATION.md        # user-facing API docs (required deliverable)
├── scenes/                 # 4 required output .ppm files
│   ├── sphere.ppm
│   ├── plane_cube.ppm
│   ├── all_objects.ppm
│   └── all_objects_alt_cam.ppm
├── crates/
│   ├── rt/                 # binary: parses CLI flags, builds World, writes .ppm
│   │   └── src/main.rs
│   ├── world/              # public API; aggregates camera + objects + lights + sky
│   │   └── src/lib.rs
│   ├── objects/            # Sphere, Cube, Plane, Cylinder + Object trait
│   │   └── src/{lib.rs, sphere.rs, cube.rs, plane.rs, cylinder.rs}
│   ├── raycaster/          # Raycaster: trace_with_reflections + occluded (shadow rays)
│   │   └── src/lib.rs
│   ├── camera/             # Camera + Frustum (culling, near→far sort)
│   │   └── src/lib.rs
│   ├── lighting/           # LightSource types, shadow + brightness math
│   │   └── src/lib.rs
│   ├── image/              # framebuffer (2D grid) + .ppm (P3) writer
│   │   └── src/lib.rs
│   ├── scenes/             # scenes-as-code: each file builds & returns a World
│   │   └── src/{lib.rs, sphere.rs, plane_cube.rs, all_objects.rs, all_objects_alt_cam.rs}
│   └── liveview/           # optional real-time window for debugging
│       └── src/lib.rs
└── tests/                  # workspace-level e2e: random scenes @ tiny resolutions, assert no crash
```

## Crate Responsibilities & Public Surface

### `objects` (trait + 4 impls + shared core types)

Shared core types live here (lowest layer of the dep graph):

```rust
pub struct Color { r: f32, g: f32, b: f32 }   // linear; gamma applied in write_ppm
// + ops: Add, Mul<f32>, Mul<Color>; constants BLACK/WHITE; Color::rgb(r,g,b)

pub struct Ray { origin: Vec3, dir: Vec3 }    // dir is normalized by Ray::new
impl Ray { pub fn new(origin, dir) -> Self;  pub fn at(t: f32) -> Vec3; }

pub struct Hit {
    distance: f32,        // t along the ray (renamed from `t`)
    point: Vec3,          // world-space hit point
    normal: Vec3,         // outward-facing, unit length
    material: Material,   // copied in by the shape
}

pub struct Material { color: Color, reflective: bool }   // no IOR, no roughness
impl Material { pub fn matte(c) -> Self; pub fn mirror(c) -> Self; }

/// Mirror reflection helper (free function — same math for every shape).
pub fn reflect(incoming: Vec3, normal: Vec3) -> Vec3;
```

```rust
pub trait Object {
    /// Ray-vs-shape intersection. Nearest forward hit (`distance > 0`), else None.
    /// The shape stamps in its own `material` on the returned `Hit`.
    fn intersect(&self, ray: Ray) -> Option<Hit>;
    /// World-space anchor point.
    fn position(&self) -> Vec3;
    /// Move the object before rendering.
    fn set_position(&mut self, p: Vec3);
    /// World-space bounding sphere `(center, radius)` for frustum culling.
    fn bounding_sphere(&self) -> (Vec3, f32);
}
pub struct Sphere   { center, radius, material }
pub struct Cube     { center, size, rotation, material }
pub struct Plane    { point, normal, material }
pub struct Cylinder { base, axis, radius, height, material }
```
Reflection is **not** on the trait — it's a free `reflect(incoming, normal)` since the math is identical for every shape. The raycaster decides when to call it (when `hit.material.reflective`) and is responsible for spawning the next ray and bounded recursion.

### `raycaster`
```rust
pub struct Raycaster;
impl Raycaster {
    /// Nearest forward hit, with full mirror-bounce chain up to `max_depth`.
    /// Returns `(tint, terminal_hit)`: tint is the product of mirror colors along the
    /// chain; terminal_hit is the first non-reflective hit (or None for sky/depth).
    pub fn trace_with_reflections(ray: Ray, objects: &[Box<dyn Object>], max_depth: u32)
        -> (Color, Option<Hit>);
    /// Shadow test: true if anything blocks `ray` within `max_distance`. Early-exits.
    pub fn occluded(ray: Ray, max_distance: f32, objects: &[Box<dyn Object>]) -> bool;
    // `trace` (nearest hit, no reflections) is private — internal helper only.
}
```
No camera knowledge — depends only on `objects`. Used by both the render pipeline (primary rays) and `lighting` (shadow rays).

### `camera`
```rust
pub struct Camera { position, look_at, up, fov, width, height }  // aspect derived from width/height
impl Camera {
    /// Construct from position, look-at target, up hint, vertical FOV (degrees), pixel dimensions.
    /// Aspect ratio is derived as width / height — cannot be set independently.
    pub fn new(position, look_at, up, fov, width: u32, height: u32) -> Self;
    /// Primary ray through pixel `(x, y)`. Uses the camera's stored width/height.
    pub fn ray_for_pixel(&self, x: u32, y: u32) -> Ray;
    /// Full-screen view frustum (4 side planes from the screen corners).
    pub fn frustum(&self) -> Frustum;
    /// Sub-frustum covering the pixel rect `[x0,x1) × [y0,y1)`.
    pub fn tile_frustum(&self, x0:u32, y0:u32, x1:u32, y1:u32) -> Frustum;
}
pub struct FrustumPlane { normal: Vec3, d: f32 }   // half-space; named to avoid clash with objects::Plane
pub struct Frustum { planes: [FrustumPlane; 4], origin: Vec3, forward: Vec3 }
impl Frustum {
    /// True if a sphere `(center, radius)` is at least partially inside all 4 side planes.
    pub fn contains(&self, center: Vec3, radius: f32) -> bool;
    /// Signed distance from the camera origin to `p` along the forward axis (sort key).
    pub fn depth(&self, p: Vec3) -> f32;
    /// Cull objects to those visible in this frustum, sorted near→far. Returns indices.
    pub fn cull_and_sort(&self, objects: &[Box<dyn Object>]) -> Vec<usize>;
}
// Pipeline:
//   1. frustum.cull_and_sort(objects) — culls off-screen, sorts near→far.
//   2. Per pixel: Camera::ray_for_pixel → Raycaster::trace_with_reflections.
//      Terminal hit (or sky) is shaded; mirrors contribute their tint along the chain.
//   3. Lighting calls Raycaster::occluded for each shadow ray.
// Optional tile loop: per-tile sub-frustum + per-tile object list for big scenes.
```
Requires `Object::bounding_sphere()` for culling. Uses the free `objects::reflect` helper for bounces.

### `lighting`
```rust
pub enum LightSource { Point { pos, intensity, color }, Directional { dir, intensity, color } }
/// Final shaded color at `hit`: sums each light's contribution (Lambert + distance falloff)
/// and casts a shadow ray per light against `objects` to mask occluded contributions.
pub fn shade(hit: &Hit, lights: &[LightSource], objects: &[Box<dyn Object>]) -> Color;
```

### `image`
```rust
pub struct Framebuffer { w, h, pixels: Vec<Color> }
impl Framebuffer {
    /// Allocate a `w × h` framebuffer initialized to black.
    pub fn new(w: u32, h: u32) -> Self;
    /// Write color `c` to pixel `(x, y)`. Panics if out of bounds.
    pub fn set(&mut self, x: u32, y: u32, c: Color);
    /// Serialize to PPM P3 (ASCII) into `out` — header + one `r g b` triple per pixel.
    pub fn write_ppm<W: Write>(&self, out: W) -> io::Result<()>;
}
```

### `world` (the glue + public API)
```rust
pub struct World {
    pub camera: Camera,
    pub objects: Vec<Box<dyn Object>>,
    pub lights: Vec<LightSource>,
    pub sky: SkyProperties,   // color/gradient when ray escapes
}
impl World {
    /// New empty world (no objects, no lights, default sky) with the given camera.
    pub fn new(camera: Camera) -> Self;
    /// Add any `Object` impl to the scene. Boxes it internally for dynamic dispatch.
    pub fn add_object<O: Object + 'static>(&mut self, o: O);
    /// Add a light source to the scene.
    pub fn add_light(&mut self, l: LightSource);
    /// Render the scene at `w × h`: build frustum, cull+sort objects, raycast every
    /// pixel, shade hits via `lighting`, fall back to `sky` on miss. Returns the framebuffer.
    pub fn render(&self, w: u32, h: u32) -> Framebuffer;
}
```
Re-exports everything users need: `pub use objects::*; pub use camera::*; pub use lighting::*; pub use image::*;`

### `scenes` (scenes-as-code)
Each scene is a Rust file exposing a single `pub fn build() -> World`. A registry in
`lib.rs` maps scene names (CLI strings) to those builders so `rt` can pick one at runtime.

```rust
// crates/scenes/src/sphere.rs
pub fn build() -> World {
    let cam = Camera::new(vec3(0.0, 1.0, 5.0), Vec3::ZERO, Vec3::Y, 60.0, 4.0/3.0);
    let mut w = World::new(cam);
    w.add_object(Sphere::new(Vec3::ZERO, 1.0, 0.0, Material::matte(rgb(0.8,0.3,0.3))));
    w.add_light(LightSource::Point { pos: vec3(2.0,4.0,2.0), intensity: 1.0, color: rgb(1,1,1) });
    w
}

// crates/scenes/src/lib.rs
pub mod sphere; pub mod plane_cube; pub mod all_objects; pub mod all_objects_alt_cam;

/// Look up a scene builder by CLI name. `None` if no such scene is registered.
pub fn get(name: &str) -> Option<fn() -> World> {
    match name {
        "sphere"               => Some(sphere::build),
        "plane_cube"           => Some(plane_cube::build),
        "all_objects"          => Some(all_objects::build),
        "all_objects_alt_cam"  => Some(all_objects_alt_cam::build),
        _ => None,
    }
}
/// All registered scene names (for `--list-scenes` and help text).
pub fn names() -> &'static [&'static str];
```
Adding a scene = new `.rs` file + one line in `get` and `names`. No parser, no JSON, type-checked.

### `liveview` (debug viewer, optional)
```rust
/// Open a window and render the world in real time, re-invoking the camera each frame.
/// Arrow keys / WASD orbit the camera; ESC quits. Resolution kept low for interactivity.
pub fn run(world: World, w: u32, h: u32);
```
Minimal: one window, blit the `Framebuffer` from `World::render` each frame, handle a few key events to mutate `world.camera`. Built on `minifb` (or `pixels`) — no GPU work, just a CPU buffer to a window. Compiled in only when the `liveview` cargo feature is enabled to avoid pulling a windowing dep into batch renders.

### `rt` (binary)
- Parses CLI: `--scene <name>`, `--out <file.ppm>`, `--width`, `--height`, `--live` (open `liveview` instead of writing PPM), `--list-scenes`, bonus flags (`-t`, `-r`, …).
- Resolves the scene via `scenes::get(name).unwrap()()` to build the `World`.
- Either calls `world.render(w,h).write_ppm(...)` or hands the world to `liveview::run(...)`.

## Shared Math
`glam` (Vec3, Mat4) added to every crate's `Cargo.toml`. No custom math crate.

## Dependency Graph (strict, no cycles)
```
rt        ──► world, scenes, liveview (liveview behind --live / cargo feature)
scenes    ──► world          (each scene file builds a World in pure Rust)
liveview  ──► world          (renders World every frame to a window)
world     ──► objects, raycaster, camera, lighting, image
camera    ──► objects        (needs Ray/Hit/Object trait)
raycaster ──► objects        (trace_with_reflections + occluded)
lighting  ──► objects, raycaster (uses Raycaster::occluded for shadow rays)
image     ──► (none, just std + glam for Color)
objects   ──► (none, just glam)
```
`Ray`, `Hit`, `Color`, `Material`, and the `reflect` helper all live in `objects` (lowest layer) so camera/lighting can depend on them without a circular dep.

## Parallel Work Split
- Person A: `objects` (4 shapes + trait + Ray/Hit).
- Person B: `camera` + `Raycaster`.
- Person C: `lighting` + sky/shadow.
- Person D: `image` (PPM writer + framebuffer).
- Person E: `world` API + `rt` binary + the 4 demo scenes + `DOCUMENTATION.md`.
- Person F: `scenes` (Rust scene builders) + `liveview` (debug window).

## Testing Strategy
- **Unit tests** live *inside each crate* (`#[cfg(test)] mod tests` next to the code). Each crate owner is responsible for testing their own public API. No shared unit-test infra.
- **Workspace e2e** in `tests/`: a fuzz-style harness that generates N random `World`s (random object counts/positions/materials/cameras/lights within sane bounds, fixed seed for reproducibility) and renders each at 100×100. Pass = no panic, no NaN pixels, PPM round-trips. Goal is crash/edge-case coverage, not visual correctness.

## Milestones
1. Lock shared types (`Ray`, `Hit`, `Color`, `Material`, `Object` trait) — unblocks everyone.
2. Each crate lands with its own local unit tests.
3. Workspace e2e harness green: 100+ random scenes render without panicking.
4. Produce the 4 required 800×600 `.ppm` scenes.
5. Bonus features behind CLI flags.
