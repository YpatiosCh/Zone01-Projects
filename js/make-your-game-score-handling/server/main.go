package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	

)

func main() {
	fmt.Printf("🎯 ARKANOID GAME SERVER - REFACTORED EDITION\n")
	fmt.Printf("====================================================\n")

	// Initialize scores file if it doesn't exist
	initializeScoresFile()

	// Setup API routes
	http.HandleFunc("/api/scores", HandleScores)

	// Serve static files with proper MIME types
	http.Handle("/", LogRequest(CustomFileServer("../")))

	// Server configuration
	port := "8080"
	
	// Startup information
	fmt.Printf("🎮 Game URL:    http://localhost:%s/game.html\n", port)
	fmt.Printf("📊 Scores API:  http://localhost:%s/api/scores\n", port)
	fmt.Printf("📁 Scores file: %s\n", SCORES_FILE)
	fmt.Printf("📄 Per page:    %d scores\n", SCORES_PER_PAGE)
	fmt.Printf("⚡ Features:    Pagination, Ranking, MIME types\n")
	fmt.Printf("🚀 Server starting on port %s...\n", port)
	fmt.Printf("====================================================\n")

	// Start server
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("❌ Server failed to start: %v", err)
	}
}

// initializeScoresFile creates an empty scores file if it doesn't exist
func initializeScoresFile() {
	if _, err := os.Stat(SCORES_FILE); os.IsNotExist(err) {
		if err := SaveScores([]Score{}); err != nil {
			fmt.Printf("⚠️  Warning: Could not create scores file: %v\n", err)
		} else {
			fmt.Printf("📄 Created empty scores file: %s\n", SCORES_FILE)
		}
	} else {
		// Load and validate existing scores
		scores, err := LoadScores()
		if err != nil {
			fmt.Printf("⚠️  Warning: Could not load existing scores: %v\n", err)
		} else {
			fmt.Printf("✅ Loaded existing scores: %d entries\n", len(scores))
		}
	}
}