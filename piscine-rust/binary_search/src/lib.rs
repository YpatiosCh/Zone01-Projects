pub fn binary_search(sorted_list: &[i32], target: i32) -> Option<usize> {
    sorted_list.binary_search(&target).ok()
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        let sorted_list = vec![1, 3, 5, 7, 9, 11, 13];
        let target1 = 7;
        let target2 = 18;

        assert_eq!(binary_search(&sorted_list, target1), Some(3));
        assert_eq!(binary_search(&sorted_list, target2), None);
    }
}
