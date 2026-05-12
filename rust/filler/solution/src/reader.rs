use std::io::BufRead;

/// A thin wrapper around any [`BufRead`] source that yields input one line at
/// a time, reusing an internal buffer to avoid per-line allocations on the
/// read side.
///
/// The wrapper is generic over the underlying reader so it can be driven by
/// `stdin` in production and by an in-memory [`std::io::Cursor`] in tests.
pub struct Reader<R: BufRead> {
    reader: R,
    buf: String,
}

impl<R: BufRead> Reader<R> {
    /// Construct a new [`Reader`] that will read from `reader`.
    ///
    /// The internal buffer starts empty and grows on demand.
    pub fn new(reader: R) -> Self {
        Reader {
            reader,
            buf: String::new(),
        }
    }

    /// Read the next line from the underlying source, with the trailing
    /// newline (and any `\r` before it) stripped.
    ///
    /// Returns:
    /// * `Some(line)` on a successful read (including for empty lines);
    /// * `None` on clean EOF *or* on any I/O error
    pub fn read_line(&mut self) -> Option<String> {
        self.buf.clear();
        match self.reader.read_line(&mut self.buf) {
            Ok(0) => None,
            Ok(_) => Some(self.buf.trim_end().to_string()),
            Err(_) => None,
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::io::Cursor;

    #[test]
    fn reads_lines_one_by_one() {
        let input = "hello\nworld\n";
        let cursor = Cursor::new(input);

        let mut reader = Reader::new(cursor);

        assert_eq!(reader.read_line(), Some("hello".to_string()));
        assert_eq!(reader.read_line(), Some("world".to_string()));
        assert_eq!(reader.read_line(), None);
    }

    #[test]
    fn trims_newlines() {
        let input = "abc\r\ndef\n";
        let cursor = Cursor::new(input);

        let mut reader = Reader::new(cursor);

        assert_eq!(reader.read_line(), Some("abc".to_string()));
        assert_eq!(reader.read_line(), Some("def".to_string()));
        assert_eq!(reader.read_line(), None);
    }

    #[test]
    fn handles_empty_lines() {
        let input = "\nhello\n\n";
        let cursor = Cursor::new(input);

        let mut reader = Reader::new(cursor);

        assert_eq!(reader.read_line(), Some("".to_string()));
        assert_eq!(reader.read_line(), Some("hello".to_string()));
        assert_eq!(reader.read_line(), Some("".to_string()));
        assert_eq!(reader.read_line(), None);
    }

    #[test]
    fn empty_input_returns_none() {
        let input = "";
        let cursor = Cursor::new(input);

        let mut reader = Reader::new(cursor);

        assert_eq!(reader.read_line(), None);
    }
}
