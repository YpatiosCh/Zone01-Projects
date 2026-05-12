# rt

The executable entry point. Parses command-line arguments and routes to one of three launch modes: interactive menu, liveview window, or headless render-to-file.

---

## Usage

```
rt [--menu] [--live] [--render] [--scene <name>] [--out <file>]
   [--width N] [--height N] [--list-scenes] [-h | --help]
```

Default behavior (no arguments): opens the startup GUI menu.

---

## Launch modes

### `--menu` (default)

Opens the main menu GUI. From there the user can navigate to the scene picker, then liveview, and back. Pressing `Esc` at any point eventually reaches the main menu `Exit` button.

Flow:

```
main menu ──[Start]──▶ scene picker ──[scene]──▶ liveview
    ▲                       │                        │
    └──────[ExitToMenu]─────┘         [Esc/Close]───┘
```

### `--live`

Skips the main menu and opens liveview directly for the named scene (defaults to `sphere` if `--scene` is not given). If the user presses `Esc` inside liveview, they are dropped into the main menu loop rather than the application closing.

```sh
rt --live --scene all_objects
```

### `--render`

Renders the named scene to a file without opening any window. Exits immediately after writing the files.

```sh
rt --render --scene plane_four_spheres --out renders/my_scene.ppm
rt --render --scene sphere --width 1920 --height 1080
```

Output files: both `.ppm` (ASCII) and `.png` (compressed) are always written. If `--out` is not given, the output path defaults to `scenes_out/<scene_name>.ppm` (and `.png` alongside it). The output directory is created automatically if it does not exist.

---

## Options

| Flag | Default | Description |
|---|---|---|
| `--scene <name>` | `sphere` | Scene to render or open in liveview |
| `--out <path>` | `scenes_out/<scene>.ppm` | PPM output path (also implies `--render`) |
| `--width N` | `800` | Image or window width in pixels |
| `--height N` | `600` | Image or window height in pixels |
| `--list-scenes` | — | Print all available scene names and exit |
| `--menu` | — | Open the startup menu (this is the default) |
| `--live` | — | Open liveview directly |
| `--render` | — | Render to file |
| `-h` / `--help` | — | Print usage and exit |

---

## Internal structure

### `LaunchMode`

```rust
enum LaunchMode { Menu, Live, Render }
```

Determined by parsing args sequentially. `--out` implies `Render`. Flags appearing later override earlier ones.

### `run_main_menu_loop`

Loops between `menu::run` → `menu::run_scene_menu` → `liveview::run` until either `MenuAction::Exit` or `LiveviewAction::ExitApp` is returned. Handles the case where liveview returns `OpenSceneMenu` by looping back to the scene picker without returning to the main menu.

### `run_live_loop`

Loops `liveview::run` → `menu::run_scene_menu` for inline scene switching. Returns `LiveLoopExit::MainMenu` when the user backs out to the main menu, or `LiveLoopExit::ExitApp` when the window is closed.

### `scene_builder`

A thin wrapper around `scenes::get` that converts `None` into a formatted error string listing available scene names.

---

## Exit codes

| Code | Meaning |
|---|---|
| `0` | Success |
| `1` | Runtime error (window open failed, file write failed) |
| `2` | Invalid arguments (unknown flag, bad value, unknown scene) |

---

## Examples

```sh
# Open the GUI (default)
rt

# Render a specific scene to default path (scenes_out/sphere.ppm)
rt --render --scene sphere

# Render at 4K
rt --render --scene all_objects --width 3840 --height 2160 --out renders/4k.ppm

# Open liveview directly for the reflective plane scene
rt --live --scene plane_reflective

# See all scene names
rt --list-scenes
```
