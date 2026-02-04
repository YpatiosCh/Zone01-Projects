package handlers

import (
	"forum/internal/middleware"
	"forum/internal/models"
	"forum/internal/services"
	"html/template"
	"net/http"
)

type postHandlerData struct {
	User          *models.User
	SinglePost    models.SinglePost
	ValidationErr string
}

func PostHandler(postService *services.PostService, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user := middleware.GetUser(r)
		postID := r.PathValue("post_id")

		post, err := postService.GetSinglePost(postID)
		if err != nil {
			RenderError(tmpl, w, http.StatusInternalServerError, "Failed to get the post")
			return
		}

		data := postHandlerData{
			User:          user,
			SinglePost:    post,
			ValidationErr: "",
		}

		tmpl.ExecuteTemplate(w, "post.html", data)
	}
}
