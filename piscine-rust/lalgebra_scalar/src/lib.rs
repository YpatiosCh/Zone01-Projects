use std::ops::{Add, Div, Mul, Sub};

pub trait Scalar:
    Add<Output = Self>
    + Sub<Output = Self>
    + Mul<Output = Self>
    + Div<Output = Self>
    + Sized
    + Copy
    + Clone
    + std::fmt::Debug
{
    type Item;

    fn zero() -> Self::Item;
    fn one() -> Self::Item;
}

macro_rules! impl_scalar {
    ($t:ty, $zero:expr, $one:expr) => {
        impl Scalar for $t {
            type Item = $t;

            fn zero() -> Self::Item {
                $zero
            }

            fn one() -> Self::Item {
                $one
            }
        }
    };
}

impl_scalar!(u32, 0, 1);
impl_scalar!(u64, 0, 1);
impl_scalar!(i32, 0, 1);
impl_scalar!(i64, 0, 1);
impl_scalar!(f32, 0.0, 1.0);
impl_scalar!(f64, 0.0, 1.0);
