pub mod areas_volumes;

pub use areas_volumes::{GeometricalShapes, GeometricalVolumes};

use areas_volumes::*;

pub fn area_fit(
    (x, y): (usize, usize),
    kind: GeometricalShapes,
    times: usize,
    (a, b): (usize, usize),
) -> bool {
    let container = rectangle_area(x, y) as f64;

    let shape_area: f64 = match kind {
        GeometricalShapes::Square => square_area(a) as f64,
        GeometricalShapes::Rectangle => rectangle_area(a, b) as f64,
        GeometricalShapes::Triangle => triangle_area(a, b),
        GeometricalShapes::Circle => circle_area(a),
    };

    shape_area * times as f64 <= container
}

pub fn volume_fit(
    (x, y, z): (usize, usize, usize),
    kind: GeometricalVolumes,
    times: usize,
    (a, b, c): (usize, usize, usize),
) -> bool {
    let container = parallelepiped_volume(x, y, z) as f64;

    let shape_volume: f64 = match kind {
        GeometricalVolumes::Cube => cube_volume(a) as f64,
        GeometricalVolumes::Sphere => sphere_volume(a),
        GeometricalVolumes::Cone => cone_volume(a, b),
        GeometricalVolumes::TriangularPyramid => triangular_pyramid_volume(a as f64, b),
        GeometricalVolumes::Parallelepiped => parallelepiped_volume(a, b, c) as f64,
    };

    shape_volume * times as f64 <= container
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_area_fit() {
        // 4 squares of side 2 fit exactly in a 4x4 area (4*4 == 16)
        assert!(area_fit((4, 4), GeometricalShapes::Square, 4, (2, 0)));
        // 5 squares of side 2 do not fit in a 4x4 area (5*4 == 20 > 16)
        assert!(!area_fit((4, 4), GeometricalShapes::Square, 5, (2, 0)));
    }

    #[test]
    fn test_volume_fit() {
        // 8 cubes of side 1 fit exactly in a 2x2x2 box (8*1 == 8)
        assert!(volume_fit(
            (2, 2, 2),
            GeometricalVolumes::Cube,
            8,
            (1, 0, 0)
        ));
        // 9 cubes of side 1 do not fit in a 2x2x2 box (9*1 == 9 > 8)
        assert!(!volume_fit(
            (2, 2, 2),
            GeometricalVolumes::Cube,
            9,
            (1, 0, 0)
        ));
    }
}
