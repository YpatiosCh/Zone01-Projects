pub fn number_logic(num: u32) -> bool {
    let str = num.to_string();
    let len = str.len();

    num == str
        .chars()
        .map(|c| c.to_digit(10).unwrap_or(0).pow(len as u32))
        .sum()
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        assert_eq!(number_logic(9), true);

        assert_eq!(number_logic(10), false);
    }
}
