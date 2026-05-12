pub fn add_curry(n: i32) -> impl Fn(i32) -> i32 {
    move |y| n + y
}

pub fn twice<F>(f: F) -> impl Fn(i32) -> i32
where
    F: Fn(i32) -> i32,
{
    move |x| f(f(x))
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        let twenty = twice(add_curry(10));

        let twenty_five = twenty(5);

        assert_eq!(25, twenty_five);
    }
}
