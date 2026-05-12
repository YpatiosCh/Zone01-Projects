use lalgebra_scalar::Scalar;
use std::ops::Add;

#[derive(Debug, PartialEq)]
pub struct Vector<T: Scalar>(pub Vec<T>);

impl<T: Scalar<Item = T>> Add for Vector<T> {
    type Output = Option<Self>;

    fn add(self, rhs: Self) -> Option<Self> {
        if self.0.len() != rhs.0.len() {
            return None;
        }

        Some(Vector(
            self.0
                .into_iter()
                .zip(rhs.0.into_iter())
                .map(|(a, b)| a + b)
                .collect(),
        ))
    }
}

impl<T: Scalar<Item = T>> Vector<T> {
    pub fn dot(self, rhs: Self) -> Option<T> {
        if self.0.len() != rhs.0.len() {
            return None;
        }

        Some(
            self.0
                .into_iter()
                .zip(rhs.0.into_iter())
                .map(|(a, b)| a * b)
                .fold(T::zero(), |acc, x| acc + x),
        )
    }
}
