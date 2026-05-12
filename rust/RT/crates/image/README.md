# image

Pixel buffer (`Framebuffer`) and file export in PPM and PNG formats. This crate is the final stage of the render pipeline — the renderer writes colors into a `Framebuffer` and then the framebuffer is saved to disk.

---

## `Framebuffer`

A flat row-major array of `Color` values.

```rust
use image::Framebuffer;
use objects::Color;

let mut fb = Framebuffer::new(800, 600); // initialized to black
fb.set(x, y, Color::rgb(1.0, 0.5, 0.0));
```

### Layout

Pixels are stored row-by-row. Pixel `(x, y)` maps to index `y * width + x`. Row 0 is the top of the image.

### `set`

```rust
fb.set(x: u32, y: u32, c: Color)
```

Writes color `c` at pixel `(x, y)`. Panics if `x >= width` or `y >= height`. The renderer calls this once per pixel after computing the shaded color.

---

## Saving images

### PPM (Portable Pixmap)

```rust
fb.save_ppm("output.ppm")?;
```

Writes the **P3 ASCII** format:

```
P3
800 600
255
r g b
r g b
...
```

Each channel is converted from `f32 [0,1]` to `u8 [0,255]` by clamping and rounding. PPM files are uncompressed and human-readable. They can be opened in most image viewers and converted with tools like ImageMagick.

### PNG

```rust
fb.save_png("output.png")?;
```

Writes a compressed PNG using the `image` crate's `PngEncoder`. Channel conversion is the same as PPM (`color_channel_to_u8`: clamp to `[0,1]`, multiply by 255, round). PNG is lossless and produces much smaller files than PPM.

### Writing to a stream

Both formats can also write to any `Write` implementor (e.g. an in-memory buffer):

```rust
let mut buffer = Vec::new();
fb.write_ppm(&mut buffer)?;

let mut buffer = Vec::new();
fb.write_png(&mut buffer)?;
```

---

## `png_path_for`

```rust
pub fn png_path_for<P: AsRef<Path>>(path: P) -> PathBuf
```

A convenience function that replaces or adds the `.png` extension on a path:

```rust
png_path_for("scenes_out/sphere.ppm") // → "scenes_out/sphere.png"
png_path_for("output")                // → "output.png"
```

Used by `rt` and `liveview` to derive the PNG output path from the PPM path so both files are always saved together.

---

## Color conversion

All export paths share `color_channel_to_u8`:

```rust
fn color_channel_to_u8(channel: f32) -> u8 {
    (channel.clamp(0.0, 1.0) * 255.0).round() as u8
}
```

Values below 0 clamp to 0 and values above 1 clamp to 255. The render pipeline may produce slightly out-of-range values due to high-intensity lights or accumulated tints; the clamp here is the final safety net.
