package handlers

import (
	// Import the package to fetch data about artists
	"groupie-tracker/filter"
	"groupie-tracker/models" // Import the package defining models like Artist
	"net/http"               // Package for handling HTTP requests and responses
	"text/template"          // Package for rendering templates
)

// Global variable for pre-parsed HTML templates
// The "template.Must" ensures that parsing the templates is completed at startup,
// and it will panic if there are errors in template files.
var Templates = template.Must(template.ParseGlob("templates/*.html"))

// HomeHandler handles requests to the home page ("/")
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the requested URL path is "/" (the home page)
	if r.URL.Path != "/" {
		// If not, return a 404 Not Found response with an error message
		RenderErrorTemplate(w, http.StatusNotFound, "Page Not Found")
		return
	}

	availableLocations := filter.GetUniqueLocations(Artists)
	// Prepare the data to pass to the HTML template
	// The data contains a slice of Artist objects
	data := struct {
		Artists            []models.Artist
		AvailableLocations []filter.LocationsByCountry
	}{
		Artists:            Artists, // Use all artists, not filtered
		AvailableLocations: availableLocations,
	}

	// Render the "index.html" template with the prepared data
	err := Templates.ExecuteTemplate(w, "index.html", data)
	// If there's an error rendering the template, return a 500 Internal Server Error response
	if err != nil {
		RenderErrorTemplate(w, http.StatusInternalServerError, "could not execute template")
		return
	}
}
