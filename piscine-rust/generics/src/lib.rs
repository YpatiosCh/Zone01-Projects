pub fn identity<T>(v: T) -> T {
    v
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        assert_eq!("Hello", identity("Hello"));
        assert_eq!(4, identity(4));
    }
}
