pub fn first_subword(mut s: String) -> String {
    let mut chars = s.char_indices();
    chars.next(); // skip first character

    let end = chars
        .find(|(_, c)| c.is_uppercase() || *c == '_') //  returns Option<(usize, char)> (index, character) |(_, char)| for each char find the one that meets criteria c.is_uppercase or c == "_"
        .map(|(i, _)| i) // then from the return on the finf method, which is (index, char) => "map" return only the index
        .unwrap_or(s.len()); // if there is an index return it otherwise (if any of above returned None) return length of s

    s.truncate(end);
    s
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        let s1 = "helloWorld";
        let s2 = "snake_case";
        let s3 = "CamelCase";
        let s4 = "just";

        assert_eq!(first_subword(s1.to_owned()), "hello".to_string());
        assert_eq!(first_subword(s2.to_owned()), "snake".to_string());
        assert_eq!(first_subword(s3.to_owned()), "Camel".to_string());
        assert_eq!(first_subword(s4.to_owned()), "just".to_string());
    }
}
