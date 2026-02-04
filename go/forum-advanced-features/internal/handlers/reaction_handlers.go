package handlers

import (
	"forum/internal/middleware"
	"forum/internal/services"
	"html/template"
	"net/http"
	"strings"
)

// PostReactionHandler handles liking/disliking post
func PostReactionHandler(service services.Services, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only allow POST method
		if r.Method != http.MethodPost {
			RenderError(tmpl, w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		// get user
		user := middleware.GetUser(r)

		// Get post ID and reaction type from URL
		ID := r.PathValue("id")
		reactionType := r.PathValue("reaction_type")
		// get the table name from url
		tables := strings.Split(r.URL.Path, "/")
		table := tables[1]

		// Toggle the reaction
		err := service.Reaction().ToggleReaction(user.ID, table, ID, reactionType)
		if err != nil {
			RenderError(tmpl, w, http.StatusInternalServerError, "Failed to toggle reaction: "+err.Error())
			return
		}

		// Get the referer URL to redirect back to the same page
		referer := r.Header.Get("Referer")
		if strings.Contains(referer, "posts") {
			// If the referer is a post page, redirect to the post page without query parameters
			referer = strings.Split(referer, "?")[0]
		}

		// Redirect back to the referring page
		http.Redirect(w, r, referer, http.StatusSeeOther)
	}
}
