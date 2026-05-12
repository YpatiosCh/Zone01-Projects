pub fn initials(names: Vec<&str>) -> Vec<String> {
    let mut v: Vec<String> = Vec::with_capacity(names.len());

    for name in names {
        let split: Vec<&str> = name.split(" ").collect();
        let first: String = split[0][0..1].to_uppercase();
        let second: String = split[1][0..1].to_uppercase();
        v.push(format!("{}. {}.", first, second));
    }
    v
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        let names = vec!["Harry Potter", "Someone Else", "J. L.", "Barack Obama"];
        let correct: Vec<String> = vec![
            "H. P.".to_string(),
            "S. E.".to_string(),
            "J. L.".to_string(),
            "B. O.".to_string(),
        ];
        assert_eq!(initials(names), correct);
    }
}
