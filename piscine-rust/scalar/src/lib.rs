pub fn sum(a: u8, b: u8) -> u8 {
    a + b
}

pub fn diff(a: i16, b: i16) -> i16 {
    a - b
}

pub fn pro(a: i8, b: i8) -> i8 {
    a * b
}

pub fn quo(a: f32, b: f32) -> f32 {
    a / b
}

pub fn rem(a: f32, b: f32) -> f32 {
    a % b
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        let sum = sum(2, 2);
        assert_eq!(sum, 4);

        let diff = diff(12, 2);
        assert_eq!(diff, 10);

        let quo: f32 = quo(10., 2.);
        assert_eq!(quo, 5.);

        let rem: f32 = rem(10., 2.);
        assert_eq!(rem, 0.);
    }
}
