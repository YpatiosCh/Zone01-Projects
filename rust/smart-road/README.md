## Smart Road

### Overview

Smart Road is an autonomous vehicle intersection simulation written in Rust with SDL2. Vehicles cross an 18×18 grid intersection from all four directions — no traffic lights, no human drivers. A reservation-based scheduler coordinates them entirely, issuing each vehicle a complete, collision-free route schedule the moment it spawns.

---

### How it works

#### Scheduling — `Book`

The core of the simulation is `Book`, a shared intersection scheduler. When a vehicle spawns, `Book` assigns it a full sequence of `(tick, cell)` pairs covering every grid cell from entry to exit. Once issued, the schedule is immutable — the vehicle just follows it.

Each grid cell has a sorted list of time reservations. Before assigning a cell to a new vehicle, `Book` scans that list to find the earliest gap wide enough to fit without overlap. This means cars can slot into quiet periods between existing reservations rather than always queuing behind the last one.

The 6×6 central zone where all four traffic directions overlap gets special treatment:

- **Variable speed** (preferred): the car enters immediately and traverses each core cell as early as possible, adjusting speed cell by cell. If any step would stall the car for too long, this strategy is rejected.
- **Sprint fallback**: the car waits outside the core until it can cross all cells at `MIN_TICKS_PER_CELL` (maximum speed) with no mid-intersection gaps.

After the core entry time is fixed, approach timestamps are smoothed so the car decelerates gradually from spawn rather than rushing to the boundary and stopping.

Reservations extend `RESERVATION_LOOKAHEAD` cells ahead of a car's current position, guaranteeing a safe gap between vehicles sharing any segment of road.

#### Vehicles and routes

Each vehicle is assigned one of 12 routes — right turn, straight, or left turn from each of the four cardinal directions. Routes are chosen randomly at spawn. The vehicle's pixel position at any tick is interpolated from its schedule, and the sprite rotates progressively through turns.

Speed is implicit in the schedule: the ticks-per-cell interval between consecutive entries determines how fast the car appears to move. The range is `MIN_TICKS_PER_CELL` (fastest) to `MAX_TICKS_PER_CELL` (slowest stall threshold), both configurable in `config.rs`.

---

### Controls

| Key | Action |
|-----|--------|
| `Arrow Up` | Spawn vehicle from south (heading north) |
| `Arrow Down` | Spawn vehicle from north (heading south) |
| `Arrow Right` | Spawn vehicle from west (heading east) |
| `Arrow Left` | Spawn vehicle from east (heading west) |
| `R` | Toggle auto-spawn (random directions, continuous) |
| `--auto` flag | Start in auto-spawn mode from launch |
| `--highlight` flag | Draw bounding boxes around all vehicles; red outline on close call |
| `Esc` | End simulation and show statistics |

---

### Statistics

Displayed on exit:

- **Vehicles through intersection** — total completed crossings
- **Close calls** — unique pairs of vehicles where one car's next scheduled cell was occupied by another (one cell apart, approaching)
- **Collisions** — unique pairs of vehicles whose rendered positions came within `COLLISION_THRESHOLD_PX` of each other; shown as a semi-transparent red fill over the car sprite
- **Slowest / fastest cell crossing** — min and max instantaneous speed across all vehicles
- **Shortest / longest transit** — time from entry to exit for the quickest and slowest vehicles
- **Total runtime** — simulation duration in ticks and seconds

Close calls and collisions are also printed to the terminal as they occur, with coordinates.

### System dependencies

This project uses SDL2 and the optional image/ttf extensions. On a development machine (and in CI) you must install the system development packages so the Rust `sdl2` crates can link to the native libraries.

Common install commands:

- Debian / Ubuntu (example):

```bash
sudo apt-get update
sudo apt-get install -y build-essential pkg-config libsdl2-dev libsdl2-image-dev libsdl2-ttf-dev
```

- Fedora / RHEL (dnf):

```bash
sudo dnf install -y @development-tools pkgconfig SDL2-devel SDL2_image-devel SDL2_ttf-devel
```

- Arch Linux:

```bash
sudo pacman -Syu --noconfirm base-devel pkgconf sdl2 sdl2_image sdl2_ttf
```

- Alpine:

```bash
sudo apk add --no-cache build-base pkgconfig sdl2-dev sdl2_image sdl2_ttf
```

- Windows (MSVC) using vcpkg:

```powershell
cd C:\Users\<user>\vcpkg
./vcpkg install sdl2:x64-windows sdl2-image:x64-windows sdl2-ttf:x64-windows
```

Then build from your project folder with:

```powershell
$env:VCPKG_ROOT='C:\Users\<user>\vcpkg'
$env:VCPKGRS_DYNAMIC='1'
cargo build
```

For windows using vcpkg, change this line in Cargo.toml:
```
sdl2 = { version = "0.37", features = ["image", "ttf"] }
```
to:
```
sdl2 = { version = "0.37", features = ["image", "ttf", "use-vcpkg"] }
```

CI note: ensure your CI image installs these packages before running `cargo build`. For example, on GitHub Actions you can run an `apt-get install` step on `ubuntu-latest` runners.

If you cannot install system packages system-wide, you can build and install SDL2 into a local prefix and set `PKG_CONFIG_PATH` and `LD_LIBRARY_PATH` appropriately before running `cargo build`.
