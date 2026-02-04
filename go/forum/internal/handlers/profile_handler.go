package handlers

import (
	"forum/internal/middleware"
	"forum/internal/models"
	"forum/internal/services"
	"html/template"
	"net/http"
	"strconv"
)

type profileData struct {
	User                *models.User
	UserPosts           []models.Post
	UserPostsPage       int
	UserPostsPages      int
	LikedPosts          []models.Post
	LikedPostsPage      int
	LikedPostsPages     int
	CommentedPosts      []models.Post
	CommentedPostsPage  int
	CommentedPostsPages int
}

func ProfileHandler(postService *services.PostService, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// method
		if r.Method != http.MethodGet {
			RenderError(tmpl, w, http.StatusMethodNotAllowed, "This method is not allowed")
			return
		}

		// get the user from context (logged in user)
		user := middleware.GetUser(r)
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Use the logged-in user's ID
		userID := user.ID

		// Get page parameters for different sections
		userPostsPage, likedPostsPage, commentedPostsPage := ValidatePagesProfile(w, r, tmpl)
		if userPostsPage == 0 { // Error occurred in validation
			return
		}

		// Get user's posts with pagination
		userPosts, userPostsPages, err := postService.GetPaginatedUserPosts(userID, userPostsPage, 5)
		if err != nil {
			RenderError(tmpl, w, http.StatusInternalServerError, "Failed to get user posts")
			return
		}
		if userPostsPage > userPostsPages && userPostsPages > 0 {
			// Redirect to the last valid user_posts_page
			query := r.URL.Query()
			query.Set("user_posts_page", strconv.Itoa(userPostsPages))
			redirectURL := "/profile?" + query.Encode()
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return
		}

		// Get liked posts with pagination
		likedPosts, likedPostsPages, err := postService.GetPaginatedUserLikedPosts(userID, likedPostsPage, 2)
		if err != nil {
			RenderError(tmpl, w, http.StatusInternalServerError, "Failed to get liked posts")
			return
		}
		if likedPostsPage > likedPostsPages && likedPostsPages > 0 {
			// Redirect to the last valid liked_posts_page
			query := r.URL.Query()
			query.Set("liked_posts_page", strconv.Itoa(likedPostsPages))
			redirectURL := "/profile?" + query.Encode()
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return
		}

		// Get commented posts with pagination
		commentedPosts, commentedPostsPages, err := postService.GetPaginatedUserCommentedPosts(userID, commentedPostsPage, 2)
		if err != nil {
			RenderError(tmpl, w, http.StatusInternalServerError, "Failed to get commented posts")
			return
		}
		if commentedPostsPage > commentedPostsPages && commentedPostsPages > 0 {
			// Redirect if commentedPostsPage is out of range
			query := r.URL.Query()
			query.Set("commented_posts_page", strconv.Itoa(commentedPostsPages))
			redirectURL := "/profile?" + query.Encode()
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return
		}

		data := profileData{
			User:                user,
			UserPosts:           userPosts,
			UserPostsPage:       userPostsPage,
			UserPostsPages:      userPostsPages,
			LikedPosts:          likedPosts,
			LikedPostsPage:      likedPostsPage,
			LikedPostsPages:     likedPostsPages,
			CommentedPosts:      commentedPosts,
			CommentedPostsPage:  commentedPostsPage,
			CommentedPostsPages: commentedPostsPages,
		}

		tmpl.ExecuteTemplate(w, "profile.html", data)
	}
}

// ValidatePagesProfile validates page parameters for profile page
func ValidatePagesProfile(w http.ResponseWriter, r *http.Request, tmpl *template.Template) (int, int, int) {
	var userPostsPage, likedPostsPage, commentedPostsPage int = 1, 1, 1

	if pageStr := r.URL.Query().Get("user_posts_page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			userPostsPage = p
		} else {
			RenderError(tmpl, w, http.StatusBadRequest, "Not valid 'user_posts_page' number")
			return 0, 0, 0
		}
	}

	if pageStr := r.URL.Query().Get("liked_posts_page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			likedPostsPage = p
		} else {
			RenderError(tmpl, w, http.StatusBadRequest, "Not valid 'liked_posts_page' number")
			return 0, 0, 0
		}
	}

	if pageStr := r.URL.Query().Get("commented_posts_page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			commentedPostsPage = p
		} else {
			RenderError(tmpl, w, http.StatusBadRequest, "Not valid 'commented_posts_page' number")
			return 0, 0, 0
		}
	}

	return userPostsPage, likedPostsPage, commentedPostsPage
}
