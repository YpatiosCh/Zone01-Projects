use std::{fmt, str::FromStr};

#[derive(PartialEq, Eq, Hash, Clone, Copy)]
pub enum Antigen {
    A,
    AB,
    B,
    O,
}

#[derive(PartialEq, Eq, Hash, Clone, Copy)]
pub enum RhFactor {
    Positive,
    Negative,
}

#[derive(PartialEq, Eq, Hash, Clone, Copy)]
pub struct BloodType {
    pub antigen: Antigen,
    pub rh_factor: RhFactor,
}

impl FromStr for BloodType {
    type Err = ();

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        let s = s.trim();

        let (antigen_str, rh_factor) = if s.ends_with('+') {
            (&s[..s.len() - 1], RhFactor::Positive)
        } else if s.ends_with('-') {
            (&s[..s.len() - 1], RhFactor::Negative)
        } else {
            return Err(());
        };

        let antigen = match antigen_str {
            "A" => Antigen::A,
            "B" => Antigen::B,
            "AB" => Antigen::AB,
            "O" => Antigen::O,
            _ => return Err(()),
        };

        Ok(BloodType { antigen, rh_factor })
    }
}

impl fmt::Debug for BloodType {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        let antigen = match self.antigen {
            Antigen::A => "A",
            Antigen::B => "B",
            Antigen::AB => "AB",
            Antigen::O => "O",
        };
        let rh = match self.rh_factor {
            RhFactor::Positive => "+",
            RhFactor::Negative => "-",
        };
        write!(f, "{}{}", antigen, rh)
    }
}

impl BloodType {
    pub fn can_receive_from(self, other: Self) -> bool {
        let rh_ok = self.rh_factor == RhFactor::Positive || other.rh_factor == RhFactor::Negative;

        let antigen_ok = match (self.antigen, other.antigen) {
            _ if self.antigen == other.antigen => true,
            (_, Antigen::O) => true,
            (Antigen::AB, _) => true,
            _ => false,
        };

        rh_ok && antigen_ok
    }

    pub fn donors(self) -> Vec<Self> {
        Self::all()
            .into_iter()
            .filter(|&bt| self.can_receive_from(bt))
            .collect()
    }

    pub fn recipients(self) -> Vec<Self> {
        Self::all()
            .into_iter()
            .filter(|&bt| bt.can_receive_from(self))
            .collect()
    }

    fn all() -> Vec<BloodType> {
        let antigens = [Antigen::A, Antigen::B, Antigen::AB, Antigen::O];
        let rh_factors = [RhFactor::Negative, RhFactor::Positive];

        antigens
            .iter()
            .flat_map(|&antigen| {
                rh_factors
                    .iter()
                    .map(move |&rh_factor| BloodType { antigen, rh_factor })
            })
            .collect()
    }
}
