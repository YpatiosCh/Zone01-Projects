#[derive(Debug, Clone, Eq, PartialEq)]
pub struct GameSession {
    pub id: u32,
    pub p1: (String, u32),
    pub p2: (String, u32),
    pub nb_games: u32,
}

impl GameSession {
    pub fn new(id: u32, p1_name: String, p2_name: String, nb_games: u32) -> GameSession {
        GameSession {
            id,
            p1: (p1_name, 0),
            p2: (p2_name, 0),
            nb_games,
        }
    }

    pub fn read_winner(&self) -> Option<&(String, u32)> {
        if self.p1.1 == self.p2.1 {
            None
        } else if self.p1.1 > self.p2.1 {
            Some(&self.p1)
        } else {
            Some(&self.p2)
        }
    }

    pub fn update_score(&mut self, user_name: &str) {
        let wins_to_end = self.nb_games / 2 + 1;
        if self.p1.1 >= wins_to_end || self.p2.1 >= wins_to_end {
            return;
        }

        if self.p1.0 == user_name {
            self.p1.1 += 1;
        } else if self.p2.0 == user_name {
            self.p2.1 += 1;
        }
    }

    pub fn delete(self) -> String {
        format!("game deleted: id -> {:?}", self.id)
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        let mut game = GameSession::new(0, String::from("Joao"), String::from("Susana"), 5);
        println!("{:?}", game.read_winner());

        game.update_score("Joao");
        game.update_score("Joao");
        game.update_score("Susana");
        game.update_score("Susana");
        println!("{:?}", game.read_winner());

        game.update_score("Joao");
        // This one will not count because it already 5 games played, the `nb_games`
        game.update_score("Susana");

        println!("{:?}", game.read_winner());

        println!("{:?}", game.delete());

        //game.read_winner();
        // This will give an error as the game was dropped with `delete` and no longer exists
    }
}
