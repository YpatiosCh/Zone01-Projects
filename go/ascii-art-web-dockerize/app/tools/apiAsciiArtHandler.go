package tools

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Handler for generating ASCII art via API and returning the result in JSON format
func ApiAsciiArtHandler(w http.ResponseWriter, r *http.Request) {
	// If the request method is not POST, return a 405 error
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed) // 405 Method Not Allowed
		return
	}

	// Create a struct to hold the incoming JSON data
	var requestData struct {
		Text   string `json:"text"`
		Banner string `json:"banner"`
	}

	// Parse the JSON request body into the struct
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest) // 400 Bad Request
		return
	}

	// Check if the required fields (text or banner) are missing
	if requestData.Text == "" || requestData.Banner == "" {
		http.Error(w, "Missing text or banner", http.StatusBadRequest) // 400 Bad Request
		return
	}

	// Construct the file path to the chosen banner
	fontFile := fmt.Sprintf("banners/%s.txt", requestData.Banner)

	// Call the function to generate ASCII art
	result, err := PrintAsciiArt(requestData.Text, fontFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating ASCII art: %v", err), http.StatusInternalServerError) // 500 Internal Server Error
		return
	}

	// Prepare the response data (ASCII art as an array of strings)
	response := struct {
		AsciiArt []string `json:"ascii_art"`
	}{
		AsciiArt: result,
	}

	// Set the response content type to JSON and encode the response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError) // 500 Internal Server Error
	}
}
