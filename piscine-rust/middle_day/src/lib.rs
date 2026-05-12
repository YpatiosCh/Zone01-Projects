use chrono::{Datelike, NaiveDate, Weekday};

pub fn middle_day(year: u32) -> Option<Weekday> {
    let is_leap = is_leap_year(year);

    // Leap years have 366 days (even) → no single middle day
    if is_leap {
        return None;
    }

    // Non-leap years have 365 days → middle day is day 183
    let date = NaiveDate::from_yo_opt(year as i32, 183)?;

    Some(date.weekday())
}

fn is_leap_year(year: u32) -> bool {
    let y = year as i32;
    (y % 4 == 0 && y % 100 != 0) || (y % 400 == 0)
}

#[cfg(test)]
mod tests {
    use super::*;
    use chrono::Weekday;

    #[test]
    fn test_1022() {
        assert_eq!(middle_day(1022), Some(Weekday::Tue));
    }

    #[test]
    fn test_leap_year_returns_none() {
        assert_eq!(middle_day(2000), None);
        assert_eq!(middle_day(2024), None);
    }

    #[test]
    fn test_non_leap_year() {
        assert!(middle_day(2023).is_some());
    }
}
