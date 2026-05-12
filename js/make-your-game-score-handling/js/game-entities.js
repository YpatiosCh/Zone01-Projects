// Game entities logic - ball, paddle, and bricks
// Contains all movement, physics, and collision logic for game objects

// Update paddle position and movement
function updatePaddle(deltaTime) {
    let moved = false;
    
    // Check left arrow movement
    if (isKeyPressed('ArrowLeft')) {
        game.paddle.x -= game.paddle.speed * deltaTime;
        moved = true;
    }
    // Check right arrow movement
    if (isKeyPressed('ArrowRight')) {
        game.paddle.x += game.paddle.speed * deltaTime;
        moved = true;
    }
    
    // Keep paddle within game bounds
    if (game.paddle.x < 0) game.paddle.x = 0;
    if (game.paddle.x > GAME_WIDTH - game.paddle.width) {
        game.paddle.x = GAME_WIDTH - game.paddle.width;
    }
    
    // Move stuck ball with paddle
    if (game.ball.stuck && moved) {
        game.ball.x = game.paddle.x + game.paddle.width / 2;
    }
}

// Update ball position and physics
function updateBall(deltaTime) {
    // Handle stuck ball state
    if (game.ball.stuck) {
        game.ball.x = game.paddle.x + game.paddle.width / 2;
        game.ball.y = game.paddle.y - game.ball.radius - 5;
        return;
    }
    
    // Update ball position based on velocity
    game.ball.x += game.ball.vx * deltaTime;
    game.ball.y += game.ball.vy * deltaTime;
}

// Launch ball from paddle with random angle
function launchBall() {
    const angle = -Math.PI / 2 + (Math.random() - 0.5) * 0.5;
    game.ball.vx = Math.cos(angle) * game.ball.speed;
    game.ball.vy = Math.sin(angle) * game.ball.speed;
    game.ball.stuck = false;
}

// Reset ball to paddle position
function resetBall() {
    game.ball.x = game.paddle.x + game.paddle.width / 2;
    game.ball.y = game.paddle.y - game.ball.radius - 5;
    game.ball.vx = 0;
    game.ball.vy = 0;
    game.ball.stuck = true;
}

// Main collision detection system optimized for performance
function checkCollisions() {
    if (game.ball.stuck) return;
    
    // Wall collision detection and response
    if (game.ball.x - game.ball.radius <= 0) {
        game.ball.x = game.ball.radius;
        game.ball.vx = Math.abs(game.ball.vx);
    }
    if (game.ball.x + game.ball.radius >= GAME_WIDTH) {
        game.ball.x = GAME_WIDTH - game.ball.radius;
        game.ball.vx = -Math.abs(game.ball.vx);
    }
    if (game.ball.y - game.ball.radius <= 0) {
        game.ball.y = game.ball.radius;
        game.ball.vy = Math.abs(game.ball.vy);
    }
    
    // Bottom wall collision results in ball lost
    if (game.ball.y + game.ball.radius > GAME_HEIGHT) {
        ballLost();
        return;
    }
    
    // Paddle collision detection and angle calculation
    if (game.ball.vy > 0 && 
        game.ball.x + game.ball.radius >= game.paddle.x &&
        game.ball.x - game.ball.radius <= game.paddle.x + game.paddle.width &&
        game.ball.y + game.ball.radius >= game.paddle.y &&
        game.ball.y - game.ball.radius <= game.paddle.y + game.paddle.height) {
        
        // Position ball above paddle
        game.ball.y = game.paddle.y - game.ball.radius;
        
        // Calculate bounce angle based on hit position
        const hitPos = (game.ball.x - (game.paddle.x + game.paddle.width / 2)) / (game.paddle.width / 2);
        const maxAngle = Math.PI / 3;
        const angle = hitPos * maxAngle;
        
        // Set new ball velocity based on angle
        game.ball.vx = Math.sin(angle) * game.ball.speed;
        game.ball.vy = -Math.abs(Math.cos(angle) * game.ball.speed);
    }
    
    // Brick collision detection using optimized active bricks list
    for (let i = game.activeBricksList.length - 1; i >= 0; i--) {
        const brick = game.activeBricksList[i];
        
        // Safety check for destroyed bricks
        if (brick.destroyed) {
            game.activeBricksList.splice(i, 1);
            continue;
        }
        
        // AABB collision detection
        if (game.ball.x + game.ball.radius >= brick.x &&
            game.ball.x - game.ball.radius <= brick.x + brick.width &&
            game.ball.y + game.ball.radius >= brick.y &&
            game.ball.y - game.ball.radius <= brick.y + brick.height) {
            
            // Reverse ball Y velocity for simple bounce
            game.ball.vy = -game.ball.vy;
            
            // Mark brick as destroyed
            brick.destroyed = true;
            
            // Remove from active collision list immediately
            game.activeBricksList.splice(i, 1);
            
            // Provide immediate visual feedback without DOM removal
            if (brick.element) {
                brick.element.style.opacity = '0';
                brick.element.style.pointerEvents = 'none';
                bricksToRemove.push(brick.element);
            }
            
            // Update score and active brick count
            game.score += brick.points;
            game.activeBricks--;
            markHudForUpdate('score');
            
            // Only process one brick collision per frame
            break;
        }
    }
}

// Rebuild active bricks list for optimized collision detection
function rebuildActiveBricksList() {
    game.activeBricksList = [];
    for (let i = 0; i < game.bricks.length; i++) {
        const brick = game.bricks[i];
        if (!brick.destroyed) {
            game.activeBricksList.push(brick);
        }
    }
}

// Handle ball lost event
function ballLost() {
    game.lives--;
    markHudForUpdate('lives');
    
    // Check for game over condition
    if (game.lives <= 0) {
        gameOver();
    } else {
        resetBall();
    }
}