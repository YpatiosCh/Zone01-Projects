use rand::Rng;

#[derive(Debug, Clone, Copy, PartialEq)]
pub enum Suit {
    Heart,
    Diamond,
    Spade,
    Club,
}

#[derive(Debug, Clone, Copy, PartialEq)]
pub enum Rank {
    Ace,
    King,
    Queen,
    Jack,
    Number(u8),
}

impl Suit {
    pub fn random() -> Suit {
        Suit::translate(rand::thread_rng().gen_range(1..=4))
    }

    pub fn translate(value: u8) -> Suit {
        match value {
            1 => Suit::Heart,
            2 => Suit::Diamond,
            3 => Suit::Spade,
            _ => Suit::Club,
        }
    }
}

impl Rank {
    pub fn random() -> Rank {
        Rank::translate(rand::thread_rng().gen_range(1..=13))
    }

    pub fn translate(value: u8) -> Rank {
        match value {
            1 => Rank::Ace,
            11 => Rank::Jack,
            12 => Rank::Queen,
            13 => Rank::King,
            n => Rank::Number(n),
        }
    }
}

#[derive(Debug, PartialEq)]
pub struct Card {
    pub suit: Suit,
    pub rank: Rank,
}

pub fn winner_card(card: &Card) -> bool {
    matches!((&card.suit, &card.rank), (Suit::Spade, Rank::Ace))
}

#[cfg(test)]
mod tests {
    use super::*;

    // --- Suit::translate ---

    #[test]
    fn suit_translate_heart() {
        assert!(matches!(Suit::translate(1), Suit::Heart));
    }

    #[test]
    fn suit_translate_diamond() {
        assert!(matches!(Suit::translate(2), Suit::Diamond));
    }

    #[test]
    fn suit_translate_spade() {
        assert!(matches!(Suit::translate(3), Suit::Spade));
    }

    #[test]
    fn suit_translate_club() {
        assert!(matches!(Suit::translate(4), Suit::Club));
    }

    // --- Rank::translate ---

    #[test]
    fn rank_translate_ace() {
        assert!(matches!(Rank::translate(1), Rank::Ace));
    }

    #[test]
    fn rank_translate_numbers() {
        for n in 2..=10 {
            assert!(matches!(Rank::translate(n), Rank::Number(x) if x == n));
        }
    }

    #[test]
    fn rank_translate_jack() {
        assert!(matches!(Rank::translate(11), Rank::Jack));
    }

    #[test]
    fn rank_translate_queen() {
        assert!(matches!(Rank::translate(12), Rank::Queen));
    }

    #[test]
    fn rank_translate_king() {
        assert!(matches!(Rank::translate(13), Rank::King));
    }

    // --- winner_card ---

    #[test]
    fn winner_card_ace_of_spades() {
        let card = Card {
            suit: Suit::Spade,
            rank: Rank::Ace,
        };
        assert!(winner_card(&card));
    }

    #[test]
    fn winner_card_ace_of_hearts_is_not_winner() {
        let card = Card {
            suit: Suit::Heart,
            rank: Rank::Ace,
        };
        assert!(!winner_card(&card));
    }

    #[test]
    fn winner_card_spade_but_not_ace() {
        let card = Card {
            suit: Suit::Spade,
            rank: Rank::King,
        };
        assert!(!winner_card(&card));
    }

    #[test]
    fn winner_card_number_spade_is_not_winner() {
        let card = Card {
            suit: Suit::Spade,
            rank: Rank::Number(7),
        };
        assert!(!winner_card(&card));
    }

    // --- random (sanity checks, not deterministic) ---

    #[test]
    fn suit_random_returns_valid_suit() {
        for _ in 0..100 {
            assert!(matches!(
                Suit::random(),
                Suit::Heart | Suit::Diamond | Suit::Spade | Suit::Club
            ));
        }
    }

    #[test]
    fn rank_random_returns_valid_rank() {
        for _ in 0..100 {
            assert!(matches!(
                Rank::random(),
                Rank::Ace | Rank::Number(_) | Rank::Jack | Rank::Queen | Rank::King
            ));
        }
    }
}
