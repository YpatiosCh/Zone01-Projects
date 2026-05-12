#[derive(Debug, Clone, Copy)]
pub struct Circle {
    pub center: Point,
    pub radius: f64,
}

impl Circle {
    pub fn new(x: f64, y: f64, radius: f64) -> Self {
        Circle {
            center: Point(x, y),
            radius: radius,
        }
    }

    pub fn area(&self) -> f64 {
        std::f64::consts::PI * self.radius * self.radius
    }

    pub fn diameter(&self) -> f64 {
        self.radius * 2.0
    }

    pub fn intersect(&self, other: Circle) -> bool {
        self.center.distance(other.center) <= (self.radius + other.radius)
    }
}

#[derive(Debug, Clone, Copy)]
pub struct Point(pub f64, pub f64);

impl Point {
    pub fn distance(&self, other: Point) -> f64 {
        let dx = self.0 - other.0;
        let dy = self.1 - other.1;
        (dx * dx + dy * dy).sqrt()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        let circle: Circle = Circle::new(500., 500., 150.);
        let circle1 = Circle {
            center: Point(80.0, 115.0),
            radius: 30.0,
        };
        print!("{:?}", circle);

        assert_eq!(70685.83470577035, circle.area());
        assert_eq!(300., circle.diameter());
        assert_eq!(false, circle.intersect(circle1));

        let point_a = Point(1.0, 1.0);
        let point_b = Point(0.0, 0.0);

        assert_eq!(1.4142135623730951, point_a.distance(point_b));
    }
}
