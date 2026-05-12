pub enum Security {
    Unknown,
    Message,
    Warning,
    NotFound,
    UnexpectedUrl,
}

pub fn fetch_data(server: Result<&str, &str>, security_level: Security) -> String {
    match security_level {
        Security::Unknown => server.unwrap().to_string(),
        Security::Message => server.expect("ERROR: program stops").to_string(),
        Security::Warning => match server {
            Ok(url) => url.to_string(),
            Err(_) => "WARNING: check the server".to_string(),
        },
        Security::NotFound => match server {
            Ok(url) => url.to_string(),
            Err(msg) => format!("Not found: {}", msg.to_string()),
        },
        Security::UnexpectedUrl => match server {
            Ok(url) => panic!("{}", url),
            Err(msg) => msg.to_string(),
        },
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    // ===== UNKNOWN ======
    #[test]
    fn ok_unknown() {
        assert_eq!(
            fetch_data(Ok("server1.com"), Security::Unknown),
            "server1.com"
        );
    }

    #[test]
    #[should_panic]
    fn err_unknown() {
        fetch_data(Err("error"), Security::Unknown);
    }

    // ===== MESSAGE ======
    #[test]
    fn ok_message() {
        assert_eq!(fetch_data(Ok("hello"), Security::Message), "hello");
    }

    #[test]
    #[should_panic(expected = "ERROR: program stops")]
    fn err_message() {
        assert_eq!(
            fetch_data(Err("some error"), Security::Message),
            "ERROR: program stops"
        );
    }

    // ===== WARNING ======
    #[test]
    fn ok_warning() {
        assert_eq!(fetch_data(Ok("hello"), Security::Warning), "hello");
    }

    #[test]
    fn err_warning() {
        assert_eq!(
            fetch_data(Err("some error"), Security::Warning),
            "WARNING: check the server"
        );
    }

    // ===== NOT FOUND ======
    #[test]
    fn ok_notfound() {
        assert_eq!(fetch_data(Ok("hello"), Security::NotFound), "hello");
    }

    #[test]
    fn err_notfound() {
        assert_eq!(fetch_data(Err("er"), Security::NotFound), "Not found: er");
    }

    // ===== UNEXPECTED URL ======

    #[test]
    #[should_panic(expected = "hello")]
    fn ok_unexpected() {
        fetch_data(Ok("hello"), Security::UnexpectedUrl);
    }

    #[test]
    fn err_unexpected() {
        assert_eq!(fetch_data(Err("err"), Security::UnexpectedUrl), "err");
    }
}
