# Arkanoid Game - High Performance JavaScript Implementation

## 🎯 Project Overview

This is a high-performance Arkanoid (brick breaker) game built with **vanilla JavaScript** achieving consistent **60 FPS** without any frame drops. The project demonstrates advanced performance optimization techniques using DOM manipulation without Canvas or frameworks.

## 🎮 Live Demo

Start the server and play at: `http://localhost:8080/game.html`

## 📋 Project Requirements

The game was built to demonstrate proficiency in:

- **requestAnimationFrame** - Smooth 60fps game loop
- **Event Loop** - Efficient JavaScript execution
- **FPS Management** - Consistent frame rate without drops
- **DOM Manipulation** - High-performance element updates
- **Anti-Jank Techniques** - Eliminating stutter/frame drops
- **Transform/Opacity** - Hardware-accelerated CSS properties
- **Browser Performance Tools** - Firefox & Chrome DevTools optimization
- **Rendering Pipeline** - Layout, Painting, Compositing optimization

## 🏗️ Architecture

### File Structure

```
├── game.html          # Main game HTML
├── styles/
│   └── game.css       # Optimized CSS with GPU acceleration
├── js/
│   ├── game.js        # Core game logic & performance optimizations
│   ├── levels.js      # Level data and brick creation
│   └── physics.js     # Collision detection utilities
├── server/
│   ├── main.go        # Go backend for high scores API
│   └── go.mod         # Go module definition
└── scores/
    └── scores.json    # Persistent score storage
```

## 🚀 Performance Optimizations

### 1. **60 FPS Game Loop**

```javascript
function gameLoop(currentTime) {
  if (!gameRunning) return;

  // Delta time clamping prevents large jumps
  const deltaTime = Math.min((currentTime - lastTime) / 1000, 1 / 30);
  lastTime = currentTime;

  update(deltaTime);
  render();

  animationId = requestAnimationFrame(gameLoop);
}
```

**Key Features:**

- Delta time clamping prevents spiral of death
- Consistent physics regardless of frame rate
- Proper cleanup with `cancelAnimationFrame`

### 2. **DOM Update Caching**

```javascript
const hudCache = {
  score: -1,
  lives: -1,
  level: -1,
  ballTransform: "",
  paddleTransform: "",
};

// Only update DOM when values actually change
if (game.score !== hudCache.score) {
  hudCache.score = game.score;
  domUpdater.updateScore(game.score);
}
```

**Benefits:**

- Eliminates redundant DOM writes
- Reduces layout thrashing
- Maintains 60fps during intensive updates

### 3. **Hardware-Accelerated CSS**

```css
#ball,
#paddle {
  will-change: transform;
  backface-visibility: hidden;
  transform: translate3d(0, 0, 0); /* Force GPU layer */
}

.brick {
  will-change: opacity;
  transition: opacity 100ms ease-out;
}
```

**GPU Optimization:**

- `translate3d()` forces hardware acceleration
- `will-change` hints optimize compositor layers
- Opacity animations bypass layout/paint

### 4. **Efficient Collision System**

```javascript
// Fast brick removal: swap-and-pop instead of splice
if (i < game.activeBricksList.length - 1) {
  game.activeBricksList[i] =
    game.activeBricksList[game.activeBricksList.length - 1];
}
game.activeBricksList.pop();

// Batch DOM cleanup in separate frame
brickPool.destroyed.push(brick);
brickPool.needsCleanup = true;
```

**Optimizations:**

- O(1) array removal vs O(n) splice
- Deferred DOM cleanup prevents frame drops
- Active brick list for faster iteration

### 5. **Memory Management**

```javascript
const domUpdater = {
  // Pre-validated element references
  ball: null,
  paddle: null,

  init() {
    // Cache all references once
    this.ball = elements.ball;
    this.paddle = elements.paddle;

    // Fail fast validation
    if (!this.ball || !this.paddle) {
      throw new Error("Critical elements missing");
    }
  },
};
```

**Memory Efficiency:**

- Pre-cached DOM references
- No repeated `getElementById` calls
- Fail-fast validation prevents runtime errors

### 6. **Input Handling Optimization**

```javascript
// Continuous movement tracking
const keys = {};

document.addEventListener("keydown", (e) => {
  keys[e.code] = true; // Track state

  // Discrete events (prevent repeat firing)
  if (!keys[e.code + "_pressed"]) {
    keys[e.code + "_pressed"] = true;
    handleGameInput(e.code);
  }
});

// Update loop uses direct key state
if (keys["ArrowLeft"]) {
  game.paddle.x -= game.paddle.speed * deltaTime;
}
```

**Input Benefits:**

- Smooth continuous movement
- Prevents key repeat issues
- Frame-rate independent motion

## 🎯 Anti-Jank Techniques

### Transform-Only Animations

- **Avoid:** Changing `left`, `top`, `width`, `height`
- **Use:** `transform: translate3d()` for position changes
- **Result:** Animations run on compositor thread

### Batch DOM Operations

```javascript
function cleanupDestroyedBricks() {
  // Batch remove multiple elements
  for (let brick of brickPool.destroyed) {
    brick.element.parentNode.removeChild(brick.element);
  }
  brickPool.destroyed.length = 0; // Clear pool
}
```

### Minimize Layout Triggers

- Cache element dimensions
- Use `transform` and `opacity` only
- Avoid reading layout properties during animations

## 🛠️ Development Tools Usage

### Chrome DevTools Optimization

1. **Performance Tab:** Record gameplay to identify frame drops
2. **Rendering Tab:** Enable "FPS Meter" and "Frame Rendering Stats"
3. **Layers Tab:** Verify GPU acceleration layers
4. **Memory Tab:** Check for memory leaks during gameplay

### Firefox DevTools

1. **Performance Tool:** Analyze call stack and frame timing
2. **Inspector:** Verify CSS transforms and will-change properties
3. **Console:** Monitor performance warnings

## 🎮 Game Features

### Core Mechanics

- **5 Levels** with increasing difficulty
- **3 Brick Types:** Weak (10 pts), Normal (25 pts), Strong (50 pts)
- **Physics:** Realistic ball bouncing with paddle angle influence
- **Lives System:** 3 lives with ball respawn
- **Scoring:** Points + time-based + level completion bonuses

### Technical Features

- **Persistent High Scores** via Go backend API
- **Real-time HUD** with cached updates
- **Keyboard-Only Controls** for consistent input
- **State Machine** for proper game flow management
- **Brick Pool Management** for efficient memory usage

## 🌐 Backend API

### Go Server Features

```go
// CORS-enabled REST API
func handleScores(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    // ... handle GET/POST for scores
}
```

**API Endpoints:**

- `GET /api/scores` - Fetch top 5 scores
- `POST /api/scores` - Submit new score
- Static file serving for game assets

## 🚀 Running the Project

### Prerequisites

- Go 1.25+ for backend server
- Modern web browser (Chrome/Firefox recommended)

### Setup

```bash
# 1. Start the Go server
cd server/
go run main.go

# 2. Open browser
# Navigate to: http://localhost:8080/game.html

# 3. Play the game!
# - Enter your name
# - Use arrow keys to move paddle
# - SPACE to launch ball/pause
# - ENTER to continue through menus
```

## 📈 Performance Metrics

### Target Performance

- **Frame Rate:** Consistent 60 FPS
- **Frame Time:** ~16.67ms per frame
- **Memory:** Stable usage, no leaks
- **Jank:** Zero stuttering during gameplay

### Measurement Tools

- Chrome: `performance.mark()` and Performance tab
- Firefox: Performance tool and frame rate monitor
- Manual: Visual smoothness testing

## 🎯 Learning Outcomes

This project demonstrates mastery of:

- High-performance JavaScript game development
- DOM optimization techniques
- CSS GPU acceleration
- Browser rendering pipeline understanding
- Memory management best practices
- Real-world performance debugging

## 🔧 Browser Compatibility

**Optimized for:**

- Chrome 90+
- Firefox 88+
- Safari 14+
- Edge 90+

**Performance Features Used:**

- `requestAnimationFrame`
- CSS `transform` and `opacity` animations
- Hardware acceleration via `translate3d()`
- Modern ES6+ JavaScript features

## 💡 Key Takeaways

1. **Transform > Position:** Always use CSS transforms for animations
2. **Cache Everything:** DOM queries are expensive, cache references
3. **Batch Operations:** Group DOM changes to minimize reflows
4. **Profile Early:** Use browser tools to identify bottlenecks
5. **Hardware Acceleration:** Force GPU layers for smooth animations
6. **Memory Awareness:** Clean up resources and avoid leaks
7. **Delta Time:** Physics should be frame-rate independent

---

_This project showcases how to build performant web games using vanilla JavaScript while maintaining 60 FPS without any frameworks or Canvas API._
