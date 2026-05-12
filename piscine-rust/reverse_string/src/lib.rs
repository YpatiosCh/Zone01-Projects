pub fn rev_str(input: &str) -> String {
    input.chars().rev().collect()
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        let result: String = rev_str("Hello, World!");
        assert_eq!(result, "!dlroW ,olleH".to_string());
    }
}
