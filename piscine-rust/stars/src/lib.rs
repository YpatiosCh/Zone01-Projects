pub fn stars(n: u32) -> String {
    let how_many = 2_usize.pow(n);

    "*".repeat(how_many as usize)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        assert_eq!(stars(1), "**");
        assert_eq!(stars(4), "****************");
    }
}
