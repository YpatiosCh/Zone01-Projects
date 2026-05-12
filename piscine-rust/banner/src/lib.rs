use std::{collections::HashMap, num::ParseFloatError};

#[derive(Debug)]
pub struct Flag {
    pub short_hand: String,
    pub long_hand: String,
    pub desc: String,
}

impl Flag {
    pub fn opt_flag(name: &str, desc: &str) -> Self {
        Self {
            short_hand: format!("-{}", name.chars().next().unwrap()),
            long_hand: format!("--{}", name.to_string()),
            desc: desc.to_string(),
        }
    }
}

pub type Callback = fn(&str, &str) -> Result<String, ParseFloatError>;

#[derive(Debug)]
pub struct FlagsHandler {
    pub flags: HashMap<String, Callback>,
}

impl FlagsHandler {
    pub fn add_flag(&mut self, flag: Flag, func: Callback) {
        self.flags.insert(flag.short_hand, func);
        self.flags.insert(flag.long_hand, func);
    }

    pub fn exec_func(&self, input: &str, argv: &[&str]) -> Result<String, String> {
        let callback = self.flags.get(input);
        match callback {
            Some(func) => func(argv[0], argv[1]).map_err(|e| e.to_string()),
            None => Err(format!("no flag {} found", input)),
        }
    }
}

pub fn div(a: &str, b: &str) -> Result<String, ParseFloatError> {
    let a_fl: f64 = a.parse()?;
    let b_fl: f64 = b.parse()?;
    Ok((a_fl / b_fl).to_string())
}

pub fn rem(a: &str, b: &str) -> Result<String, ParseFloatError> {
    let a_fl: f64 = a.parse()?;
    let b_fl: f64 = b.parse()?;
    Ok((a_fl % b_fl).to_string())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        let mut handler = FlagsHandler {
            flags: HashMap::new(),
        };

        let d = Flag::opt_flag("division", "divides the values, formula (a / b)");
        let r = Flag::opt_flag(
            "remainder",
            "remainder of the division between two values, formula (a % b)",
        );

        println!("FLAG 1: {:?}", d);

        handler.add_flag(d, div);
        handler.add_flag(r, rem);

        println!("Handler: {:?}", handler);

        assert_eq!(
            handler.exec_func("-d", &["1.0", "2.0"]),
            Ok("0.5".to_string())
        );

        assert_eq!(
            handler.exec_func("-r", &["2.0", "2.0"]),
            Ok("0".to_string())
        );

        assert_eq!(
            handler.exec_func("--division", &["a", "2.0"]),
            Err("invalid float literal".to_string())
        );

        assert_eq!(
            handler.exec_func("--remainder", &["2.0", "fd"]),
            Err("invalid float literal".to_string())
        );
    }
}
