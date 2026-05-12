# menu

GUI menus for the ray tracer. Provides the main menu (Start / Credits / Exit) and the scene selection screen. Both menus block until the user makes a choice and return a typed action value to the caller.

---

## Main menu — `menu::run`

```rust
pub fn run(width: u32, height: u32) -> Result<MenuAction, WindowError>
```

Opens the main menu window and returns when the user clicks a button or presses `Esc`.

```rust
match menu::run(800, 600)? {
    MenuAction::Start => { /* open scene selection */ }
    MenuAction::Exit  => { /* quit */ }
}
```

### Layout

```
┌────────────────────────────────┐
│         RAY TRACING            │  ← title (scale 4)
│     CPU renderer project       │  ← subtitle (scale 2)
│                                │
│          [ Start ]             │  ← green button
│         [ Credits ]            │  ← yellow button
│          [ Exit ]              │  ← red button
└────────────────────────────────┘
```

The Credits button shows an overlay screen with contributor names. Pressing `Esc` or clicking `Back` returns to the main menu. The main menu itself returns `MenuAction::Exit` on `Esc`.

---

## Scene selection — `menu::run_scene_menu`

```rust
pub fn run_scene_menu(width: u32, height: u32) -> Result<SceneMenuAction, WindowError>
```

Opens the scene picker and returns when the user selects a scene or exits.

```rust
match menu::run_scene_menu(800, 600)? {
    SceneMenuAction::Scene(name) => { /* load scene by name */ }
    SceneMenuAction::ExitToMainMenu => { /* return to main menu */ }
}
```

### Layout

```
┌────────────────────────────────┐
│           SCENES               │
│  Choose the next liveview scene│
│                                │
│         [ Sphere ]             │
│        [ Cylinder ]            │
│          [ Plane ]             │
│    [ Plane + 4 Spheres ]       │
│    [ Plane Reflective ]        │
│      [ Plane + Cube ]          │
│       [ All Objects ]          │
│  [ All Objects (Alt Cam) ]     │
│                                │
│          [ Exit ]              │
└────────────────────────────────┘
```

Each scene button has its own fill color and text color for visual variety.

---

## Return types

### `MenuAction`

```rust
pub enum MenuAction {
    Start, // proceed to scene selection
    Exit,  // quit the application
}
```

### `SceneMenuAction`

```rust
pub enum SceneMenuAction {
    Scene(&'static str), // name of the selected scene (e.g. "plane_four_spheres")
    ExitToMainMenu,      // user backed out
}
```

The `&'static str` scene names match the keys in `scenes::REGISTRY`, so they can be passed directly to `scenes::get`.

---

## Scene table

The list of scenes shown in the menu is a static table at the top of `lib.rs`:

```rust
static SCENES: &[(&str, &str, u32, u32)] = &[
    // (scene_name, button_label, fill_color_hex, text_color_hex)
    ("sphere",            "Sphere",              0x5FBF77, 0x0F1A14),
    ("cylinder",          "Cylinder",            0x6CB4D9, 0x0A1720),
    ("plane",             "Plane",               0xB7B7B7, 0x171717),
    ("plane_four_spheres","Plane + 4 Spheres",   0x9D8DF1, 0x111022),
    ("plane_reflective",  "Plane Reflective",    0xE9C46A, 0x241A05),
    ("plane_cube",        "Plane + Cube",        0xE07B54, 0x2A0F09),
    ("all_objects",       "All Objects",         0x5FBF77, 0x0F1A14),
    ("all_objects_alt_cam","All Objects (Alt Cam)",0x3DA88A, 0x071410),
];
```

To add a scene: add a row here and add the corresponding `build` function to `crates/scenes/src/`.

---

## Layout helpers

The drawing layer has three private helpers that produce `Rect` values for button placement:

- `centered_button(width, y, btn_w, btn_h)` — centers a button horizontally at a given Y position.
- `footer_button(width, height, btn_w, btn_h)` — places a button 48 pixels above the bottom edge of the window.
- `draw_button(window, rect, label, fill, text)` — fills a `Rect` with a solid color, draws a 1-pixel border, and centers the label text inside.

Buttons are tested for mouse hits using `Rect::contains(mx, my)`, with the pixel position coming from `window.take_left_click()`.

---

## Minimum window size

Both `run` and `run_scene_menu` enforce a minimum window size of **640×480**. Smaller values passed in are clamped up silently.
