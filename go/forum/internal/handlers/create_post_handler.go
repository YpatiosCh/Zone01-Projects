package handlers

import (
	"fmt"
	"forum/internal/middleware"
	"forum/internal/models"
	"forum/internal/services"
	"forum/internal/utils/validation"
	"html/template"
	"net/http"
	"strings"
)

type createPostData struct {
	User          *models.User
	Categories    []models.Category
	ValidationErr string
	SuccessMsg    string
	FormData      map[string]string
}

// CreatePostHandler handles the create post page
func CreatePostHandler(postService *services.PostService, categoryService *services.CategoryService, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the current user
		user := middleware.GetUser(r)
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Get all categories for the form
		categories, err := categoryService.GetAllCategories()
		if err != nil {
			RenderError(tmpl, w, http.StatusInternalServerError, "Failed to get categories")
			return
		}

		// Common data for all responses
		data := createPostData{
			User:       user,
			Categories: categories,
			FormData:   make(map[string]string),
		}

		// If POST request, handle form submission
		if r.Method == http.MethodPost {
			// Parse form
			if err := r.ParseForm(); err != nil {
				data.ValidationErr = "Invalid form data"
				tmpl.ExecuteTemplate(w, "create-post.html", data)
				return
			}

			// Get form values
			title := r.FormValue("title")
			content := r.FormValue("content")
			categoryIDs := r.Form["categories"] // Gets multiple selected values

			// Save form data for re-displaying on error
			data.FormData["title"] = title
			data.FormData["content"] = content
			// Sanitize and validate input
			// title = validation.SanitizeInput(title)
			// content = validation.SanitizeInput(content)

			// Validate the post using our validation utilities
			errorStruct := validation.ValidatePost(title, content, categoryIDs)
			if errorStruct.Error() != nil {
				// If it's a validation error, display it to the user
				data.ValidationErr = strings.Join(errorStruct.Error(), " & ")
				fmt.Println(data.ValidationErr)
				tmpl.ExecuteTemplate(w, "create-post.html", data)
				return

			}

			// Create post in the database
			postID, err := postService.CreatePost(user.ID, title, content, categoryIDs)
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
