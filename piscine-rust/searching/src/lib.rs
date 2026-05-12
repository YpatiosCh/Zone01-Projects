pub fn search(array: &[i32], key: i32) -> Option<usize> {
    array.iter().rposition(|&n| n == key)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        let ar = [1, 3, 4, 6, 8, 9, 11, 8];
        let f = search(&ar, 8);
        assert_eq!(f, Some(7));
    }
}
