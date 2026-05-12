pub fn diagonals(player: char, table: [[char; 3]; 3]) -> bool {
    table[0][0] == player && table[1][1] == player && table[2][2] == player
        || table[0][2] == player && table[1][1] == player && table[2][0] == player
}

pub fn horizontal(player: char, table: [[char; 3]; 3]) -> bool {
    table.iter().any(|row| row.iter().all(|&c| c == player))
}

pub fn vertical(player: char, table: [[char; 3]; 3]) -> bool {
    (0..3).any(|col| (0..3).all(|row| table[row][col] == player))
}

pub fn tic_tac_toe(table: [[char; 3]; 3]) -> String {
    for player in ['X', 'O'] {
        if horizontal(player, table) || vertical(player, table) || diagonals(player, table) {
            return format!("player {} won", player);
        }
    }
    "tie".to_string()
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        assert_eq!(
            tic_tac_toe([['O', 'X', 'O'], ['O', 'P', 'X'], ['X', '#', 'X']]),
            "tie".to_string()
        );

        assert_eq!(
            tic_tac_toe([['X', 'O', 'O'], ['X', 'O', 'O'], ['#', 'O', 'X']]),
            "player O won".to_string()
        );

        assert_eq!(
            tic_tac_toe([['O', 'O', 'X'], ['O', 'X', 'O'], ['X', '#', 'X']]),
            "player X won".to_string()
        );
    }
}
