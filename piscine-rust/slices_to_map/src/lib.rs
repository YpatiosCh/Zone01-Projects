use std::collections::HashMap;

pub fn slices_to_map<'a, T: Eq + std::hash::Hash, U>(
    keys: &'a [T],
    values: &'a [U],
) -> HashMap<&'a T, &'a U> {
    keys.iter().zip(values.iter()).collect()
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        let keys = ["Olivia", "Liam", "Emma", "Noah", "James"];
        let values = [1, 3, 23, 5, 2];
        println!("{:?}", slices_to_map(&keys, &values));
    }
}
