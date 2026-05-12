use std::collections::HashMap;

pub fn mean(list: &[i32]) -> f64 {
    list.iter().sum::<i32>() as f64 / list.len() as f64
}

pub fn median(list: &[i32]) -> i32 {
    let mut sorted = list.to_vec();
    sorted.sort();
    let mid = sorted.len() / 2;
    if sorted.len() % 2 == 0 {
        (sorted[mid - 1] + sorted[mid]) / 2
    } else {
        sorted[mid]
    }
}

pub fn mode(list: &[i32]) -> i32 {
    let mut freq: HashMap<i32, usize> = HashMap::new();
    for &n in list {
        *freq.entry(n).or_insert(0) += 1;
    }
    freq.into_iter().max_by_key(|&(_, count)| count).unwrap().0
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        let v = [4, 7, 5, 2, 5, 1, 3];

        assert_eq!(mean(&v), 3.857142857142857);
        assert_eq!(median(&v), 4);
        assert_eq!(mode(&v), 5);
    }
}
