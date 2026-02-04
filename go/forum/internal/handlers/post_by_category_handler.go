package handlers

import (
	"fmt"
	"forum/internal/middleware"
	"forum/internal/models"
	"forum/internal/services"
	"html/template"
	"net/http"
)

type PostsByCategoryData struct {
	User       *models.User
	Categories []models.Category
	TopPosts   []models.Post
	NewPosts   []models.Post
	Page       int
	TotalPages int
}

// PostsByCategory handles the showcasing of all posts of a specific category.
func PostsByCategory(postService *services.PostService, categoryService *services.CategoryService, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// method
		if r.Method != http.MethodGet {
			RenderError(tmpl, w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		// get user
		user := middleware.GetUser(r)

		// get categories
		categories, err := categoryService.GetAllCategories()
		if err != nil {
			RenderError(tmpl, w, http.StatusInternalServerError, "Failed to get categories")
			return
		}

		// Get category ID from URL
		categoryID := r.PathValue("category_id")

		// Validate that the category exists
		categoryExists := false
		for _, cat := range categories {
			if cat.ID == categoryID {
				categoryExists = true
				break
			}
		}
		if !categoryExists {
			RenderError(tmpl, w, http.StatusNotFound, "Category not found")
			return
		}

		// Set default page to 1
		page := 1
		if pageParam := ValidatePages(w, r, tmpl); pageParam > 0 {
			page = pageParam
		}

		// Get paginated posts for this category by engagement
		topPosts, totalPages, err := postService.GetPaginatedPostsByCategoryByEngagement(categoryID, page, 5)
		if err != nil {
			RenderError(tmpl, w, http.StatusInternalServerError, "Failed to get category posts")
			return
		}

		// Handle empty category case
		if totalPages == 0 {
			// Get empty newest posts too
			newPosts, _, err := postService.GetPaginatedPostsByCategoryNewest(categoryID, page, 5)
			if err != nil {
				RenderError(tmpl, w, http.StatusInternalServerError, "Failed to get newest category posts")
				return
			}

			data := PostsByCategoryData{
				User:       user,
				Categories: categories,
				TopPosts:   topPosts,
				NewPosts:   newPosts,
				Page:       1, // Always show page 1 when empty
				TotalPages: 1, // Avoid pagination display issues
			}

			tmpl.ExecuteTemplate(w, "index.html", data)
			return
		}

		// Handle out-of-range pages - redirect to last valid page
		if page > totalPages {
			redirectURL := fmt.Sprintf("/categories/%s/posts?page=%d", categoryID, totalPages)
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return
		}

		// Handle negative or zero pages - redirect to first page
		if page < 1 {
			redirectURL := fmt.Sprintf("/categories/%s/posts?page=1", categoryID)
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return
		}

		// Get newest posts for this category
		newPosts, _, err := postService.GetPaginatedPostsByCategoryNewest(categoryID, page, 5)
		if err != nil {
			RenderError(tmpl, w, http.StatusInternalServerError, "Failed to get newest category posts")
			return
		}

		data := PostsByCategoryData{
			User:       user,
			Categories: categories,
			TopPosts:   topPosts,
			NewPosts:   newPosts,
			Page:       page,
			TotalPages: totalPages,
		}

		tmpl.ExecuteTemplate(w, "index.html", data)
	}
}
