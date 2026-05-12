use std::cell::RefCell;
use std::rc::Rc;

pub struct Tracker {
    pub messages: RefCell<Vec<String>>,
    value: RefCell<usize>,
    max: usize,
}

impl Tracker {
    pub fn new(max: usize) -> Self {
        Tracker {
            messages: Default::default(),
            value: Default::default(),
            max,
        }
    }

    pub fn set_value<T>(&self, value: &Rc<T>) {
        let count = Rc::strong_count(value);

        if count > self.max {
            self.messages
                .borrow_mut()
                .push("Error: You can't go over your quota!".to_string());
            return;
        }

        *self.value.borrow_mut() = count;
        let percentage = count * 100 / self.max;

        if percentage >= 70 {
            self.messages.borrow_mut().push(format!(
                "Warning: You have used up over {percentage}% of your quota!"
            ));
        }
    }

    pub fn peek<T>(&self, value: &Rc<T>) {
        let count = Rc::strong_count(value);
        let percentage = count * 100 / self.max;
        self.messages.borrow_mut().push(format!(
            "Info: This value would use {}% of your quota",
            percentage
        ));
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        let tracker = Tracker::new(1);

        let hello = "Helloooo".to_string();

        let one = Rc::new(1);
        let two = Rc::new(&hello);
        let three = Rc::new(&hello);

        tracker.set_value(&one);
        tracker.set_value(&two);
        tracker.set_value(&three);

        tracker
            .messages
            .borrow()
            .iter()
            .for_each(|msg| println!("{msg}"));
    }
}
