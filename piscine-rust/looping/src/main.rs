use std::io;

fn main() {
    start_riddle();
}

fn start_riddle() {
    let riddle: &str = "I am the beginning of the end, and the end of time and space. I am essential to creation, and I surround every place. What am I?";
    let answer: &str = "The letter e";
    let mut tries: u32 = 0;

    loop {
        println!("{}", riddle);

        let mut input = String::new();
        io::stdin()
            .read_line(&mut input)
            .expect("Failed to read line");
        let input = input.trim();

        tries += 1;

        if input == answer {
            println!("Number of trials: {}", tries);
            break;
        }
    }
}
