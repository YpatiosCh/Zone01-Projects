package handlers

import (
	"fmt"
	"forum/internal/middleware"
	"forum/internal/models"
	"forum/internal/services"
	"html/template"
	"net/http"
	"strconv"
)

func ProfileHandler(service services.Services, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// method
		if r.Method != http.MethodGet {
			RenderError(tmpl, w, http.StatusMethodNotAllowed, "This method is not allowed")
			return
		}

		// get real user
		realUser := middleware.GetUser(r)

		// get notifications for the user
		notifications, err := service.Notify().GetUnreadNotificationCount(realUser.ID)
		if err != nil {
			RenderError(tmpl, w, http.StatusInternalServerError, "Failed to get notifications")
			return
		}

		// get user from query parameters
		username := r.PathValue("username")
		if !service.User().UserExist(username) {
			redirectUrl := "/profile/" + realUser.Username
			http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
			return
		}

		var userHolder string
		if realUser.Username != username {
			userHolder = username
		} else {
			userHolder = "You"
		}

		// Get page parameters for different sections
		userPostsPage, likedPostsPage, commentedPostsPage, dislikedPostsPage := ValidatePagesProfile(w, r, tmpl)
		if userPostsPage == 0 { // Error occurred in validation
			return
		}

		// Get user's posts with pagination
		userPosts, userPostsPages, err := service.Post().GetPaginatedUserPosts(username, userPostsPage, 2)
		if err != nil {
			RenderError(tmpl, w, http.StatusInternalServerError, err.Error())
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
		likedPosts, likedPostsPages, err := service.Post().GetPaginatedUserLikedPosts(username, likedPostsPage, 2)
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
		commentedPosts, commentedPostsPages, err := service.Post().GetPaginatedUserCommentedPosts(username, commentedPostsPage, 2)
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

		// Get disliked posts with pagination
		dislikedPosts, dislikedPostPages, err := service.Post().GetPaginatedUserDislikedPosts(username, dislikedPostsPage, 2)
		if err != nil {
			RenderError(tmpl, w, http.StatusInternalServerError, "Failed to get commented posts")
			return
		}
		if dislikedPostsPage > dislikedPostPages && dislikedPostPages > 0 {
			// Redirect if dislikedPostPage is out of range
			query := r.URL.Query()
			query.Set("disliked_posts_page", strconv.Itoa(dislikedPostPages))
			redirectURL := "/profile?" + query.Encode()
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return
		}

		// Get user's likes/dislikes for the posts
		for i, post := range userPosts {
			liked, disliked, err := service.Reaction().UserLikedDislikedPost(realUser.ID, post.ID)
			if err != nil {
				RenderError(tmpl, w, http.StatusInternalServerError, "Failed to check post likes/dislikes")
				return
			}
			userPosts[i].IsLiked = liked
			userPosts[i].IsDisliked = disliked
		}
		for i, post := range likedPosts {
			liked, disliked, err := service.Reaction().UserLikedDislikedPost(realUser.ID, post.ID)
			if err != nil {
				RenderError(tmpl, w, http.StatusInternalServerError, "Failed to check post likes/dislikes")
				return
			}
			likedPosts[i].IsLiked = liked
			likedPosts[i].IsDisliked = disliked
		}
		for i, post := range commentedPosts {
			liked, disliked, err := service.Reaction().UserLikedDislikedPost(realUser.ID, post.ID)
			if err != nil {
				RenderError(tmpl, w, http.StatusInternalServerError, "Failed to check post likes/dislikes")
				return
			}
			commentedPosts[i].IsLiked = liked
			commentedPosts[i].IsDisliked = disliked
		}
		for i, post := range dislikedPosts {
			liked, disliked, err := service.Reaction().UserLikedDislikedPost(realUser.ID, post.ID)
			if err != nil {
				RenderError(tmpl, w, http.StatusInternalServerError, "Failed to check post likes/disliked")
				return
			}
			dislikedPosts[i].IsLiked = liked
			dislikedPosts[i].IsDisliked = disliked
		}

		data := models.ProfileData{
			User:                realUser,
			UserPosts:           userPosts,
			UserPostsPage:       userPostsPage,
			UserPostsPages:      userPostsPages,
			LikedPosts:          likedPosts,
			LikedPostsPage:      likedPostsPage,
			LikedPostsPages:     likedPostsPages,
			CommentedPosts:      commentedPosts,
			CommentedPostsPage:  commentedPostsPage,
			CommentedPostsPages: commentedPostsPages,
			DislikedPosts:       dislikedPosts,
			DislikedPostsPage:   dislikedPostsPage,
			DislikedPostsPages:  dislikedPostPages,
			UserHolder:          userHolder,
			Notifications:       notifications,
		}

		err = tmpl.ExecuteTemplate(w, "profile.html", data)
		if err != nil {
			fmt.Println("can't execute profile template", err)
		}
	}
}

// ValidatePagesProfile validates page parameters for profile page
func ValidatePagesProfile(w http.ResponseWriter, r *http.Request, tmpl *template.Template) (int, int, int, int) {
	var userPostsPage, likedPostsPage, commentedPostsPage, dislikedPostsPage int = 1, 1, 1, 1

	if pageStr := r.URL.Query().Get("user_posts_page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			userPostsPage = p
		} else {
			RenderError(tmpl, w, http.StatusBadRequest, "Not valid 'user_posts_page' number")
			return 0, 0, 0, 0
		}
	}

	if pageStr := r.URL.Query().Get("liked_posts_page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			likedPostsPage = p
		} else {
			RenderError(tmpl, w, http.StatusBadRequest, "Not valid 'liked_posts_page' number")
			return 0, 0, 0, 0
		}
	}

	if pageStr := r.URL.Query().Get("commented_posts_page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			commentedPostsPage = p
		} else {
			RenderError(tmpl, w, http.StatusBadRequest, "Not valid 'commented_posts_page' number")
			return 0, 0, 0, 0
		}
	}

	if pageStr := r.URL.Query().Get("disliked_posts_page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			dislikedPostsPage = p
		} else {
			RenderError(tmpl, w, http.StatusBadRequest, "Not valid 'commented_posts_page' number")
			return 0, 0, 0, 0
		}
	}

	return userPostsPage, likedPostsPage, commentedPostsPage, dislikedPostsPage
}
