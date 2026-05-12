pub fn num_to_ordinal(x: u32) -> String {
    format!(
        "{}{}",
        x,
        if matches!(x % 100, 11..=13) {
            "th"
        } else {
            match x % 10 {
                1 => "st",
                2 => "nd",
                3 => "rd",
                _ => "th",
            }
        }
    )
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        assert_eq!(num_to_ordinal(1), "1st");
        assert_eq!(num_to_ordinal(11), "11th");
    }
}
