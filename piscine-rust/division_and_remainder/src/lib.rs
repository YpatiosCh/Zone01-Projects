pub fn divide(x: i32, y: i32) -> (i32, i32) {
    let division = x / y;
    let remainder = x % y;
    (division, remainder)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        let result: (i32, i32) = divide(9, 4);
        assert_eq!(result, (2, 1));
    }
}
