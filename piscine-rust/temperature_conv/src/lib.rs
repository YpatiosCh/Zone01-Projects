pub fn fahrenheit_to_celsius(f: f64) -> f64 {
    (f - 32.0) / (9.0 / 5.0)
}

pub fn celsius_to_fahrenheit(c: f64) -> f64 {
    c * (9.0 / 5.0) + 32.0
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        let c: f64 = fahrenheit_to_celsius(-459.67);
        assert_eq!(c, -273.15);

        let f: f64 = celsius_to_fahrenheit(0.0);
        assert_eq!(f, 32.0);
    }
}
