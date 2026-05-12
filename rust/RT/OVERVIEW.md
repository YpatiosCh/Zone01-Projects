# RT — Project Overview

A ray tracer in Rust. It generates a 2D image of a 3D scene by shooting rays from a virtual camera, tracking what each ray hits, and coloring pixels based on materials, lights, and shadows.

## The pieces

- **objects** — The 3D shapes (sphere, cube, plane, cylinder) and the shared building blocks they all use: rays, hit records, colors, and materials. This is the foundation everything else depends on.

- **raycaster** — Walks rays through the scene and bounces them off mirrors. Exposes `trace_with_reflections` (full mirror-bounce chain for a primary ray) and `occluded` (fast shadow-ray test). Used by both `camera` (via the render pipeline) and `lighting` (shadow rays).

- **camera** — Knows where the viewer is and which direction they're looking. Shoots rays into the scene through each pixel. Also handles a "frustum" — the cone-shaped volume the camera can actually see — used to cull and sort objects that aren't on screen.

- **lighting** — Light sources (point lights, directional lights) and the math that decides how bright a hit point should be. Casts shadow rays via `raycaster::Raycaster::occluded` to figure out what's in shadow.

- **image** — A 2D grid of pixels and the code that writes it out to a `.ppm` image file.

- **world** — The glue. Holds the camera, the list of objects, the lights, and the sky. Exposes the main `render()` function that orchestrates everything. This is what users of the library actually touch.

- **scenes** — Pre-built example worlds, written as plain Rust functions. Each scene is one file that builds and returns a `World`. The required deliverables (one sphere, plane+cube, all four shapes, alt camera angle) live here.

- **liveview** — An optional debug window that re-renders the scene every frame so you can move the camera around in real time. Useful for spotting bugs visually.

- **rt** — The command-line program. Picks a scene, sets the resolution, and either writes a PPM file or opens the live viewer.

- **tests/** — End-to-end fuzz tests that generate random scenes at 100×100 to make sure nothing crashes on weird inputs.

## How a render flows

`rt` → picks a scene from `scenes` → builds a `World` → `world.render()` → camera shoots rays → raycaster finds nearest hit (bouncing off mirrors as needed) → lighting shades the hit → image collects pixels → write PPM.
