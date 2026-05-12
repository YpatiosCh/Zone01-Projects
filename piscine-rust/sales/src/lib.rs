#[derive(Debug, Clone, PartialEq)]
pub struct Store {
    pub products: Vec<(String, f32)>,
}

impl Store {
    pub fn new(products: Vec<(String, f32)>) -> Self {
        Store { products }
    }
}

#[derive(Debug, Clone, PartialEq)]
pub struct Cart {
    pub items: Vec<(String, f32)>,
    pub receipt: Vec<f32>,
}

impl Cart {
    pub fn new() -> Self {
        Cart {
            items: Default::default(),
            receipt: Default::default(),
        }
    }

    pub fn insert_item(&mut self, store: &Store, item: String) {
        if let Some(product) = store.products.iter().find(|(i, _)| *i == item) {
            self.items.push(product.clone());
        }
    }

    pub fn generate_receipt(&mut self) -> Vec<f32> {
        let mut prices: Vec<f32> = self.items.iter().map(|(_, p)| *p).collect();
        prices.sort_by(|a, b| a.partial_cmp(b).unwrap());

        let n = prices.len();
        let free_count = n / 3;

        let free_total: f32 = prices.iter().take(free_count).sum();
        let full_total: f32 = prices.iter().sum();
        let paid_total = full_total - free_total;
        let ratio = paid_total / full_total;

        let receipt: Vec<f32> = prices
            .iter()
            .map(|&p| (p * ratio * 100.0).round() / 100.0)
            .collect();

        self.receipt = receipt.clone();
        receipt
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        let store = Store::new(vec![
            (String::from("product A"), 1.23),
            (String::from("product B"), 23.1),
            (String::from("product C"), 3.12),
        ]);

        let mut cart = Cart::new();
        cart.insert_item(&store, String::from("product A"));
        cart.insert_item(&store, String::from("product B"));
        cart.insert_item(&store, String::from("product C"));

        let receipt = cart.generate_receipt();

        assert_eq!(receipt, [1.17, 2.98, 22.06]);
    }
}
