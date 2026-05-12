// Physics utilities and collision detection helpers
// Contains all physics calculations and collision detection functions

// Convert game coordinates to screen coordinates
function gameToScreen(gameX, gameY) {
    return {
        x: gameX + HUD_WIDTH,  // Add HUD offset
        y: gameY
    };
}

// Convert screen coordinates to game coordinates
function screenToGame(screenX, screenY) {
    return {
        x: screenX - HUD_WIDTH,  // Remove HUD offset
        y: screenY
    };
}

// Check collision between circle and rectangle
function circleRectCollision(circleX, circleY, radius, rectX, rectY, rectW, rectH) {
    // Find closest point on rectangle to circle center
    const closestX = Math.max(rectX, Math.min(circleX, rectX + rectW));
    const closestY = Math.max(rectY, Math.min(circleY, rectY + rectH));
    
    // Calculate distance between circle center and closest point
    const distanceX = circleX - closestX;
    const distanceY = circleY - closestY;
    const distanceSquared = (distanceX * distanceX) + (distanceY * distanceY);
    
    return distanceSquared < (radius * radius);
}

// Determine which side of rectangle was hit for better bounce physics
function getCollisionSide(circleX, circleY, rectX, rectY, rectW, rectH) {
    const rectCenterX = rectX + rectW / 2;
    const rectCenterY = rectY + rectH / 2;
    
    const deltaX = circleX - rectCenterX;
    const deltaY = circleY - rectCenterY;
    
    // Determine collision side based on angle
    if (Math.abs(deltaX) > Math.abs(deltaY)) {
        return deltaX > 0 ? 'right' : 'left';
    } else {
        return deltaY > 0 ? 'bottom' : 'top';
    }
}

// Normalize vector to unit length
function normalizeVector(x, y) {
    const length = Math.sqrt(x * x + y * y);
    if (length === 0) return { x: 0, y: 0 };
    return { x: x / length, y: y / length };
}

// Set vector to specific magnitude
function setVectorMagnitude(x, y, magnitude) {
    const normalized = normalizeVector(x, y);
    return {
        x: normalized.x * magnitude,
        y: normalized.y * magnitude
    };
}

// Calculate paddle bounce angle based on hit position
function calculatePaddleBounce(ballX, paddleX, paddleWidth, ballSpeed) {
    // Calculate hit position from -1 to 1
    const hitPos = ((ballX - paddleX) - paddleWidth / 2) / (paddleWidth / 2);
    
    // Calculate bounce angle with maximum of 60 degrees
    const maxAngle = Math.PI / 3;
    const angle = hitPos * maxAngle;
    
    // Return new velocity components
    return {
        vx: Math.sin(angle) * ballSpeed,
        vy: -Math.abs(Math.cos(angle) * ballSpeed)
    };
}

// Simple AABB collision detection between two rectangles
function rectCollision(rect1, rect2) {
    return rect1.x < rect2.x + rect2.width &&
           rect1.x + rect1.width > rect2.x &&
           rect1.y < rect2.y + rect2.height &&
           rect1.y + rect1.height > rect2.y;
}

// Clamp value between minimum and maximum bounds
function clamp(value, min, max) {
    return Math.min(Math.max(value, min), max);
}

// Linear interpolation between two values
function lerp(a, b, t) {
    return a + (b - a) * t;
}