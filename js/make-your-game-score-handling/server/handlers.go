package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	
)

// EnableCORS sets up CORS headers for API requests
func EnableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
}

// HandleScores manages both GET (retrieve) and POST (submit) score requests
func HandleScores(w http.ResponseWriter, r *http.Request) {
	EnableCORS(w)

	// Handle preflight requests
	if r.Method == "OPTIONS" {
		return
	}

	switch r.Method {
	case "GET":
		handleGetScores(w, r)
	case "POST":
		handlePostScore(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleGetScores retrieves paginated scores
func handleGetScores(w http.ResponseWriter, r *http.Request) {
	// Parse page parameter
	pageStr := r.URL.Query().Get("page")
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// Get paginated scores
	scoresResponse, err := GetPaginatedScores(page)
	if err != nil {
		fmt.Printf("❌ Error loading scores: %v\n", err)
		http.Error(w, "Could not load scores", http.StatusInternalServerError)
		return
	}

	// Send response
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(scoresResponse); err != nil {
		fmt.Printf("❌ Error encoding scores response: %v\n", err)
		http.Error(w, "Could not encode response", http.StatusInternalServerError)
		return
	}

	fmt.Printf("📊 Scores served: page %d/%d (%d total scores)\n", 
		scoresResponse.CurrentPage, scoresResponse.TotalPages, scoresResponse.TotalScores)
}

// handlePostScore submits a new score and returns ranking information
func handlePostScore(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var newScore Score
	if err := json.NewDecoder(r.Body).Decode(&newScore); err != nil {
		fmt.Printf("❌ Invalid JSON: %v\n", err)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if strings.TrimSpace(newScore.Name) == "" {
		http.Error(w, "Player name is required", http.StatusBadRequest)
		return
	}
	if newScore.Score <= 0 {
		http.Error(w, "Score must be positive", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(newScore.Time) == "" {
		http.Error(w, "Time is required", http.StatusBadRequest)
		return
	}

	// Sanitize name (remove extra whitespace, limit length)
	newScore.Name = strings.TrimSpace(newScore.Name)
	if len(newScore.Name) > 20 {
		newScore.Name = newScore.Name[:20]
	}

	// Add score and get ranking information
	submitResponse, err := AddScoreWithRanking(newScore)
	if err != nil {
		fmt.Printf("❌ Error saving score: %v\n", err)
		http.Error(w, "Could not save score", http.StatusInternalServerError)
		return
	}

	// Send success response with ranking
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(submitResponse); err != nil {
		fmt.Printf("❌ Error encoding submit response: %v\n", err)
		http.Error(w, "Could not encode response", http.StatusInternalServerError)
		return
	}

	fmt.Printf("✅ Score added: %s - %d points (%s) - Rank: %d/%d (%.1f%%)\n", 
		newScore.Name, newScore.Score, newScore.Time, 
		submitResponse.PlayerRank, submitResponse.TotalScores, submitResponse.Percentage)
}