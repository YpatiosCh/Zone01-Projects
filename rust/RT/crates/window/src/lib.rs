use font8x8::{BASIC_FONTS, UnicodeFonts};
use minifb::{Key, KeyRepeat, MouseButton, MouseMode, Window, WindowOptions};

pub use minifb::Error as WindowError;
pub use minifb::Key as WindowKey;

#[cfg(target_os = "linux")]
use raw_window_handle::{HasDisplayHandle, HasWindowHandle, RawDisplayHandle, RawWindowHandle};

#[cfg(windows)]
use windows_sys::Win32::Foundation::{HWND, POINT, RECT};
#[cfg(windows)]
use windows_sys::Win32::Graphics::Gdi::ClientToScreen;
#[cfg(windows)]
use windows_sys::Win32::Graphics::Gdi::ScreenToClient;
#[cfg(windows)]
use windows_sys::Win32::UI::WindowsAndMessaging::{ClipCursor, GetClientRect, GetCursorPos, SetCursorPos};

#[cfg(target_os = "linux")]
use x11_dl::xlib;

#[derive(Clone, Copy, Debug)]
pub struct Rect {
    pub x: usize,
    pub y: usize,
    pub w: usize,
    pub h: usize,
}

impl Rect {
    pub fn contains(&self, px: usize, py: usize) -> bool {
        px >= self.x && px < self.x + self.w && py >= self.y && py < self.y + self.h
    }
}

pub struct AppWindow {
    window: Window,
    buffer: Vec<u32>,
    width: usize,
    height: usize,
    left_was_down: bool,
    cursor_confined: bool,
    #[cfg(target_os = "linux")]
    xlib: Option<xlib::Xlib>,
}

impl AppWindow {
    pub fn new(title: &str, width: usize, height: usize) -> Result<Self, WindowError> {
        let mut window = Window::new(title, width, height, WindowOptions::default())?;
        window.set_target_fps(60);
        Ok(Self {
            window,
            buffer: vec![0; width * height],
            width,
            height,
            left_was_down: false,
            cursor_confined: false,
            #[cfg(target_os = "linux")]
            xlib: xlib::Xlib::open().ok(),
        })
    }

    pub fn is_open(&self) -> bool {
        self.window.is_open()
    }

    pub fn is_active(&mut self) -> bool {
        self.window.is_active()
    }

    pub fn is_key_pressed(&self, key: Key) -> bool {
        self.window.is_key_pressed(key, KeyRepeat::No)
    }

    pub fn is_key_down(&self, key: Key) -> bool {
        self.window.is_key_down(key)
    }

    pub fn set_cursor_visible(&mut self, visible: bool) {
        self.window.set_cursor_visibility(visible);
    }

    pub fn center_cursor(&mut self) {
        #[cfg(windows)]
        {
            let hwnd = self.window.get_window_handle() as HWND;
            if hwnd.is_null() {
                return;
            }

            let mut rect = RECT {
                left: 0,
                top: 0,
                right: 0,
                bottom: 0,
            };
            if unsafe { GetClientRect(hwnd, &mut rect) } == 0 {
                return;
            }

            let mut center = POINT {
                x: (rect.right - rect.left) / 2,
                y: (rect.bottom - rect.top) / 2,
            };
            if unsafe { ClientToScreen(hwnd, &mut center) } == 0 {
                return;
            }

            unsafe {
                SetCursorPos(center.x, center.y);
            }
        }

        #[cfg(target_os = "linux")]
        {
            let Some(xlib) = &self.xlib else {
                return;
            };
            let Some((display, window)) = self.linux_x11_handles() else {
                return;
            };

            unsafe {
                (xlib.XWarpPointer)(
                    display,
                    0,
                    window,
                    0,
                    0,
                    0,
                    0,
                    (self.width / 2) as i32,
                    (self.height / 2) as i32,
                );
                (xlib.XFlush)(display);
            }
        }
    }

    pub fn mouse_delta_from_center(&mut self) -> Option<(f32, f32)> {
        let (mouse_x, mouse_y) = self.native_mouse_position()?;
        let center_x = self.width as f32 * 0.5;
        let center_y = self.height as f32 * 0.5;
        let dx = mouse_x - center_x;
        let dy = mouse_y - center_y;

        self.center_cursor();
        Some((dx, dy))
    }

    pub fn confine_cursor(&mut self, confined: bool) {
        #[cfg(windows)]
        {
            if confined == self.cursor_confined {
                return;
            }

            if confined {
                let hwnd = self.window.get_window_handle() as HWND;
                if hwnd.is_null() {
                    return;
                }

                let mut rect = RECT {
                    left: 0,
                    top: 0,
                    right: 0,
                    bottom: 0,
                };
                let ok = unsafe { GetClientRect(hwnd, &mut rect) };
                if ok == 0 {
                    return;
                }

                let mut top_left = POINT { x: rect.left, y: rect.top };
                let mut bottom_right = POINT { x: rect.right, y: rect.bottom };
                if unsafe { ClientToScreen(hwnd, &mut top_left) } == 0 {
                    return;
                }
                if unsafe { ClientToScreen(hwnd, &mut bottom_right) } == 0 {
                    return;
                }

                let clip_rect = RECT {
                    left: top_left.x,
                    top: top_left.y,
                    right: bottom_right.x,
                    bottom: bottom_right.y,
                };

                if unsafe { ClipCursor(&clip_rect) } != 0 {
                    self.cursor_confined = true;
                }
            } else if unsafe { ClipCursor(std::ptr::null()) } != 0 {
                self.cursor_confined = false;
            }
        }

        #[cfg(target_os = "linux")]
        {
            if confined == self.cursor_confined {
                return;
            }

            let Some(xlib) = &self.xlib else {
                return;
            };
            let Some((display, window)) = self.linux_x11_handles() else {
                return;
            };

            if confined {
                let grab_result = unsafe {
                    (xlib.XGrabPointer)(
                        display,
                        window,
                        xlib::True,
                        0,
                        xlib::GrabModeAsync,
                        xlib::GrabModeAsync,
                        window,
                        0,
                        xlib::CurrentTime,
                    )
                };

                if grab_result == xlib::GrabSuccess {
                    self.cursor_confined = true;
                    self.center_cursor();
                }
            } else {
                unsafe {
                    (xlib.XUngrabPointer)(display, xlib::CurrentTime);
                    (xlib.XFlush)(display);
                }
                self.cursor_confined = false;
            }
        }

        #[cfg(not(any(windows, target_os = "linux")))]
        {
            self.cursor_confined = confined;
        }
    }

    pub fn clear(&mut self, color: u32) {
        self.buffer.fill(color);
    }

    pub fn fill_rect(&mut self, rect: Rect, color: u32) {
        let x_end = (rect.x + rect.w).min(self.width);
        let y_end = (rect.y + rect.h).min(self.height);
        for y in rect.y.min(self.height)..y_end {
            let row_start = y * self.width;
            for x in rect.x.min(self.width)..x_end {
                self.buffer[row_start + x] = color;
            }
        }
    }

    pub fn stroke_rect(&mut self, rect: Rect, color: u32) {
        if rect.w == 0 || rect.h == 0 {
            return;
        }

        self.fill_rect(Rect { x: rect.x, y: rect.y, w: rect.w, h: 1 }, color);
        self.fill_rect(Rect { x: rect.x, y: rect.y + rect.h.saturating_sub(1), w: rect.w, h: 1 }, color);
        self.fill_rect(Rect { x: rect.x, y: rect.y, w: 1, h: rect.h }, color);
        self.fill_rect(Rect { x: rect.x + rect.w.saturating_sub(1), y: rect.y, w: 1, h: rect.h }, color);
    }

    pub fn draw_text(&mut self, x: usize, y: usize, text: &str, color: u32, scale: usize) {
        if scale == 0 {
            return;
        }

        for (index, ch) in text.chars().enumerate() {
            if let Some(glyph) = BASIC_FONTS.get(ch) {
                let glyph_x = x + index * 8 * scale;
                for (row, bits) in glyph.iter().enumerate() {
                    for col in 0..8 {
                        if ((bits >> col) & 1) == 1 {
                            let pixel_x = glyph_x + col * scale;
                            let pixel_y = y + row * scale;
                            self.fill_rect(Rect { x: pixel_x, y: pixel_y, w: scale, h: scale }, color);
                        }
                    }
                }
            }
        }
    }

    pub fn draw_text_centered(&mut self, rect: Rect, text: &str, color: u32, scale: usize) {
        let (text_width, text_height) = measure_text(text, scale);
        let x = rect.x + rect.w.saturating_sub(text_width) / 2;
        let y = rect.y + rect.h.saturating_sub(text_height) / 2;
        self.draw_text(x, y, text, color, scale);
    }

    pub fn mouse_position(&self) -> Option<(usize, usize)> {
        let (x, y) = self.window.get_mouse_pos(MouseMode::Clamp)?;
        if x.is_sign_negative() || y.is_sign_negative() {
            return None;
        }
        Some((x as usize, y as usize))
    }

    pub fn take_left_click(&mut self) -> Option<(usize, usize)> {
        let is_down = self.window.get_mouse_down(MouseButton::Left);
        let click = is_down && !self.left_was_down;
        self.left_was_down = is_down;
        if click {
            self.mouse_position()
        } else {
            None
        }
    }

    pub fn set_pixels(&mut self, pixels: &[u32]) {
        assert_eq!(pixels.len(), self.buffer.len(), "pixel buffer size must match the window size");
        self.buffer.copy_from_slice(pixels);
    }

    pub fn present(&mut self) -> Result<(), WindowError> {
        self.window.update_with_buffer(&self.buffer, self.width, self.height)
    }
}

impl AppWindow {
    fn native_mouse_position(&self) -> Option<(f32, f32)> {
        #[cfg(windows)]
        {
            let hwnd = self.window.get_window_handle() as HWND;
            if hwnd.is_null() {
                return None;
            }

            let mut point = POINT { x: 0, y: 0 };
            if unsafe { GetCursorPos(&mut point) } == 0 {
                return None;
            }
            if unsafe { ScreenToClient(hwnd, &mut point) } == 0 {
                return None;
            }

            Some((point.x as f32, point.y as f32))
        }

        #[cfg(target_os = "linux")]
        {
            let Some(xlib) = &self.xlib else {
                return None;
            };
            let Some((display, window)) = self.linux_x11_handles() else {
                return None;
            };

            let mut root: xlib::Window = 0;
            let mut root_x: i32 = 0;
            let mut root_y: i32 = 0;
            let mut child: xlib::Window = 0;
            let mut child_x: i32 = 0;
            let mut child_y: i32 = 0;
            let mut mask: u32 = 0;

            let inside = unsafe {
                (xlib.XQueryPointer)(
                    display,
                    window,
                    &mut root,
                    &mut child,
                    &mut root_x,
                    &mut root_y,
                    &mut child_x,
                    &mut child_y,
                    &mut mask,
                )
            } != xlib::False;

            if inside {
                Some((child_x as f32, child_y as f32))
            } else {
                None
            }
        }

        #[cfg(not(any(windows, target_os = "linux")))]
        {
            self.window.get_mouse_pos(MouseMode::Discard)
        }
    }
}

#[cfg(target_os = "linux")]
impl AppWindow {
    fn linux_x11_handles(&self) -> Option<(*mut xlib::Display, xlib::Window)> {
        let display_handle = self.window.display_handle().ok()?.as_raw();
        let window_handle = self.window.window_handle().ok()?.as_raw();

        let RawDisplayHandle::Xlib(display) = display_handle else {
            return None;
        };
        let RawWindowHandle::Xlib(window) = window_handle else {
            return None;
        };

        let display_ptr = display.display?.as_ptr() as *mut xlib::Display;
        Some((display_ptr, window.window))
    }
}

impl Drop for AppWindow {
    fn drop(&mut self) {
        if self.cursor_confined {
            self.confine_cursor(false);
        }
    }
}

pub fn measure_text(text: &str, scale: usize) -> (usize, usize) {
    (text.chars().count() * 8 * scale, 8 * scale)
}
