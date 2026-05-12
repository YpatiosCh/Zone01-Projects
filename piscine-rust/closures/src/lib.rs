pub fn first_fifty_even_square() -> Vec<i32> {
    (1..)
        .filter(|n| n % 2 == 0)
        .take(50)
        .map(|n| n * n)
        .collect()
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        let even = first_fifty_even_square();
        println!("{:?}", even);
    }
}
