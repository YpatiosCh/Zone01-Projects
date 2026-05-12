
## Two small details worth knowing

**1. The Manhattan shortcut.** Filler maps are open rectangles with no
walls. On an open grid, BFS distance is always `|Δrow| + |Δcol|`. So when
we're scoring a placement, instead of running a whole new BFS from every
candidate piece layout, we just use that formula directly. That's the
`((r as i32 - pr).abs() + (c as i32 - pc).abs())` line in
`score_placement`.

**2. The tiebreaker.** The real score is actually a pair:
`(voronoi_gain, -enemy_dist_sum)`. The second number is "sum of
`enemy_dist` over all the cells the piece would cover" — smaller is
better, meaning the piece lands deeper into enemy territory. It only
matters when `voronoi_gain` is the same across multiple placements (common
late in the game), and it nudges the bot toward *attacking* the enemy
rather than piling up in a corner.

---

## Code map

If you want to match this to the source:

| Concept                             | Function                              | File          |
|------------------------------------|---------------------------------------|---------------|
| Main game loop                      | `run`                                 | `filler.rs`   |
| Is this placement legal?            | `is_valid_placement`                  | `strategy.rs` |
| List all legal placements           | `enumerate_placements`                | `strategy.rs` |
| Build `my_dist` and `enemy_dist`    | `bfs_distance_maps` (+ `bfs_from`)    | `strategy.rs` |
| Score one placement                 | `score_placement`                     | `strategy.rs` |
| Pick the best placement             | `pick_best_placement`                 | `strategy.rs` |

---

## Full code pipeline — one turn, 5×5 board, piece `OO`

This walks through **every function call** on the same tiny example, with
concrete values at every step. Starting point: the parser has already
produced `board` and `piece`. We pick up from there.

### Setup

Board (5×5, rows × cols):

```
.....
.@...
.....
...$.
.....
```

After parsing, `board` is a `Board` struct with:
- `width = 5`, `height = 5`
- `cells[1][1] = Mine`, `cells[3][3] = Opponent`, everything else `Empty`

Piece `OO` (1 row tall, 2 cols wide):
- `width = 2`, `height = 1`
- `cells[0] = [true, true]`

---

### 1 — `run()` in `filler.rs`

After the two `parse_*` calls, this is the code that fires:

```rust
let placements = strategy::enumerate_placements(&board, &piece);
let chosen = strategy::pick_best_placement(&board, &piece, &placements);

match chosen {
    Some(pos) => writeln!(writer, "{} {}", pos.col, pos.row).unwrap(),
    None      => writeln!(writer, "0 0").unwrap(),
}
```

Two function calls, then print `col row`. Let's open each one.

---

### 2 — `enumerate_placements(&board, &piece)`

Source:

```rust
pub fn enumerate_placements(board: &Board, piece: &Piece) -> Vec<Position> {
    let mut out = Vec::new();
    for row in 0..board.height {          // 0..5
        for col in 0..board.width {       // 0..5
            let pos = Position { row, col };
            if is_valid_placement(board, piece, pos) {
                out.push(pos);
            }
        }
    }
    out
}
```

It's a brute-force double loop: try every board cell as the top-left
corner of the piece, keep the ones that pass the legality check. On a 5×5
board that's 25 candidate positions.

For each position it calls `is_valid_placement`, which decides yes/no:

```rust
fn is_valid_placement(board: &Board, piece: &Piece, pos: Position) -> bool {
    let mut overlap = 0u32;
    for pr in 0..piece.height {
        for pc in 0..piece.width {
            if !piece.cells[pr][pc] { continue; }   // transparent → ignore
            let br = pos.row + pr;
            let bc = pos.col + pc;
            if br >= board.height || bc >= board.width {
                return false;                       // off-board
            }
            match board.cells[br][bc] {
                Cell::Opponent => return false,     // banned overlap
                Cell::Mine     => overlap += 1,     // anchor candidate
                Cell::Empty    => {}                // fine
            }
        }
    }
    overlap == 1                                    // rule 3: exactly one
}
```

Let's hand-trace the interesting positions:

- `pos = (0,0)`: piece cells `(0,0),(0,1)` land on board `(0,0),(0,1)`.
  Both Empty → `overlap = 0` → returns **false**.
- `pos = (1,0)`: lands on `(1,0),(1,1)`. `(1,0)` Empty, `(1,1)` Mine →
  `overlap = 1` → returns **true** ✅
- `pos = (1,1)`: lands on `(1,1),(1,2)`. `(1,1)` Mine, `(1,2)` Empty →
  `overlap = 1` → returns **true** ✅
- `pos = (1,4)`: `(1,4)` is in bounds but `(1,5)` is out → returns
  **false**.
- `pos = (3,2)`: `(3,2)` Empty, `(3,3)` Opponent → returns **false**.
- Every other position: piece doesn't touch (1,1), so `overlap = 0` →
  false.

Only two positions pass, so:

```
enumerate_placements returns  [Position{1,0}, Position{1,1}]
```

---

### 3 — `pick_best_placement(&board, &piece, &placements)`

Source:

```rust
pub fn pick_best_placement(
    board: &Board, piece: &Piece, placements: &[Position],
) -> Option<Position> {
    let mut best_pos = *placements.first()?;              // (a)

    let (my_dist, enemy_dist) = bfs_distance_maps(board); // (b)

    let mut best_score =
        score_placement(&my_dist, &enemy_dist, board, piece, &best_pos);   // (c)

    for pos in &placements[1..] {
        let score = score_placement(&my_dist, &enemy_dist, board, piece, pos);
        if score > best_score {                            // (d)
            best_score = score;
            best_pos = *pos;
        }
    }

    Some(best_pos)
}
```

What happens:

- **(a)** Grabs the first candidate as the initial "best so far":
  `best_pos = Position{1, 0}`. If `placements` was empty, `?` returns
  `None` and `run()` prints `0 0`.
- **(b)** Builds the two distance maps *once* (not per candidate — they
  don't change between candidates on the same turn).
- **(c)** Scores `best_pos`.
- **(d)** Loops the remaining candidates, replacing `best_pos` whenever
  it finds a *strictly* higher score. Strict `>` means ties keep the
  earlier (row-major) candidate.

Let's open the two sub-functions.

---

### 4 — `bfs_distance_maps(&board)`

Source:

```rust
pub fn bfs_distance_maps(board: &Board) -> (Vec<Vec<u32>>, Vec<Vec<u32>>) {
    let my_dist    = bfs_from(board, |c| c == Cell::Mine);
    let enemy_dist = bfs_from(board, |c| c == Cell::Opponent);
    (my_dist, enemy_dist)
}
```

Two BFSes. Same function, different "source predicate".

`bfs_from`:

```rust
fn bfs_from(board: &Board, is_source: impl Fn(Cell) -> bool) -> Vec<Vec<u32>> {
    let mut dist  = vec![vec![u32::MAX; board.width]; board.height];
    let mut queue = VecDeque::new();

    // Phase 1: seed every source cell at distance 0
    for r in 0..board.height {
        for c in 0..board.width {
            if is_source(board.cells[r][c]) {
                dist[r][c] = 0;
                queue.push_back((r, c));
            }
        }
    }

    // Phase 2: expand outward
    while let Some((r, c)) = queue.pop_front() {
        let next = dist[r][c] + 1;
        for (dr, dc) in NEIGHBORS {     // (-1,0),(1,0),(0,-1),(0,1)
            let nr = r as isize + dr;
            let nc = c as isize + dc;
            if nr < 0 || nc < 0 { continue; }
            let nr = nr as usize;
            let nc = nc as usize;
            if nr >= board.height || nc >= board.width { continue; }
            if next < dist[nr][nc] {    // only update if we found a shorter path
                dist[nr][nc] = next;
                queue.push_back((nr, nc));
            }
        }
    }

    dist
}
```

**First call** — sources are `Cell::Mine`. Only `(1,1)`:

- Phase 1 seeds `dist[1][1] = 0`, queue = `[(1,1)]`.
- Phase 2 pops `(1,1)`, pushes its 4 neighbours with distance 1, then
  expands outward one ring at a time. After the whole BFS finishes:

```
my_dist:
 2  1  2  3  4
 1  0  1  2  3
 2  1  2  3  4
 3  2  3  4  5
 4  3  4  5  6
```

**Second call** — sources are `Cell::Opponent`. Only `(3,3)`:

```
enemy_dist:
 6  5  4  3  4
 5  4  3  2  3
 4  3  2  1  2
 3  2  1  0  1
 4  3  2  1  2
```

Both are just Manhattan distance to the single source on this board
(there are no obstacles), which is what BFS reduces to on an open grid.
`bfs_distance_maps` returns this tuple.

---

### 5 — `score_placement` for `pos = (1, 0)` (placement A)

Source:

```rust
pub fn score_placement(
    my_dist: &[Vec<u32>], enemy_dist: &[Vec<u32>],
    board: &Board, piece: &Piece, pos: &Position,
) -> (i32, i32) {
    // Part A — collect board coordinates of every filled piece cell,
    //          and sum enemy_dist over them (excluding the Mine anchor).
    let mut piece_coords: Vec<(i32, i32)> = Vec::with_capacity(piece.width * piece.height);
    let mut enemy_dist_sum: i32 = 0;

    for pr in 0..piece.height {
        for pc in 0..piece.width {
            if !piece.cells[pr][pc] { continue; }
            let br = pos.row + pr;
            let bc = pos.col + pc;
            piece_coords.push((br as i32, bc as i32));

            if board.cells[br][bc] == Cell::Mine { continue; }   // skip anchor
            let ed = enemy_dist[br][bc];
            if ed != u32::MAX {
                enemy_dist_sum += ed as i32;
            }
        }
    }

    // Part B — count the board cells that would flip to mine.
    let mut gain: i32 = 0;
    for r in 0..board.height {
        for c in 0..board.width {
            let ed = enemy_dist[r][c];
            if ed == u32::MAX    { continue; }   // opponent unreachable
            if my_dist[r][c] < ed { continue; }  // already strictly mine
            for &(pr, pc) in &piece_coords {
                let d = ((r as i32 - pr).abs() + (c as i32 - pc).abs()) as u32;
                if d < ed {
                    gain += 1;
                    break;          // one reaching piece cell is enough
                }
            }
        }
    }

    (gain, -enemy_dist_sum)
}
```

#### Part A — `piece_coords` and `enemy_dist_sum`

Walk the piece's filled cells, map each to the board, push the board
coord into `piece_coords`, and sum `enemy_dist` over the non-anchor ones.

For `pos = (1, 0)` the piece covers `(1,0)` and `(1,1)`:
- `(1,0)`: Empty → push `(1,0)`, `enemy_dist[1][0] = 5`, `sum = 5`.
- `(1,1)`: Mine → push `(1,1)`, skip the sum (it's the anchor).

Result: `piece_coords = [(1,0), (1,1)]`, `enemy_dist_sum = 5`.

#### Part B — `gain`

Iterate all 25 board cells. For each:

1. Skip if `enemy_dist[r][c] == u32::MAX` (can't happen here — opponent
   is reachable from every cell).
2. Skip if `my_dist[r][c] < enemy_dist[r][c]` (already strictly mine —
   nothing to gain).
3. Otherwise, try each piece coord and check if Manhattan distance to it
   is strictly less than `enemy_dist[r][c]`. If any succeeds, the cell
   flips: `gain += 1` and `break`.

Cells that pass the "not strictly mine" filter (condition 2): all the
non-M squares from the Voronoi picture — 16 cells total. Going through
them for placement A:

| cell  | enemy_dist | dist from (1,0) | dist from (1,1) | flips? |
|-------|-----------|-----------------|-----------------|--------|
| (0,3) | 3         | 4               | 3               | no     |
| (0,4) | 4         | 5               | 4               | no     |
| (1,3) | 2         | 3               | 2               | no     |
| (1,4) | 3         | 4               | 3               | no     |
| (2,2) | 2         | 3               | 2               | no     |
| (2,3) | 1         | 4               | 3               | no     |
| (2,4) | 2         | 5               | 4               | no     |
| (3,0) | 3         | **2**           | 3               | **✓**  |
| (3,1) | 2         | 3               | 2               | no     |
| (3,2) | 1         | 4               | 3               | no     |
| (3,4) | 1         | 6               | 5               | no     |
| (4,0) | 4         | **3**           | 4               | **✓**  |
| (4,1) | 3         | 4               | 3               | no     |
| (4,2) | 2         | 5               | 4               | no     |
| (4,3) | 1         | 6               | 5               | no     |
| (4,4) | 2         | 7               | 6               | no     |

**gain = 2.**

`score_placement` returns `(gain, -enemy_dist_sum) = (2, -5)`.

Back in `pick_best_placement`, that becomes `best_score`.

---

### 6 — `score_placement` for `pos = (1, 1)` (placement B)

Same function, same logic, different `pos`. The piece covers `(1,1)` and
`(1,2)`:
- `(1,1)`: Mine → push `(1,1)`, skip the sum.
- `(1,2)`: Empty → push `(1,2)`, `enemy_dist[1][2] = 3`, `sum = 3`.

`piece_coords = [(1,1), (1,2)]`, `enemy_dist_sum = 3`.

Gain table:

| cell  | enemy_dist | dist from (1,1) | dist from (1,2) | flips? |
|-------|-----------|-----------------|-----------------|--------|
| (0,3) | 3         | 3               | **2**           | **✓**  |
| (0,4) | 4         | 4               | **3**           | **✓**  |
| (1,3) | 2         | 2               | **1**           | **✓**  |
| (1,4) | 3         | 3               | **2**           | **✓**  |
| (2,2) | 2         | 2               | **1**           | **✓**  |
| (2,3) | 1         | 3               | 2               | no     |
| (2,4) | 2         | 4               | 3               | no     |
| (3,0) | 3         | 3               | 4               | no     |
| (3,1) | 2         | 2               | 3               | no     |
| (3,2) | 1         | 3               | 2               | no     |
| (3,4) | 1         | 5               | 4               | no     |
| (4,0) | 4         | 4               | 5               | no     |
| (4,1) | 3         | 3               | 4               | no     |
| (4,2) | 2         | 4               | 3               | no     |
| (4,3) | 1         | 5               | 4               | no     |
| (4,4) | 2         | 6               | 5               | no     |

**gain = 5.**

Returns `(5, -3)`.

---

### 7 — back inside `pick_best_placement`

```rust
for pos in &placements[1..] {
    let score = score_placement(...);
    if score > best_score {
        best_score = score;
        best_pos = *pos;
    }
}
```

Comparison: `(5, -3) > (2, -5)`? Tuple compare in Rust is lexicographic
— first element decides. `5 > 2` → true. So:
- `best_score = (5, -3)`
- `best_pos = Position{1, 1}`

No more candidates. `pick_best_placement` returns `Some(Position{1, 1})`.

---

### 8 — back in `run()`, print and flush

```rust
match chosen {
    Some(pos) => writeln!(writer, "{} {}", pos.col, pos.row).unwrap(),
    None      => writeln!(writer, "0 0").unwrap(),
}
writer.flush().unwrap();
```

`pos = Position{row: 1, col: 1}` → prints `1 1\n` (engine wants
`X Y` = `col row`). `flush()` makes sure it actually leaves our stdout
buffer before we go around the loop for the next turn.

The main loop then goes back to parsing the next board/piece the engine
sends us.

---

### Call tree summary

```
run()  [filler.rs]
 ├── parser::parse_board     (not shown)
 ├── parser::parse_piece     (not shown)
 ├── enumerate_placements(&board, &piece)        → Vec<Position>
 │    └── is_valid_placement(&board, &piece, pos)  ×25 calls
 ├── pick_best_placement(&board, &piece, &placements)  → Option<Position>
 │    ├── bfs_distance_maps(&board)                 → (my_dist, enemy_dist)
 │    │    ├── bfs_from(&board, is Mine)             → my_dist
 │    │    └── bfs_from(&board, is Opponent)         → enemy_dist
 │    └── score_placement(...)                        ×2 calls (one per candidate)
 └── writeln!  (col row)
```

That's the entire per-turn path the code takes, start to finish.

---

## Worked example — map00, seed 42, first turn

Real game. Map00 is 20×15, I'm player 1, starts at (2, 9):

```
    01234567890123456789
000 ....................
001 ....................
002 .........@..........
003 ....................
 …
012 .........$..........
013 ....................
014 ....................
```

### The piece I get

The engine prints `Piece 4 3` (4 wide, 3 tall) followed by:

```
OO..
OO..
O...
```

Where does the list `(0,0) (0,1) (1,0) (1,1) (2,0)` come from? It's just
the `(row, col)` index of every `O` inside that 4×3 grid:

```
            col 0   col 1   col 2   col 3
row 0:       O       O       .       .    → (0,0), (0,1) are O
row 1:       O       O       .       .    → (1,0), (1,1) are O
row 2:       O       .       .       .    → (2,0) is O
```

Shape-wise: a 2×2 square with one extra cell hanging off the bottom-left.

### Finding legal placements

A placement is a `pos = (row, col)` saying "put the piece's top-left
corner HERE". The only legal anchor point right now is my single Mine at
board `(2, 9)`, so one of the five `O`s of the piece has to sit on top of
it. Which one? We try each:

| Which piece `O` is the anchor? | `pos = (2 − pr, 9 − pc)` | Board squares covered           |
|--------------------------------|--------------------------|--------------------------------|
| (0,0)                          | (2, 9)                   | (2,9) (2,10) (3,9) (3,10) (4,9) |
| (0,1)                          | (2, 8)                   | (2,8) (2,9) (3,8) (3,9) (4,8)   |
| (1,0)                          | (1, 9)                   | (1,9) (1,10) (2,9) (2,10) (3,9) |
| (1,1)                          | (1, 8)                   | (1,8) (1,9) (2,8) (2,9) (3,8)   |
| (2,0)                          | (0, 9)                   | (0,9) (0,10) (1,9) (1,10) (2,9) |

All five are legal (empty board, enemy far away).

### Scoring

Today I only own (2,9), enemy only owns (12,9). So:
- `my_dist[r][c] = |r − 2| + |c − 9|`
- `enemy_dist[r][c] = |r − 12| + |c − 9|`

A cell is "not strictly mine" (flippable) when `my_dist >= enemy_dist`.
Solve: `|r−2| ≥ |r−12|` → `r ≥ 7`. So rows 7..14 = **8 rows × 20 cols =
160 flippable cells**.

Comparing the two most southerly placements:

**`pos = (2, 9)`** — piece goes straight down, deepest filled square at
board `(4, 9)`:
- Row 7, column `c`: from (4,9), my reach is `3+|c-9|`; enemy reach is
  `5+|c-9|`. I win by 2 everywhere → **all 20 cells in row 7 flip to me**.
- Row 8: my reach `4+|c-9|`, enemy reach `4+|c-9|`. Tie — no flip.
- **Gain = 20.**

**`pos = (2, 8)`** — piece shifted one column left, deepest square at
`(4, 8)`:
- Row 7: my reach `3+|c-8|` vs enemy `5+|c-9|`. I win for every `c` → 20.
- Row 8: my reach `4+|c-8|` vs enemy `4+|c-9|`. I win whenever
  `|c-8| < |c-9|`, i.e. `c = 0..8` → 9 cells.
- **Gain = 29.**

`pos = (2, 8)` wins. The bot prints the answer as `X Y` = `col row`:

```
-> Answer (@): 8 2
```

The intuition in one sentence: shifting the piece's tail one column to the
left puts it closer to the middle of the board, which lets it beat the
enemy to an extra 9 cells on row 8 that the straight-down placement could
only tie on.

---

## Worked example continued — turn 2

After my first move and the opponent's reply, the board looks like:

```
002 ........aa..........
003 ........aa..........
004 ........a...........
 …
012 .........ss.........
```

My Mines: **(2,8), (2,9), (3,8), (3,9), (4,8)**.
Enemy: **(12,9), (12,10)**.

New piece, 3×3:

```
OOO
OOO
.OO
```

That's 8 filled cells in a 9-cell grid — very dense.

### Why "one legal origin per Mine cell" is WRONG now

On turn 1 I had a single Mine, so each filled piece cell gave exactly one
candidate origin and they never collided. Now my 5 Mines sit packed
together — most positions that put a piece cell on one of my Mines put
**another** piece cell on one of my other Mines at the same time. That
violates the "exactly one overlap" rule.

Going through them one by one:

| Anchor cell | # of legal origins | Why |
|---|---|---|
| (3,8) | **0** | Boxed in: Mines at (2,8), (3,9), (4,8) are its neighbours. With a piece this dense, *every* orientation hits at least one of them — rule 3 fails. |
| (2,8) | 1 | Only `pos = (0,6)` works (piece dangles up-and-left). |
| (2,9) | 2 | `pos = (1,9)` and `pos = (0,8)`. |
| (3,9) | 1 | `pos = (3,9)` (piece goes right and down into empty board). |
| (4,8) | 3 | `pos = (4,6), (4,7), (4,8)`. Lots of room because nothing south of row 4 is occupied yet. |

**Total: 7 legal placements.** You were right: (3,8) literally cannot
anchor this piece.

### Picking the winner

The 3 placements anchored on (4,8) all reach deepest south (down to row 6
or 7), which is where the contested frontier is. Among those three,
`pos = (4, 7)` has the best spread — its filled squares straddle columns
7–9, covering the most "still leaning enemy" cells.

Covered squares: `(4,7) (4,8) (4,9) (5,7) (5,8) (5,9) (6,8) (6,9)`.
Anchor: `(4,8)`. Everything else lands on Empty.

```
-> Answer (@): 7 4
```

---

## Worked example continued — turn 3

After the opponent's second reply:

```
 …
006 ........@@..........
 …
011 ..........ss........
012 .........$ss........
```

My 13 Mines span rows 2–6. Enemy has 6 cells in rows 11–12. The contested
frontier has shifted down toward row 8–9.

New piece, 1×3:
```
O
O
.
```

Just a vertical 2-cell stripe (the third row is blank).

With such a tiny piece, there are lots of legal anchors — essentially any
of my Mines on the southern edge can host it. The scorer picks the one
that pushes the front line furthest south:

Chosen move: `pos = (6, 9)`. The anchor is my existing Mine at (6,9); the
fresh piece cell lands at (7,9), which sits right at the old midline and
re-opens the southern front.

```
-> Answer (@): 9 6
```

---

## Properties of this bot

- **Deterministic.** Same board + same piece → same move. Debug-friendly.
- **Single-ply.** It never simulates the opponent's reply — that'd need an
  opponent model. The Voronoi gain metric already captures most of what a
  deeper search would find, so it's a deliberate simplicity/quality trade.
- **Open-grid assumption.** The "Manhattan distance == BFS distance"
  shortcut inside `score_placement` is only exact because Filler maps
  have no walls. If walls were added, the inner loop would over-estimate
  our reach; the fix would be a real BFS seeded from the piece cells.
- **Forfeit.** If no legal placement exists, we print `0 0` and keep
  looping. The engine stops calling us; the opponent may keep scoring.

---

## If you only remember one sentence

> For every legal placement, count how many empty squares would flip from
> "closer to the enemy" to "closer to me" if I played it. Pick the one
> that flips the most.

Everything else is either plumbing or a tiebreaker.
