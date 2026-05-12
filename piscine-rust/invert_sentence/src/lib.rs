pub fn invert_sentence(string: &str) -> String {
    let mut segments: Vec<(&str, &str)> = Vec::new();
    let mut rest = string;

    while !rest.is_empty() {
        let word_start = rest
            .find(|c: char| !c.is_whitespace())
            .unwrap_or(rest.len());
        let spaces = &rest[..word_start];
        rest = &rest[word_start..];

        let word_end = rest.find(|c: char| c.is_whitespace()).unwrap_or(rest.len());
        let word = &rest[..word_end];
        rest = &rest[word_end..];

        segments.push((spaces, word));
    }

    // trailing spaces end up as a ("  ", "") segment — keep them at the end
    let trailing = if segments.last().map(|(_, w)| w.is_empty()).unwrap_or(false) {
        segments.pop().map(|(s, _)| s).unwrap_or("")
    } else {
        ""
    };

    let mut words: Vec<&str> = segments.iter().map(|(_, w)| *w).rev().collect();

    let mut result: String = segments
        .iter()
        .map(|(spaces, _)| format!("{}{}", spaces, words.remove(0)))
        .collect();

    result.push_str(trailing);
    result
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_simple() {
        assert_eq!(invert_sentence("Rust is Awesome"), "Awesome is Rust");
    }

    #[test]
    fn test_preserve_whitespace() {
        assert_eq!(
            invert_sentence("    word1     word2  "),
            "    word2     word1  "
        );
    }

    #[test]
    fn test_single_word() {
        assert_eq!(invert_sentence("Hello, World!"), "World! Hello,");
    }
}
