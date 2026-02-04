package routes

import (
	"forum/internal/config"
	"forum/internal/handlers"
	"forum/internal/middleware"
	"forum/internal/services"
	"forum/internal/utils"
	"html/template"
	"net/http"
)

// SetUpRoutes sets up the routes for the application
func SetUpRoutes(middleware *middleware.Middleware, services services.Services, config *config.AppConfig) http.Handler {

	// parse templates
	templates := template.Must(template.New("").Funcs(utils.TemplateFuncs()).ParseGlob("templates/*.html"))

	// set up routes
	mux := http.NewServeMux()

	// serve static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	mux.Handle("/", handlers.HomePage(services, templates))
	mux.Handle("/signup", handlers.RegisterUser(services, templates))
	mux.Handle("/login", handlers.LoginUser(services, templates))
	mux.Handle("/categories/{category_id}/posts", handlers.PostsByCategory(services, templates))
	mux.Handle("/posts/{post_id}", handlers.PostHandler(services, templates))

	// OAuth routes
	mux.Handle("/auth/google", handlers.AuthHandler(services, services.OAuth().GetGoogleAuthURL))
	mux.Handle("/auth/google/callback", handlers.GoogleCallbackHandler(services, templates))
	mux.Handle("/auth/github", handlers.AuthHandler(services, services.OAuth().GetGitHubAuthURL))
	mux.Handle("/auth/github/callback", handlers.GitHubCallbackHandler(services, templates))

	// OAuth username selection routes
	mux.Handle("/auth/username", handlers.OAuthUsernameHandler(templates))
	mux.Handle("/auth/set-username", handlers.SetUsernameHandler(services, templates))

	// protected routes
	mux.Handle("/logout", middleware.RequireAuth(handlers.LogoutUser(services, templates)))
	mux.Handle("/profile/{username}", middleware.RequireAuth(handlers.ProfileHandler(services, templates)))
	mux.Handle("/create-post", middleware.RequireAuth(handlers.CreatePostHandler(services, templates)))
	mux.Handle("/posts/{post_id}/edit", middleware.RequireAuth(handlers.EditPostHandler(services, templates)))
	mux.Handle("/posts/{post_id}/delete", middleware.RequireAuth(handlers.DeletePost(services, templates)))
	mux.Handle("/posts/{post_id}/comments", middleware.RequireAuth(handlers.CreateComment(services, templates)))
	mux.Handle("/post/{id}/{reaction_type}", middleware.RequireAuth(handlers.PostReactionHandler(services, templates)))
	mux.Handle("/comment/{id}/delete", middleware.RequireAuth(handlers.DeleteComment(services, templates)))
	mux.Handle("/comment/{id}/edit", middleware.RequireAuth(handlers.UpdateCommentHandler(services, templates)))
	mux.Handle("/comment/{id}/{reaction_type}", middleware.RequireAuth(handlers.PostReactionHandler(services, templates)))
	mux.Handle("/notifications", middleware.RequireAuth(handlers.NotificationsHandler(services, templates)))
	mux.Handle("/notifications/{notification_id}/delete", middleware.RequireAuth(handlers.DeleteNotificationHandler(services, templates)))
	mux.Handle("/notifications/delete", middleware.RequireAuth(handlers.DeleteAllNotifications(services, templates)))
	mux.Handle("/notifications/mark-all", middleware.RequireAuth(handlers.MarkAllAsRead(services, templates)))
	return middleware.Auth(mux)
}
