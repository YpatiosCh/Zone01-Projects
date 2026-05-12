pub fn delete_prefix<'s>(prefix: &str, s: &'s str) -> Option<&'s str> {
    s.strip_prefix(prefix)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        assert_eq!(
            delete_prefix("ab", "abcdefghijklmnop"),
            Some("cdefghijklmnop")
        );
        assert_eq!(delete_prefix("x", "abcdefghijklmnop"), None);
    }
}
