use std::collections::HashMap;

pub fn is_permutation(s1: &str, s2: &str) -> bool {
    if s1.len() != s2.len() {
        return false;
    }

    let mut counts: HashMap<char, i32> = HashMap::new();

    for c in s1.chars() {
        *counts.entry(c).or_insert(0) += 1;
    }

    for c in s2.chars() {
        *counts.entry(c).or_insert(0) -= 1;
    }

    counts.values().all(|&v| v == 0)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        let word = "thought";
        let word1 = "thougth";

        assert_eq!(true, is_permutation(word1, word));
    }
}
