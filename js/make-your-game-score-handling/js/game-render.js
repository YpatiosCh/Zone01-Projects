// Rendering and HUD update system
// Handles all DOM updates and visual feedback with performance optimization

// DOM elements cache for fast access
const elements = {};

// HUD cache system with dirty flags for batched updates
const hudCache = {
    score: -1,
    lives: -1,
    level: -1,
    gameTime: -1,
    ballTransform: '',
    paddleTransform: '',
    hudNeedsUpdate: false,
    hudUpdatesPending: new Set()
};

// HUD update batching configuration
let hudUpdateTimer = 0;
const HUD_UPDATE_INTERVAL = 100; // Update HUD every 100ms instead of every frame

// Initialize and cache all DOM elements - UPDATED
function initElements() {
    // List of required DOM element IDs - UPDATED for unified modal system
    const elementIds = [
        'playerSection', 'gameContainer', 'playerNameInput', 'highscoresList',
        'playername', 'ball', 'paddle', 'bricksContainer', 'score', 'lives',
        'level', 'gameTime', 'gameMessages', 
        // UNIFIED MODAL ELEMENTS - ADDED
        'gameModal', 'modalTitle', 'modalText', 'modalContent', 'modalOptions', 'modalHint'
    ];
    
    // Cache each element and check for missing ones
    for (const id of elementIds) {
        const element = document.getElementById(id);
        if (!element) {
            console.error(`Missing element: ${id}`);
            return false;
        }
        elements[id] = element;
    }
    
    return true;
}

// STATE-BASED RENDERING SYSTEM - renders different content based on game state - UPDATED
function renderByState(currentTime) {
    switch(game.state) {
        case STATES.MENU:
            renderMainMenu(currentTime);
            break;
            
        case STATES.SCENARIO:
        case STATES.PAUSED:
        case STATES.LEVEL_COMPLETE:
        case STATES.GAME_OVER:
        case STATES.ALL_LEVELS_COMPLETE: // NEW state
            renderGameplay(currentTime); // Show game background
            renderModalOverlay(currentTime); // UNIFIED modal overlay
            break;
            
        case STATES.PLAYING:
            renderGameplay(currentTime);
            break;
    }
}

// Main menu rendering at 60fps
function renderMainMenu(currentTime) {
    // Update HUD with batching for performance
    updateHUD(currentTime);
    
    // Future: Add menu highlight animations here
}

// Gameplay rendering function at 60fps
function renderGameplay(currentTime) {
    // Update HUD with batching for performance
    updateHUD(currentTime);
    
    // Update ball position with GPU acceleration
    const ballX = (game.ball.x - game.ball.radius) | 0;
    const ballY = (game.ball.y - game.ball.radius) | 0;
    const ballTransform = `translate3d(${ballX}px, ${ballY}px, 0)`;
    
    if (ballTransform !== hudCache.ballTransform) {
        hudCache.ballTransform = ballTransform;
        elements.ball.style.transform = ballTransform;
    }
    
    // Update paddle position with GPU acceleration
    const paddleX = game.paddle.x | 0;
    const paddleY = game.paddle.y | 0;
    const paddleTransform = `translate3d(${paddleX}px, ${paddleY}px, 0)`;
    
    if (paddleTransform !== hudCache.paddleTransform) {
        hudCache.paddleTransform = paddleTransform;
        elements.paddle.style.transform = paddleTransform;
    }
}

// UNIFIED modal overlay rendering at 60fps - ADDED
function renderModalOverlay(currentTime) {
    // Modal animations are handled by updateModal() in game-core.js
    // This function is reserved for future modal-specific rendering needs
    
    // Future: Add modal-specific visual effects here
    // For example: particle effects, background dimming variations, etc.
}

// Mark HUD elements for batched updates
function markHudForUpdate(element) {
    hudCache.hudNeedsUpdate = true;
    hudCache.hudUpdatesPending.add(element);
}

// Batched HUD updates for performance optimization
function updateHUD(currentTime) {
    // Only update HUD every 100ms not every frame
    if (currentTime - hudUpdateTimer < HUD_UPDATE_INTERVAL) return;
    
    if (!hudCache.hudNeedsUpdate) return;
    
    // Get pending updates set
    const updates = hudCache.hudUpdatesPending;
    
    // Update score display if changed
    if (updates.has('score') && game.score !== hudCache.score) {
        hudCache.score = game.score;
        elements.score.textContent = `Score: ${game.score}`;
    }
    
    // Update lives display if changed
    if (updates.has('lives') && game.lives !== hudCache.lives) {
        hudCache.lives = game.lives;
        elements.lives.textContent = `Lives: ${game.lives}`;
    }
    
    // Update level display if changed
    if (updates.has('level') && game.level !== hudCache.level) {
        hudCache.level = game.level;
        elements.level.textContent = `Level: ${game.level}`;
    }
    
    // Update game time display if changed
    if (updates.has('gameTime') && game.gameTime !== hudCache.gameTime) {
        hudCache.gameTime = game.gameTime;
        const minutes = Math.floor(game.gameTime / 60);
        const seconds = game.gameTime % 60;
        const timeStr = minutes.toString().padStart(2, '0') + ':' + seconds.toString().padStart(2, '0');
        elements.gameTime.textContent = `Time: ${timeStr}`;
    }
    
    // Clear update flags after processing
    hudCache.hudNeedsUpdate = false;
    hudCache.hudUpdatesPending.clear();
    hudUpdateTimer = currentTime;
}

// Clean up destroyed bricks after frame rendering
function cleanupDestroyedBricks() {
    // Only process if we have bricks queued for removal
    if (bricksToRemove.length === 0) return;
    
    // Batch remove all destroyed brick elements
    for (let i = 0; i < bricksToRemove.length; i++) {
        const element = bricksToRemove[i];
        if (element && element.parentNode) {
            element.parentNode.removeChild(element);
        }
    }
    
    // Clear the removal queue
    bricksToRemove.length = 0;
}

// Display high scores in the main menu
function displayHighScores(scores) {
    if (!scores || scores.length === 0) {
        elements.highscoresList.innerHTML = '<div class="no-scores">No scores yet! Be the first to play!</div>';
        return;
    }
    
    // Build HTML for top 5 scores
    let html = '<div class="top-scores">';
    scores.slice(0, 5).forEach((score, index) => {
        const rank = index + 1;
        html += `
            <div class="score-entry">
                <span class="rank">${rank}</span>
                <span class="name">${score.name}</span>
                <span class="score">${score.score.toLocaleString()}</span>
                <span class="time">${score.time}</span>
            </div>
        `;
    });
    html += '</div>';
    elements.highscoresList.innerHTML = html;
}