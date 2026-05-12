use ::std::collections::HashMap;

pub fn bigger(h: HashMap<&str, i32>) -> i32 {
    h.into_values().max().unwrap_or(0)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        let hash = HashMap::from_iter([
            ("Daniel", 122),
            ("Ashley", 333),
            ("Katie", 334),
            ("Robert", 14),
        ]);

        assert_eq!(bigger(hash), 334);
    }
}
