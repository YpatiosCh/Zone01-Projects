use glam::{Quat, Vec3};
use std::fs;
use std::path::PathBuf;
use world::{Color, LightSource, World};
use window::{AppWindow, WindowKey};

const MOVE_SPEED: f32 = 0.25;
const VERTICAL_MOVE_SPEED: f32 = 0.20;
const MOUSE_SENSITIVITY: f32 = 0.004;
const MAX_PITCH: f32 = 1.45;
const RENDER_SCALE: usize = 3;
const CAMERA_CLEARANCE: f32 = 0.1;
const LIGHT_MOVE_SPEED: f32 = 0.25;
const LIGHT_ROTATE_SPEED: f32 = 0.03;
const SCREENSHOT_DIR: &str = "screenshots";

pub enum LiveviewAction {
	OpenSceneMenu,
	ExitApp,
}

fn finish_liveview(window: &mut AppWindow, action: LiveviewAction) -> LiveviewAction {
	window.confine_cursor(false);
	window.set_cursor_visible(true);
	action
}

/// Open a window and render `world` in real time. `WASD` moves the camera, mouse motion
/// rotates the view, and `Esc` closes the liveview.
pub fn run(mut world: World, w: u32, h: u32) -> LiveviewAction {
	let width = w.max(640) as usize;
	let height = h.max(480) as usize;
	let render_width = (width / RENDER_SCALE).max(320);
	let render_height = (height / RENDER_SCALE).max(240);

	let mut window = match AppWindow::new("RayTracing Liveview", width, height) {
		Ok(window) => window,
		Err(error) => {
			eprintln!("failed to open liveview window: {error}");
			return LiveviewAction::ExitApp;
		}
	};
	window.set_cursor_visible(false);
	window.center_cursor();

	let mut camera_position = world.camera.position;
	camera_position = clamp_camera_to_floor(&world, camera_position);
	let initial_forward = (world.camera.look_at - world.camera.position).normalize();
	let mut yaw = initial_forward.z.atan2(initial_forward.x);
	let mut pitch = initial_forward.y.asin();
	let mut needs_redraw = true;
	let mut was_active = false;

	while window.is_open() {
		let is_active = window.is_active();
		window.confine_cursor(is_active);
		if is_active && !was_active {
			window.center_cursor();
		}
		was_active = is_active;

		if window.is_key_pressed(WindowKey::Escape) {
			return finish_liveview(&mut window, LiveviewAction::OpenSceneMenu);
		}

		let mut camera_changed = false;
		if is_active {
			if let Some((dx, dy)) = window.mouse_delta_from_center() {
				if dx != 0.0 || dy != 0.0 {
					yaw += dx * MOUSE_SENSITIVITY;
					pitch = (pitch - dy * MOUSE_SENSITIVITY).clamp(-MAX_PITCH, MAX_PITCH);
					camera_changed = true;
				}
			}
		}

		let forward = direction_from_angles(yaw, pitch);
		let forward_flat = Vec3::new(forward.x, 0.0, forward.z).normalize_or_zero();
		let right = Vec3::new(-forward.z, 0.0, forward.x).normalize_or_zero();

		if window.is_key_down(WindowKey::W) {
			camera_position += forward_flat * MOVE_SPEED;
			camera_changed = true;
		}
		if window.is_key_down(WindowKey::S) {
			camera_position -= forward_flat * MOVE_SPEED;
			camera_changed = true;
		}
		if window.is_key_down(WindowKey::A) {
			camera_position -= right * MOVE_SPEED;
			camera_changed = true;
		}
		if window.is_key_down(WindowKey::D) {
			camera_position += right * MOVE_SPEED;
			camera_changed = true;
		}
		if window.is_key_down(WindowKey::Space) {
			camera_position.y += VERTICAL_MOVE_SPEED;
			camera_changed = true;
		}
		if window.is_key_down(WindowKey::LeftCtrl) || window.is_key_down(WindowKey::RightCtrl) {
			camera_position.y -= VERTICAL_MOVE_SPEED;
			camera_changed = true;
		}

		camera_position = clamp_camera_to_floor(&world, camera_position);

		let mut lights_changed = false;
		{
			let left  = window.is_key_down(WindowKey::Left);
			let right = window.is_key_down(WindowKey::Right);
			let up    = window.is_key_down(WindowKey::Up);
			let down  = window.is_key_down(WindowKey::Down);
			if left || right || up || down {
				let dx         = if right { LIGHT_MOVE_SPEED    } else if left { -LIGHT_MOVE_SPEED    } else { 0.0 };
				let dy         = if up    { LIGHT_MOVE_SPEED    } else if down { -LIGHT_MOVE_SPEED    } else { 0.0 };
				let yaw_delta  = if right { LIGHT_ROTATE_SPEED  } else if left { -LIGHT_ROTATE_SPEED  } else { 0.0 };
				let pitch_delta = if up   { LIGHT_ROTATE_SPEED  } else if down { -LIGHT_ROTATE_SPEED  } else { 0.0 };
				for light in &mut world.lights {
					match light {
						LightSource::Point { pos, .. } => {
							pos.x += dx;
							pos.y += dy;
						}
						LightSource::Directional { dir, .. } => {
							// Left/Right: orbit around world Y (azimuth).
							*dir = Quat::from_axis_angle(Vec3::Y, yaw_delta).mul_vec3(*dir).normalize();
							// Up/Down: tilt along the horizontal axis perpendicular to the direction.
							let elev_axis = Vec3::new(-dir.z, 0.0, dir.x).normalize_or_zero();
							if elev_axis != Vec3::ZERO {
								*dir = Quat::from_axis_angle(elev_axis, pitch_delta).mul_vec3(*dir).normalize();
							}
						}
					}
				}
				lights_changed = true;
			}
		}

		let screenshot_requested = window.is_key_pressed(WindowKey::F);

		if camera_changed || lights_changed || needs_redraw || screenshot_requested {
			world.camera.set_position(camera_position);
			world.camera.set_look_at(camera_position + forward);

			if camera_changed || lights_changed || needs_redraw {
				let framebuffer = world.render(render_width as u32, render_height as u32);
				let pixels = upscale_framebuffer_to_window_pixels(
					&framebuffer.pixels,
					render_width,
					render_height,
					width,
					height,
				);
				window.set_pixels(&pixels);
				needs_redraw = false;
			}

			if screenshot_requested {
				match save_screenshot(&world, width as u32, height as u32) {
					Ok((ppm_path, png_path)) => {
						println!("saved screenshot {} and {}", ppm_path.display(), png_path.display());
					}
					Err(error) => {
						eprintln!("failed to save screenshot: {error}");
					}
				}
			}
		}

		if let Err(error) = window.present() {
			eprintln!("failed to update liveview window: {error}");
			return finish_liveview(&mut window, LiveviewAction::ExitApp);
		}
	}

	finish_liveview(&mut window, LiveviewAction::ExitApp)
}

fn direction_from_angles(yaw: f32, pitch: f32) -> Vec3 {
	let cos_pitch = pitch.cos();
	Vec3::new(cos_pitch * yaw.cos(), pitch.sin(), cos_pitch * yaw.sin()).normalize()
}

fn upscale_framebuffer_to_window_pixels(
	pixels: &[Color],
	src_width: usize,
	src_height: usize,
	dst_width: usize,
	dst_height: usize,
) -> Vec<u32> {
	let mut out = vec![0; dst_width * dst_height];

	for y in 0..dst_height {
		let src_y = (y * src_height / dst_height).min(src_height.saturating_sub(1));
		for x in 0..dst_width {
			let src_x = (x * src_width / dst_width).min(src_width.saturating_sub(1));
			let pixel = pixels[src_y * src_width + src_x];
			let r = color_channel_to_u32(pixel.r);
			let g = color_channel_to_u32(pixel.g);
			let b = color_channel_to_u32(pixel.b);
			out[y * dst_width + x] = (r << 16) | (g << 8) | b;
		}
	}

	out
}

fn color_channel_to_u32(channel: f32) -> u32 {
	(channel.clamp(0.0, 1.0) * 255.0).round() as u32
}

fn clamp_camera_to_floor(world: &World, position: Vec3) -> Vec3 {
	let floor_y = world
		.objects
		.iter()
		.filter_map(|object| object.floor_height_at(position))
		.reduce(f32::max);

	if let Some(floor_y) = floor_y {
		let min_y = floor_y + CAMERA_CLEARANCE;
		return Vec3::new(position.x, position.y.max(min_y), position.z);
	}

	position
}

fn save_screenshot(world: &World, width: u32, height: u32) -> Result<(PathBuf, PathBuf), String> {
	let ppm_path = next_screenshot_ppm_path()?;
	let png_path = world::png_path_for(&ppm_path);
	let framebuffer = world.render(width, height);

	framebuffer
		.save_ppm(&ppm_path)
		.map_err(|error| format!("failed to write {}: {error}", ppm_path.display()))?;
	framebuffer
		.save_png(&png_path)
		.map_err(|error| format!("failed to write {}: {error}", png_path.display()))?;

	Ok((ppm_path, png_path))
}

fn next_screenshot_ppm_path() -> Result<PathBuf, String> {
	let screenshot_dir = PathBuf::from(SCREENSHOT_DIR);
	fs::create_dir_all(&screenshot_dir).map_err(|error| {
		format!(
			"failed to create screenshot directory {}: {error}",
			screenshot_dir.display()
		)
	})?;

	for index in 1.. {
		let ppm_path = screenshot_dir.join(format!("screenshot_{index}.ppm"));
		let png_path = world::png_path_for(&ppm_path);
		if !ppm_path.exists() && !png_path.exists() {
			return Ok(ppm_path);
		}
	}

	Err(String::from("failed to allocate screenshot filename"))
}
