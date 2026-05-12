#[derive(Debug, Clone, Copy, PartialEq)]
pub struct Person<'name> {
    pub name: &'name str,
    pub age: u32,
}

impl<'name> Person<'name> {
    pub fn new(name: &'name str) -> Self {
        Person { name, age: 0 }
    }
}

#[cfg(test)]
mod tests {

    use super::*;

    #[test]
    fn it_works() {
        let leo = Person::new("Leo");
        assert_eq!(leo.name, "Leo");
    }

    #[test]
    fn copy() {
        let p1 = Person::new("Leo");
        let p2 = p1;

        assert_eq!(p1, p2);
    }

    #[test]
    fn name_is_borrowed_not_owned() {
        let name = String::from("Leo");
        let person = Person::new(&name);

        // can't mutate through Person
        assert_eq!(person.name, "Leo");
    }
}
