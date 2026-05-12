use std::ops::Add;

pub struct StepIterator<T> {
    pub beg: T,
    pub end: T,
    pub step: T,
}

impl<T: Add> StepIterator<T> {
    pub fn new(beg: T, end: T, step: T) -> Self {
        StepIterator { beg, end, step }
    }
}

impl<T: Add<Output = T> + Clone + PartialOrd> std::iter::Iterator for StepIterator<T> {
    type Item = T;

    fn next(&mut self) -> Option<Self::Item> {
        if self.beg > self.end {
            return None;
        }
        let current = self.beg.clone();
        self.beg = self.beg.clone() + self.step.clone();
        Some(current)
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        let mut res = Vec::new();
        for v in StepIterator::new(0, 100, 10) {
            res.push(v);
        }

        let mut res1 = Vec::new();
        for v in StepIterator::new(0, 100, 12) {
            res1.push(v);
        }

        let mut res2 = Vec::new();
        for v in StepIterator::new(0.0, 100.0, 15.0) {
            res2.push(v);
        }

        println!("{:?}", res);
        println!("{:?}", res1);
        println!("{:?}", res2);
    }
}
