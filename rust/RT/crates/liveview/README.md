# liveview

Real-time interactive 3D viewer. Opens a window, renders the scene at a reduced resolution, and re-renders on every frame where the camera or lights changed.

---

## Entry point

```rust
pub fn run(world: World, w: u32, h: u32) -> LiveviewAction
```

Opens a window of at least 640×480 and enters the render loop. Returns when the user presses `Esc` or closes the window:

- `LiveviewAction::OpenSceneMenu` — user pressed `Esc`; the caller should open the scene selection menu.
- `LiveviewAction::ExitApp` — window was closed or a rendering error occurred.

---

## Camera controls

| Input | Action |
|---|---|
| `W` | Move forward (flat, ignores pitch) |
| `S` | Move backward |
| `A` | Strafe left |
| `D` | Strafe right |
| `Space` | Move up |
| `Left Ctrl` / `Right Ctrl` | Move down |
| Mouse move | Look around (yaw + pitch) |
| `Esc` | Exit to scene menu |
| `F` | Save a screenshot |

Movement uses a **flat forward vector** (pitch ignored) so pressing `W` while looking up does not make the camera drift upward. The camera stays at the same height and moves horizontally. Vertical movement is separate (`Space`/`Ctrl`).

---

## Light controls

| Input | Action |
|---|---|
| `←` / `→` arrow keys | Move point lights left/right; orbit directional lights around Y |
| `↑` / `↓` arrow keys | Move point lights up/down; tilt directional lights |

---

## Render scale

Liveview renders at `1/RENDER_SCALE` of the window resolution (default scale = 3, so a 640×480 window renders at 213×160) and upscales the result with nearest-neighbor sampling. This keeps the frame rate interactive even on complex scenes. The low-resolution image is visibly pixelated but updates in real time.

The upscaler maps each destination pixel `(dx, dy)` to source pixel `(src_x, src_y)` using integer scaling:

```
src_x = dx * src_width / dst_width
src_y = dy * src_height / dst_height
```

Screenshots bypass the render scale and save at full window resolution.

---

## Camera representation

The camera orientation is stored as **yaw** (horizontal rotation) and **pitch** (vertical tilt) angles rather than a look-at vector. This is standard for first-person navigation because it prevents gimbal lock and makes it natural to move in the horizontal plane.

The look direction is reconstructed from angles each frame:

```rust
fn direction_from_angles(yaw: f32, pitch: f32) -> Vec3 {
    let cos_pitch = pitch.cos();
    Vec3::new(cos_pitch * yaw.cos(), pitch.sin(), cos_pitch * yaw.sin())
}
```

Pitch is clamped to `±MAX_PITCH = ±1.45 rad ≈ ±83°` to prevent the camera from flipping over at the poles.

Mouse movement is computed as the delta from the window center. After reading the delta, the cursor is immediately recentered. The cursor is hidden and confined to the window while liveview is active.

---

## Floor collision

The camera is prevented from passing through floor objects using `Object::floor_height_at`. After every movement step, all objects are queried:

```rust
fn clamp_camera_to_floor(world: &World, position: Vec3) -> Vec3 {
    let floor_y = world.objects.iter()
        .filter_map(|o| o.floor_height_at(position))
        .reduce(f32::max);

    if let Some(floor_y) = floor_y {
        Vec3::new(position.x, position.y.max(floor_y + CAMERA_CLEARANCE), position.z)
    } else {
        position
    }
}
```

`CAMERA_CLEARANCE = 0.1` keeps the camera eye slightly above the floor surface. Only `Plane` implements `floor_height_at`; all other shapes return `None`.

---

## Screenshots

Pressing `F` saves the current scene at full window resolution (bypassing the render scale) to `screenshots/screenshot_N.ppm` and `screenshots/screenshot_N.png`, where `N` is the next available index. The `screenshots/` directory is created automatically if it does not exist.

---

## Constants

| Constant | Value | Description |
|---|---|---|
| `MOVE_SPEED` | `0.25` | Units per frame for WASD movement |
| `VERTICAL_MOVE_SPEED` | `0.20` | Units per frame for vertical movement |
| `MOUSE_SENSITIVITY` | `0.004` | Radians per pixel of mouse delta |
| `MAX_PITCH` | `1.45` | Maximum pitch angle in radians (~83°) |
| `RENDER_SCALE` | `3` | Divisor applied to window dimensions for render resolution |
| `CAMERA_CLEARANCE` | `0.1` | Minimum gap between camera and floor |
| `LIGHT_MOVE_SPEED` | `0.25` | Units per frame for point light translation |
| `LIGHT_ROTATE_SPEED` | `0.03` | Radians per frame for directional light rotation |
