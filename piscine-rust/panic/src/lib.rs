use std::fs::File;

pub fn open_file(s: &str) -> File {
    let file = File::open(s);
    match file {
        Ok(this) => this,
        Err(err) => panic!("{}", err),
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::fs;

    #[test]
    #[should_panic]
    fn should_panic() {
        let filename = "created.txt";
        File::create(filename).unwrap();

        println!("{:?}", open_file(filename));

        fs::remove_file(filename).unwrap();

        // this should panic!
        open_file(filename);
    }

    #[test]
    fn should_open() {
        let filename = "test.txt";
        File::create(filename).unwrap();

        let file = open_file(filename);
        println!("{:?}", file);

        fs::remove_file(filename).unwrap();
    }
}
