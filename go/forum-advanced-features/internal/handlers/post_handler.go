package handlers

import (
	"forum/internal/middleware"
	"forum/internal/models"
	"forum/internal/services"
	"forum/internal/utils/image"
	"forum/internal/utils/parser"
	"forum/internal/utils/validation"
	"html/template"
	"mime/multipart"
	"net/http"
	"strings"
)

func PostHandler(service services.Services, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user := middleware.GetUser(r)
		postID := r.PathValue("post_id")

		post, err := service.Post().GetSinglePost(postID)
		if err != nil {
			RenderError(tmpl, w, http.StatusInternalServerError, "Failed to get the post")
			return
		}

		var notifications int
		var commentHighlight string
		var userToHighlight string
		// check if user is logged in
		if user != nil {
			userToHighlight = r.URL.Query().Get("u")
			if userToHighlight != "" {
				if userToHighlight == "You" {
					userToHighlight = user.Username
				}
			}
			// check if the user came from notification page
			notificationID := r.URL.Query().Get("n")
			if notificationID != "" {
				// mark the notification as read
				err := service.Notify().MarkNotificationAsRead(notificationID, user.ID)
				if err != nil {
					RenderError(tmpl, w, http.StatusInternalServerError, err.Error())
					return
				}
			}

			commentID := r.URL.Query().Get("c")
			if commentID != "" {
				commentHighlight = commentID
			}

			// get notifications for the user
			notificationCount, err := service.Notify().GetUnreadNotificationCount(user.ID)
			if err != nil {
				RenderError(tmpl, w, http.StatusInternalServerError, "Failed to get notifications")
				return
			}
			notifications = notificationCount
			// check if user has liked or disliked the post
			liked, disliked, err := service.Reaction().UserLikedDislikedPost(user.ID, postID)
			if err != nil {
				RenderError(tmpl, w, http.StatusInternalServerError, "Failed to check like/dislike status")
				return
			}
			post.IsLiked = liked
			post.IsDisliked = disliked

			// check if the user has liked or disliked any of the comments
			for i, comment := range post.Comments {
				liked, disliked, err := service.Reaction().UserLikedDislikedComment(user.ID, comment.ID)
				if err != nil {
					RenderError(tmpl, w, http.StatusInternalServerError, "Failed to check comment like/dislike status")
					return
				}
				post.Comments[i].IsLiked = liked
				post.Comments[i].IsDisliked = disliked
			}
		}

		data := models.PostHandlerData{
			User:             user,
			SinglePost:       post,
			ValidationErr:    "",
			Notifications:    notifications,
			CommentHighlight: commentHighlight,
			UserHighlight:    userToHighlight,
		}

		tmpl.ExecuteTemplate(w, "post.html", data)
	}
}

// CreatePostHandler handles the create post page
func CreatePostHandler(service services.Services, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the current user
		user := middleware.GetUser(r)
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Get all categories for the form
		categories, err := service.Category().GetAllCategories()
		if err != nil {
			RenderError(tmpl, w, http.StatusInternalServerError, "Failed to get categories")
			return
		}

		notifications, err := service.Notify().GetUnreadNotificationCount(user.ID)
		if err != nil {
			RenderError(tmpl, w, http.StatusInternalServerError, "Failed to get notifications")
			return
		}

		// Common data for all responses
		data := models.CreatePostData{
			User:          user,
			Categories:    categories,
			FormData:      make(map[string]string),
			Notifications: notifications,
		}

		// If POST request, handle form submission
		if r.Method == http.MethodPost {
			var code int
			// Parse multipart form (for file uploads)
			data, code, err = parser.ParseValuesToCreatePost(r, data, categories, service.Config())

			if err != nil {
				data.ValidationErr = err.Error()
				w.WriteHeader(code)
				tmpl.ExecuteTemplate(w, "create-post.html", data)
				return
			}
			// prepare the image and validate
			var imageFile multipart.File
			var imageHeader *multipart.FileHeader
			var hasImage bool

			// Form Image
			imageFile, imageHeader, err = r.FormFile("image")
			if err != nil {
				if err.Error() != "http: no such file" {
					RenderError(tmpl, w, http.StatusInternalServerError, "Failed to form file")
				}

			}

			hasImage, code, err = validation.ValidateImage(imageHeader, imageFile, service.Config().MaxImageSize)
			if err != nil {
				RenderError(tmpl, w, code, err.Error())
			}

			var createPostData = models.PostData{
				UserID:      user.ID,
				Title:       data.FormData["title"],
				Content:     data.FormData["content"],
				CategoryIDs: data.FormCategories,
				HasImage:    hasImage,
				ImageFile:   imageFile,
				ImageHeader: imageHeader,
			}

			// Create post in the database
			postID, err := service.Post().CreatePost(createPostData)
			if err != nil {
				data.ValidationErr = "Failed to create post. Please try again."
				tmpl.ExecuteTemplate(w, "create-post.html", data)
				return
			}

			// Redirect to the newly created post
			http.Redirect(w, r, "/posts/"+postID, http.StatusSeeOther)
			return
		}

		// For GET request, just render the create post page
		tmpl.ExecuteTemplate(w, "create-post.html", data)
	}
}

func DeletePost(service services.Services, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// check method
		if r.Method != http.MethodPost {
			RenderError(tmpl, w, http.StatusMethodNotAllowed, "this method is not allowed :)")
			return
		}

		// get user
		user := middleware.GetUser(r)

		if user == nil {
			RenderError(tmpl, w, http.StatusUnauthorized, "You are not authorized to do that")
		}

		// get polst id from path
		postID := r.PathValue("post_id")

		// delete the post
		err, statusCode := service.Post().DeletePost(postID, user.ID)
		if err != nil {
			RenderError(tmpl, w, statusCode, err.Error())
		}

		referer := r.Header.Get("Referer")
		if strings.Contains(referer, "posts") {
			referer = "/"
		}

		http.Redirect(w, r, referer, http.StatusSeeOther)

	}
}

// EditPostHandler handles the edit post page
func EditPostHandler(service services.Services, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the current user
		user := middleware.GetUser(r)
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Get post ID from URL
		postID := r.PathValue("post_id")

		// Get the post
		post, err := service.Post().GetSinglePost(postID)
		if err != nil {
			RenderError(tmpl, w, http.StatusNotFound, "Post not found")
			return
		}

		// Check if user owns this post
		if post.Username != user.Username {
			RenderError(tmpl, w, http.StatusForbidden, "You can only edit your own posts")
			return
		}

		// Get all categories for the form
		categories, err := service.Category().GetAllCategories()
		if err != nil {
			RenderError(tmpl, w, http.StatusInternalServerError, "Failed to get categories")
			return
		}

		// Common data for all responses
		data := models.EditPostData{
			User:       user,
			Categories: categories,
			Post:       post,
			FormData:   make(map[string]string),
		}

		// If GET request, pre-populate form data
		if r.Method == http.MethodGet {
			data.FormData["title"] = post.Title
			data.FormData["content"] = post.Content
			// Clear any existing categories
			data.FormCategories = make([]string, 0)
			// Add the post's current categories
			for _, category := range post.Categories {
				data.FormCategories = append(data.FormCategories, category.ID)
			}
			tmpl.ExecuteTemplate(w, "edit.html", data)
			return
		}

		// If POST request, handle form submission
		if r.Method == http.MethodPost {
			var code int
			removeImage, code, data, err := parser.ParseValuesToEdit(r, data, categories, service.Config())
			if err != nil {
				data.ValidationErr = err.Error()
				w.WriteHeader(code)
				tmpl.ExecuteTemplate(w, "edit.html", data)
				return
			}
			var imageFile multipart.File
			var imageHeader *multipart.FileHeader
			var hasNewImage bool

			imageFile, imageHeader, err = r.FormFile("image")
			if err != nil {
				if err.Error() != "http: no such file" {
					RenderError(tmpl, w, http.StatusInternalServerError, err.Error())
				}

			}
			var errCode int
			hasNewImage, errCode, err = validation.ValidateImage(imageHeader, imageFile, service.Config().MaxImageSize)
			if err != nil {
				RenderError(tmpl, w, errCode, err.Error())
				return
			}

			if !removeImage {
				if hasNewImage {
					if post.HasImage {
						err := image.RemoveImageLocaly(service.Config().UploadDir, post.Image)
						if err != nil {
							RenderError(tmpl, w, http.StatusInternalServerError, "Failed to delete old image")
							return
						}
					}

					//Save new image locally
					filepath, err := image.SaveImageLocally(service.Config().UploadDir, imageHeader, imageFile)
					if err != nil {
						RenderError(tmpl, w, http.StatusInternalServerError, "Failed to save new image")
						return
					}
					//updating the post Struct and table
					err = service.Post().UpdatePost(postID, user.ID, data.FormData["title"], data.FormData["content"], data.FormCategories, true, filepath)
					if err != nil {
						RenderError(tmpl, w, http.StatusInternalServerError, "Failed to update post")
						return
					}
				} else {
					//updating the post Struct and table
					err = service.Post().UpdatePost(postID, user.ID, data.FormData["title"], data.FormData["content"], data.FormCategories, post.HasImage, post.Image)
					if err != nil {
						RenderError(tmpl, w, http.StatusInternalServerError, "Failed to update post")
						return
					}
				}

			} else if removeImage {
				//updating the post Struct and table
				err = service.Post().UpdatePost(postID, user.ID, data.FormData["title"], data.FormData["content"], data.FormCategories, false, "")
				if err != nil {
					RenderError(tmpl, w, http.StatusInternalServerError, "Failed to update post")
					return
				}
			}
			// Redirect to the updated post
			http.Redirect(w, r, "/posts/"+postID, http.StatusSeeOther)
			return
		}

		// Method not allowed
		RenderError(tmpl, w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}
