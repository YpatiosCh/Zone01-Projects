pub fn parse_into_boxed(s: String) -> Vec<Box<u32>> {
    s.split_whitespace()
        .map(|str| {
            let number = if str.ends_with('k') {
                let num: f64 = str[..str.len() - 1].parse().unwrap();
                (num * 1000.0) as u32
            } else {
                str.parse().unwrap()
            };
            Box::new(number)
        })
        .collect()
}

pub fn into_unboxed(a: Vec<Box<u32>>) -> Vec<u32> {
    a.into_iter().map(|b| *b).collect()
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::mem;

    #[test]
    fn it_works() {
        let s = "5.5k 8.9k 32".to_owned();

        let boxed = parse_into_boxed(s);
        assert_eq!(5500 as u32, *boxed[0]);
        println!("Element size: {:?} bytes", mem::size_of_val(&boxed[0]));

        let unboxed = into_unboxed(boxed);
        assert_eq!(5500, unboxed[0]);
        println!("Element size: {:?} bytes", mem::size_of_val(&unboxed[0]));
    }
}
