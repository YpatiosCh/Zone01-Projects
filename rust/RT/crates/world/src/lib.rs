pub use camera::*;
pub use image::*;
pub use lighting::*;
pub use objects::*;
pub use raycaster::*;

use glam::Vec3;
use std::f32::consts::TAU;

// Radius of the disk used to spread point light samples (world units).
const AREA_LIGHT_RADIUS: f32 = 1.0;
// Half-arc of the elevation sweep used to spread directional light samples (radians).
const DIRECTIONAL_SPREAD: f32 = 0.2;

pub struct World {
    pub camera: Camera,
    pub objects: Vec<Box<dyn Object>>,
    pub lights: Vec<LightSource>,
    pub sky: SkyProperties,
}

impl World {
    /// New empty world with the given camera (no objects, no lights, default sky).
    pub fn new(camera: Camera) -> Self {
        Self {
            camera,
            objects: Vec::new(),
            lights: Vec::new(),
            sky: SkyProperties::default(),
        }
    }
    /// Add any `Object` impl; boxed internally for dynamic dispatch.
    pub fn add_object<O: Object + 'static>(&mut self, o: O) {
        self.objects.push(Box::new(o));
    }
    /// Register a light source.
    pub fn add_light(&mut self, l: LightSource) {
        self.lights.push(l);
    }

    /// Register `light_num` slightly varied copies of `base` to approximate a soft area light.
    ///
    /// For `light_num == 1` the light is added unchanged. For larger counts, samples are
    /// distributed evenly around the base:
    /// - **Point** — positions spread on a small horizontal disk (`AREA_LIGHT_RADIUS`).
    /// - **Directional** — directions spread across a narrow cone (`DIRECTIONAL_SPREAD`).
    ///
    /// Intensity is split equally across all samples so total scene energy is preserved.
    pub fn add_lights(&mut self, base: LightSource, light_num: u32) {
        if light_num == 0 {
            return;
        }
        if light_num == 1 {
            self.lights.push(base);
            return;
        }

        match base {
            LightSource::Point {
                pos,
                intensity,
                color,
            } => {
                let per = intensity / light_num as f32;
                for i in 0..light_num {
                    let angle = TAU * i as f32 / light_num as f32;
                    let offset = Vec3::new(
                        angle.cos() * AREA_LIGHT_RADIUS,
                        0.0,
                        angle.sin() * AREA_LIGHT_RADIUS,
                    );
                    self.lights.push(LightSource::Point {
                        pos: pos + offset,
                        intensity: per,
                        color,
                    });
                }
            }
            LightSource::Directional {
                dir,
                intensity,
                color,
            } => {
                let per = intensity / light_num as f32;
                let dir_n = dir.normalize();
                // Elevation axis: component of world-up perpendicular to the light direction.
                // Perturbing dir_n along this axis raises or lowers the light's altitude angle.
                let y_perp = (Vec3::Y - dir_n * dir_n.dot(Vec3::Y)).normalize();
                for i in 0..light_num {
                    // Spread linearly from -half to +half so the center sample is unchanged.
                    let t = if light_num == 1 {
                        0.0
                    } else {
                        (i as f32 / (light_num - 1) as f32) - 0.5
                    };
                    let new_dir = (dir_n + y_perp * t * DIRECTIONAL_SPREAD).normalize();
                    self.lights.push(LightSource::Directional {
                        dir: new_dir,
                        intensity: per,
                        color,
                    });
                }
            }
        }
    }
    /// Render at `w × h`. Mirror tint from `trace_with_reflections` is still applied.
    pub fn render(&self, w: u32, h: u32) -> Framebuffer {
        let mut fb = Framebuffer::new(w, h);
        let camera = Camera::new(
            self.camera.position,
            self.camera.look_at,
            self.camera.up,
            self.camera.fov,
            w,
            h,
        );
        for y in 0..h {
            for x in 0..w {
                let ray = camera.ray_for_pixel(x, y);
                let (tint, terminal_hit) = Raycaster::trace_with_reflections(ray, &self.objects, 8);
                let color = match terminal_hit {
                    Some(hit) => shade(&hit, &self.lights, &self.objects) * tint,
                    None => self.sky.color * tint,
                };
                fb.set(x, y, color);
            }
        }
        fb
    }
}
