package tools

import "net/http"

// Home handler for displaying the form
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the URL path is incorrect (not "/"), return 404 error
	if r.URL.Path != "/" {
		RenderErrorTemplate(w, http.StatusNotFound, "Page Not Found")
		return
	}

	// Check if the HTTP method is not GET (only allow GET requests for the home page)
	if r.Method != http.MethodGet {
		RenderErrorTemplate(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	// Render the home page (index.html)
	if err := templates.ExecuteTemplate(w, "index.html", nil); err != nil {
		RenderErrorTemplate(w, http.StatusInternalServerError, "Error rendering template")
	}
}
