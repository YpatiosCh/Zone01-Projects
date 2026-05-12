pub fn arrange_phrase(phrase: &str) -> String {
    let mut words: Vec<&str> = phrase.split_whitespace().collect();
    words.sort_by_key(|w| w.chars().find(|c| c.is_numeric()).unwrap());
    words
        .iter()
        .map(|w| w.chars().filter(|c| !c.is_numeric()).collect::<String>())
        .collect::<Vec<String>>()
        .join(" ")
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        assert_eq!(arrange_phrase("is2 Thi1s T4est 3a"), "This is a Test");
    }
}
