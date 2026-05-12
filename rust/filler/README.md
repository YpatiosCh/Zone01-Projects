# Filler — Rust bot

A territorial-game bot written in Rust for the 01-edu Filler engine. It
reads a board and a piece from stdin every turn and replies with the
coordinates of the best place to drop the piece.

The bot beats every reference opponent bundled with the engine (`bender`,
`wall_e`, `h2_d2`, `terminator`) on every shipped map.

---

## The game in one paragraph

Filler is a two-player land-grab on a rectangular grid. Each player owns a
starting cell. On every turn the engine hands you a randomly shaped piece
and you must place it so that **exactly one** of its filled cells overlaps
your existing territory (the "anchor"), no filled cell overlaps your
opponent, and nothing leaves the board. You score one point per filled
piece cell placed. When neither player can place any more, whoever covers
the bigger area wins. Full spec: see the
[01-edu subject](https://github.com/01-edu/public/tree/master/subjects/filler).

---

## How the bot plays

**One-ply Voronoi territory heuristic.** Every turn:

1. Enumerate every legal placement of the current piece.
2. Run two multi-source BFSes over the board — one from our cells, one
   from the opponent's — to produce `my_dist` and `enemy_dist` maps.
3. Score each placement by counting how many board cells would **flip to
   our Voronoi region** if we played it (i.e. how many cells we'd reach
   faster than the opponent after the move). Ties broken by "piece lands
   closer to the enemy".
4. Pick the highest-scoring placement.

No search tree, no opponent model, no learned weights — the whole
decision comes out of comparing two distance maps. A detailed walkthrough
with worked examples and a function-by-function trace lives in
[`algo.md`](./algo.md).

---

## Project layout

```
filler-rust/
├── solution/                 ← the bot (this is what you build)
│   ├── Cargo.toml
│   └── src/
│       ├── main.rs            ← entry point, calls filler::run
│       ├── lib.rs             ← module declarations
│       ├── filler.rs          ← main game loop (I/O + orchestration)
│       ├── reader.rs          ← line-buffered stdin wrapper
│       ├── parser.rs          ← parses player header, board, piece
│       ├── models.rs          ← Board / Piece / Position / Cell / PlayerChars
│       └── strategy.rs        ← legality + BFS + scoring + picking
├── game/                      ← engine & reference bots (opaque binaries)
│   ├── linux_game_engine
│   ├── linux_robots/          ← bender, wall_e, h2_d2, terminator
│   └── maps/                  ← map00, map01, map02
├── Makefile                   ← convenience targets
├── Dockerfile                 ← reproducible run environment
└── algo.md                    ← deep technical write-up of the strategy
```

### Module responsibilities

| File           | Responsibility                                                      |
|----------------|---------------------------------------------------------------------|
| `filler.rs`    | Main per-turn loop: read, parse, decide, write                      |
| `reader.rs`    | `BufRead` wrapper yielding one trimmed line at a time               |
| `parser.rs`    | Converts engine text blocks into typed `Board` / `Piece` values     |
| `models.rs`    | Typed representation: `Board`, `Cell`, `Piece`, `Position`, `PlayerChars` |
| `strategy.rs`  | Placement enumeration, BFS distance maps, scoring, selection        |

Board characters are **normalised** during parsing: `@`/`a` and `$`/`s`
become `Cell::Mine` or `Cell::Opponent` depending on which player we are,
so downstream code never cares about raw engine chars.

---

## Building and running

### Local (Linux)

```bash
make build                   # cargo build --release
make run                     # runs vs. bender on map00
make test                    # cargo test
```

### Pick a different opponent or map

```bash
./game/linux_game_engine \
    -f ./game/maps/map01 \
    -p1 ./target/release/solution \
    -p2 ./game/linux_robots/terminator
```

Engine flags:

| flag          | meaning                                  |
|---------------|------------------------------------------|
| `-f <path>`   | map file                                 |
| `-p1 <path>`  | player 1 binary                          |
| `-p2 <path>`  | player 2 binary                          |
| `-s <int>`    | random seed (reproducible games)         |
| `-t <sec>`    | timeout per move (default 10)            |
| `-q`          | quiet mode                               |
| `-r`          | throttled / animated mode                |

### Docker

The project ships a `Dockerfile` that bundles the engine, maps, and
reference robots on top of `rust:1.86-slim-bookworm`:

```bash
make build-docker            # docker build -t filler .
make run-docker              # mounts ./solution and drops into a shell
```

Inside the container:

```bash
cd /filler/solution
cargo build --release
cd ..
./linux_game_engine -f maps/map00 -p1 solution/target/release/solution -p2 linux_robots/bender
```

---

## Testing

Unit tests cover the three subsystems that matter for correctness:

- **Reader** — line splitting, trailing-newline handling, empty input.
- **Parser** — player header (`p1` / `p2` / malformed), board dimensions
  and cell normalisation, piece dimensions and `.` / `O` charset.
- **Strategy** — BFS distance maps on small grids, the Voronoi-gain
  scoring, the anchor-overlap legality rule, opponent-overlap rejection,
  and end-to-end picker behaviour including the invade-vs-retreat
  tiebreaker.

Run them with:

```bash
make test
```

---

For the actual algorithmic details — how scoring works, a full worked
example on map00, and a function-by-function trace of one turn — read
[`algo.md`](./algo.md).
