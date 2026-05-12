use lalgebra_scalar::Scalar;
use matrix::Matrix;
use std::ops::{Add, Mul, Sub};

#[derive(Debug, Eq, PartialEq, Clone, Copy)]
pub struct Wrapper<const W: usize, const H: usize, T>(pub Matrix<W, H, T>);

impl<const W: usize, const H: usize, T: Scalar<Item = T>> From<[[T; W]; H]> for Wrapper<W, H, T> {
    fn from(arr: [[T; W]; H]) -> Self {
        Wrapper(Matrix(arr))
    }
}

impl<const W: usize, const H: usize, T: Scalar<Item = T>> Add for Wrapper<W, H, T> {
    type Output = Self;

    fn add(self, rhs: Self) -> Self {
        let mut result: [[_; W]; H] = [[T::zero(); W]; H];
        for i in 0..H {
            for j in 0..W {
                result[i][j] = self.0.0[i][j] + rhs.0.0[i][j];
            }
        }
        Wrapper(Matrix(result))
    }
}

impl<const W: usize, const H: usize, T: Scalar<Item = T>> Sub for Wrapper<W, H, T> {
    type Output = Self;

    fn sub(self, rhs: Self) -> Self {
        let mut result = [[T::zero(); W]; H];
        for i in 0..H {
            for j in 0..W {
                result[i][j] = self.0.0[i][j] - rhs.0.0[i][j];
            }
        }
        Wrapper(Matrix(result))
    }
}

impl<const S: usize, T: Scalar<Item = T>> Mul for Wrapper<S, S, T> {
    type Output = Self;

    fn mul(self, rhs: Self) -> Self {
        let mut result = [[T::zero(); S]; S];
        for i in 0..S {
            for j in 0..S {
                for k in 0..S {
                    result[i][j] = result[i][j] + self.0.0[i][k] * rhs.0.0[k][j];
                }
            }
        }
        Wrapper(Matrix(result))
    }
}
