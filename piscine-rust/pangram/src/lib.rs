pub fn is_pangram(s: &str) -> bool {
    ('a'..='z').all(|c| s.to_lowercase().contains(c))
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        assert_eq!(
            is_pangram("the quick brown fox jumps over the lazy dog!"),
            true
        );

        assert_eq!(is_pangram("this is not a pangram!"), false);
    }
}
