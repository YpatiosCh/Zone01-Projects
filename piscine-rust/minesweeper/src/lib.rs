pub fn solve_board(minefield: &[&str]) -> Vec<String> {
    if minefield.is_empty() {
        return vec![];
    }

    let rows = minefield.len();
    let cols = minefield[0].len();

    let is_mine = |r: usize, c: usize| minefield[r].as_bytes()[c] == b'*';

    let count_adjacent = |r: usize, c: usize| -> u8 {
        let mut count = 0;
        for dr in -1i32..=1 {
            for dc in -1i32..=1 {
                if dr == 0 && dc == 0 {
                    continue;
                }
                let nr = r as i32 + dr;
                let nc = c as i32 + dc;
                if nr >= 0 && nr < rows as i32 && nc >= 0 && nc < cols as i32 {
                    if is_mine(nr as usize, nc as usize) {
                        count += 1;
                    }
                }
            }
        }
        count
    };

    (0..rows)
        .map(|r| {
            (0..cols)
                .map(|c| {
                    if is_mine(r, c) {
                        '*'
                    } else {
                        let n = count_adjacent(r, c);
                        if n == 0 {
                            ' '
                        } else {
                            (b'0' + n) as char
                        }
                    }
                })
                .collect()
        })
        .collect()
}
