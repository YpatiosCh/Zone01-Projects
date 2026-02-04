package handlers

import (
	"encoding/json"
	"groupie-tracker/filter"
	"groupie-tracker/models"
	"net/http"
)

// FilterArtistsHandler is the HTTP handler for filtering artists based on request parameters
func FilterArtistsHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the form data from the request body
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Create a Filters object and extract filter parameters from the request
	filters := filter.Filters{
		CreationDateMin: filter.GetIntParam(r, "CreationDateMin", 1960), // Default min year: 1960
		CreationDateMax: filter.GetIntParam(r, "CreationDateMax", 2024), // Default max year: 2024
		FirstAlbumMin:   filter.GetIntParam(r, "albumMin", 1960),        // Default min first album year: 1960
		FirstAlbumMax:   filter.GetIntParam(r, "albumMax", 2024),        // Default max first album year: 2024
		MemberMin:       filter.GetIntParam(r, "memberMin", 1),          // Default min members: 1
		MemberMax:       filter.GetIntParam(r, "memberMax", 8),          // Default max members: 8
		Locations:       r.Form["locations"],                            // List of location filters
	}

	// Filter the artists based on the specified filters
	filteredArtists := FilterArtists(Artists, filters)

	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode the filtered artists list into JSON and send as the response
	json.NewEncoder(w).Encode(filteredArtists)
}

// FilterArtists applies the filters to the list of artists and returns the filtered result
func FilterArtists(artists []models.Artist, filters filter.Filters) []models.Artist {
	var filtered []models.Artist

	// Loop through each artist and apply filters
	for _, artist := range artists {
		// Check if the artist's creation date is within the specified range
		if !filter.InRange(artist.CreationDate, filters.CreationDateMin, filters.CreationDateMax) {
			continue // Skip the artist if it doesn't match the creation date filter
		}

		// Check if the artist's first album date is within the specified range
		if !filter.InAlbumRange(artist.FirstAlbum, filters.FirstAlbumMin, filters.FirstAlbumMax) {
			continue // Skip the artist if it doesn't match the first album date filter
		}

		// Check if the artist's member count is within the specified range
		if !filter.InRange(len(artist.Members), filters.MemberMin, filters.MemberMax) {
			continue // Skip the artist if it doesn't match the member count filter
		}

		// Check if the artist's locations match the specified location filters
		if !filter.CheckLocations(artist.LocationsData.Locations, filters.Locations) {
			continue // Skip the artist if it doesn't match the location filter
		}

		// If all filters passed, add the artist to the filtered list
		filtered = append(filtered, artist)
	}

	// Return the filtered list of artists
	return filtered
}
