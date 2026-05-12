#[macro_export]
macro_rules! hash_map {
    () => {
        std::collections::HashMap::new()
    };

    ($($key: expr => $val: expr),+ $(,)?) => {
        {
            let mut map = std::collections::HashMap::new();
            $(map.insert($key, $val);)+
            map
        }
    };
}

#[cfg(test)]
mod tests {
    use std::collections::HashMap;

    #[test]
    fn empty_map() {
        let m: HashMap<u32, u32> = hash_map![];
        assert!(m.is_empty());
    }

    #[test]
    fn single_entry() {
        let m = hash_map!["key" => 42];
        assert_eq!(m.len(), 1);
        assert_eq!(m["key"], 42);
    }

    #[test]
    fn multiple_entries() {
        let m = hash_map!['a' => 1, 'b' => 2, 'c' => 3];
        assert_eq!(m.len(), 3);
        assert_eq!(m[&'a'], 1);
        assert_eq!(m[&'b'], 2);
        assert_eq!(m[&'c'], 3);
    }
}
