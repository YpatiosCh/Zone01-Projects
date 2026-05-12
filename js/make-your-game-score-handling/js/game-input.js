// Input handling system
// Manages all keyboard input and game control logic

// Input state tracking object
const keys = {};

// Setup all input event listeners
function setupInput() {
    // Handle key press events
    document.addEventListener('keydown', (e) => {
        // Allow normal input field behavior
        if (e.target.tagName === 'INPUT') {
            if (e.code === 'Space' || e.code === 'Enter') {
                e.preventDefault();
                handleGameInput(e.code);
            }
            return;
        }
        
        // Prevent default browser behavior
        e.preventDefault();
        keys[e.code] = true;
        
        // Handle discrete key events only once per press
        if (!keys[e.code + '_pressed']) {
            keys[e.code + '_pressed'] = true;
            
            // Modal states handle all keys as discrete events - UPDATED
            if (game.state === STATES.SCENARIO || 
                game.state === STATES.PAUSED || 
                game.state === STATES.LEVEL_COMPLETE ||
                game.state === STATES.GAME_OVER ||
                game.state === STATES.ALL_LEVELS_COMPLETE) {
                handleGameInput(e.code);
            }
            // Playing state only handles non-movement keys as discrete
            else if (e.code !== 'ArrowLeft' && e.code !== 'ArrowRight') {
                handleGameInput(e.code);
            }
        }
    });
    
    // Handle key release events
    document.addEventListener('keyup', (e) => {
        if (e.target.tagName === 'INPUT') return;
        
        e.preventDefault();
        keys[e.code] = false;
        keys[e.code + '_pressed'] = false;
    });
}

// Main input handler for discrete key events - UPDATED
function handleGameInput(keyCode) {
    switch (game.state) {
        case STATES.MENU:
            if (keyCode === 'Space' || keyCode === 'Enter') {
                startGame();
            }
            break;
            
        case STATES.SCENARIO:
            if (keyCode === 'Space' || keyCode === 'Enter') {
                selectModalOption(); // UNIFIED
            }
            break;
            
        case STATES.PLAYING:
            if (keyCode === 'Space') {
                if (game.ball.stuck) {
                    launchBall();
                } else {
                    pauseGame();
                }
            }
            break;
            
        case STATES.PAUSED:
        case STATES.GAME_OVER:
        case STATES.ALL_LEVELS_COMPLETE: // NEW state
            if (keyCode === 'ArrowLeft' || keyCode === 'ArrowRight') {
                navigateModal(keyCode); // UNIFIED navigation
            } else if (keyCode === 'Space' || keyCode === 'Enter') {
                selectModalOption(); // UNIFIED selection
            }
            break;
            
        case STATES.LEVEL_COMPLETE:
            if (keyCode === 'Space' || keyCode === 'Enter') {
                selectModalOption(); // UNIFIED
            }
            break;
    }
}

// Check if a key is currently pressed
function isKeyPressed(keyCode) {
    return keys[keyCode] === true;
}

// Check if a key was just pressed this frame
function isKeyJustPressed(keyCode) {
    return keys[keyCode + '_pressed'] === true;
}