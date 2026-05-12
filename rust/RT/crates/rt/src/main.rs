use std::env;
use std::fs;
use std::path::PathBuf;
use std::process::ExitCode;

use liveview::LiveviewAction;
use menu::{MenuAction, SceneMenuAction};

const DEFAULT_WIDTH: u32 = 800;
const DEFAULT_HEIGHT: u32 = 600;

enum LaunchMode {
    Menu,
    Live,
    Render,
}

enum LiveLoopExit {
    MainMenu,
    ExitApp,
}

fn main() -> ExitCode {
    let args: Vec<String> = env::args().collect();
    let mut scene_name = String::from("sphere");
    let mut out: Option<PathBuf> = None;
    let mut width: u32 = DEFAULT_WIDTH;
    let mut height: u32 = DEFAULT_HEIGHT;
    let mut launch_mode = LaunchMode::Menu;
    let mut list = false;

    let mut i = 1;
    while i < args.len() {
        match args[i].as_str() {
            "--scene" => {
                let Some(v) = args.get(i + 1) else {
                    eprintln!("--scene requires a value");
                    return ExitCode::from(2);
                };
                scene_name = v.clone();
                i += 2;
            }
            "--out" => {
                let Some(v) = args.get(i + 1) else {
                    eprintln!("--out requires a path");
                    return ExitCode::from(2);
                };
                out = Some(PathBuf::from(v));
                launch_mode = LaunchMode::Render;
                i += 2;
            }
            "--width" => {
                let Some(v) = args.get(i + 1).and_then(|s| s.parse().ok()) else {
                    eprintln!("--width requires a positive integer");
                    return ExitCode::from(2);
                };
                width = v;
                i += 2;
            }
            "--height" => {
                let Some(v) = args.get(i + 1).and_then(|s| s.parse().ok()) else {
                    eprintln!("--height requires a positive integer");
                    return ExitCode::from(2);
                };
                height = v;
                i += 2;
            }
            "--menu" => {
                launch_mode = LaunchMode::Menu;
                i += 1;
            }
            "--live" => {
                launch_mode = LaunchMode::Live;
                i += 1;
            }
            "--render" => {
                launch_mode = LaunchMode::Render;
                i += 1;
            }
            "--list-scenes" => {
                list = true;
                i += 1;
            }
            "--help" | "-h" => {
                print_help();
                return ExitCode::SUCCESS;
            }
            other => {
                eprintln!("unknown argument: {other}");
                return ExitCode::from(2);
            }
        }
    }

    if list {
        for n in scenes::names() {
            println!("{n}");
        }
        return ExitCode::SUCCESS;
    }

    match launch_mode {
        LaunchMode::Menu => match run_main_menu_loop(width, height) {
            Ok(()) => return ExitCode::SUCCESS,
            Err(error) => {
                eprintln!("{error}");
                return ExitCode::from(1);
            }
        },
        LaunchMode::Live => match run_live_loop(&scene_name, width, height) {
            Ok(LiveLoopExit::ExitApp) => return ExitCode::SUCCESS,
            Ok(LiveLoopExit::MainMenu) => match run_main_menu_loop(width, height) {
                Ok(()) => return ExitCode::SUCCESS,
                Err(error) => {
                    eprintln!("{error}");
                    return ExitCode::from(1);
                }
            },
            Err(error) => {
                eprintln!("{error}");
                return ExitCode::from(1);
            }
        },
        LaunchMode::Render => {}
    }

    let Some(builder) = scenes::get(&scene_name) else {
        eprintln!("unknown scene: {scene_name}");
        eprintln!("available: {}", scenes::names().join(", "));
        return ExitCode::from(2);
    };
    let world = builder();

    let out = out.unwrap_or_else(|| PathBuf::from(format!("scenes_out/{scene_name}.ppm")));
    if let Some(parent) = out.parent().filter(|parent| !parent.as_os_str().is_empty()) {
        if let Err(e) = fs::create_dir_all(parent) {
            eprintln!(
                "failed to create output directory {}: {e}",
                parent.display()
            );
            return ExitCode::from(1);
        }
    }
    let fb = world.render(width, height);
    if let Err(e) = fb.save_ppm(&out) {
        eprintln!("failed to write {}: {e}", out.display());
        return ExitCode::from(1);
    }

    let png_out = world::png_path_for(&out);
    if let Err(e) = fb.save_png(&png_out) {
        eprintln!("failed to write {}: {e}", png_out.display());
        return ExitCode::from(1);
    }

    println!(
        "wrote {} and {} ({}x{})",
        out.display(),
        png_out.display(),
        width,
        height
    );
    ExitCode::SUCCESS
}

fn print_help() {
    println!("rt — ray tracer");
    println!();
    println!("Usage: rt [--menu] [--live] [--render] [--scene <name>] [--out <file.ppm>] [--width N] [--height N] [--list-scenes]");
    println!();
    println!("Options:");
    println!("  --scene <name>   Scene to render (default: sphere)");
    println!("  --out <path>     Output PPM file (default: scenes_out/<scene>.ppm)");
    println!("  --width  N       Image width  in pixels (default: {DEFAULT_WIDTH})");
    println!("  --height N       Image height in pixels (default: {DEFAULT_HEIGHT})");
    println!("  --menu           Open the startup menu (default when no mode is given)");
    println!("  --live           Open the liveview window directly");
    println!("  --render         Render to files instead of opening the menu");
    println!("  --list-scenes    Print available scene names");
    println!("  -h, --help       Show this help");
}

fn run_main_menu_loop(width: u32, height: u32) -> Result<(), String> {
    loop {
        match menu::run(width, height)
            .map_err(|error| format!("failed to open menu window: {error}"))?
        {
            MenuAction::Start => {
                let selected_scene = match menu::run_scene_menu(width, height)
                    .map_err(|error| format!("failed to open scene menu window: {error}"))?
                {
                    SceneMenuAction::Scene(name) => Some(name),
                    SceneMenuAction::ExitToMainMenu => None,
                };

                let Some(selected_scene) = selected_scene else {
                    continue;
                };

                match run_live_loop(selected_scene, width, height)? {
                    LiveLoopExit::MainMenu => continue,
                    LiveLoopExit::ExitApp => return Ok(()),
                }
            }
            MenuAction::Exit => return Ok(()),
        }
    }
}

fn run_live_loop(initial_scene: &str, width: u32, height: u32) -> Result<LiveLoopExit, String> {
    let mut scene_name = initial_scene.to_string();

    loop {
        let builder = scene_builder(&scene_name)?;
        match liveview::run(builder(), width, height) {
            LiveviewAction::ExitApp => return Ok(LiveLoopExit::ExitApp),
            LiveviewAction::OpenSceneMenu => {
                match menu::run_scene_menu(width, height)
                    .map_err(|error| format!("failed to open scene menu window: {error}"))?
                {
                    SceneMenuAction::Scene(name) => scene_name = name.to_string(),
                    SceneMenuAction::ExitToMainMenu => return Ok(LiveLoopExit::MainMenu),
                }
            }
        }
    }
}

fn scene_builder(name: &str) -> Result<fn() -> world::World, String> {
    scenes::get(name).ok_or_else(|| {
        format!(
            "unknown scene: {name}\navailable: {}",
            scenes::names().join(", ")
        )
    })
}
