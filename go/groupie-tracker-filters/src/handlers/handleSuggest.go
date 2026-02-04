package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// SuggestHandler provides autocomplete suggestions based on a query.
func SuggestHandler(w http.ResponseWriter, r *http.Request) {
	// Extract and clean the query string.
	query := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("q")))
	if query == "" {
		json.NewEncoder(w).Encode([]SearchSuggestion{})
		return
	}

	// Perform a grouped search based on the query.
	results := searchArtistsGrouped(Artists, query)

	// Generate autocomplete suggestions from the search results.
	suggestions := generateSuggestionsFromResults(results)

	// Respond with the suggestions as JSON.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(suggestions)
}

// generateSuggestionsFromResults creates autocomplete suggestions from search results.
func generateSuggestionsFromResults(results SearchResult) []SearchSuggestion {
	var suggestions []SearchSuggestion
	seen := make(map[string]bool) // Track unique suggestions to avoid duplicates.

	// Add artist suggestions.
	for _, match := range results.Artists {
		key := match.Name + "|artist"
		if !seen[key] {
			suggestions = append(suggestions, SearchSuggestion{
				Text: match.Name + " - artist/band",
				Type: "artist/band",
				URL:  fmt.Sprintf("/artist/%d", match.Artist.ID),
			})
			seen[key] = true
		}
	}

	// Add member suggestions.
	for _, match := range results.Members {
		key := match.Member + "|member|" + match.Artist.Name
		if !seen[key] {
			suggestions = append(suggestions, SearchSuggestion{
				Text: match.Artist.Name + " - member: " + match.Member,
				Type: "member",
				URL:  fmt.Sprintf("/artist/%d", match.Artist.ID),
			})
			seen[key] = true
		}
	}

	// Add location suggestions.
	for _, match := range results.Locations {
		key := match.Location + "|location|" + match.Artist.Name
		if !seen[key] {
			suggestions = append(suggestions, SearchSuggestion{
				Text: match.Artist.Name + " - location: " + match.Location,
				Type: "location",
				URL:  fmt.Sprintf("/artist/%d", match.Artist.ID),
			})
			seen[key] = true
		}
	}

	// Add creation date suggestions.
	for _, match := range results.CreationDates {
		key := fmt.Sprint(match.CreationDate) + "|date|" + match.Artist.Name
		if !seen[key] {
			suggestions = append(suggestions, SearchSuggestion{
				Text: match.Artist.Name + " - creation date: " + fmt.Sprint(match.CreationDate),
				Type: "creation date",
				URL:  fmt.Sprintf("/artist/%d", match.Artist.ID),
			})
			seen[key] = true
		}
	}

	// Add first album suggestions.
	for _, match := range results.FirstAlbums {
		key := match.FirstAlbum + "|album|" + match.Artist.Name
		if !seen[key] {
			suggestions = append(suggestions, SearchSuggestion{
				Text: match.Artist.Name + " - first album: " + match.FirstAlbum,
				Type: "first album",
				URL:  fmt.Sprintf("/artist/%d", match.Artist.ID),
			})
			seen[key] = true
		}
	}

	return suggestions
}
