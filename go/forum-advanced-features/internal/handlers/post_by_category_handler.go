package handlers

import (
	"fmt"
	"forum/internal/middleware"
	"forum/internal/models"
	"forum/internal/services"
	"html/template"
	"net/http"
)

// PostsByCategory handles the showcasing of all posts of a specific category.
func PostsByCategory(service services.Services, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// method
		if r.Method != http.MethodGet {
			RenderError(tmpl, w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		// get user
		user := middleware.GetUser(r)

		// get categories
		categories, err := service.Category().GetAllCategories()
		if err != nil {
			RenderError(tmpl, w, http.StatusInternalServerError, "Failed to get categories")
			return
		}

		var notifications int
		if user != nil {
			notifCount, err := service.Notify().GetUnreadNotificationCount(user.ID)
			if err != nil {
				RenderError(tmpl, w, http.StatusInternalServerError, "Failed to get notifications")
				return
			}
			notifications = notifCount
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
		topPosts, totalPages, err := service.Post().GetPaginatedPostsByCategoryByEngagement(categoryID, page, 5)
		if err != nil {
			RenderError(tmpl, w, http.StatusInternalServerError, "Failed to get category posts")
			return
		}

		// Handle empty category case
		if totalPages == 0 {
			// Get empty newest posts too
			newPosts, _, err := service.Post().GetPaginatedPostsByCategoryNewest(categoryID, page, 5)
			if err != nil {
				RenderError(tmpl, w, http.StatusInternalServerError, "Failed to get newest category posts")
				return
			}

			data := models.HomeHandlerData{
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
		newPosts, _, err := service.Post().GetPaginatedPostsByCategoryNewest(categoryID, page, 5)
		if err != nil {
			RenderError(tmpl, w, http.StatusInternalServerError, "Failed to get newest category posts")
			return
		}

		// check if user is logged in and get their likes/dislikes for the posts
		if user != nil {
			for i, post := range topPosts {
				liked, disliked, err := service.Reaction().UserLikedDislikedPost(user.ID, post.ID)
				if err != nil {
					RenderError(tmpl, w, http.StatusInternalServerError, "Failed to get user post engagement")
					return
				}
				topPosts[i].IsLiked = liked
				topPosts[i].IsDisliked = disliked
			}
			for i, post := range newPosts {
				liked, disliked, err := service.Reaction().UserLikedDislikedPost(user.ID, post.ID)
				if err != nil {
					RenderError(tmpl, w, http.StatusInternalServerError, "Failed to get user post engagement")
					return
				}
				newPosts[i].IsLiked = liked
				newPosts[i].IsDisliked = disliked
			}
		}

		data := models.HomeHandlerData{
			User:          user,
			Categories:    categories,
			TopPosts:      topPosts,
			NewPosts:      newPosts,
			Page:          page,
			TotalPages:    totalPages,
			Notifications: notifications,
		}

		tmpl.ExecuteTemplate(w, "index.html", data)
	}
}
