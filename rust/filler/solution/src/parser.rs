use crate::models::*;

/// Parse the first line emitted by the game engine and work out which player
/// we are.
///
/// The engine prints something like `$$$ exec p1 : [./bot]` or
/// `$$$ exec p2 : [./bot]`. Only the `p1`/`p2` token matters: player 1 always
/// owns `('@', 'a')` and player 2 always owns `('$', 's')`. The returned
/// [`PlayerChars`] is oriented from *our* point of view, so `my_chars` is
/// whichever pair belongs to us and `ops_chars` is the other pair.
///
/// # Errors
/// Returns `Err` if the line does not contain `exec p`, or if the character
/// following `exec p` is neither `'1'` nor `'2'`.
pub fn parse_player(line: &str) -> Result<PlayerChars, String> {
    let start = "$$$ exec p";
    let idx = line
        .find(start)
        .ok_or("malformed player line: missing '$$$ exec p'")?;
    let player_char = line.as_bytes()[idx + start.len()];
    let player_number = match player_char {
        b'1' => 1,
        b'2' => 2,
        _ => return Err(format!("malformed player line: expected p1 or p2")),
    };

    if player_number == 1 {
        Ok(PlayerChars {
            my_chars: ('@', 'a'),
            ops_chars: ('$', 's'),
        })
    } else {
        Ok(PlayerChars {
            my_chars: ('$', 's'),
            ops_chars: ('@', 'a'),
        })
    }
}

/// Parse one `Anfield` block from an iterator of input lines into a [`Board`].
///
/// The expected block layout (as emitted by the Filler engine) is:
/// ```text
/// Anfield <width> <height>:
///     0123456789...            <- column ruler, ignored
/// 000 ..........                <- `<row-number> <row-of-cells>`
/// 001 ....@.....
/// ...
/// ```
///
/// The raw board characters are normalised against `players` so the returned
/// [`Board::cells`] only contains `Mine` / `Opponent` / `Empty` variants.
///
/// # Errors
/// Returns `Err` if the header, column ruler, or any row line is missing or
/// malformed — in particular if a row's width does not match the header.
pub fn parse_board<I: Iterator<Item = String>>(
    lines: &mut I,
    players: &PlayerChars,
) -> Result<Board, String> {
    // get the board header
    let header = lines.next().ok_or("malformed board: missing header")?;
    let (width, height) = parse_dimensions(&header, "Anfield")?;

    // skip the columns line
    lines
        .next()
        .ok_or("malformed board: missing columns line")?;

    // create a vector to store the cells
    let mut cells = Vec::with_capacity(height);

    for _ in 0..height {
        // get the row line
        let row_line = lines.next().ok_or("malformed board: missing row")?;
        let row_line = row_line.trim_end();

        // skip the first column (the row number)
        let (_, grid_str) = row_line.split_once(' ').ok_or("malformed board row")?;

        // convert the row string to a vector of cells
        let row: Vec<Cell> = grid_str.chars().map(|c| char_to_cell(c, players)).collect();

        // check if the row has the correct width
        if row.len() != width {
            return Err(format!(
                "malformed board: expected width {}, got {}",
                width,
                row.len()
            ));
        }

        cells.push(row);
    }

    Ok(Board {
        width,
        height,
        cells,
    })
}

/// Parse one `Piece` block from an iterator of input lines into a [`Piece`].
///
/// The expected block layout is:
/// ```text
/// Piece <width> <height>:
/// ..O..
/// .OOO.
/// ```
/// where `O` marks an occupied cell and `.` marks a transparent cell. The
/// returned [`Piece::cells`] is a boolean mask of the same shape, with `true`
/// at every `O`.
///
/// # Errors
/// Returns `Err` if the header is missing/malformed, if a row is missing, if
/// a row's width does not match the header, or if a row contains any
/// character other than `.` or `O`.
pub fn parse_piece<I: Iterator<Item = String>>(lines: &mut I) -> Result<Piece, String> {
    // get the piece header
    let header = lines.next().ok_or("malformed piece: missing header")?;
    let (width, height) = parse_dimensions(&header, "Piece")?;

    // create a vector to store the cells
    let mut cells: Vec<Vec<bool>> = Vec::with_capacity(height);

    for _ in 0..height {
        // get the row line
        let row_line = lines.next().ok_or("malformed piece: missing row")?;
        let row_line = row_line.trim_end();

        // convert the row string to a vector of booleans
        let mut row = Vec::with_capacity(width);

        // iterate over the characters in the row line
        for c in row_line.chars() {
            match c {
                '.' => row.push(false),
                'O' => row.push(true),
                some_char => {
                    return Err(format!("malformed piece: invalid character: {}", some_char));
                }
            }
        }

        // check if the row has the correct width
        if row.len() != width {
            return Err(format!(
                "malformed piece: expected width {}, got {}",
                width,
                row.len()
            ));
        }

        cells.push(row);
    }

    Ok(Piece {
        width,
        height,
        cells,
    })
}

/// Parse a header line of the form `<tag> <width> <height>:` — shared by both
/// the board (`tag = "Anfield"`) and the piece (`tag = "Piece"`) sections.
///
/// The trailing colon on the height token is tolerated and stripped before
/// parsing. Returns the `(width, height)` pair on success.
///
/// # Errors
/// Returns `Err` when the line does not start with `tag`, does not have at
/// least three whitespace-separated tokens, or contains non-numeric
/// width/height values.
fn parse_dimensions(line: &str, tag: &str) -> Result<(usize, usize), String> {
    let parts: Vec<&str> = line.split_whitespace().collect();
    if parts.len() < 3 || parts[0] != tag {
        return Err(format!("malformed board: missing {}", tag));
    }
    let width: usize = parts[1]
        .parse()
        .map_err(|_| "malformed board: invalid width".to_string())?;
    let height: usize = parts[2]
        .trim_end_matches(':')
        .parse()
        .map_err(|_| "malformed board: invalid height".to_string())?;
    Ok((width, height))
}

/// Map a single raw board character to the player-relative [`Cell`] variant.
///
/// Both the uppercase and lowercase forms of each player's characters are
/// matched, so a newly-placed anchor cell and an older occupied cell resolve
/// to the same variant. Anything that isn't one of the four player characters
/// is treated as empty territory.
fn char_to_cell(c: char, players: &PlayerChars) -> Cell {
    match c {
        c if c == players.ops_chars.0 || c == players.ops_chars.1 => Cell::Opponent,
        c if c == players.my_chars.0 || c == players.my_chars.1 => Cell::Mine,
        _ => Cell::Empty,
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::io::{BufRead, Cursor};

    /// Turn a raw multi-line string into an iterator of owned lines suitable
    /// for feeding into the parsing functions under test.
    fn make_lines(input: &str) -> impl Iterator<Item = String> {
        Cursor::new(input).lines().map(|l| l.unwrap())
    }

    #[test]
    fn test_parse_player() {
        let player_line = "$$$ exec p1";
        let players = parse_player(player_line).expect("parse failed");
        assert_eq!(players.my_chars, ('@', 'a'));
        assert_eq!(players.ops_chars, ('$', 's'));
    }

    #[test]
    fn test_parse_player_2() {
        let player_line = "$$$ exec p2";
        let players = parse_player(player_line).expect("parse failed");
        assert_eq!(players.my_chars, ('$', 's'));
        assert_eq!(players.ops_chars, ('@', 'a'));
    }

    #[test]
    #[should_panic]
    fn test_parse_p3_panic() {
        let player_line = "@@@ exec p3";
        parse_player(player_line).unwrap();
    }

    #[test]
    fn parse_board_basic() {
        let input = "\
Anfield 3 4:
    012
000 ...
001 .@.
002 .$.
003 ...
";

        let players = PlayerChars {
            my_chars: ('@', 'a'),
            ops_chars: ('$', 's'),
        };

        let mut lines = input.lines().map(|s| s.to_string());
        let board = parse_board(&mut lines, &players).expect("parse failed");

        assert_eq!(board.width, 3);
        assert_eq!(board.height, 4);
        assert_eq!(board.cells[0][0], Cell::Empty);
        assert_eq!(board.cells[0][1], Cell::Empty);
        assert_eq!(board.cells[0][2], Cell::Empty);
        assert_eq!(board.cells[1][0], Cell::Empty);
        assert_eq!(board.cells[1][1], Cell::Mine);
        assert_eq!(board.cells[1][2], Cell::Empty);
        assert_eq!(board.cells[2][0], Cell::Empty);
        assert_eq!(board.cells[2][1], Cell::Opponent);
        assert_eq!(board.cells[2][2], Cell::Empty);
        assert_eq!(board.cells[3][0], Cell::Empty);
        assert_eq!(board.cells[3][1], Cell::Empty);
        assert_eq!(board.cells[3][2], Cell::Empty);
    }

    #[test]
    fn parse_board_rejects_wrong_width() {
        let input = "\
Anfield 3 2:
    012
000 ....
001 .@.
";

        let mut lines = make_lines(input);

        let players = PlayerChars {
            my_chars: ('@', 'a'),
            ops_chars: ('$', 's'),
        };

        let result = parse_board(&mut lines, &players);

        assert!(result.is_err());
    }

    #[test]
    fn parse_piece_basic() {
        let input = "\
Piece 3 2:
.O.
OOO
";

        let mut lines = make_lines(input);

        let piece = parse_piece(&mut lines).unwrap();

        assert_eq!(piece.width, 3);
        assert_eq!(piece.height, 2);

        assert_eq!(piece.cells[0], vec![false, true, false]);
        assert_eq!(piece.cells[1], vec![true, true, true]);
    }

    #[test]
    fn parse_piece_rejects_wrong_width() {
        let input = "\
Piece 3 1:
.OOO
";

        let mut lines = make_lines(input);

        let result = parse_piece(&mut lines);

        assert!(result.is_err());
    }

    #[test]
    fn parse_piece_missing_header() {
        let input = "";

        let mut lines = make_lines(input);

        let result = parse_piece(&mut lines);

        assert!(result.is_err());
    }

    #[test]
    fn parse_piece_rejects_invalid_characters() {
        let input = "\
Piece 2 1:
.X
";

        let mut lines = make_lines(input);

        let result = parse_piece(&mut lines);

        assert!(result.is_err());
    }
}
