// All menu systems and UI navigation
// Handles main menu, unified modal system, and all game state transitions

// UNIFIED MODAL SYSTEM - ADDED
function showModal(modalType, customData = {}) {
    const config = MODAL_CONFIGS[modalType];
    if (!config) {
        console.error(`Unknown modal type: ${modalType}`);
        return;
    }
    
    // Set current modal for animation system
    game.currentModal = modalType;
    
    // Update modal title
    elements.modalTitle.textContent = customData.title || config.text;
    
    // Update modal text content
    if (customData.content) {
        elements.modalText.style.display = 'none';
        elements.modalContent.style.display = 'block';
        elements.modalContent.innerHTML = `<p>${customData.content.replace(/\n\n/g, '</p><p>')}</p>`;
    } else {
        elements.modalText.style.display = 'block';
        elements.modalContent.style.display = 'none';
        elements.modalText.textContent = config.text;
    }
    
    // Update options based on configuration
    if (config.optionCount === 1) {
        elements.modalOptions.style.display = 'none';
    } else {
        elements.modalOptions.style.display = 'block';
        updateModalOptions(config);
    }
    
    // Update keyboard hint
    elements.modalHint.innerHTML = customData.hint || config.hint;
    
    // Show modal
    elements.gameModal.style.display = 'block';
    
    // Reset selection for navigable modals
    if (config.navigation) {
        game.selectedOption = 0;
        updateModalSelection();
    }
}

// Update modal options based on configuration - ADDED
function updateModalOptions(config) {
    if (!config.options) return;
    
    let optionsHTML = '<div class="modal-options">';
    config.options.forEach((option, index) => {
        const selectedClass = index === 0 ? 'selected' : '';
        optionsHTML += `<div class="modal-option ${selectedClass}">${option}</div>`;
    });
    optionsHTML += '</div>';
    
    elements.modalOptions.innerHTML = optionsHTML;
}

// Update visual selection in modal options - ADDED
function updateModalSelection() {
    const config = MODAL_CONFIGS[game.currentModal];
    if (!config || !config.navigation) return;
    
    const options = document.querySelectorAll('.modal-option');
    options.forEach((option, index) => {
        option.classList.toggle('selected', index === game.selectedOption);
    });
}

// Hide modal - ADDED
function hideModal() {
    elements.gameModal.style.display = 'none';
    game.currentModal = null;
}

// Navigate modal options - UNIFIED
function navigateModal(direction) {
    const config = MODAL_CONFIGS[game.currentModal];
    if (!config || !config.navigation) return;
    
    if (direction === 'ArrowLeft') {
        game.selectedOption = Math.max(0, game.selectedOption - 1);
    } else if (direction === 'ArrowRight') {
        game.selectedOption = Math.min(config.optionCount - 1, game.selectedOption + 1);
    }
    updateModalSelection();
}

// Execute selected modal option - UNIFIED
function selectModalOption() {
    const config = MODAL_CONFIGS[game.currentModal];
    if (!config) return;
    
    switch (game.currentModal) {
        case 'SCENARIO':
            hideScenario();
            break;
            
        case 'PAUSE':
            if (game.selectedOption === 0) {
                resumeGame();
            } else if (game.selectedOption === 1) {
                restartGame();
            } else {
                showMainMenu();
            }
            break;
            
        case 'LEVEL_COMPLETE':
            nextLevel();
            break;
            
        case 'GAME_OVER':
        case 'ALL_LEVELS_COMPLETE':
            if (game.selectedOption === 0) {
                startGame();
            } else {
                showMainMenu();
            }
            break;
    }
}

// Start a new game from main menu
function startGame() {
    const playerName = elements.playerNameInput.value.trim();
    if (!playerName) {
        alert('Please enter your name!');
        return;
    }
    
    // Initialize game state
    game.playerName = playerName;
    game.score = 0;
    game.lives = 3;
    game.level = 1;
    resetTimer();  // Reset timer to 0
    startTimer();  // Start the 1-second interval
    
    // Reset HUD cache and mark for updates
    hudCache.score = -1;
    hudCache.lives = -1;
    hudCache.level = -1;
    hudCache.gameTime = -1;
    markHudForUpdate('score');
    markHudForUpdate('lives');
    markHudForUpdate('level');
    
    // Update UI and show game
    elements.playername.textContent = `Player: ${playerName}`;
    elements.playerSection.style.display = 'none';
    elements.gameContainer.style.display = 'block';
    
    // Start with level scenario
    showScenario(1);
}

// Show main menu screen
function showMainMenu() {
    game.state = STATES.MENU;
    
    // Hide modal and show menu elements
    hideModal();
    clearAllThemes(); 
    elements.playerSection.style.display = 'block';
    elements.gameContainer.style.display = 'none';
    
    // Clear any pending brick removals
    bricksToRemove.length = 0;
    
    // Reset lastSubmitResult for fresh game
    game.lastSubmitResult = null;
    
    // Load and display high scores with new pagination system
    loadHighScores(1); // Always start from page 1
}
// Show level scenario popup - UPDATED
function showScenario(levelNumber) {
    const story = getLevelStory(levelNumber);
    if (story) {
        game.state = STATES.SCENARIO;
        elements.gameContainer.style.display = 'block';
        showModal('SCENARIO', {
            title: story.title,
            content: story.text
        });
    } else {
        // No scenario available start level directly
        startLevel(levelNumber);
    }
}

// Hide scenario and start level - UPDATED
function hideScenario() {
    hideModal();
    
    // Load level data and initialize bricks
    const levelResult = loadLevel(game.level);
    if (levelResult) {
        game.bricks = levelResult.bricks;
        game.activeBricks = levelResult.activeBricks;
        rebuildActiveBricksList();
    }
    
    // Reset ball and start gameplay
    resetBall();
    applyLevelTheme(game.level); 
    game.state = STATES.PLAYING;
    startGameLoop();
}

// Start level without scenario
function startLevel(levelNumber) {
    // Load level data and initialize bricks
    const levelResult = loadLevel(levelNumber);
    if (levelResult) {
        game.bricks = levelResult.bricks;
        game.activeBricks = levelResult.activeBricks;
        rebuildActiveBricksList();
    }
    
    // Reset ball and start gameplay
    resetBall();
    applyLevelTheme(levelNumber);
    game.state = STATES.PLAYING;
    startGameLoop();
}

// Pause game and show pause menu - UPDATED
function pauseGame() {
    pauseTimer();  // Stop the timer
    game.state = STATES.PAUSED;
    showModal('PAUSE');
}

// Resume game from pause - UPDATED
function resumeGame() {
    startTimer();  // Resume the timer
    game.state = STATES.PLAYING;
    hideModal();
}

// Restart game from level 1 with same player name - UPDATED
function restartGame() {
    // Reset game state but keep player name
    game.score = 0;
    game.lives = 3;
    game.level = 1;
    resetTimer();  // Reset timer but don't start yet (scenario will start it)
    
    // Reset HUD cache and mark for updates
    hudCache.score = -1;
    hudCache.lives = -1;
    hudCache.level = -1;
    hudCache.gameTime = -1;
    markHudForUpdate('score');
    markHudForUpdate('lives');
    markHudForUpdate('level');
    markHudForUpdate('gameTime');
    
    // Hide modal and start from level 1
    hideModal();
    clearAllThemes();
    showScenario(1);
}

// Handle level completion - UPDATED
function levelComplete() {
    game.state = STATES.LEVEL_COMPLETE;
    game.score += 500; // Level completion bonus
    markHudForUpdate('score');
    showModal('LEVEL_COMPLETE');
}

// Advance to next level - UPDATED
function nextLevel() {
    game.level++;
    markHudForUpdate('level');
    hideModal();
    
    // Check if all levels completed - UPDATED
    if (game.level > getTotalLevels()) {
        allLevelsComplete();
        return;
    }
    
    // Show next level scenario
    showScenario(game.level);
}

// Handle all levels completion - NEW FUNCTION
async function allLevelsComplete() {
    pauseTimer();  // Stop the timer
    game.state = STATES.ALL_LEVELS_COMPLETE;
    
    // Submit final score and wait for ranking response
    await submitScore();
    
    // Get personalized score message with ranking
    const scoreMessage = getScoreDisplayMessage();
    
    // Show victory screen with personalized ranking message
    const minutes = Math.floor(game.gameTime / 60);
    const seconds = game.gameTime % 60;
    const timeString = `${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`;
    
    showModal('ALL_LEVELS_COMPLETE', {
        content: `🎉 Congratulations, Captain ${game.playerName}! You have successfully escaped the Cosmic Brickfield!<br><br>${scoreMessage}<br>Total Time: ${timeString}<br><br>The galaxy salutes your courage and skill!`
    });
}

// Handle game over state - UPDATED
async function gameOver() {
    pauseTimer();  // Stop the timer
    game.state = STATES.GAME_OVER;
    
    // Submit final score and wait for ranking response
    await submitScore();
    
    // Get personalized score message with ranking
    const scoreMessage = getScoreDisplayMessage();
    
    // Show game over screen with personalized ranking message
    showModal('GAME_OVER', {
        content: `💥 Your ship was destroyed in the Cosmic Brickfield, Captain ${game.playerName}...<br><br>${scoreMessage}<br><br>Try again to escape the galaxy's most dangerous sector!`
    });
}