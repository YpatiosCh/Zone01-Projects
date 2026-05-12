pub fn add_curry(n: i32) -> impl Fn(i32) -> i32 {
    move |y| n + y
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        let add10 = add_curry(10);
        let added = add10(5);

        assert_eq!(15, added);
    }
}
