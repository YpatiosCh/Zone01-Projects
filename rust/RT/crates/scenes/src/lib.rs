pub mod all_objects;
pub mod all_objects_alt_cam;
pub mod cube;
pub mod cylinder;
pub mod plane;
pub mod plane_cube;
pub mod plane_four_spheres;
pub mod plane_reflective;
pub mod sphere;

use world::World;

static REGISTRY: &[(&str, fn() -> World)] = &[
    ("sphere", sphere::build),
    ("cylinder", cylinder::build),
    ("cube", cube::build),
    ("plane", plane::build),
    ("plane_four_spheres", plane_four_spheres::build),
    ("plane_reflective", plane_reflective::build),
    ("plane_cube", plane_cube::build),
    ("all_objects", all_objects::build),
    ("all_objects_alt_cam", all_objects_alt_cam::build),
];

pub fn get(name: &str) -> Option<fn() -> World> {
    REGISTRY.iter().find(|(n, _)| *n == name).map(|(_, f)| *f)
}

pub fn names() -> Vec<&'static str> {
    REGISTRY.iter().map(|(n, _)| *n).collect()
}
