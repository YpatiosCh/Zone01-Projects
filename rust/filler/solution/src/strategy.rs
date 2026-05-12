use std::collections::VecDeque;

use crate::models::{Board, Cell, Piece, Position};

/// Four-way orthogonal offsets used by every BFS in this module.
const NEIGHBORS: [(isize, isize); 4] = [(-1, 0), (1, 0), (0, -1), (0, 1)];

/// Run a 4-connected breadth-first expansion from every cell the predicate
/// accepts, producing a distance map in which each reachable cell holds the
/// length of the shortest path back to the nearest source. Unreachable cells
/// retain `u32::MAX`.
///
/// The BFS treats all cell types as passable; the only distinction between
/// them is whether they act as a source (distance 0) or not.
fn bfs_from(board: &Board, is_source: impl Fn(Cell) -> bool) -> Vec<Vec<u32>> {
    let mut dist = vec![vec![u32::MAX; board.width]; board.height];
    let mut queue: VecDeque<(usize, usize)> = VecDeque::new();

    for r in 0..board.height {
        for c in 0..board.width {
            if is_source(board.cells[r][c]) {
                dist[r][c] = 0;
                queue.push_back((r, c));
            }
        }
    }

    while let Some((r, c)) = queue.pop_front() {
        let next = dist[r][c] + 1;
        for (dr, dc) in NEIGHBORS {
            let nr = r as isize + dr;
            let nc = c as isize + dc;
            if nr < 0 || nc < 0 {
                continue;
            }
            let nr = nr as usize;
            let nc = nc as usize;
            if nr >= board.height || nc >= board.width {
                continue;
            }
            if next < dist[nr][nc] {
                dist[nr][nc] = next;
                queue.push_back((nr, nc));
            }
        }
    }

    dist
}

/// Compute the pair of distance maps that underpins the Voronoi territory
/// heuristic: `(my_dist, enemy_dist)` where each entry is the BFS distance
/// from the nearest Mine cell and the nearest Opponent cell respectively.
///
/// Comparing the two maps cell-by-cell yields a cheap proxy for ownership:
/// wherever `my_dist < enemy_dist` the cell is likely ours to claim, and the
/// signed gap `enemy_dist - my_dist` quantifies how strongly it leans our way.
pub fn bfs_distance_maps(board: &Board) -> (Vec<Vec<u32>>, Vec<Vec<u32>>) {
    let my_dist = bfs_from(board, |c| c == Cell::Mine);
    let enemy_dist = bfs_from(board, |c| c == Cell::Opponent);
    (my_dist, enemy_dist)
}

/// Score a candidate placement as `(voronoi_gain, -enemy_dist_sum)` so that
/// `max` picks the best move lexicographically.
///
/// * `voronoi_gain` — the number of board cells that would flip from the
///   opponent's Voronoi region (or contested) into ours if this placement
///   went ahead. A cell `(r,c)` counts when it is currently not-strictly-ours
///   (`my_dist[r][c] >= enemy_dist[r][c]`) AND some filled piece cell sits
///   at a Manhattan distance strictly less than `enemy_dist[r][c]`, so the
///   new Mine source reaches it before the enemy does. This is the exact
///   territorial delta of the move, not a proxy.
/// * `enemy_dist_sum` — sum of `enemy_dist` over the piece's non-anchor
///   cells. Negated in the returned tuple so "minimise this" becomes
///   "maximise" lexicographically. Breaks ties toward aggressive moves that
///   push the piece into enemy-leaning terrain (the original "rush-the-enemy"
///   heuristic), which matters a lot in the endgame once `voronoi_gain` is
///   frequently zero.
///
/// Board cells where `enemy_dist` is `u32::MAX` (opponent unreachable) are
/// skipped so they neither inflate `voronoi_gain` nor overflow the sum.
/// BFS distances on an obstacle-free rectangular grid equal Manhattan
/// distance, which is why the inner loop uses the cheap closed form instead
/// of a second BFS per candidate.
pub fn score_placement(
    my_dist: &[Vec<u32>],
    enemy_dist: &[Vec<u32>],
    board: &Board,
    piece: &Piece,
    pos: &Position,
) -> (i32, i32) {
    let mut piece_coords: Vec<(i32, i32)> = Vec::with_capacity(piece.width * piece.height);
    let mut enemy_dist_sum: i32 = 0;

    for pr in 0..piece.height {
        for pc in 0..piece.width {
            if !piece.cells[pr][pc] {
                continue;
            }
            let br = pos.row + pr;
            let bc = pos.col + pc;
            piece_coords.push((br as i32, bc as i32));

            if board.cells[br][bc] == Cell::Mine {
                continue;
            }
            let ed = enemy_dist[br][bc];
            if ed != u32::MAX {
                enemy_dist_sum += ed as i32;
            }
        }
    }

    let mut gain: i32 = 0;
    for r in 0..board.height {
        for c in 0..board.width {
            let ed = enemy_dist[r][c];
            if ed == u32::MAX {
                continue;
            }
            if my_dist[r][c] < ed {
                continue;
            }
            for &(pr, pc) in &piece_coords {
                let d = ((r as i32 - pr).abs() + (c as i32 - pc).abs()) as u32;
                if d < ed {
                    gain += 1;
                    break;
                }
            }
        }
    }

    (gain, -enemy_dist_sum)
}

/// Check whether placing `piece` with its top-left corner at `pos` satisfies
/// the Filler placement rules: every filled piece cell must fall inside the
/// board, no filled piece cell may land on an Opponent cell, and exactly one
/// filled piece cell must land on one of our Mine cells (the "anchor"
/// overlap). Returns `true` only when all three conditions hold.
fn is_valid_placement(board: &Board, piece: &Piece, pos: Position) -> bool {
    let mut overlap = 0u32;
    for pr in 0..piece.height {
        for pc in 0..piece.width {
            if !piece.cells[pr][pc] {
                continue;
            }
            let br = pos.row + pr;
            let bc = pos.col + pc;
            if br >= board.height || bc >= board.width {
                return false;
            }
            match board.cells[br][bc] {
                Cell::Opponent => return false,
                Cell::Mine => overlap += 1,
                Cell::Empty => {}
            }
        }
    }
    overlap == 1
}

/// Enumerate every legal origin for the current piece on the current board by
/// trying every `(row, col)` in row-major order and keeping the ones that
/// pass the Filler placement rules. O(W·H·|piece_cells|) — simple and more
/// than fast enough for the map sizes this bot plays on.
pub fn enumerate_placements(board: &Board, piece: &Piece) -> Vec<Position> {
    let mut out = Vec::new();
    for row in 0..board.height {
        for col in 0..board.width {
            let pos = Position { row, col };
            if is_valid_placement(board, piece, pos) {
                out.push(pos);
            }
        }
    }
    out
}

/// Select the highest-scoring placement from a list of legal candidates.
///
/// Each candidate is scored via [`score_placement`] against a single pair of
/// BFS distance maps (computed once up-front, not per candidate). The scan
/// uses a strict `>` comparison so ties fall to whichever candidate appears
/// first in `placements` — callers that want row-major tie-breaking should
/// pass a row-major-ordered slice (which is what [`enumerate_placements`]
/// produces).
///
/// Returns `None` only when `placements` is empty, which the caller should
/// interpret as "no legal move this turn" (forfeit).
pub fn pick_best_placement(
    board: &Board,
    piece: &Piece,
    placements: &[Position],
) -> Option<Position> {
    let mut best_pos = *placements.first()?;

    let (my_dist, enemy_dist) = bfs_distance_maps(board);

    let mut best_score = score_placement(&my_dist, &enemy_dist, board, piece, &best_pos);

    for pos in &placements[1..] {
        let score = score_placement(&my_dist, &enemy_dist, board, piece, pos);
        if score > best_score {
            best_score = score;
            best_pos = *pos;
        }
    }

    Some(best_pos)
}

#[cfg(test)]
mod tests {
    use super::*;

    fn board_from_rows(rows: &[&str]) -> Board {
        let height = rows.len();
        let width = rows[0].chars().count();
        let cells: Vec<Vec<Cell>> = rows
            .iter()
            .map(|r| {
                r.chars()
                    .map(|c| match c {
                        '@' | 'a' => Cell::Mine,
                        '$' | 's' => Cell::Opponent,
                        _ => Cell::Empty,
                    })
                    .collect()
            })
            .collect();
        Board {
            width,
            height,
            cells,
        }
    }

    fn piece_from_rows(rows: &[&str]) -> Piece {
        let height = rows.len();
        let width = rows[0].chars().count();
        let cells: Vec<Vec<bool>> = rows
            .iter()
            .map(|r| r.chars().map(|c| c == 'O').collect())
            .collect();
        Piece {
            width,
            height,
            cells,
        }
    }

    #[test]
    fn bfs_distance_maps_basic() {
        let board = board_from_rows(&[
            "@....", //
            ".....", //
            ".....", //
            ".....", //
            "....$", //
        ]);
        let (my_dist, enemy_dist) = bfs_distance_maps(&board);

        assert_eq!(my_dist[0][0], 0);
        assert_eq!(my_dist[4][4], 8);
        assert_eq!(enemy_dist[4][4], 0);
        assert_eq!(enemy_dist[0][0], 8);
        assert_eq!(my_dist[2][2], 4);
        assert_eq!(enemy_dist[2][2], 4);
    }

    #[test]
    fn bfs_distance_maps_multiple_sources() {
        let board = board_from_rows(&[
            "@...@", //
            ".....", //
            ".....", //
        ]);
        let (my_dist, _) = bfs_distance_maps(&board);
        assert_eq!(my_dist[0][2], 2);
        assert_eq!(my_dist[2][2], 4);
    }

    #[test]
    fn score_placement_counts_voronoi_flips() {
        // Single row; Mine at column 0, Opponent at column 7.
        // Current Voronoi: enemy owns columns 4..=7 (columns 4,5,6 are
        // strictly closer to the opponent; column 7 is the opponent itself).
        // A 3-cell horizontal piece anchored at column 0 extends us to
        // column 2. Column 4 has enemy_dist=3; Manhattan distance from the
        // piece's tip (0,2) to (0,4) is 2 < 3, so that cell flips to ours.
        let board = board_from_rows(&["@......$"]);
        let piece = piece_from_rows(&["OOO"]);
        let (my_dist, enemy_dist) = bfs_distance_maps(&board);
        let (gain, _) = score_placement(
            &my_dist,
            &enemy_dist,
            &board,
            &piece,
            &Position { row: 0, col: 0 },
        );
        assert_eq!(gain, 1, "expected exactly (0,4) to flip to our side");
    }

    #[test]
    fn score_placement_prefers_invasion_over_retreat() {
        // Mine at (0,2), Opponent at (2,5). Two valid placements with a
        // horizontal 2-cell piece: one extends left (away from enemy), the
        // other extends right (toward enemy).
        let board = board_from_rows(&[
            "..@.....", //
            "........", //
            ".....$..", //
        ]);
        let piece = piece_from_rows(&["OO"]);
        let (my_dist, enemy_dist) = bfs_distance_maps(&board);

        let retreat = Position { row: 0, col: 1 };
        let invade = Position { row: 0, col: 2 };

        let s_retreat = score_placement(&my_dist, &enemy_dist, &board, &piece, &retreat);
        let s_invade = score_placement(&my_dist, &enemy_dist, &board, &piece, &invade);

        assert!(
            s_invade > s_retreat,
            "expected invasion {:?} to outscore retreat {:?}",
            s_invade,
            s_retreat,
        );
    }

    #[test]
    fn enumerate_placements_single_cell_piece() {
        let board = board_from_rows(&[
            "@..", //
            "...", //
            "..$", //
        ]);
        let piece = piece_from_rows(&["O"]);
        let placements = enumerate_placements(&board, &piece);
        assert_eq!(placements, vec![Position { row: 0, col: 0 }]);
    }

    #[test]
    fn enumerate_placements_respects_anchor_rule() {
        let board = board_from_rows(&[
            "@...", //
            "....", //
            "...$", //
        ]);
        let piece = piece_from_rows(&[
            "OO", //
            ".O", //
        ]);
        let placements = enumerate_placements(&board, &piece);
        assert!(placements.contains(&Position { row: 0, col: 0 }));
        assert!(!placements.contains(&Position { row: 1, col: 2 }));
        for p in &placements {
            assert!(is_valid_placement(&board, &piece, *p));
        }
    }

    #[test]
    fn enumerate_placements_rejects_opponent_overlap() {
        let board = board_from_rows(&[
            "@$.", //
            "...", //
        ]);
        let piece = piece_from_rows(&["OO"]);
        let placements = enumerate_placements(&board, &piece);
        assert!(placements.iter().all(|p| p.row != 0 || p.col != 0));
    }

    #[test]
    fn pick_best_placement_empty_list_returns_none() {
        let board = board_from_rows(&["@.$"]);
        let piece = piece_from_rows(&["O"]);
        assert!(pick_best_placement(&board, &piece, &[]).is_none());
    }

    #[test]
    fn pick_best_placement_prefers_invasion() {
        // Two valid placements from the same Mine anchor; the one extending
        // toward the opponent should win on Voronoi gain or the aggressive
        // tiebreaker.
        let board = board_from_rows(&[
            "..@.....", //
            "........", //
            ".....$..", //
        ]);
        let piece = piece_from_rows(&["OO"]);
        let candidates = [
            Position { row: 0, col: 1 }, // extends left, away
            Position { row: 0, col: 2 }, // extends right, toward enemy
        ];
        let best = pick_best_placement(&board, &piece, &candidates).unwrap();
        assert_eq!(best, Position { row: 0, col: 2 });
    }

    #[test]
    fn is_valid_placement_rejects_multi_mine_overlap() {
        // Two Mine cells side by side; a 2-wide piece would overlap both.
        let board = board_from_rows(&["@@..", "...."]);
        let piece = piece_from_rows(&["OO"]);
        // pos (0,0) covers both Mines → overlap = 2, must be rejected.
        assert!(!is_valid_placement(
            &board,
            &piece,
            Position { row: 0, col: 0 }
        ));
    }

    #[test]
    fn is_valid_placement_rejects_off_board() {
        // 3x3 board with Mine on the right edge.
        let board = board_from_rows(&["..@", "...", "..."]);
        let piece = piece_from_rows(&["OO"]);
        // pos (0,2) anchors on the Mine but the second cell is at col=3, off-board.
        assert!(!is_valid_placement(
            &board,
            &piece,
            Position { row: 0, col: 2 }
        ));
        // Same idea vertically.
        let tall = piece_from_rows(&["O", "O", "O"]);
        assert!(!is_valid_placement(
            &board,
            &tall,
            Position { row: 2, col: 2 }
        ));
    }
}
