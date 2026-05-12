use ::std::collections::HashMap;

pub fn is_anagram(s1: &str, s2: &str) -> bool {
    if s1.len() != s2.len() {
        return false;
    }

    let mut counts: HashMap<char, i32> = HashMap::new();

    let s1 = s1.to_lowercase();
    let s2 = s2.to_lowercase();

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
        let s1 = "listen";
        let s2 = "silent";

        assert_eq!(true, is_anagram(s1, s2));
    }
}
