/// Character pairs used by the Filler game engine to identify each player's
/// cells on the board.
///
/// Each player owns two characters: an uppercase "anchor" placed on the turn a
/// piece is laid down (e.g. `@` or `$`), and a lowercase form (e.g. `a` or `s`)
/// used for every subsequent cell of that player. Keeping both forms lets the
/// parser recognise any cell belonging to a given player regardless of when it
/// was placed.
#[derive(Debug, Clone)]
pub struct PlayerChars {
    pub my_chars: (char, char),
    pub ops_chars: (char, char),
}

/// The game board ("Anfield"), a rectangular grid of [`Cell`]s.
///
/// `cells` is indexed as `cells[row][col]`, with `row` in `0..height` and
/// `col` in `0..width`.
#[derive(Debug, Clone)]
pub struct Board {
    pub width: usize,
    pub height: usize,
    pub cells: Vec<Vec<Cell>>,
}

/// The possible states of a single cell on the board, normalised relative to
/// the current player (so downstream logic never has to care about the raw
/// engine characters).
#[derive(Debug, Clone, Copy, PartialEq)]
pub enum Cell {
    Empty,
    Mine,
    Opponent,
}

/// A tetromino-like piece handed to the player each turn.
///
/// `cells[row][col] == true` means the piece occupies that offset; `false`
/// means the cell is transparent (does not constrain placement).
#[derive(Debug, Clone)]
pub struct Piece {
    pub width: usize,
    pub height: usize,
    pub cells: Vec<Vec<bool>>,
}

/// A 2D coordinate on the board, expressed as `(row, col)`.
///
/// Used to describe where the top-left corner of a piece should be placed.
#[derive(Debug, Clone, Copy, PartialEq)]
pub struct Position {
    pub row: usize,
    pub col: usize,
}
