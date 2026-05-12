#[derive(Debug)]
pub struct Numbers<'a> {
    numbers: &'a [u32],
}

impl<'a> Numbers<'a> {
    pub fn new(numbers: &'a [u32]) -> Self {
        Numbers { numbers }
    }

    pub fn list(&self) -> &[u32] {
        self.numbers
    }

    pub fn latest(&self) -> Option<u32> {
        if self.numbers.len() == 0 {
            return None;
        }

        Some(self.numbers[self.numbers.len() - 1])
    }

    pub fn highest(&self) -> Option<u32> {
        self.numbers.iter().max().map(|n| *n)
    }

    pub fn highest_three(&self) -> Vec<u32> {
        let mut res = self.numbers.to_vec();
        res.sort_by(|a, b| b.cmp(a));
        res.into_iter().take(3).collect()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        let result = [30, 500, 20, 70];
        let n = Numbers::new(&result);

        assert_eq!(n.list(), [30, 500, 20, 70]);

        assert_eq!(n.highest(), Some(500));
        assert_eq!(n.latest(), Some(70));

        assert_eq!(n.highest_three(), [500, 70, 30]);
    }
}
