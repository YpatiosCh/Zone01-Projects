pub mod boss;
pub mod member;

pub use boss::Boss;
pub use member::{Member, Role};

use std::collections::{HashMap, HashSet};

#[derive(PartialEq, Debug, Clone)]
pub struct Mob {
    pub name: String,
    pub boss: Boss,
    pub members: HashMap<String, Member>,
    pub cities: HashSet<String>,
    pub wealth: u64,
}

impl Mob {
    pub fn recruit(&mut self, member: (&str, u32)) {
        self.members
            .insert(member.0.to_string(), Member::new(member.1));
    }

    pub fn attack(&mut self, enemy: &mut Mob) {
        let (winner, loser) = if self.combat_power() <= enemy.combat_power() {
            // lost
            (enemy, self)
        } else {
            // win
            (self, enemy)
        };

        let youngest_age = loser.members.values().map(|m| m.age).min().unwrap_or(0);

        loser.members.retain(|_, m| m.age > youngest_age);

        if loser.members.is_empty() {
            loser.surrender_to(winner);
        }
    }

    pub fn steal(&mut self, enemy: &mut Mob, value: u64) {
        // how much are we taking?
        let this_much = value.min(enemy.wealth);

        self.wealth += this_much;
        enemy.wealth -= this_much;
    }

    pub fn conquer_city(&mut self, mobs: &[&Mob], wanted_city: String) {
        if !mobs
            .iter()
            .flat_map(|m| &m.cities)
            .any(|c| *c == wanted_city)
        {
            self.cities.insert(wanted_city);
        }
    }

    fn combat_power(&self) -> u32 {
        self.members
            .values()
            .map(|m| match m.role {
                Role::Associate => 1,
                Role::Soldier => 2,
                Role::Caporegime => 3,
                Role::Underboss => 4,
            })
            .sum()
    }

    fn surrender_to(&mut self, enemy: &mut Mob) {
        // give cities
        enemy.cities.extend(self.cities.drain());

        // give wealth
        enemy.wealth += self.wealth;
        self.wealth = 0;
    }
}
