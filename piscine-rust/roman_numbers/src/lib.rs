use crate::RomanDigit::*;

#[derive(Copy, Clone, Debug, PartialEq, Eq)]
pub enum RomanDigit {
    Nulla,
    I,
    V,
    X,
    L,
    C,
    D,
    M,
}

#[derive(Clone, Debug, PartialEq, Eq)]
pub struct RomanNumber(pub Vec<RomanDigit>);

impl From<u32> for RomanNumber {
    fn from(mut value: u32) -> Self {
        if value == 0 {
            return RomanNumber(vec![Nulla]);
        }

        let digits = [
            (1000, M),
            (900, C),
            (900, M),
            (500, D),
            (400, C),
            (400, D),
            (100, C),
            (90, X),
            (90, C),
            (50, L),
            (40, X),
            (40, L),
            (10, X),
            (9, I),
            (9, X),
            (5, V),
            (4, I),
            (4, V),
            (1, I),
        ];

        let mut result = vec![];

        let mut i = 0;
        while i < digits.len() {
            let (val, digit) = digits[i];
            if value >= val {
                // subtractive pairs: 900=CM, 400=CD, 90=XC, 40=XL, 9=IX, 4=IV
                // they're stored as two consecutive entries with the same value
                if i + 1 < digits.len() && digits[i + 1].0 == val {
                    result.push(digit);
                    result.push(digits[i + 1].1);
                    value -= val;
                    i += 2;
                } else {
                    result.push(digit);
                    value -= val;
                }
            } else {
                i += 1;
            }
        }

        RomanNumber(result)
    }
}
