pub fn char_length(s: &str) -> usize {
    s.chars().count()
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        assert_eq!(1, char_length("❤"));
        assert_eq!(3, char_length("形聲字"));
        assert_eq!(6, char_length("change"));
        assert_eq!(1, char_length("😍"));
    }
}
