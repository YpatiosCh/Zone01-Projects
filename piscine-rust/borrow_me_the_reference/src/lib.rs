pub fn delete_and_backspace(s: &mut String) {
    let mut result = Vec::new();
    let mut skip = 0;

    for c in s.chars().rev() {
        match c {
            '-' => skip += 1,
            '+' => {
                result.pop();
            }
            _ => {
                if skip > 0 {
                    skip -= 1;
                } else {
                    result.push(c);
                }
            }
        }
    }

    *s = result.iter().rev().collect();
}

pub fn do_operations(v: &mut [String]) {
    for s in v.iter_mut() {
        if let Some(i) = s.find('+') {
            let left: i32 = s[..i].parse().unwrap();
            let right: i32 = s[i + 1..].parse().unwrap();
            *s = (left + right).to_string();
        } else if let Some(i) = s.find('-') {
            let left: i32 = s[..i].parse().unwrap();
            let right: i32 = s[i + 1..].parse().unwrap();
            *s = (left - right).to_string();
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        let mut a = "bpp--o+er+++sskroi-++lcw".to_owned();
        delete_and_backspace(&mut a);
        let a_result = "borrow".to_string();

        let mut b = [
            "2+2".to_owned(),
            "3+2".to_owned(),
            "10-3".to_owned(),
            "5+5".to_owned(),
        ];
        do_operations(&mut b);
        let b_result = [
            "4".to_string(),
            "5".to_string(),
            "7".to_string(),
            "10".to_string(),
        ];

        assert_eq!(a, a_result);
        assert_eq!(b, b_result);
    }
}
