package tools

import (
	"net/http"
)

// renderErrorTemplate renders the errors.html template with a status code and message
func RenderErrorTemplate(w http.ResponseWriter, statusCode int, message string) {
	// Set the status code in the HTTP response
	w.WriteHeader(statusCode)

	// Define the data to pass to the template (status code and message)
	data := struct {
		StatusCode int
		Message    string
	}{
		StatusCode: statusCode,
		Message:    message,
	}

	// Render the error page using the errors.html template
	if err := templates.ExecuteTemplate(w, "errors.html", data); err != nil {
		http.Error(w, "Error rendering error page", http.StatusInternalServerError) // Handle template rendering errors
	}
}
