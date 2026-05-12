#[derive(Debug, PartialEq, Clone)]
pub struct Boss {
    pub name: String,
    pub age: u32,
}

impl Boss {
    pub fn new(name: &str, age: u32) -> Boss {
        Boss {
            name: name.to_owned(),
            age,
        }
    }
}
