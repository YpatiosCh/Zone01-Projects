package handlers

import (
	"fmt"
	"forum/internal/middleware"
	"forum/internal/models"
	"forum/internal/services"
	"html/template"
	"net/http"
)

// HomePage handles the homepage of the forum. It passes as data all categories, filtered posts by most engagement and by newest.
// Also passes a user data(if any) to handle the navbar accordingly
func HomePage(service services.Services, tmpl *template.Template) http.HandlerFunc {
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
		categories, err := service.Category().GetAllCategories()
		if err != nil {
			RenderError(tmpl, w, http.StatusInternalServerError, "Failed to get categories")
			return
		}

		// Get notification count for logged-in users
		notifications := 0
		if user != nil {
			notificationCount, err := service.Notify().GetUnreadNotificationCount(user.ID)
			if err != nil {
				RenderError(tmpl, w, http.StatusInternalServerError, "Failed to get notifications")
				return
			}
			notifications = notificationCount
		}
		// Set default page to 1
		page := 1
		if pageParam := ValidatePages(w, r, tmpl); pageParam > 0 {
			page = pageParam
		}

		// Get top posts by engagement with pagination
		topPosts, totalPages, err := service.Post().GetPaginatedAllPostsByEngagement(page, 5)
		if err != nil {
			RenderError(tmpl, w, http.StatusInternalServerError, err.Error())
			return
		}

		// Handle empty database case
		if totalPages == 0 {
			newPosts, _, err := service.Post().GetPaginatedAllNewestPosts(page, 5)
			if err != nil {
				RenderError(tmpl, w, http.StatusInternalServerError, err.Error())
				return
			}

			data := models.HomeHandlerData{
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
		newPosts, _, err := service.Post().GetPaginatedAllNewestPosts(page, 5)
		if err != nil {
			RenderError(tmpl, w, http.StatusInternalServerError, "Failed to get newest posts")
			return
		}

		// check if user is logged in and get their likes/dislikes for the posts
		if user != nil {
			for i, post := range topPosts {
				// Check if the user has liked or disliked the post
				liked, disliked, err := service.Reaction().UserLikedDislikedPost(user.ID, post.ID)
				if err != nil {
					RenderError(tmpl, w, http.StatusInternalServerError, "Failed to check post likes/dislikes")
					return
				}
				topPosts[i].IsLiked = liked
				topPosts[i].IsDisliked = disliked
			}

			for i, post := range newPosts {
				// Check if the user has liked or disliked the post
				liked, disliked, err := service.Reaction().UserLikedDislikedPost(user.ID, post.ID)
				if err != nil {
					RenderError(tmpl, w, http.StatusInternalServerError, "Failed to check post likes/dislikes")
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

		if err := tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
			RenderError(tmpl, w, http.StatusInternalServerError, "Error: "+err.Error())
			return
		}
	}
}
