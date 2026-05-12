#[derive(Debug, Eq, PartialEq, Clone)]
pub struct Light {
    pub alias: String,
    pub brightness: u8,
}

impl Light {
    pub fn new(alias: &str) -> Self {
        Light {
            alias: alias.to_string(),
            brightness: 0,
        }
    }
}

pub fn change_brightness(lights: &mut [Light], alias: &str, value: u8) {
    if let Some(light) = lights.iter_mut().find(|l| l.alias == alias) {
        light.brightness = value;
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        let mut lights = ["living_room", "bedroom", "rest_room"].map(Light::new);
        change_brightness(&mut lights, "living_room", 200);

        assert_eq!(200, lights[0].brightness);
    }
}
