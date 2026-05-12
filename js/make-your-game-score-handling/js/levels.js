// Level data and brick management system
// Contains all level definitions, story data, and brick creation logic

// Brick layout constants
const BRICK_WIDTH = 80;
const BRICK_HEIGHT = 30;
const BRICK_SPACING_X = 85;
const BRICK_SPACING_Y = 35;
const BRICK_START_X = 10;  
const BRICK_START_Y = 50;

// Tile type definitions for level grids
const TILES = {
    EMPTY: 0,
    WEAK: 1,     // Light blue brick, 1 hit, 10 points
    NORMAL: 2    // Orange brick, 1 hit, 10 points  
};

// Pre-calculated tile configurations for performance optimization
const TILE_CONFIG = {
    [TILES.WEAK]: { 
        color: '#87CEEB', 
        health: 1, 
        points: 10,
        styleString: 'background-color: #87CEEB; border: 1px solid rgba(255, 255, 255, 0.3);'
    },
    [TILES.NORMAL]: { 
        color: '#FF8C00', 
        health: 1, 
        points: 20,
        styleString: 'background-color: #FF8C00; border: 1px solid rgba(255, 255, 255, 0.3);'
    }
};

// Level definitions with story and brick layout data
const LEVELS = {
    1: {
        story: {
            title: "Sector 1: The Outer Rim",
            text: `Captain, we've detected the first barrier field ahead. These energy blocks are weakly shielded - your plasma cannon should break through easily.

Mission Objective: Clear the basic energy barriers blocking our path. This is a simple warm-up before we encounter heavier defenses.

Good luck, pilot. The ship's survival depends on you.`
        },
        grid: [
            [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
            [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
            [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
            [0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0],
            [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]
        ]
    },
    
    2: {
        story: {
            title: "Sector 2: The Defense Grid",
            text: `Well done on the outer barriers, Captain. Now we're facing a more organized defense pattern.

Intelligence reports show reinforced energy blocks mixed with standard barriers. The orange blocks have enhanced shielding but remain vulnerable to direct hits.

Stay sharp - the enemy is learning our tactics.`
        },
        grid: [
            [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
            [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
            [0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0],
            [0, 0, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 0, 0],
            [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]
        ]
    },
    
    3: {
        story: {
            title: "Sector 3: Heavy Fortifications",
            text: `Captain, we're approaching the enemy's heavy fortifications. Sensors detect red-class barrier blocks with double-layer shielding.

Warning: These crimson barriers require two direct hits to breach. Conserve your ammunition and maintain precise targeting.

The enemy is getting desperate. Push through!`
        },
        grid: [
            [0, 0, 0, 0, 0, 0, 2, 2, 0, 0, 0, 0, 0, 0],
            [0, 0, 0, 0, 0, 2, 2, 2, 2, 0, 0, 0, 0, 0],
            [0, 0, 0, 0, 2, 2, 2, 2, 2, 2, 0, 0, 0, 0],
            [0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0],
            [0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0]
        ]

    },
    
    4: {
        story: {
            title: "Sector 4: The Core Approach",
            text: `Excellent work, Captain! We're now approaching the enemy's core defense systems. This formation shows tactical sophistication.

The pyramid structure ahead contains their most advanced barrier technology. Expect mixed resistance levels and concentrated firepower positions.

We're almost through to their command center. Stay focused!`
        },
        grid: [
            [0, 0, 2, 0, 2, 0, 2, 2, 0, 2, 0, 2, 0, 0],
            [0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0],
            [2, 2, 2, 2, 0, 0, 2, 2, 0, 0, 2, 2, 2, 2],
            [1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1],
            [2, 0, 0, 0, 2, 2, 2, 2, 2, 2, 0, 0, 0, 2]
        ]
    },
    
    5: {
        story: {
            title: "Sector 5: Final Barrier Matrix",
            text: `This is it, Captain! The final barrier matrix protecting their command center. 

The enemy has deployed their most complex defense pattern - a checkerboard formation designed to confuse targeting systems. Every shot must count.

Break through this final defense and we can escape the Cosmic Brickfield. The entire galaxy believes in you!`
        },
        grid: [
            [2, 1, 2, 1, 2, 1, 2, 2, 1, 2, 1, 2, 1, 2],
            [1, 2, 1, 2, 0, 2, 1, 1, 2, 0, 2, 1, 2, 1],
            [2, 1, 0, 1, 2, 1, 2, 2, 1, 2, 1, 0, 1, 2],
            [1, 2, 1, 2, 1, 0, 1, 1, 0, 1, 2, 1, 2, 1],
            [2, 1, 2, 1, 2, 1, 0, 0, 1, 2, 1, 2, 1, 2]
        ]

    }
};

// Load level data and create brick elements
function loadLevel(levelNumber) {
    const levelData = LEVELS[levelNumber];
    if (!levelData) {
        console.error(`Level ${levelNumber} not found`);
        return null;
    }
    
    // Clear existing bricks with fast innerHTML clearing
    const bricksContainer = document.getElementById('bricksContainer');
    bricksContainer.innerHTML = '';
    
    const bricks = [];
    let activeBricks = 0;
    
    // Use document fragment for batched DOM creation
    const fragment = document.createDocumentFragment();
    
    // Create bricks from level grid data
    for (let row = 0; row < levelData.grid.length; row++) {
        for (let col = 0; col < levelData.grid[row].length; col++) {
            const tileType = levelData.grid[row][col];
            
            if (tileType !== TILES.EMPTY) {
                const tile = createTile(tileType, col, row);
                if (tile) {
                    bricks.push(tile);
                    fragment.appendChild(tile.element);
                    activeBricks++;
                }
            }
        }
    }
    
    // Single DOM append for better performance
    bricksContainer.appendChild(fragment);
    
    console.log(`Level ${levelNumber} loaded: ${activeBricks} bricks`);
    
    return {
        number: levelNumber,
        story: levelData.story,
        bricks: bricks,
        activeBricks: activeBricks
    };
}

// Create individual brick element with optimized DOM operations
function createTile(type, col, row) {
    const config = TILE_CONFIG[type];
    if (!config) return null;
    
    // Calculate brick position on screen
    const x = BRICK_START_X + (col * BRICK_SPACING_X);
    const y = BRICK_START_Y + (row * BRICK_SPACING_Y);
    
    // Create brick data object
    const tile = {
        type: type,
        x: x,
        y: y,
        width: BRICK_WIDTH,
        height: BRICK_HEIGHT,
        health: config.health,
        maxHealth: config.health,
        points: config.points,
        destroyed: false,
        element: null
    };
    
    // Create DOM element with optimized styling
    const element = document.createElement('div');
    element.className = 'brick';
    
    // Single style assignment using cssText for performance
    element.style.cssText = `
        position: absolute;
        left: ${x}px;
        top: ${y}px;
        width: ${BRICK_WIDTH}px;
        height: ${BRICK_HEIGHT}px;
        ${config.styleString}
        border-radius: 4px;
        contain: layout paint;
        backface-visibility: hidden;
        will-change: opacity;
    `;
    
    tile.element = element;
    return tile;
}

// Get story data for specific level
function getLevelStory(levelNumber) {
    const levelData = LEVELS[levelNumber];
    return levelData ? levelData.story : null;
}

// Check if level exists
function hasLevel(levelNumber) {
    return LEVELS.hasOwnProperty(levelNumber);
}

// Get total number of available levels
function getTotalLevels() {
    return Object.keys(LEVELS).length;
}

// Get tile configuration data
function getTileConfig(tileType) {
    return TILE_CONFIG[tileType];
}