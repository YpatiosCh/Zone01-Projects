use window::{measure_text, AppWindow, Rect, WindowError, WindowKey};

pub enum MenuAction {
    Start,
    Exit,
}

pub enum SceneMenuAction {
    Scene(&'static str),
    ExitToMainMenu,
}

// (scene_name, button_label, fill_color, text_color)
// To add a scene to the menu, add one line here.
static SCENES: &[(&str, &str, u32, u32)] = &[
    ("sphere",            "Sphere",           0x5FBF77, 0x0F1A14),
    ("cylinder",          "Cylinder",          0x6CB4D9, 0x0A1720),
    ("plane",             "Plane",             0xB7B7B7, 0x171717),
    ("plane_four_spheres","Plane + 4 Spheres", 0x9D8DF1, 0x111022),
    ("plane_reflective",  "Plane Reflective",  0xE9C46A, 0x241A05),
    ("plane_cube",        "Plane + Cube",      0xE07B54, 0x2A0F09),
    ("all_objects",         "All Objects",         0x5FBF77, 0x0F1A14),
    ("all_objects_alt_cam", "All Objects (Alt Cam)", 0x3DA88A, 0x071410),
];

pub fn run(width: u32, height: u32) -> Result<MenuAction, WindowError> {
    let width = width.max(640) as usize;
    let height = height.max(480) as usize;
    let mut window = AppWindow::new("RayTracing", width, height)?;
    let mut show_credits = false;

    while window.is_open() {
        if show_credits {
            draw_credits(&mut window, width, height);
            if window.is_key_pressed(WindowKey::Escape) {
                show_credits = false;
            }
            if let Some((mx, my)) = window.take_left_click() {
                let back = footer_button(width, height, 160, 56);
                if back.contains(mx, my) {
                    show_credits = false;
                }
            }
        } else {
            let buttons = draw_main_menu(&mut window, width, height);
            if window.is_key_pressed(WindowKey::Escape) {
                return Ok(MenuAction::Exit);
            }
            if let Some((mx, my)) = window.take_left_click() {
                if buttons.0.contains(mx, my) {
                    return Ok(MenuAction::Start);
                }
                if buttons.1.contains(mx, my) {
                    show_credits = true;
                }
                if buttons.2.contains(mx, my) {
                    return Ok(MenuAction::Exit);
                }
            }
        }

        window.present()?;
    }

    Ok(MenuAction::Exit)
}

pub fn run_scene_menu(width: u32, height: u32) -> Result<SceneMenuAction, WindowError> {
    let width = width.max(640) as usize;
    let height = height.max(480) as usize;
    let mut window = AppWindow::new("Select Scene", width, height)?;

    while window.is_open() {
        let (scene_buttons, exit_button) = draw_scene_menu(&mut window, width, height);
        if window.is_key_pressed(WindowKey::Escape) {
            return Ok(SceneMenuAction::ExitToMainMenu);
        }
        if let Some((mx, my)) = window.take_left_click() {
            for (rect, name) in &scene_buttons {
                if rect.contains(mx, my) {
                    return Ok(SceneMenuAction::Scene(name));
                }
            }
            if exit_button.contains(mx, my) {
                return Ok(SceneMenuAction::ExitToMainMenu);
            }
        }

        window.present()?;
    }

    Ok(SceneMenuAction::ExitToMainMenu)
}

fn draw_main_menu(window: &mut AppWindow, width: usize, _height: usize) -> (Rect, Rect, Rect) {
    window.clear(0x101820);

    let title_box = Rect {
        x: width / 2 - 220,
        y: 70,
        w: 440,
        h: 80,
    };
    window.draw_text_centered(title_box, "RAY TRACING", 0xF4EBD0, 4);

    let subtitle = Rect {
        x: width / 2 - 180,
        y: 160,
        w: 360,
        h: 32,
    };
    window.draw_text_centered(subtitle, "CPU renderer project", 0x9FB3C8, 2);

    let start = centered_button(width, 250, 260, 60);
    let credits = centered_button(width, 330, 260, 60);
    let exit = centered_button(width, 410, 260, 60);

    draw_button(window, start, "Start", 0x5FBF77, 0x0F1A14);
    draw_button(window, credits, "Credits", 0xE9C46A, 0x241A05);
    draw_button(window, exit, "Exit", 0xE76F51, 0x2A0F09);

    (start, credits, exit)
}

fn draw_credits(window: &mut AppWindow, width: usize, height: usize) {
    window.clear(0x131A22);

    let title = Rect {
        x: width / 2 - 180,
        y: 70,
        w: 360,
        h: 60,
    };
    window.draw_text_centered(title, "CREDITS", 0xF4EBD0, 4);

    let lines = [
        "RayTracing workspace",
        "Menu scaffold: menu crate",
        "Window scaffold: window crate",
        "Liveview will be wired next",
    ];

    let mut y = 180;
    for line in lines {
        let (text_w, _) = measure_text(line, 2);
        let x = width / 2 - text_w / 2;
        window.draw_text(x, y, line, 0xC7D3E0, 2);
        y += 42;
    }

    let back = footer_button(width, height, 160, 56);
    draw_button(window, back, "Back", 0x6C8EAD, 0x0D1620);
}

fn draw_scene_menu(
    window: &mut AppWindow,
    width: usize,
    height: usize,
) -> (Vec<(Rect, &'static str)>, Rect) {
    window.clear(0x111923);

    let title_box = Rect {
        x: width / 2 - 180,
        y: 30,
        w: 360,
        h: 56,
    };
    window.draw_text_centered(title_box, "SCENES", 0xF4EBD0, 4);

    let subtitle = Rect {
        x: width / 2 - 230,
        y: 93,
        w: 460,
        h: 24,
    };
    window.draw_text_centered(subtitle, "Choose the next liveview scene", 0xA9BACD, 2);

    let btn_w = 360;
    let btn_h = 40;
    let start_y = 133;
    let step = 46;

    let mut scene_buttons = Vec::new();
    for (i, (name, label, fill, text)) in SCENES.iter().enumerate() {
        let rect = centered_button(width, start_y + i * step, btn_w, btn_h);
        draw_button(window, rect, label, *fill, *text);
        scene_buttons.push((rect, *name));
    }

    let exit = footer_button(width, height, 200, 48);
    draw_button(window, exit, "Exit", 0xE76F51, 0x2A0F09);

    (scene_buttons, exit)
}

fn draw_button(window: &mut AppWindow, rect: Rect, label: &str, fill: u32, text: u32) {
    window.fill_rect(rect, fill);
    window.stroke_rect(rect, 0xF4EBD0);
    window.draw_text_centered(rect, label, text, 2);
}

fn centered_button(width: usize, y: usize, w: usize, h: usize) -> Rect {
    Rect {
        x: width / 2 - w / 2,
        y,
        w,
        h,
    }
}

fn footer_button(width: usize, height: usize, w: usize, h: usize) -> Rect {
    Rect {
        x: width / 2 - w / 2,
        y: height.saturating_sub(h + 48),
        w,
        h,
    }
}
