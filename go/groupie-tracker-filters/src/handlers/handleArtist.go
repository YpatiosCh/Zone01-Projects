package handlers

import (
	// Import the package for fetching artist data
	"groupie-tracker/models" // Import the models package which contains the Artist struct
	"log"                    // Package for logging errors or messages
	"net/http"               // Package for handling HTTP requests and responses
	"strconv"                // Package for converting strings to integers
	"strings"                // Package for string manipulation
)

// HandleArtist is an HTTP handler function that processes requests for artist details
func HandleArtist(w http.ResponseWriter, r *http.Request) {
	// Extract the artist ID from the URL path by trimming the "/artist/" prefix
	id := strings.TrimPrefix(r.URL.Path, "/artist/")

	// If the ID is empty, return a 400 Bad Request response with an error message
	if id == "" {
		RenderErrorTemplate(w, http.StatusBadRequest, "Invalid artist ID")
		return
	}

	// Convert the extracted ID from string to integer
	idInt, err := strconv.Atoi(id)
	// If conversion fails or the ID is non-positive, return a 400 Bad Request response
	if err != nil || idInt <= 0 {
		RenderErrorTemplate(w, http.StatusBadRequest, "Invalid artist ID")
		return
	}

	// Declare a variable to hold the selected artist
	var artistSelected *models.Artist

	// If the provided artist ID exceeds the number of available artists, return a 400 Bad Request response
	if idInt > len(Artists) {
		RenderErrorTemplate(w, http.StatusBadRequest, "Artist ID does not exist")
		return
	}

	// Iterate through the list of artists to find the one matching the given ID
	for _, artist := range Artists {
		if artist.ID == idInt {
			artistSelected = &artist // Assign the matched artist to the artistSelected variable
			break
		}
	}

	// Prepare the data to be passed to the HTML template
	data := struct {
		Artist *models.Artist // Embed the artist data into the structure
	}{
		Artist: artistSelected,
	}

	// Render the artist details using the "artist.html" template
	err = Templates.ExecuteTemplate(w, "artist.html", data)
	// If there's an error rendering the template, log the error and return a 500 Internal Server Error response
	if err != nil {
		log.Printf("Error rendering template: %v", err) // Log the specific error for debugging
		RenderErrorTemplate(w, http.StatusInternalServerError, "could not execute template")
		return
	}
}
