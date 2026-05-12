// GAME THEMES SYSTEM
// Manages background themes for each level

// Theme class names for each level
const LEVEL_THEMES = {
    1: 'level-1-theme',  // Starfield theme for "The Outer Rim"
    2: 'level-2-theme',  // Future: Energy grid theme
    3: 'level-3-theme',  // Future: Nebula theme
    4: 'level-4-theme',  // Future: Plasma field theme
    5: 'level-5-theme'   // Future: Command center theme
};

// Apply theme based on current level
function applyLevelTheme(levelNumber) {
    // Remove all existing theme classes from body
    Object.values(LEVEL_THEMES).forEach(theme => {
        document.body.classList.remove(theme);
    });
    
    // Add the current level's theme class
    const currentTheme = LEVEL_THEMES[levelNumber];
    if (currentTheme) {
        document.body.classList.add(currentTheme);
        console.log(`🎨 Applied theme: ${currentTheme} for level ${levelNumber}`);
    } else {
        console.warn(`⚠️ No theme found for level ${levelNumber}`);
    }
}

// Remove all themes (for menu screens and clean states)
function clearAllThemes() {
    Object.values(LEVEL_THEMES).forEach(theme => {
        document.body.classList.remove(theme);
    });
    console.log('🎨 All themes cleared');
}

// Get current active theme (for debugging)
function getCurrentTheme() {
    for (const [level, theme] of Object.entries(LEVEL_THEMES)) {
        if (document.body.classList.contains(theme)) {
            return { level: parseInt(level), theme };
        }
    }
    return null;
}

// Preload all themes for better performance
function preloadThemes() {
    console.log('🎨 Preloading themes...');
    
    // Temporarily apply each theme to force CSS parsing
    Object.values(LEVEL_THEMES).forEach(theme => {
        document.body.classList.add(theme);
        // Force a reflow to parse CSS and animations
        document.body.offsetHeight;
        document.body.classList.remove(theme);
    });
    
    console.log('✅ All themes preloaded for optimal performance');
}

// Theme system initialization
function initThemeSystem() {
    console.log('🎨 Theme system initialized');
    preloadThemes();
}