use uuid::Uuid;
pub type Timestamp = u64; // milliseconds or ticks

#[derive(Clone, Copy, PartialEq)]
pub struct CellPosition {
    pub cell: Cell,
    pub dir: Direction,
}

// TODO: degrees
#[derive(Clone, Copy, PartialEq, Debug)]
pub enum Direction {
    North,
    South,
    East,
    West,
}

impl Direction {
    pub fn index(self) -> i32 {
        match self {
            Direction::North => 0,
            Direction::East => 1,
            Direction::South => 2,
            Direction::West => 3,
        }
    }
}

#[derive(Clone, Copy, PartialEq, Debug)]
pub enum Origin {
    North,
    South,
    East,
    West,
}

pub enum Routes {
    NorthL,
    NorthR,
    NorthS,
    SouthL,
    SouthR,
    SouthS,
    EastL,
    EastR,
    EastS,
    WestL,
    WestR,
    WestS,
}

#[derive(Clone, Copy, PartialEq, Eq, Hash)]
pub struct Cell {
    pub x: i32,
    pub y: i32,
}

pub type ToRender = Vec<(Uuid, RenderPosition)>;

#[derive(Clone, Copy, PartialEq, Debug)]
pub enum Highlight {
    None,
    CloseCall,
    Collision,
}

#[derive(Clone, Copy, PartialEq, Debug)]
pub struct RenderPosition {
    pub angle: f32,
    pub x: f32,
    pub y: f32,
    pub highlight: Highlight,
}
