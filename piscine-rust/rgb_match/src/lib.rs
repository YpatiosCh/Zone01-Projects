#[derive(Debug, Clone, Copy, PartialEq)]
pub struct Color {
    pub r: u8,
    pub g: u8,
    pub b: u8,
    pub a: u8,
}

impl Color {
    pub fn swap(mut self, first: u8, second: u8) -> Color {
        [&mut self.r, &mut self.g, &mut self.b, &mut self.a]
            .into_iter()
            .for_each(|v| {
                if *v == first {
                    *v = second
                } else if *v == second {
                    *v = first;
                }
            });

        self
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        let c = Color {
            r: 255,
            g: 200,
            b: 10,
            a: 30,
        };

        assert_eq!(
            c.swap(c.r, c.b),
            Color {
                r: 10,
                g: 200,
                b: 255,
                a: 30
            }
        );
    }
}
