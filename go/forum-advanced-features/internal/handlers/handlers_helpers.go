package handlers

import (
	"html/template"
	"net/http"
	"strconv"
)

// ValidatePages checks the "page" query parameter in the request URL.
// If the parameter is valid, it returns the page number; otherwise, it renders an error.
func ValidatePages(w http.ResponseWriter, r *http.Request, tmpl *template.Template) int {
	var selectedPage int
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			selectedPage = p
		} else if err != nil {
			RenderError(tmpl, w, http.StatusBadRequest, "Not valid 'page number'")
			return 0
		}
	}
	return selectedPage
}
