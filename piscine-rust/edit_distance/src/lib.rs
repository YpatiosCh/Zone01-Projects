pub fn edit_distance(source: &str, target: &str) -> usize {
    let s: Vec<char> = source.chars().collect();
    let t: Vec<char> = target.chars().collect();
    let m = s.len();
    let n = t.len();

    // dp[i][j] = edit distance between s[0..i] and t[0..j]
    let mut dp = vec![vec![0usize; n + 1]; m + 1];

    // base cases: transforming to/from empty string
    for i in 0..=m {
        dp[i][0] = i;
    } // delete all chars from source
    for j in 0..=n {
        dp[0][j] = j;
    } // insert all chars from target

    for i in 1..=m {
        for j in 1..=n {
            dp[i][j] = if s[i - 1] == t[j - 1] {
                // chars match, no operation needed
                dp[i - 1][j - 1]
            } else {
                1 + dp[i - 1][j - 1] // substitution
                    .min(dp[i - 1][j]) // deletion
                    .min(dp[i][j - 1]) // insertion
            };
        }
    }

    dp[m][n]
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        assert_eq!(2, edit_distance("alignment", "assignment"));
    }
}
