//! # Matrix
//!
//! A generic, fixed-size matrix library built on top of `lalgebra_scalar`.
//!
//! Matrix dimensions are const generics, meaning sizes are known at compile time,
//! allowing stack allocation and zero-cost abstractions.
//!
//! ## Example
//!
//! ```
//! use matrix::Matrix;
//!
//! let m = Matrix::<3, 3, f64>::identity();
//! let z = Matrix::<4, 3, f64>::zero();
//! ```

use lalgebra_scalar::Scalar;

/// A matrix of width `W` and height `H` containing elements of type `T`.
///
/// Internally represented as a 2D array `[[T; W]; H]`, fully stack-allocated.
/// `T` must implement [`Scalar`] to support arithmetic operations.
///
/// # Type Parameters
///
/// * `W` - Number of columns (compile-time constant)
/// * `H` - Number of rows (compile-time constant)
/// * `T` - Element type, must implement [`Scalar`]
///
/// # Example
///
/// ```
/// use matrix::Matrix;
///
/// let m = Matrix([[1, 2, 3], [4, 5, 6]]);
/// ```
#[derive(Debug, Eq, PartialEq, Clone, Copy)]
pub struct Matrix<const W: usize, const H: usize, T>(pub [[T; W]; H]);

impl<const W: usize, const H: usize, T: Scalar<Item = T>> Matrix<W, H, T> {
    /// Returns a `W x H` matrix with all elements set to `T::zero()`.
    ///
    /// # Example
    ///
    /// ```
    /// use matrix::Matrix;
    ///
    /// let m = Matrix::<3, 2, f64>::zero();
    /// assert_eq!(m, Matrix([[0.0, 0.0, 0.0], [0.0, 0.0, 0.0]]));
    /// ```
    pub fn zero() -> Self {
        Matrix([[T::zero(); W]; H])
    }
}

impl<const S: usize, T: Scalar<Item = T> + PartialEq> Matrix<S, S, T> {
    /// Returns the `S x S` identity matrix.
    ///
    /// The identity matrix has `T::one()` on the main diagonal and
    /// `T::zero()` everywhere else. Multiplying any matrix by the identity
    /// matrix returns the original matrix unchanged.
    ///
    /// Only defined for square matrices — `W` and `H` must be equal.
    ///
    /// # Example
    ///
    /// ```
    /// use matrix::Matrix;
    ///
    /// let m = Matrix::<3, 3, u32>::identity();
    /// assert_eq!(m, Matrix([[1, 0, 0], [0, 1, 0], [0, 0, 1]]));
    /// ```
    pub fn identity() -> Self {
        let mut m = [[T::zero(); S]; S];
        for i in 0..S {
            m[i][i] = T::one();
        }
        Matrix(m)
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn zero_integer_matrix() {
        let m = Matrix::<4, 3, u32>::zero();
        assert_eq!(m, Matrix([[0, 0, 0, 0], [0, 0, 0, 0], [0, 0, 0, 0]]));
    }

    #[test]
    fn zero_float_matrix() {
        let m = Matrix::<3, 4, f64>::zero();
        assert_eq!(
            m,
            Matrix([
                [0.0, 0.0, 0.0],
                [0.0, 0.0, 0.0],
                [0.0, 0.0, 0.0],
                [0.0, 0.0, 0.0]
            ])
        );
    }

    #[test]
    fn identity_4x4() {
        let m = Matrix::<4, 4, u32>::identity();
        assert_eq!(
            m,
            Matrix([[1, 0, 0, 0], [0, 1, 0, 0], [0, 0, 1, 0], [0, 0, 0, 1]])
        );
    }

    #[test]
    fn identity_1x1() {
        let m = Matrix::<1, 1, u32>::identity();
        assert_eq!(m, Matrix([[1]]));
    }

    #[test]
    fn identity_diagonal_is_one() {
        let m = Matrix::<3, 3, f64>::identity();
        for i in 0..3 {
            assert_eq!(m.0[i][i], 1.0);
        }
    }

    #[test]
    fn identity_off_diagonal_is_zero() {
        let m = Matrix::<3, 3, f64>::identity();
        for i in 0..3 {
            for j in 0..3 {
                if i != j {
                    assert_eq!(m.0[i][j], 0.0);
                }
            }
        }
    }

    #[test]
    fn matrix_from_literal() {
        let m = Matrix([[0; 4]; 3]);
        assert_eq!(m, Matrix([[0, 0, 0, 0], [0, 0, 0, 0], [0, 0, 0, 0]]));
    }
}
