package routes

import (
	"database/sql"
	"forum/internal/handlers"
	"forum/internal/middleware"
	"forum/internal/repository"
	"forum/internal/services"
	"forum/internal/utils"
	"html/template"
	"net/http"
)

// SetUpRoutes sets up the routes for the application
func SetUpRoutes(db *sql.DB) http.Handler {
	// set up database manager
	dbManager := repository.NewManager(db)

	// set up services
	authService := services.NewAuthService(dbManager)
	postService := services.NewPostService(dbManager)
	categoryService := services.NewCategoryService(dbManager)
	commentService := services.NewCommentService(dbManager)
	reactionService := services.NewReactionService(dbManager)

	// set up middleware
	middleware := middleware.NewMiddleware(dbManager)

	// parse templates
	templates := template.Must(template.New("").Funcs(utils.TemplateFuncs()).ParseGlob("templates/*.html"))

	// set up routes
	mux := http.NewServeMux()

	// serve static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	mux.Handle("/", handlers.HomePage(postService, categoryService, templates))
	mux.Handle("/signup", handlers.RegisterUser(authService, templates))
	mux.Handle("/login", handlers.LoginUser(authService, templates))
	mux.Handle("/categories/{category_id}/posts", handlers.PostsByCategory(postService, categoryService, templates))
	mux.Handle("/posts/{post_id}", handlers.PostHandler(postService, templates))

	// protected routes
	mux.Handle("/logout", middleware.RequireAuth(handlers.LogoutUser(authService, templates)))
	mux.Handle("/profile", middleware.RequireAuth(handlers.ProfileHandler(postService, templates)))
	mux.Handle("/create-post", middleware.RequireAuth(handlers.CreatePostHandler(postService, categoryService, templates)))
	mux.Handle("/posts/{post_id}/comments", middleware.RequireAuth(handlers.CreateComment(commentService, postService, templates)))
	mux.Handle("/post/{id}/{reaction_type}", middleware.RequireAuth(handlers.PostReactionHandler(reactionService, templates)))
	mux.Handle("/comment/{id}/{reaction_type}", middleware.RequireAuth(handlers.PostReactionHandler(reactionService, templates)))
	return middleware.Auth(mux)
}
