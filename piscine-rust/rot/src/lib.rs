pub fn rotate(input: &str, key: i8) -> String {
    input
        .chars()
        .map(|c| {
            if c.is_ascii_alphabetic() {
                let base = if c.is_uppercase() { b'A' } else { b'a' };
                let shifted = ((c as i16 - base as i16 + key as i16).rem_euclid(26)) as u8;
                (base + shifted) as char
            } else {
                c
            }
        })
        .collect()
}
