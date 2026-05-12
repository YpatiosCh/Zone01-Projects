pub fn is_empty(v: &str) -> bool {
    v.len() == 0
}

pub fn is_ascii(v: &str) -> bool {
    v.is_ascii()
}

pub fn contains(v: &str, pat: &str) -> bool {
    v.contains(pat)
}

pub fn split_at(v: &str, index: usize) -> (&str, &str) {
    (&v[..index], &v[index..])
}

pub fn find(v: &str, pat: char) -> usize {
    v.chars().position(|c| c == pat).unwrap_or(0)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        // IS EMPTY ?
        assert_eq!(false, is_empty("Hello"));
        assert_eq!(true, is_empty(""));

        // IS ASCII ?
        assert_eq!(true, is_ascii("rust"));
        assert_eq!(false, is_ascii("γεια κοσμε"));

        // CONTAINS ?
        assert_eq!(true, contains("hello", "el"));
        assert_eq!(false, contains("hello", "world"));

        // SPLIT AT
        assert_eq!(("ru", "st"), split_at("rust", 2));

        // FIND
        assert_eq!(1, find("rust", 'u'));
        assert_eq!(0, find("rust", 'k'));
    }
}
