# Book — The Intersection Scheduler

## What problem does it solve?

Multiple autonomous vehicles cross the same 18×18 grid intersection. They share cells. Without coordination they would collide. `Book` is the single authority that gives each vehicle a **full, collision-free schedule** the moment it spawns. The vehicle then just follows that schedule — no further negotiation needed.

---

## The grid and lanes

```
  col:  0 1 2 3 4 5 | 6  7  8  9 10 11 | 12 13 14 15 16 17
                    |                   |
  road starts at 6 (ROAD_START), ends at 11 (ROAD_END = ROAD_START + 5)
```

The six road columns/rows are named by their lane role:

| Variable | Value | Used by |
|----------|-------|---------|
| `nr` | 6 | North entry right-turn lane / East entry right-turn row |
| `ns` | 7 | North entry straight lane / East entry straight row |
| `nl` | 8 | North entry left-turn lane / East entry left-turn row |
| `sl` | 9 | South entry left-turn lane / West entry left-turn row |
| `ss` | 10 | South entry straight lane / West entry straight row |
| `sr` | 11 | South entry right-turn lane / West entry right-turn row |

Traffic entering from the **North** travels **south** (↓) and uses columns 6, 7, 8.
Traffic entering from the **South** travels **north** (↑) and uses columns 9, 10, 11.
Traffic entering from the **East** travels **west** (←) and uses rows 6, 7, 8.
Traffic entering from the **West** travels **east** (→) and uses rows 9, 10, 11.

The **intersection core** is the 6×6 square where all four directions overlap: cells with both `x` and `y` in `[ROAD_START, ROAD_END]` (i.e. `[6, 11]`). This region gets special treatment during scheduling (see Phase 2 below).

---

## Types used

### `Timestamp` (`u64`)
A game tick counter. Incremented once per frame in the main loop. Everything is scheduled in ticks.

### `Cell { x: i32, y: i32 }`
One grid square. `x` is the column (left = 0), `y` is the row (top = 0). `i32` because exit cells go off-grid into negative coordinates.

### `Direction`
The direction the vehicle is **facing** while occupying a cell: `North`, `South`, `East`, `West`. Used by the renderer to orient the car sprite. Changes at the pivot cell of a turn.

### `CellPosition { cell: Cell, dir: Direction }`
One step in a schedule: a grid square plus the car's facing direction there.

### `Routes`
An enum with 12 variants that fully describes where a car comes from and where it's going:

| Variant | Entry direction | Maneuver |
|---------|-----------------|----------|
| `NorthS` | North (going south) | Straight |
| `NorthR` | North | Right turn (exits west) |
| `NorthL` | North | Left turn (exits east) |
| `SouthS` | South (going north) | Straight |
| `SouthR` | South | Right turn (exits east) |
| `SouthL` | South | Left turn (exits west) |
| `EastS` | East (going west) | Straight |
| `EastR` | East | Right turn (exits north) |
| `EastL` | East | Left turn (exits south) |
| `WestS` | West (going east) | Straight |
| `WestR` | West | Right turn (exits south) |
| `WestL` | West | Left turn (exits north) |

### `MIN_TICKS_PER_CELL` (`u64 = 20`)
The minimum ticks to traverse one cell at maximum speed. Used as the lower bound when propagating arrival times forward and as the unit for all speed multipliers (1×, 2×, 3×).

---

## The internal booking table

```rust
struct Reservation {
    start: Timestamp,  // tick when this car enters the cell
    end:   Timestamp,  // tick when the car arrives at cell+RESERVATION_LOOKAHEAD
}

pub struct Book {
    cells: HashMap<(i32, i32), Vec<Reservation>>,
}
```

Each cell maps to a **sorted list of reservations** (sorted by `start`). When a new car books a cell, `earliest_free` scans this list to find the earliest slot that fits without overlap — including gaps *between* existing reservations, not just after the last one.

This is the key design choice: instead of storing only "when is this cell next free" (a single timestamp), we keep the full recent history, so a car that arrives in the middle of a quiet period can exploit the gap rather than waiting for the last reservation to expire.

The lists are kept compact by `prune()`, which drops reservations whose `end ≤ now` at the start of every `book()` call.

---

## How to call it

```rust
let mut book = Book::new();

if book.can_spawn(&route, now) {
    let schedule: Vec<(Timestamp, CellPosition)> = book.book(&route, now);
    // store schedule, create vehicle, etc.
}
```

### `can_spawn(&route, now) -> bool`

Checks whether the **first two on-grid cells** of the route are free at `now`. "Free" means `earliest_free` returns `now` unchanged — no existing reservation blocks entry at that tick. Read-only; registers nothing.

### `book(&route, now) -> Vec<(Timestamp, CellPosition)>`

Schedules a new car and returns its complete position schedule: one `(tick, cell_position)` pair per cell in travel order. Times are strictly increasing. The last entry is off-grid (one cell beyond the canvas edge).

Remove the vehicle when `now >= schedule.last().unwrap().0`.

---

## What happens inside `book()`

### Step 1 — Derive the cell path

`path_for_route(&route)` returns an ordered `Vec<CellPosition>` covering every cell from entry edge to `RESERVATION_LOOKAHEAD` off-grid exit cells. These trailing cells are never booked but are required so that `times[i + RESERVATION_LOOKAHEAD]` is always a valid index for the last on-grid cell.

### Step 2 — Prune stale reservations

Reservations with `end ≤ now` are removed from all buckets.

### Step 3 — Locate the core segment

The path is split into three regions:

```
[ approach cells ] [ core cells ] [ exit cells ] [ off-grid ]
  0 .. cs-1          cs .. ce       ce+1 ..
```

`cs` and `ce` are the first and last path indices where `is_core_cell` returns true. Each region is scheduled differently.

### Phase 1 — Approach

```
for i in 0..cs:
    times[i] = earliest_free(cell_i, times[i], 2*MIN_TICKS_PER_CELL)
    times[i+1] = times[i] + MIN_TICKS_PER_CELL
```

The gap-check duration is `APPROACH_RESERVATION_TICKS` (= `RESERVATION_LOOKAHEAD * MIN_TICKS_PER_CELL`) — the exact size of an approach reservation (`[times[i], times[i + RESERVATION_LOOKAHEAD]]` with no padding, at minimum speed). This means a new car can slot into any gap wider than `APPROACH_RESERVATION_TICKS` between prior bookings on approach cells, rather than waiting until after the last booking ends.

After this phase, `times[cs]` is the earliest tick the car could physically reach the core entry.

### Phase 2 — Core traversal

Two strategies are tried and the better one is used.

#### Variable speed (preferred)

The car enters the core immediately at `times[cs]` (no extra waiting outside) and traverses each cell greedily:

```
var_core[0] = earliest_free(core_cell_0, times[cs], var_dur)
for k in 1 .. n_core:
    var_core[k] = earliest_free(core_cell_k, var_core[k-1] + MIN_TICKS_PER_CELL, var_dur)
    if var_core[k] - var_core[k-1] > MAX_TICKS_PER_CELL:  // stall cap: 3×MIN
        reject variable speed, try sprint
```

`var_dur = RESERVATION_LOOKAHEAD * MAX_TICKS_PER_CELL` matches the worst-case registered window for a core cell: cell `i` is held until `times[i + RESERVATION_LOOKAHEAD]`, and at maximum slow speed that end tick is at most `t + RESERVATION_LOOKAHEAD * MAX_TICKS_PER_CELL`. Using the same expression for the gap check ensures it never passes for a slot too narrow to hold the actual registration.

Each inter-cell step is at least `MIN_TICKS_PER_CELL` (the car always moves). If any step would exceed `MAX_TICKS_PER_CELL = 3 * MIN_TICKS_PER_CELL`, the car would look like it is stalling inside the intersection — variable speed is rejected and the sprint fallback is used instead.

#### Sprint fallback

The car waits outside the core until it can traverse all core cells at `MIN_TICKS_PER_CELL` without any mid-intersection gaps. `sprint_entry` finds the earliest valid entry time via a fixpoint loop:

```
entry = times[cs]
repeat (up to n_core times):
    for each core cell k:
        t_at_k = entry + k * MIN_TICKS_PER_CELL
        actual = earliest_free(core_cell_k, t_at_k, RESERVATION_LOOKAHEAD * MIN_TICKS_PER_CELL)
        if actual > t_at_k:
            entry = max(entry, actual - k * MIN_TICKS_PER_CELL)
    stop if entry did not change
```

The fixpoint is necessary because pushing `entry` later for cell `k` can create a new conflict at an earlier cell `k'` — a single left-to-right pass would miss it. The loop converges in at most `n_core` iterations because `entry` is monotone non-decreasing and bounded.

### Phase 3 — Exit

All cells after the core are scheduled at maximum speed:

```
times[i+1] = times[i] + MIN_TICKS_PER_CELL   for i from ce onward
```

### Phase 4 — Smooth approach

After the core entry time is fixed, the approach times are redistributed backward so the car decelerates gradually from its spawn point rather than rushing to the core entry and waiting there.

```
smooth_times(&mut times[..=cs])
```

`smooth_times` scans for "wait points" — consecutive entries where `times[i] - times[i-1] > MIN_TICKS_PER_CELL`. For each such point it linearly interpolates all intermediate times between the previous anchor and this wait point.

Example: `[0, 20, 40, 60, 120, 140]` → `[0, 30, 60, 90, 120, 140]`. The car now decelerates uniformly from spawn, arriving at the core at exactly tick 120 with no stopping.

### Step 4 — Register reservations

```
for each on-grid cell i:
    end = times[i + RESERVATION_LOOKAHEAD]
    insert Reservation { start: times[i], end }
    into cells[(x,y)], maintaining sort order by start
```

`times[i + RESERVATION_LOOKAHEAD]` is always valid because the path ends with `RESERVATION_LOOKAHEAD` off-grid cells, so the index never goes out of bounds.

### Step 5 — Return the schedule

The `times` array is zipped with the `path` array and returned. All but one off-grid cell are removed (`RESERVATION_LOOKAHEAD - 1` elements truncated), so callers always get exactly one off-grid entry — the tick when the vehicle has fully exited.

---

## How `earliest_free` works

```
fn earliest_free(cell, min_t, duration) -> Timestamp:
    candidate = min_t
    for each Reservation R in cells[cell] (sorted by R.start):
        if candidate + duration <= R.start:
            STOP — the window [candidate, candidate+duration] fits before R
        if candidate >= R.end:
            SKIP — already past this reservation
        candidate = R.end   — inside R's window, jump to when R expires
    return candidate
```

`duration` is the size of the window we are trying to place. Using the correct duration is critical:

| Context | Duration used | Why |
|---------|--------------|-----|
| Approach gap check | `APPROACH_RESERVATION_TICKS` (= `RESERVATION_LOOKAHEAD * MIN_TICKS_PER_CELL`) | Exact: `[times[i], times[i + RESERVATION_LOOKAHEAD]]` at MIN speed |
| Variable-speed core | `RESERVATION_LOOKAHEAD * MAX_TICKS_PER_CELL` | Matches worst-case registered window at MAX speed |
| Sprint at speed `s` | `RESERVATION_LOOKAHEAD * s` | Matches registered window at uniform speed `s` |
| `can_spawn` | `APPROACH_RESERVATION_TICKS` | Approach cells only |

An overestimate for `duration` is always safe (the placed reservation will never overlap an existing one) but may cause the car to miss a gap it could have used. An underestimate is unsafe: the gap check may pass, but the actual registration would overlap the next reservation.

### Example walkthrough

Reservations on cell (7,5): `R1 = [100, 160]`, `R2 = [200, 260]`. `duration = 40`.

**New car arrives at t=70:**
- R1: `70 + 40 = 110 > 100` — gap too small. `70 < 160` → `candidate = 160`
- R2: `160 + 40 = 200 ≤ 200` — fits before R2. STOP.
- Result: **160** (slotted in the gap between R1 and R2)

**New car arrives at t=170:**
- R1: `170 + 40 = 210 > 100` — gap fails. `170 ≥ 160` → SKIP
- R2: `170 + 40 = 210 > 200` — gap too small. `170 < 260` → `candidate = 260`
- Result: **260** (must wait until after R2)

**New car arrives at t=50:**
- R1: `50 + 40 = 90 ≤ 100` — fits before R1. STOP.
- Result: **50** (goes before both reservations)

---

## Safety guarantees

### No same-cell collision

`earliest_free` ensures every new car's entry time is outside all existing reservation windows for that cell. Each cell is booked independently, so no two cars share a cell at the same tick.

### Two-cell safe distance

`end = times[i+2]` means the cell is held until the car has moved two cells further. The next car cannot enter until then, guaranteeing a minimum two-cell gap on any shared segment.

### Gap reuse does not break safety

When a car slots into a gap *before* an existing reservation (gap check passed), the placed window `[candidate, candidate + duration]` ends at or before `R.start`. At registration, `end = times[i+2]` which is ≤ `candidate + duration` (the duration estimate is an upper bound). So `end ≤ R.start` — no overlap.

### Crossing routes

When routes cross at a core cell, the first car's reservation blocks the second car's `earliest_free` call for that cell. The cell-by-cell independence handles all crossing points automatically.

## What the caller must do

1. **Create one `Book` for the whole simulation.** All vehicles go through the same instance.

2. **Call `can_spawn` first.** If it returns false, defer spawning. Don't call `book` when the approach lanes are occupied.

3. **Call `book.book(&route, now)` the moment a vehicle spawns.** The returned schedule is fixed — it does not update if conditions change later.

4. **Store the returned schedule.** The rendering system reads it to interpolate pixel positions.

5. **Remove a vehicle when `now >= schedule.last().unwrap().0`.** That timestamp is the first off-grid tick — the car is off canvas.

6. **Do not call `book` from multiple threads.** `Book` takes `&mut self`; it is not thread-safe.

---

## What `Book` does NOT do

- Move vehicles. It only assigns timestamps to cell positions.
- Know about pixel positions. That is the renderer's job.
- Track vehicle IDs. Pair the returned schedule with an ID yourself.
- Re-schedule. A schedule is immutable once returned. If a route changes, book a new one.

---

## Quick reference

```
Book::new()                          create the shared booking table
book.can_spawn(&route, now)          true if first 2 approach cells are free at now
book.book(&route, now)               schedule a car, get back its position schedule
schedule.last().unwrap().0           tick when the car has fully exited; remove it then
schedule[i].0                        tick when car arrives at the i-th cell
schedule[i].1.cell                   the Cell at that step
schedule[i].1.dir                    the Direction the car faces there
```
