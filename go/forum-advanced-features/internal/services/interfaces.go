package services

import (
	"forum/internal/config"
	"forum/internal/models"
	"forum/internal/repository"
)

type Services interface {
	User() UserInterface
	Reaction() ReactionInterface
	Comment() CommentInterface
	Post() PostInterface
	Auth() AuthInterface
	OAuth() OAuthInterface
	Category() CategoryInterface
	Notify() NotificationInterface
	Config() *config.AppConfig
}

type ServiceContainer struct {
	userService         UserInterface
	reactionService     ReactionInterface
	commentService      CommentInterface
	postService         PostInterface
	authService         AuthInterface
	oauthService        OAuthInterface
	categoryService     CategoryInterface
	notificationService NotificationInterface
	appConfig           *config.AppConfig
}

func NewServiceContainer(repo *repository.Manager, config *config.AppConfig) *ServiceContainer {
	// Initialize services in dependency order
	authService := NewAuthService(repo)
	categoryService := NewCategoryService(repo)
	userService := NewUserService(repo)
	notificationService := NewNotificationService(repo)
	reactionService := NewReactionService(repo, notificationService)
	commentService := NewCommentService(repo, reactionService, notificationService)
	oauthService := NewOAuthService(repo, config)
	postService := NewPostService(repo, reactionService, categoryService, userService, commentService, config)

	appConfig := config

	return &ServiceContainer{
		userService:         userService,
		reactionService:     reactionService,
		commentService:      commentService,
		authService:         authService,
		oauthService:        oauthService,
		categoryService:     categoryService,
		notificationService: notificationService,
		postService:         postService,
		appConfig:           appConfig,
	}
}

// Confic returns the app configuration implemented
func (s *ServiceContainer) Config() *config.AppConfig {
	return s.appConfig
}

// User returns the UserInterface implementation
// for accessing user-related services.
func (s *ServiceContainer) User() UserInterface {
	return s.userService
}

// Reaction returns the ReactionInterface implementation
// for accessing reaction-related services.
func (s *ServiceContainer) Reaction() ReactionInterface {
	return s.reactionService
}

// Comment returns the CommentInterface implementation
// for accessing comment-related services.
func (s *ServiceContainer) Comment() CommentInterface {
	return s.commentService
}

// Post returns the PostInterface implementation
// for accessing post-related services.
func (s *ServiceContainer) Post() PostInterface {
	return s.postService
}

// Auth returns the AuthInterface implementation
// for accessing authentication-related services.
func (s *ServiceContainer) Auth() AuthInterface {
	return s.authService
}

// OAuth returns the OAuthInterface implementation
// for accessing OAuth-related services.
func (s *ServiceContainer) OAuth() OAuthInterface {
	return s.oauthService
}

// Category returns the CategoryInterface implementation
// for accessing category-related services.
func (s *ServiceContainer) Category() CategoryInterface {
	return s.categoryService
}

// Notify returns the NotificationInterface implementation
// for accessing notification-related services.
func (s *ServiceContainer) Notify() NotificationInterface {
	return s.notificationService
}

// ------ Interfaces for each service ------

// UserInterface defines methods for user-related operations.
type UserInterface interface {
	UserExist(username string) bool
}

// CategoryInterface defines methods for category-related operations.
type CategoryInterface interface {
	GetAllCategories() ([]models.Category, error)
	GetCategoriesForPost(postID string) ([]models.Category, error)
}

// CommentInterface defines methods for comment-related operations.
type CommentInterface interface {
	CreateComment(userID, postID, content string) (string, error)
	ConvertToCommentStruct(comments []map[string]interface{}) ([]models.Comment, error)
	UpdateComment(userID, commentID, content string) (int, error)
	GetSingleComment(commentId string) (string, models.Comment, error)
	DeleteComment(userID, commentID string) (int, error)
}

// ReactionInterface defines methods for handling reactions (likes/dislikes).
type ReactionInterface interface {
	ToggleReaction(userID, table, ID, reactionType string) error
	UserLikedDislikedComment(userID, commentID string) (bool, bool, error)
	UserLikedDislikedPost(userID, postID string) (bool, bool, error)
	GetReactionsForComment(commentID string) (int, int, error)
	GetReactionsForPost(postID string) (int, int, error)
}

// AuthInterface defines methods for user authentication and registration.
type AuthInterface interface {
	RegisterUser(username, email, password string) (string, Error)
	Login(email, password string) (string, Error)
	Logout(sessionID string) Error
}

// OAuthInterface defines methods for handling OAuth authentication.
type OAuthInterface interface {
	GetGoogleAuthURL(state string) string
	GetGitHubAuthURL(state string) string
	HandleGoogleCallback(code string) (string, Error)
	HandleGitHubCallback(code string) (string, Error)
	handleOAuthUser(provider string, userInfo *models.User) (string, Error)
	CheckExistingOAuthUser(provider, providerID string) (*models.User, error)
	ExchangeGoogleCode(code string) (string, error)
	GetGoogleUserInfo(token string) (*models.User, error)
	ExchangeGitHubCode(code string) (string, error)
	GetGitHubUserInfo(token string) (*models.User, error)
	CreateOAuthUser(provider, providerID, email, name, username string) (string, Error)
	LoginoAuthUser(userID string) (string, Error)
}

// PostInterface defines methods for post-related operations.
type PostInterface interface {
	GetSinglePost(postID string) (models.Post, error)
	GetPaginatedAllPostsByEngagement(page, perPage int) ([]models.Post, int, error)
	GetPaginatedAllNewestPosts(page, perPage int) ([]models.Post, int, error)
	GetPaginatedUserPosts(username string, page, perPage int) ([]models.Post, int, error)
	GetPaginatedUserLikedPosts(username string, page, perPage int) ([]models.Post, int, error)
	GetPaginatedUserCommentedPosts(username string, page, perPage int) ([]models.Post, int, error)
	GetPaginatedUserDislikedPosts(username string, page, perPage int) ([]models.Post, int, error)
	GetPaginatedPostsByCategoryByEngagement(categoryID string, page, perPage int) ([]models.Post, int, error)
	GetPaginatedPostsByCategoryNewest(categoryID string, page, perPage int) ([]models.Post, int, error)
	ConvertToPostStruct(posts []map[string]interface{}) ([]models.Post, error)
	CreatePost(data models.PostData) (string, error)
	UpdatePost(postID, userID, title, content string, categoryIDs []string, hasImage bool, filepath string) error
	DeletePost(postID string, username string) (error, int)
}

// NotificationInterface defines methods for handling notifications.
type NotificationInterface interface {
	CreateNotification(userID, notificationType, sourceUserID, message string, postID, commentID *string) (string, error)
	GetUserNotifications(userID string) ([]models.Notification, error)
	GetUnreadNotificationCount(userID string) (int, error)
	MarkNotificationAsRead(notificationID, userID string) error
	DeleteNotification(notificationID, userID string) error
}
