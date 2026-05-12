//! # Filler bot
//!
//! A Rust implementation of a bot for the 01-edu Filler game engine.
//!
//! On every turn the engine sends the current board (the "Anfield") and a
//! randomly shaped piece on stdin; the bot must reply with the coordinates of
//! a legal placement on stdout. The bot picks its move with a one-ply Voronoi
//! territory heuristic — it counts how many board cells would flip to its own
//! "closer-than-the-opponent" region for each legal placement and keeps the
//! best one. See `algo.md` in the repo root for a full walkthrough.
//!
//! Module layout:
//!
//! - [`filler`]   — the per-turn I/O loop.
//! - [`reader`]   — a `BufRead` wrapper that yields one trimmed line at a time.
//! - [`parser`]   — turns the engine's text blocks into typed values.
//! - [`models`]   — the domain types (`Board`, `Piece`, `Position`, `Cell`, …).
//! - [`strategy`] — legality check, BFS distance maps, scoring, selection.

pub mod filler;
pub mod models;
pub mod parser;
pub mod reader;
pub mod strategy;
