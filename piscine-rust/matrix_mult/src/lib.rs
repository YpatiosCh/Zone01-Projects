use std::ops::Mul;

#[derive(Debug, PartialEq, Clone)]
pub struct Matrix<T>(pub Vec<Vec<T>>);

impl<T: Clone> Matrix<T> {
    pub fn number_of_cols(&self) -> usize {
        self.0[0].len()
    }

    pub fn number_of_rows(&self) -> usize {
        self.0.len()
    }

    pub fn row(&self, n: usize) -> Vec<T> {
        self.0[n].clone()
    }

    pub fn col(&self, n: usize) -> Vec<T> {
        self.0.iter().map(|row| row[n].clone()).collect()
    }
}

impl<T> Mul for Matrix<T>
where
    T: Clone + Default + std::ops::Add<Output = T> + std::ops::Mul<Output = T>,
{
    type Output = Option<Self>;

    fn mul(self, rhs: Self) -> Option<Self> {
        // matrix multiplication requires left cols == right rows
        if self.number_of_cols() != rhs.number_of_rows() {
            return None;
        }

        let rows = self.number_of_rows();
        let cols = rhs.number_of_cols();
        let inner = self.number_of_cols();

        let mut result = vec![vec![T::default(); cols]; rows];

        for i in 0..rows {
            for j in 0..cols {
                for k in 0..inner {
                    result[i][j] =
                        result[i][j].clone() + self.0[i][k].clone() * rhs.0[k][j].clone();
                }
            }
        }

        Some(Matrix(result))
    }
}
