pub fn bubble_sort(arr: &mut [i32]) {
    arr.sort()
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        let mut v = [3, 2, 4, 5, 1, 7];
        bubble_sort(&mut v);
        let sorted = [1, 2, 3, 4, 5, 7];

        assert_eq!(sorted, v);
    }
}
