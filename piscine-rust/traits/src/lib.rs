use std::fmt;

#[derive(Debug)]
pub struct Player<'a> {
    pub name: &'a str,
    pub strength: f64,
    pub score: u32,
    pub money: u32,
    pub weapons: Vec<&'a str>,
}

pub struct Fruit {
    pub weight_in_kg: f64,
}

pub struct Meat {
    pub weight_in_kg: f64,
    pub fat_content: f64,
}

impl Player<'_> {
    pub fn eat(&mut self, food: impl Food) {
        self.strength += food.gives();
    }
}

impl fmt::Display for Player<'_> {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(
            f,
            "{}\nStrength: {}, Score: {}, Money: {}\nWeapons: {:?}",
            self.name, self.strength, self.score, self.money, self.weapons
        )
    }
}

pub trait Food {
    fn gives(&self) -> f64;
}

impl Food for Fruit {
    fn gives(&self) -> f64 {
        self.weight_in_kg * 4.0
    }
}

impl Food for Meat {
    fn gives(&self) -> f64 {
        let fat = self.weight_in_kg * self.fat_content;
        let protein = self.weight_in_kg - fat;
        (protein * 4.0) + (fat * 9.0)
    }
}
