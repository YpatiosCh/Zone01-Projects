pub fn sum(a: &[i32]) -> i32 {
    a.iter().sum()
}

pub fn thirtytwo_tens() -> [i32; 32] {
    [10; 32]
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        let thirty_two_tens = thirtytwo_tens();
        assert_eq!(32, thirty_two_tens.len());

        assert_eq!(sum(&[1, 2, 3, 4, 5, 6, 7, 8, 9, 10]), 55);
    }
}
