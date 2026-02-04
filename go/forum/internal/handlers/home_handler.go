package handlers

import (
	"fmt"
	"forum/internal/middleware"
	"forum/internal/models"
	"forum/internal/services"
	"html/template"
	"net/http"
)

type homeHandlerData struct {
	User       *models.User
	Categories []models.Category
	TopPosts   []models.Post
	NewPosts   []models.Post
	Page       int
	TotalPages int
}

// HomePage handles the homepage of the forum. It passes as data all categories, filtered posts by most engagement and by newest.
// Also passes a user data(if any) to handle the navbar accordingly
func HomePage(postService *services.PostService, categoryService *services.CategoryService, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			RenderError(tmpl, w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		// check url path
		if r.URL.Path != "/" {
			RenderError(tmpl, w, http.StatusNotFound, "The page you are looking for does not exist :)")
			return
		}
		// get the user
		user := middleware.GetUser(r)

		// get all post categories
		categories, err := categoryService.GetAllCategories()
		if err != nil {
			RenderError(tmpl, w, http.StatusInternalServerError, "Failed to get categories")
			return
		}

		// Set default page to 1
		page := 1
		if pageParam := ValidatePages(w, r, tmpl); pageParam > 0 {
			page = pageParam
		}

		// Get top posts by engagement with pagination
		topPosts, totalPages, err := postService.GetPaginatedAllPostsByEngagement(page, 5)
		if err != nil {
			RenderError(tmpl, w, http.StatusInternalServerError, "Failed to get posts")
			return
		}

		// Handle empty database case
		if totalPages == 0 {
			newPosts, _, err := postService.GetPaginatedAllNewestPosts(page, 5)
			if err != nil {
				RenderError(tmpl, w, http.StatusInternalServerError, "Failed to get newest posts")
				return
			}

			data := homeHandlerData{
				User:       user,
				Categories: categories,
				TopPosts:   topPosts,
				NewPosts:   newPosts,
				Page:       1,
				TotalPages: 1,
			}

			if err := tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
				RenderError(tmpl, w, http.StatusInternalServerError, "Error: "+err.Error())
				return
			}
			return
		}

		// Redirect out-of-range pages to last valid page
		if page > totalPages {
			// Redirect to the last valid page
			redirectURL := fmt.Sprintf("/?page=%d", totalPages)
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return
		}

		// Also handle negative or zero pages
		if page < 1 {
			http.Redirect(w, r, "/?page=1", http.StatusSeeOther)
			return
		}

		// Get newest posts with pagination
		newPosts, _, err := postService.GetPaginatedAllNewestPosts(page, 5)
		if err != nil {
			RenderError(tmpl, w, http.StatusInternalServerError, "Failed to get newest posts")
			return
		}

		data := homeHandlerData{
			User:       user,
			Categories: categories,
			TopPosts:   topPosts,
			NewPosts:   newPosts,
			Page:       page,
			TotalPages: totalPages,
		}

		if err := tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
			RenderError(tmpl, w, http.StatusInternalServerError, "Error: "+err.Error())
			return
		}
	}
}
