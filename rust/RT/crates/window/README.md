# window

Cross-platform windowing, input handling, and software rendering surface. Wraps `minifb` with additional utilities for cursor management, bitmap font text rendering, and rectangle drawing.

---

## `AppWindow`

The main type. Manages a native OS window with a 32-bit RGB pixel buffer.

```rust
use window::AppWindow;

let mut window = AppWindow::new("My Window", 800, 600)?;
```

The window targets **60 FPS** via `minifb`'s `set_target_fps`. The pixel buffer is initialized to black.

### Frame loop

```rust
while window.is_open() {
    // 1. Handle input
    // 2. Draw into buffer
    // 3. Present
    window.present()?;
}
```

`present()` calls `minifb::Window::update_with_buffer`, which flushes the pixel buffer to the screen and processes OS events.

---

## Input

### Keyboard

```rust
// One-shot press (fires once per key-down, does not repeat)
window.is_key_pressed(WindowKey::Escape)

// Held down (fires every frame while held)
window.is_key_down(WindowKey::W)
```

`WindowKey` is re-exported from `minifb::Key`. All standard keyboard keys are available.

### Mouse

```rust
// One-shot left click: returns pixel position or None
if let Some((x, y)) = window.take_left_click() { ... }

// Current mouse position (clamped to window bounds)
if let Some((x, y)) = window.mouse_position() { ... }

// Delta from window center (used by liveview for mouse look)
if let Some((dx, dy)) = window.mouse_delta_from_center() { ... }
```

`take_left_click` uses edge-detection: it returns a position only on the frame the button transitions from up to down, preventing a single click from registering multiple times.

`mouse_delta_from_center` reads the native OS cursor position and immediately recenters the cursor so the next call starts fresh. This is how the liveview first-person look works without the cursor hitting the window edge.

---

## Cursor management

```rust
window.set_cursor_visible(false); // hide cursor for liveview
window.center_cursor();            // warp cursor to window center
window.confine_cursor(true);       // trap cursor inside window
window.confine_cursor(false);      // release cursor
```

Confinement is implemented using platform APIs:

- **Windows** — `ClipCursor` with the window's client rectangle
- **Linux** — `XGrabPointer` via the `x11-dl` crate

On other platforms the flag is stored but no OS call is made. The `Drop` impl automatically releases cursor confinement when the window is destroyed.

---

## Drawing

### Clearing

```rust
window.clear(0x101820); // fill with RGB hex color
```

The color is a packed `0x00RRGGBB` u32 (same format `minifb` expects).

### Rectangles

```rust
window.fill_rect(rect, 0x5FBF77);   // filled rectangle
window.stroke_rect(rect, 0xF4EBD0); // 1-pixel border
```

`Rect` is a simple `{ x, y, w, h }` struct with a `contains(px, py) -> bool` hit test.

### Text

```rust
// Top-left aligned
window.draw_text(x, y, "Hello", 0xF4EBD0, scale);

// Centered inside a rect
window.draw_text_centered(rect, "START", 0xF4EBD0, scale);
```

Text is rendered using the `font8x8` crate's basic 8×8 bitmap font. Each character is an 8×8 pixel glyph. `scale` is an integer pixel multiplier: `scale = 2` produces 16×16 characters, `scale = 4` produces 32×32 characters.

`measure_text(text, scale)` returns `(width, height)` in pixels for layout calculations.

### Pixel buffer

```rust
// Replace the entire pixel buffer (for raycaster output)
window.set_pixels(&pixels); // pixels: &[u32], length must equal width * height
```

Each `u32` is `0x00RRGGBB`. The liveview crate packs its `Color` values into this format with `(r << 16) | (g << 8) | b`.

---

## `Rect`

```rust
pub struct Rect { pub x: usize, pub y: usize, pub w: usize, pub h: usize }
```

Used for drawing and for button hit testing in the menu. `contains(px, py)` returns true if the point is inside the half-open interval `[x, x+w) × [y, y+h)`.

---

## Platform notes

The cursor delta and confinement features require native API calls that differ between platforms:

| Feature | Windows | Linux | Other |
|---|---|---|---|
| Cursor confinement | `ClipCursor` | `XGrabPointer` | no-op |
| Cursor center warp | `SetCursorPos` | `XWarpPointer` | no-op |
| Native cursor position | `GetCursorPos` + `ScreenToClient` | `XQueryPointer` | `minifb` fallback |

`x11-dl` is loaded dynamically at runtime on Linux; if the library is unavailable, cursor features silently degrade.
