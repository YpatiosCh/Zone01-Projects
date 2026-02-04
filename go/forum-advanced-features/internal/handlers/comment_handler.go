package handlers

import (
	"forum/internal/middleware"
	"forum/internal/models"
	"forum/internal/services"
	"forum/internal/utils/parser"
	"forum/internal/utils/validation"
	"html/template"
	"net/http"
	"strings"
)

func CreateComment(service services.Services, tmpl *template.Template) http.HandlerFunc {
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
		post, err := service.Post().GetSinglePost(postID)
		if err != nil {
			RenderError(tmpl, w, http.StatusInternalServerError, "Failed to get the post")
			return
		}

		// Parse form
		if err := r.ParseForm(); err != nil {
			// Prepare data with error
			data := models.CreateCommentData{
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
			data := models.CreateCommentData{
				User:          user,
				SinglePost:    post,
				ValidationErr: strings.Join(Error.Error(), " & "),
				FormData:      map[string]string{"content": content},
			}

			// Render the template with error
			tmpl.ExecuteTemplate(w, "post.html", data)
			return
		}

		// Create comment
		_, err = service.Comment().CreateComment(user.ID, postID, content)
		if err != nil {
			// Prepare data with error
			data := models.CreateCommentData{
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

func UpdateCommentHandler(service services.Services, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		//Get and check if the user exists
		user := middleware.GetUser(r)
		if user == nil {
			RenderError(tmpl, w, http.StatusUnauthorized, "Not allowed user")
			return
		}
		//Take the comment Id from url
		commentId := r.PathValue("id")

		postID, comment, err := service.Comment().GetSingleComment(commentId)
		if err != nil {
			RenderError(tmpl, w, http.StatusInternalServerError, "failed to get single comment")
			return
		}

		data := models.EditCommentData{
			User:    user,
			Comment: comment,
			Content: "",
		}

		content, err := parser.ParseValuesToEditComment(r, data)
		if err != nil {
			RenderError(tmpl, w, http.StatusBadRequest, err.Error())
			return
		}

		errCode, err := service.Comment().UpdateComment(user.ID, commentId, content)
		if err != nil {
			RenderError(tmpl, w, errCode, err.Error())
		}
		http.Redirect(w, r, "/posts/"+postID, http.StatusSeeOther)
	}
}

func DeleteComment(service services.Services, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//Get and check if the user exists
		user := middleware.GetUser(r)
		if user == nil {
			RenderError(tmpl, w, http.StatusUnauthorized, "Not allowed user")
			return
		}
		//Take the comment Id from url
		commentId := r.PathValue("id")

		postID, _, err := service.Comment().GetSingleComment(commentId)
		if err != nil {
			RenderError(tmpl, w, http.StatusInternalServerError, "failed to get single comment")
			return
		}

		errCode, err := service.Comment().DeleteComment(user.ID, commentId)
		if err != nil {
			RenderError(tmpl, w, errCode, err.Error())
		}
		http.Redirect(w, r, "/posts/"+postID, http.StatusSeeOther)
	}
}
