const NUMBERS: &[(u64, &str)] = &[
    (0, "zero"),
    (1, "one"),
    (2, "two"),
    (3, "three"),
    (4, "four"),
    (5, "five"),
    (6, "six"),
    (7, "seven"),
    (8, "eight"),
    (9, "nine"),
    (10, "ten"),
    (11, "eleven"),
    (12, "twelve"),
    (13, "thirteen"),
    (14, "fourteen"),
    (15, "fifteen"),
    (16, "sixteen"),
    (17, "seventeen"),
    (18, "eighteen"),
    (19, "nineteen"),
    (20, "twenty"),
    (30, "thirty"),
    (40, "forty"),
    (50, "fifty"),
    (60, "sixty"),
    (70, "seventy"),
    (80, "eighty"),
    (90, "ninety"),
];

#[derive(Debug, PartialEq, Eq, Clone, Copy)]
pub enum Number {
    Ones,
    Tens,
    Hundreds,
    Thousands,
    TensThousands,
    HundredThousands,
    Million,
}

pub fn spell(n: u64) -> String {
    let num = what_number(n.to_string());

    match num {
        Number::Million => "one million".to_string(),
        Number::Ones => implement_ones(n),
        Number::Tens => implement_tens(n),
        Number::Hundreds => implement_hundreds(n),
        Number::Thousands | Number::TensThousands | Number::HundredThousands => {
            implement_thousands(n)
        }
    }
}

fn what_number(str: String) -> Number {
    match str.len() {
        1 => Number::Ones,
        2 => Number::Tens,
        3 => Number::Hundreds,
        4 => Number::Thousands,
        5 => Number::TensThousands,
        6 => Number::HundredThousands,
        7 => Number::Million,
        _ => unreachable!(),
    }
}

fn implement_ones(n: u64) -> String {
    NUMBERS
        .iter()
        .find(|(num, _)| *num == n)
        .map(|(_, name)| name.to_string())
        .unwrap_or_default()
}

fn implement_tens(n: u64) -> String {
    if let Some((_, name)) = NUMBERS.iter().find(|(num, _)| *num == n) {
        name.to_string()
    } else {
        let tens = (n / 10) * 10;
        let ones = n % 10;
        let tens_name = NUMBERS
            .iter()
            .find(|(num, _)| *num == tens)
            .map(|(_, name)| name.to_string())
            .unwrap_or_default();
        let ones_name = implement_ones(ones);

        format!("{tens_name}-{ones_name}")
    }
}

fn implement_hundreds(n: u64) -> String {
    let hundreds = n / 100;
    let remainder = n % 100;

    let hundreds_name = implement_ones(hundreds);

    if remainder == 0 {
        format!("{hundreds_name} hundred")
    } else {
        let remainder_name = implement_tens(remainder);
        format!("{hundreds_name} hundred {remainder_name}")
    }
}

fn implement_thousands(n: u64) -> String {
    let thousands = n / 1000;
    let remainder = n % 1000;
    let thousands_name = spell_below_thousand(thousands);

    if remainder == 0 {
        format!("{thousands_name} thousand")
    } else {
        let remainder_name = spell_below_thousand(remainder);
        format!("{thousands_name} thousand {remainder_name}")
    }
}

fn spell_below_thousand(n: u64) -> String {
    match what_number(n.to_string()) {
        Number::Ones => implement_ones(n),
        Number::Tens => implement_tens(n),
        Number::Hundreds => implement_hundreds(n),
        _ => unreachable!(),
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn works() {
        assert_eq!(spell(23), "twenty-three");
        assert_eq!(spell(143), "one hundred forty-three");
    }
}
