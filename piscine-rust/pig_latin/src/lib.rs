pub fn pig_latin(text: &str) -> String {
    if starts_with_consonant_and_qu(text) {
        let prefix: String = text.chars().take(3).collect();
        let rest: String = text.chars().skip(3).collect();

        format!("{}{}ay", rest, prefix)
    } else if starts_with_consonant(text) {
        let prefix: String = text
            .chars()
            .map(|c| c.to_ascii_lowercase())
            .take_while(|c| !"aeiou".contains(*c))
            .collect();

        let rest: String = text.chars().skip(prefix.len()).collect();

        format!("{}{}ay", rest, prefix)
    } else {
        format!("{text}ay")
    }
}

/// finds if text starts with consontant followed by "qu" -> "squ", "mqu" etc.
fn starts_with_consonant_and_qu(text: &str) -> bool {
    let lower = text.to_ascii_lowercase();
    let mut chars = lower.chars();

    match (chars.next(), chars.next(), chars.next()) {
        (Some(c1), Some('q'), Some('u')) => !"aeiou".contains(c1),
        _ => false,
    }
}

/// guess what it does
fn starts_with_consonant(text: &str) -> bool {
    text.chars()
        .next()
        .map(|c| c.to_ascii_lowercase())
        .map(|c| !"aeiou".contains(c))
        .unwrap_or(false)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        assert_eq!(pig_latin("igloo"), "iglooay");
        assert_eq!(pig_latin("apple"), "appleay");
        assert_eq!(pig_latin("hello"), "ellohay");
        assert_eq!(pig_latin("square"), "aresquay");
        assert_eq!(pig_latin("xenon"), "enonxay");
        assert_eq!(pig_latin("chair"), "airchay");
        assert_eq!(pig_latin("queen"), "ueenqay");
    }
}
