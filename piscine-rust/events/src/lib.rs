use colored::*;
use std::{fmt, time::Duration};

#[derive(Debug, PartialEq, Clone, Copy)]
pub enum Position {
    Top,
    Bottom,
    Center,
}

#[derive(Debug, PartialEq, Clone)]
pub struct Notification {
    pub size: u32,
    pub color: (u8, u8, u8),
    pub position: Position,
    pub content: String,
}

#[derive(Clone, Copy)]
pub enum Event<'a> {
    Remainder(&'a str),
    Registration(Duration),
    Appointment(&'a str),
    Holiday,
}

impl fmt::Display for Notification {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        let (r, g, b) = self.color;
        let colored = self.content.truecolor(r, g, b);

        write!(f, "({:?}, {}, {})", self.position, self.size, colored)
    }
}

impl<'a> Event<'a> {
    pub fn notify(self) -> Notification {
        match self {
            Event::Remainder(smth) => Notification {
                size: 50,
                color: (50, 50, 50),
                position: (Position::Bottom),
                content: smth.to_owned(),
            },
            Event::Appointment(smth) => Notification {
                size: 100,
                color: (200, 200, 3),
                position: Position::Center,
                content: smth.to_owned(),
            },
            Event::Holiday => Notification {
                size: 25,
                color: (0, 255, 0),
                position: Position::Top,
                content: "Enjoy your holiday".to_owned(),
            },
            Event::Registration(duration) => {
                let total = duration.as_secs();
                let hours = total / 3600;
                let minutes = (total % 3600) / 60;
                let seconds = total % 60;
                Notification {
                    size: 30,
                    color: (255, 2, 22),
                    position: Position::Top,
                    content: format!(
                        "You have {}H:{}M:{}S left before the registration ends",
                        hours, minutes, seconds
                    ),
                }
            }
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        // RUN WITH `cargo test -- --nocapture`

        let notif1 = Event::Remainder("Go to the doctor man").notify();
        let notif2 = Event::Appointment("Go to the doctor man").notify();
        let notif3 = Event::Registration(Duration::from_secs(49094)).notify();
        let notif4 = Event::Holiday.notify();

        println!("{notif1}");
        println!("{notif2}");
        println!("{notif3}");
        println!("{notif4}");
    }
}
