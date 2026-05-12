use std::io::Write;
use std::iter;

use crate::parser;
use crate::reader::Reader;
use crate::strategy;

/// Run the bot's main I/O loop.
///
/// The lifecycle is:
///
/// 1. Read the one-shot player header (`$$$ exec p1 : [...]` or `p2 : [...]`)
///    from stdin and figure out which pair of characters represents us.
/// 2. Loop forever:
///    - parse the next `Anfield` block into a [`crate::models::Board`],
///    - parse the next `Piece` block into a [`crate::models::Piece`],
///    - enumerate every legal placement for that piece,
///    - pick the highest-scoring one,
///    - print `col row\n` (or `0 0\n` if there is no legal move) and flush.
///
/// The loop terminates implicitly when stdin closes: the engine kills our
/// process once the game is over, so there's no explicit exit condition.
///
/// # Panics
/// The function treats every parse failure as fatal and panics with a short
/// message. This is deliberate — the engine protocol is rigid, so any
/// malformed input is a logic bug we'd rather surface loudly than silently
/// skip a turn on.
pub fn run() {
    let mut reader = Reader::new(std::io::stdin().lock());
    let mut writer = std::io::stdout().lock();

    let player_line = match reader.read_line() {
        Some(line) => line,
        None => panic!("Could not read player line"),
    };
    let players = match parser::parse_player(&player_line) {
        Ok(players) => players,
        Err(e) => panic!("Could not parse player line: {}", e),
    };

    loop {
        let mut lines = iter::from_fn(|| reader.read_line());

        let board = match parser::parse_board(&mut lines, &players) {
            Ok(board) => board,
            Err(e) => panic!("Could not read board: {}", e),
        };

        let piece = match parser::parse_piece(&mut lines) {
            Ok(piece) => piece,
            Err(e) => panic!("Could not read piece: {}", e),
        };

        let placements = strategy::enumerate_placements(&board, &piece);
        let chosen = strategy::pick_best_placement(&board, &piece, &placements);

        match chosen {
            Some(pos) => writeln!(writer, "{} {}", pos.col, pos.row).unwrap(),
            None => writeln!(writer, "0 0").unwrap(),
        }
        writer.flush().unwrap();
    }
}
