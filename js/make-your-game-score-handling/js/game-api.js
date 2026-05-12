// Enhanced Score API and storage management with pagination and ranking
// Handles communication with backend server for high scores with pagination

// Global pagination state
let currentScorePage = 1;
let totalScorePages = 1;
let isLoadingScores = false;

// Load paginated high scores from server
async function loadHighScores(page = 1) {
    if (isLoadingScores) return; // Prevent concurrent requests
    
    isLoadingScores = true;
    
    try {
        // Show loading indicator
        elements.highscoresList.innerHTML = '<div class="loading-scores">Loading scores...</div>';
        
        const response = await fetch(`http://localhost:8080/api/scores?page=${page}`);
        const data = await response.json();
        
        // Update global pagination state
        currentScorePage = data.currentPage;
        totalScorePages = data.totalPages;
        
        displayPaginatedScores(data);
        
    } catch (error) {
        console.error('Error loading high scores:', error);
        elements.highscoresList.innerHTML = '<div class="error-message">Could not load scores. Please try again.</div>';
    } finally {
        isLoadingScores = false;
    }
}

// Display paginated scores with navigation
function displayPaginatedScores(data) {
    if (!data.scores || data.scores.length === 0) {
        if (data.totalScores === 0) {
            elements.highscoresList.innerHTML = '<div class="no-scores">No scores yet! Be the first to play!</div>';
        } else {
            elements.highscoresList.innerHTML = '<div class="error-message">No scores found for this page.</div>';
        }
        return;
    }
    
    // Build enhanced HTML table
    let html = '<div class="scores-container">';
    
    // Header
    html += `
        <div class="scores-header">
            <div class="scores-title">High Scores (${data.totalScores} total)</div>
        </div>
    `;
    
    // Table
    html += '<div class="scores-table">';
    
    // Table header
    html += `
        <div class="score-row header-row">
            <span class="rank-col">Rank</span>
            <span class="name-col">Name</span>
            <span class="score-col">Score</span>
            <span class="time-col">Time</span>
        </div>
    `;
    
    // Score rows
    data.scores.forEach((score) => {
        html += `
            <div class="score-row data-row">
                <span class="rank-col rank-${score.rank}">${score.rankText}</span>
                <span class="name-col">${escapeHtml(score.name)}</span>
                <span class="score-col">${score.score.toLocaleString()}</span>
                <span class="time-col">${score.time}</span>
            </div>
        `;
    });
    
    html += '</div>'; // Close scores-table
    
    // Pagination controls
    if (data.totalPages > 1) {
        html += '<div class="pagination-controls">';
        
        // Previous button
        const prevDisabled = data.currentPage <= 1 ? 'disabled' : '';
        html += `<button class="page-btn ${prevDisabled}" onclick="navigateScores('prev')" ${prevDisabled ? 'disabled' : ''}>
                    ← Previous
                 </button>`;
        
        // Page indicator
        html += `<span class="page-info">Page ${data.currentPage}/${data.totalPages}</span>`;
        
        // Next button
        const nextDisabled = data.currentPage >= data.totalPages ? 'disabled' : '';
        html += `<button class="page-btn ${nextDisabled}" onclick="navigateScores('next')" ${nextDisabled ? 'disabled' : ''}>
                    Next →
                 </button>`;
        
        html += '</div>'; // Close pagination-controls
    }
    
    html += '</div>'; // Close scores-container
    
    elements.highscoresList.innerHTML = html;
}

// Navigate between score pages
async function navigateScores(direction) {
    if (isLoadingScores) return;
    
    let newPage = currentScorePage;
    
    if (direction === 'prev' && currentScorePage > 1) {
        newPage = currentScorePage - 1;
    } else if (direction === 'next' && currentScorePage < totalScorePages) {
        newPage = currentScorePage + 1;
    } else {
        return; // No change needed
    }
    
    await loadHighScores(newPage);
}

// Submit player score to server with enhanced response handling
async function submitScore() {
    if (game.score <= 0) {
        console.warn('Cannot submit score of 0 or negative');
        return;
    }
    
    // Format game time as MM:SS string
    const minutes = Math.floor(game.gameTime / 60);
    const seconds = game.gameTime % 60;
    const timeString = `${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`;
    
    try {
        const response = await fetch('http://localhost:8080/api/scores', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                name: game.playerName,
                score: game.score,
                time: timeString
            })
        });
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        const result = await response.json();
        
        // Store ranking information for display in modals
        if (result.status === 'success') {
            game.lastSubmitResult = {
                rank: result.playerRank,
                totalScores: result.totalScores,
                percentage: result.percentage,
                rankMessage: result.rankMessage
            };
            
            console.log(`Score submitted successfully! ${result.rankMessage}`);
            
            // Refresh the scores display to show the new score
            await loadHighScores(1); // Go to first page to potentially see the new high score
        }
        
    } catch (error) {
        console.error('Error submitting score:', error);
        // Store error info for display
        game.lastSubmitResult = {
            error: true,
            message: 'Failed to submit score. Please check your connection.'
        };
    }
}

// Utility function to escape HTML to prevent XSS
function escapeHtml(unsafe) {
    return unsafe
         .replace(/&/g, "&amp;")
         .replace(/</g, "&lt;")
         .replace(/>/g, "&gt;")
         .replace(/"/g, "&quot;")
         .replace(/'/g, "&#039;");
}

// Enhanced score display for game over / victory screens
function getScoreDisplayMessage() {
    if (!game.lastSubmitResult) {
        return `Final Score: ${game.score.toLocaleString()}`;
    }
    
    if (game.lastSubmitResult.error) {
        return `Final Score: ${game.score.toLocaleString()}<br><br>${game.lastSubmitResult.message}`;
    }
    
    return `Final Score: ${game.score.toLocaleString()}<br><br>${game.lastSubmitResult.rankMessage}`;
}

// Initialize scores on page load
document.addEventListener('DOMContentLoaded', () => {
    // Add some delay to ensure elements are ready
    setTimeout(() => {
        if (elements.highscoresList) {
            loadHighScores(1);
        }
    }, 100);
});