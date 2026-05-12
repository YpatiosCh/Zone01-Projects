#[derive(Debug, PartialEq)]
pub struct CipherError {
    pub expected: String,
}

pub fn cipher(original: &str, ciphered: &str) -> Result<(), CipherError> {
    let expected: String = original.bytes().map(|c| transform(c)).collect();

    if expected == ciphered {
        Ok(())
    } else {
        Err(CipherError { expected: expected })
    }
}

fn transform(c: u8) -> char {
    if !c.is_ascii_alphabetic() {
        return c as char;
    }

    // upper
    // 65-90
    if c >= 65 as u8 && c <= 90 as u8 {
        let mirror = 65 + 90 - c;
        return mirror as char;
    }

    // lower
    // 97-122
    if c >= 97 as u8 && c <= 122 as u8 {
        let mirror = 97 + 122 - c;
        return mirror as char;
    }

    c as char
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        assert_eq!(cipher("1Hello 2world!", "1Svool 2dliow!"), Ok(()));
        assert_eq!(
            cipher("1Hello 2world!", "svool"),
            Err(CipherError {
                expected: "1Svool 2dliow!".to_string()
            })
        );
    }
}
