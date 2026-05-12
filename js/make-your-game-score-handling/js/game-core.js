// Main game loop and state management
// Contains critical 60fps code that must stay together

// Game constants
const GAME_WIDTH = 1000;  
const GAME_HEIGHT = 800;
const HUD_WIDTH = 200;  

// Game states enumeration - UPDATED with new state
const STATES = {
    MENU: 'menu',
    SCENARIO: 'scenario',
    PLAYING: 'playing',
    PAUSED: 'paused',
    LEVEL_COMPLETE: 'levelComplete',
    GAME_OVER: 'gameOver',
    ALL_LEVELS_COMPLETE: 'allLevelsComplete'  // NEW
};

// Modal configuration constants - ADDED
const MODAL_CONFIGS = {
    SCENARIO: {
        text: "", // Dynamic - will be set from level story data
        optionCount: 1,
        navigation: false,
        action: 'continue',
        hint: 'Press <strong>SPACE</strong> or <strong>ENTER</strong> to begin mission'
    },
    PAUSE: {
        text: "Game Paused",
        optionCount: 3,
        navigation: true,
        options: ['Resume', 'Restart', 'Menu'],
        hint: 'Use <strong>Arrow Left/Right</strong> to select, <strong>SPACE</strong> or <strong>ENTER</strong> to confirm'
    },
    LEVEL_COMPLETE: {
        text: "Level Complete!",
        optionCount: 1,
        navigation: false,
        action: 'continue',
        hint: 'Press <strong>SPACE</strong> or <strong>ENTER</strong> to continue'
    },
    GAME_OVER: {
        text: "Game Over",
        optionCount: 2,
        navigation: true,
        options: ['Play Again', 'Back to Menu'],
        hint: 'Use <strong>Arrow Left/Right</strong> to select, <strong>SPACE</strong> or <strong>ENTER</strong> to confirm'
    },
    ALL_LEVELS_COMPLETE: {
        text: "Victory! All Sectors Cleared!",
        optionCount: 2,
        navigation: true,
        options: ['Play Again', 'Back to Menu'],
        hint: 'Use <strong>Arrow Left/Right</strong> to select, <strong>SPACE</strong> or <strong>ENTER</strong> to confirm'
    }
};

// Main game state object
const game = {
    state: STATES.MENU,
    score: 0,
    lives: 3,
    level: 1,
    gameTime: 0,           // Current seconds (no longer calculated)
    timerInterval: null,   // Holds the setInterval ID
    playerName: '',
    
    // Ball state properties
    ball: {
        x: GAME_WIDTH / 2,
        y: GAME_HEIGHT - 60,
        vx: 0,
        vy: 0,
        radius: 8,
        speed: 400,
        stuck: true
    },
    
    // Paddle state properties
    paddle: {
        x: (GAME_WIDTH - 120) / 2,
        y: GAME_HEIGHT - 30,
        width: 120,
        height: 20,
        speed: 600
    },
    
    // Brick management arrays for performance optimization
    bricks: [],           // All bricks reference array
    activeBricks: 0,      // Count of non-destroyed bricks
    activeBricksList: [], // Only non-destroyed bricks for collision detection
    
    // Menu navigation state
    selectedOption: 0,
    
    // Current modal state - ADDED
    currentModal: null
};

// Game loop timing variables
let gameRunning = false;
let animationId = null;
let lastTime = 0;

// Array to track bricks that need DOM removal after frame
let bricksToRemove = [];

// Simple FPS counter variables
let frameCount = 0;
let fpsStartTime = 0;

function startTimer() {
    if (game.timerInterval) return; // prevent multiple intervals
    game.timerInterval = setInterval(() => {
        game.gameTime++;
        markHudForUpdate('gameTime');
    }, 1000);
}

function pauseTimer() {
    if (game.timerInterval) {
        clearInterval(game.timerInterval);
        game.timerInterval = null;
    }
}

function resetTimer() {
    pauseTimer();
    game.gameTime = 0;
    markHudForUpdate('gameTime');
}

// UNIFIED 60FPS GAME LOOP - runs continuously for ALL states
function gameLoop(currentTime) {
    if (!gameRunning) return;
    
    // Cap delta time to prevent large time jumps
    const deltaTime = Math.min((currentTime - lastTime) / 1000, 1/30);
    lastTime = currentTime;
    
    // Update based on current game state
    updateByState(deltaTime);
    
    // Render based on current game state  
    renderByState(currentTime);
    
    // Clean up destroyed bricks outside critical path
    cleanupDestroyedBricks();
    
    // Request next frame - ALWAYS continues regardless of state
    animationId = requestAnimationFrame(gameLoop);
}

// State-based update system for 60fps across all states - UPDATED
function updateByState(deltaTime) {
    switch(game.state) {
        case STATES.MENU:
            updateMainMenu(deltaTime);
            break;
            
        case STATES.SCENARIO:
        case STATES.PAUSED:
        case STATES.LEVEL_COMPLETE:
        case STATES.GAME_OVER:
        case STATES.ALL_LEVELS_COMPLETE:
            updateModal(deltaTime);  // UNIFIED modal update
            break;
            
        case STATES.PLAYING:
            updateGameplay(deltaTime);
            break;
    }
}

// Main menu update function for 60fps animations
function updateMainMenu(deltaTime) {
    const time = performance.now() * 0.001; // Convert to seconds
    
    // Breathing title effect
    const scale = 1 + Math.sin(time * 2) * 0.02; // 2% scale variation
    const title = document.querySelector('h1');
    if (title) {
        title.style.transform = `scale(${scale})`;
    }
    
    // Glow pulse on input
    const glow = 0.5 + Math.sin(time * 3) * 0.3; // Pulse between 0.2-0.8
    elements.playerNameInput.style.boxShadow = `0 0 ${10 + glow * 10}px rgba(0, 170, 255, ${glow})`;
}

// UNIFIED modal update function for 60fps animations - ADDED
function updateModal(deltaTime) {
    const time = performance.now() * 0.001; // Convert to seconds
    
    if (!elements.gameModal) return;
    
    // Gentle breathing effect for all modals
    const scale = 1 + Math.sin(time * 1.5) * 0.01; // 1% scale variation
    elements.gameModal.style.transform = `scale(${scale})`;
    
    // Border glow animation based on modal type
    let glowColor = 'rgba(0, 170, 255, '; // Default blue
    
    if (game.state === STATES.LEVEL_COMPLETE || game.state === STATES.ALL_LEVELS_COMPLETE) {
        glowColor = 'rgba(0, 255, 0, '; // Green for success
    } else if (game.state === STATES.GAME_OVER) {
        glowColor = 'rgba(255, 68, 68, '; // Red for game over
    } else if (game.state === STATES.PAUSED) {
        glowColor = 'rgba(255, 170, 0, '; // Orange for pause
    }
    
    const glow = 0.4 + Math.sin(time * 2.5) * 0.4; // Pulse between 0-0.8
    elements.gameModal.style.borderColor = glowColor + glow + ')';
    elements.gameModal.style.boxShadow = `0 0 ${15 + glow * 15}px ${glowColor}${glow * 0.5})`;
    
    // Selection highlight breathing effect for navigable options
    const config = MODAL_CONFIGS[game.currentModal];
    if (config && config.navigation) {
        const selectedOption = document.querySelector('.modal-option.selected');
        if (selectedOption) {
            const selectionGlow = 0.6 + Math.sin(time * 4) * 0.3;
            selectedOption.style.boxShadow = `0 0 ${15 + selectionGlow * 10}px rgba(0, 170, 255, ${selectionGlow})`;
        }
    }
}

// Gameplay-specific update function
function updateGameplay(deltaTime) {
    const time = performance.now() * 0.001; // Convert to seconds
    
    // Floating title animation for 60fps visual feedback
    const titleElement = document.getElementById('floatingTitle');
    if (titleElement) {
        // Gentle floating effect
        const floatY = Math.sin(time * 1.5) * 3; // 3px vertical float
        const scale = 1 + Math.sin(time * 2.2) * 0.015; // 1.5% scale breathing
        const glow = 0.6 + Math.sin(time * 3.5) * 0.3; // Glow pulse
        
        titleElement.style.transform = `translateX(-50%) translateY(${floatY}px) scale(${scale})`;
        titleElement.style.textShadow = `0 0 ${8 + glow * 12}px rgba(0, 255, 136, ${glow})`;
        
        // Fade out slightly when ball is launched for less distraction
        const opacity = game.ball.stuck ? 1 : 0.7;
        titleElement.style.opacity = opacity;
    }
    
    // Update all game entities
    updatePaddle(deltaTime);
    updateBall(deltaTime);
    checkCollisions();
    
    // Check win condition
    if (game.activeBricks === 0) {
        levelComplete();
    }
}

// Start the main game loop
function startGameLoop() {
    gameRunning = true;
    lastTime = performance.now();
    gameLoop(lastTime);
}

// Initialize game when DOM is ready
function initGame() {
    if (!initElements()) return false;
    setupInput();
    initThemeSystem();
    game.state = STATES.MENU; 
    startGameLoop();
    loadHighScores(); 
    return true;
}

// Game initialization entry point
document.addEventListener('DOMContentLoaded', () => {
    if (!initGame()) {
        console.error('Game initialization failed');
        document.body.innerHTML = '<div style="color: red; text-align: center; padding: 50px;">Game failed to load. Please refresh the page.</div>';
    }
});