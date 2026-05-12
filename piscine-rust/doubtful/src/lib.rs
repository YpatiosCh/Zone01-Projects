pub fn doubtful(s: &mut String) {
    s.push('?');
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        let mut s: String = "Hello".to_owned();
        let answer: String = "Hello?".to_string();
        doubtful(&mut s);
        assert_eq!(s, answer);
    }
}
