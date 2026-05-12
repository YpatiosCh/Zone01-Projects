use std::fs::File;
use std::io::{self, BufWriter, Write};
use std::path::{Path, PathBuf};

use image_rs::codecs::png::PngEncoder;
use image_rs::{ColorType, ImageEncoder, ImageResult};
use objects::Color;

pub struct Framebuffer {
    pub w: u32,
    pub h: u32,
    pub pixels: Vec<Color>,
}

impl Framebuffer {
    /// Allocate a `w × h` framebuffer initialized to black.
    pub fn new(w: u32, h: u32) -> Self {
        let len = (w as usize) * (h as usize);
        Self { w, h, pixels: vec![Color::BLACK; len] }
    }
    /// Write color `c` at pixel `(x, y)`. Panics if out of bounds.
    pub fn set(&mut self, x: u32, y: u32, c: Color) {
        assert!(x < self.w, "x coordinate {x} out of bounds for width {}", self.w);
        assert!(y < self.h, "y coordinate {y} out of bounds for height {}", self.h);

        let index = (y as usize) * (self.w as usize) + (x as usize);
        self.pixels[index] = c;
    }
    /// Serialize as PPM P3 (ASCII): header + one `r g b` triple per pixel.
    pub fn write_ppm<W: Write>(&self, mut out: W) -> io::Result<()> {
        writeln!(out, "P3")?;
        writeln!(out, "{} {}", self.w, self.h)?;
        writeln!(out, "255")?;

        for pixel in &self.pixels {
            writeln!(
                out,
                "{} {} {}",
                color_channel_to_u8(pixel.r),
                color_channel_to_u8(pixel.g),
                color_channel_to_u8(pixel.b)
            )?;
        }

        Ok(())
    }

    /// Serialize as PNG: header + compressed RGB pixel data.
    pub fn write_png<W: Write>(&self, out: W) -> ImageResult<()> {
        let encoder = PngEncoder::new(out);
        let bytes = self.rgb8_bytes();
        encoder.write_image(&bytes, self.w, self.h, ColorType::Rgb8.into())
    }

    /// Save the framebuffer directly to a `.ppm` file path.
    pub fn save_ppm<P: AsRef<Path>>(&self, path: P) -> io::Result<()> {
        let file = File::create(path)?;
        self.write_ppm(BufWriter::new(file))
    }

    /// Save the framebuffer directly to a `.png` file path.
    pub fn save_png<P: AsRef<Path>>(&self, path: P) -> ImageResult<()> {
        let file = File::create(path)?;
        self.write_png(BufWriter::new(file))
    }

    fn rgb8_bytes(&self) -> Vec<u8> {
        let mut bytes = Vec::with_capacity(self.pixels.len() * 3);
        for pixel in &self.pixels {
            bytes.push(color_channel_to_u8(pixel.r));
            bytes.push(color_channel_to_u8(pixel.g));
            bytes.push(color_channel_to_u8(pixel.b));
        }
        bytes
    }
}

fn color_channel_to_u8(channel: f32) -> u8 {
    (channel.clamp(0.0, 1.0) * 255.0).round() as u8
}

pub fn png_path_for<P: AsRef<Path>>(path: P) -> PathBuf {
    path.as_ref().with_extension("png")
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn framebuffer_is_black_by_default() {
        let framebuffer = Framebuffer::new(3, 2);

        assert_eq!(framebuffer.w, 3);
        assert_eq!(framebuffer.h, 2);
        assert_eq!(framebuffer.pixels.len(), 6);
        assert!(framebuffer.pixels.iter().all(|pixel| {
            pixel.r == Color::BLACK.r && pixel.g == Color::BLACK.g && pixel.b == Color::BLACK.b
        }));
    }

    #[test]
    fn set_writes_pixel_at_row_major_index() {
        let mut framebuffer = Framebuffer::new(4, 3);
        let color = Color::rgb(0.2, 0.4, 0.6);

        framebuffer.set(2, 1, color);

        let pixel = framebuffer.pixels[6];
        assert_eq!(pixel.r, color.r);
        assert_eq!(pixel.g, color.g);
        assert_eq!(pixel.b, color.b);
    }

    #[test]
    fn write_ppm_serializes_header_and_clamped_pixels() {
        let mut framebuffer = Framebuffer::new(2, 2);
        framebuffer.set(0, 0, Color::rgb(0.0, 0.0, 0.0));
        framebuffer.set(1, 0, Color::rgb(1.0, 0.5, -0.25));
        framebuffer.set(0, 1, Color::rgb(0.25, 0.75, 0.1));
        framebuffer.set(1, 1, Color::rgb(1.5, 1.0, 1.0));

        let mut ppm = Vec::new();
        framebuffer.write_ppm(&mut ppm).unwrap();

        let ppm = String::from_utf8(ppm).unwrap();
        assert_eq!(
            ppm,
            "P3\n2 2\n255\n0 0 0\n255 128 0\n64 191 26\n255 255 255\n"
        );
    }

    #[test]
    fn write_png_starts_with_png_signature() {
        let mut framebuffer = Framebuffer::new(1, 1);
        framebuffer.set(0, 0, Color::rgb(1.0, 0.5, 0.0));

        let mut png = Vec::new();
        framebuffer.write_png(&mut png).unwrap();

        assert_eq!(&png[..8], &[137, 80, 78, 71, 13, 10, 26, 10]);
    }

    #[test]
    fn png_path_uses_same_directory_and_basename() {
        let path = png_path_for("scenes_out/sphere.ppm");
        assert_eq!(path.to_string_lossy(), "scenes_out/sphere.png");
    }
}
