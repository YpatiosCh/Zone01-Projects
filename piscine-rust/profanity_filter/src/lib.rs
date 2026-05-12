pub fn check_ms(message: &str) -> Result<&str, &str> {
    if message == "" || message.contains("stupid") {
        Err("ERROR: illegal")
    } else {
        Ok(message)
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        assert_eq!(check_ms("hello"), Ok("hello"));
        assert_eq!(check_ms("stupid monkey"), Err("ERROR: illegal"));
    }
}
