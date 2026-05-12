pub fn to_url(s: &str) -> String {
    s.replace(' ', "%20")
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        let result = to_url("Hello, world!");
        let correct = "Hello,%20world!".to_string();
        assert_eq!(result, correct);
    }
}
