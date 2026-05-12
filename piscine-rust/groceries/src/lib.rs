pub fn insert(vec: &mut Vec<String>, val: String) {
    vec.push(val);
}

pub fn at_index(slice: &[String], index: usize) -> &str {
    &slice[index]
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        let mut vec: Vec<String> = vec!["hello".to_string()];
        let last: usize = vec.len();

        insert(&mut vec, "world".to_string());

        assert_eq!(at_index(&vec, last), "world");
    }
}
