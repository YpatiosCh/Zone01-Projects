package handlers

import (
	"html/template"
	"net/http"
)

type errorHandlerData struct {
	StatusCode int
	Message    string
}

// RenderError renders the "errors.html" template with a given status code and error message
func RenderError(template *template.Template, w http.ResponseWriter, statusCode int, message string) {
	// Set the status code of the response
	w.WriteHeader(statusCode)

	// Define the data structure to pass to the template
	data := errorHandlerData{
		StatusCode: statusCode,
		Message:    message,
	}

	// Attempt to render the "errors.html" template with the data provided
	if err := template.ExecuteTemplate(w, "errors.html", data); err != nil {
		// If template execution fails, set status code and write error message
		w.WriteHeader(statusCode)
		w.Write([]byte("Error rendering error page"))
		return
	}
}
