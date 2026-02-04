package handlers

import (
	"forum/internal/middleware"
	"forum/internal/models"
	"forum/internal/services"
	"forum/internal/utils/validation"
	"html/template"
	"net/http"
	"strings"
)

func CreateComment(commentService *services.CommentService, postService *services.PostService, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only allow POST method
		if r.Method != http.MethodPost {
			RenderError(tmpl, w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		// Get authenticated user
		user := middleware.GetUser(r)
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Get post ID from URL path
		postID := r.PathValue("post_id")

		// Get the post data
		post, err := postService.GetSinglePost(postID)
		if err != nil {
			RenderError(tmpl, w, http.StatusInternalServerError, "Failed to get the post")
			return
		}

		// Parse form
		if err := r.ParseForm(); err != nil {
			// Prepare data with error
			data := struct {
				User          *models.User
				SinglePost    models.SinglePost
				ValidationErr string
			}{
				User:          user,
				SinglePost:    post,
				ValidationErr: "Failed to parse form",
			}

			// Render the template with error
			tmpl.ExecuteTemplate(w, "post.html", data)
			return
		}

		// Get comment content
		content := r.FormValue("content")

		// Validate content
		Error := validation.ValidateComment(content)
		if Error.ErrorSlice != nil {
			// Prepare data with error
			data := struct {
				User          *models.User
				SinglePost    models.SinglePost
				ValidationErr string
				FormData      map[string]string
			}{
				User:          user,
				SinglePost:    post,
				ValidationErr: strings.Join(Error.Error(), " & "),
				FormData:      map[string]string{"content": content},
			}

			// Render the template with error
			tmpl.ExecuteTemplate(w, "post.html", data)
			return
		}

		// sanitize content
		// content = validation.SanitizeInput(content)

		// Create comment
		_, err = commentService.CreateComment(user.ID, postID, content)
		if err != nil {
			// Prepare data with error
			data := struct {
				User          *models.User
				SinglePost    models.SinglePost
				ValidationErr string
				FormData      map[string]string
			}{
				User:          user,
				SinglePost:    post,
				ValidationErr: "Failed to create comment: " + err.Error(),
				FormData:      map[string]string{"content": content},
			}

			// Render the template with error
			tmpl.ExecuteTemplate(w, "post.html", data)
			return
		}

		// Redirect back to post page
		http.Redirect(w, r, "/posts/"+postID, http.StatusSeeOther)
	}
}
