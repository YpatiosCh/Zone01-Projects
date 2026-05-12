use chrono::Local;

#[derive(Debug, Eq, PartialEq)]
pub struct FormError {
    pub form_values: (&'static str, String),
    pub date: String,
    pub err: &'static str,
}

impl FormError {
    pub fn new(field_name: &'static str, field_value: String, err: &'static str) -> Self {
        FormError {
            form_values: (field_name, field_value),
            date: Local::now().format("%Y-%m-%d %H:%M:%S").to_string(),
            err: err,
        }
    }
}

#[derive(Debug, Eq, PartialEq)]
pub struct Form {
    pub name: String,
    pub password: String,
}

impl Form {
    pub fn validate(&self) -> Result<(), FormError> {
        if self.name.is_empty() {
            return Err(FormError::new(
                "name",
                self.name.clone(),
                "Username is empty",
            ));
        }
        if self.password.len() < 8 {
            return Err(FormError::new(
                "password",
                self.password.clone(),
                "Password should be at least 8 characters long",
            ));
        }

        let is_alphabetic = self.password.chars().any(|c| c.is_alphabetic());
        let is_numeric: bool = self.password.chars().any(|c| c.is_ascii_digit());
        let has_symbol = self.password.chars().any(|c| c.is_ascii_punctuation());

        if !is_alphabetic || !has_symbol || !is_numeric {
            return Err(FormError::new(
                "password",
                self.password.clone(),
                "Password should be a combination of ASCII numbers, letters and symbols",
            ));
        }
        Ok(())
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        // VALID NAME AND PASSWORD

        let form = Form {
            name: "Lee".to_owned(),
            password: "qwqwsa1dty_".to_owned(),
        };

        assert_eq!(Ok(()), form.validate());
    }

    #[test]
    fn empty_name() {
        let form = Form {
            name: "".to_owned(),
            password: "qwqwsa1dty_".to_owned(),
        };

        let result = form.validate().unwrap_err();
        let expected = FormError::new("name", "".to_string(), "Username is empty");

        assert_eq!(result.form_values, expected.form_values);
        assert_eq!(result.err, expected.err);
    }

    #[test]
    fn password_length_err() {
        let form = Form {
            name: "Hello".to_owned(),
            password: "edi".to_owned(),
        };

        let result = form.validate().unwrap_err();
        let expected = FormError::new(
            "password",
            form.password.clone(),
            "Password should be at least 8 characters long",
        );

        assert_eq!(result.form_values, expected.form_values);
        assert_eq!(result.err, expected.err);
    }

    #[test]
    fn password_alphanumericsymbol_err() {
        let not_numeric = Form {
            name: "Hello".to_owned(),
            password: "helowdfnmdncjnl!!".to_owned(),
        };

        let result = not_numeric.validate().unwrap_err();
        let expected = FormError::new(
            "password",
            not_numeric.password.clone(),
            "Password should be a combination of ASCII numbers, letters and symbols",
        );

        assert_eq!(result.form_values, expected.form_values);
        assert_eq!(result.err, expected.err);

        let not_symbolic = Form {
            name: "Hello".to_owned(),
            password: "helowdfnmdncjnl123".to_owned(),
        };

        let result = not_symbolic.validate().unwrap_err();
        let expected = FormError::new(
            "password",
            not_symbolic.password.clone(),
            "Password should be a combination of ASCII numbers, letters and symbols",
        );

        assert_eq!(result.form_values, expected.form_values);
        assert_eq!(result.err, expected.err);
    }
}
