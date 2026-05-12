pub fn get_products(arr: Vec<usize>) -> Vec<usize> {
    if arr.len() <= 1 {
        return vec![];
    }
    arr.iter()
        .enumerate()
        .map(|(i, _)| {
            arr.iter()
                .enumerate()
                .filter(|(j, _)| *j != i)
                .map(|(_, v)| v)
                .product()
        })
        .collect()
}

#[cfg(test)]
mod tests {
    use std::vec;

    use super::*;

    #[test]
    fn it_works() {
        let arr: Vec<usize> = vec![1, 7, 3, 4];
        let output = get_products(arr);

        assert_eq!(output, [84, 12, 28, 21]);

        let single: Vec<usize> = vec![2];
        let single_output: Vec<usize> = get_products(single);
        let empty: Vec<usize> = vec![];

        assert_eq!(empty, single_output);
    }
}
