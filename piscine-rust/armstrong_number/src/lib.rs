pub fn is_armstrong_number(nb: u32) -> Option<u32> {
    let digits: Vec<u32> = nb
        .to_string()
        .chars()
        .map(|c| c.to_digit(10).unwrap())
        .collect();
    let power = digits.len() as u32;
    let sum: u32 = digits.iter().map(|d| d.pow(power)).sum();
    if sum == nb { Some(nb) } else { None }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        assert_eq!(Some(0), is_armstrong_number(0));
        assert_eq!(Some(1), is_armstrong_number(1));
        assert_eq!(Some(153), is_armstrong_number(153));
        assert_eq!(Some(370), is_armstrong_number(370));
        assert_eq!(Some(371), is_armstrong_number(371));
        assert_eq!(Some(407), is_armstrong_number(407));
        assert_eq!(None, is_armstrong_number(400));
        assert_eq!(None, is_armstrong_number(198));
    }
}
