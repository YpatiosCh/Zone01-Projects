package models

import "mime/multipart"

type EditPostData struct {
	User           *User
	Categories     []Category
	Post           Post
	ValidationErr  string
	SuccessMsg     string
	FormData       map[string]string
	FormCategories []string
}

type EditCommentData struct {
	User          *User
	Comment       Comment
	ValidationErr string
	Content       string
}

type CreatePostData struct {
	User           *User
	Categories     []Category
	ValidationErr  string
	SuccessMsg     string
	FormData       map[string]string
	Notifications  int
	FormCategories []string
}

type PostData struct {
	UserID      string
	Title       string
	Content     string
	CategoryIDs []string
	HasImage    bool
	ImageFile   multipart.File
	ImageHeader *multipart.FileHeader
}

type CreateCommentData struct {
	User          *User
	SinglePost    Post
	ValidationErr string
	FormData      map[string]string
}

type ErrorHandlerData struct {
	StatusCode int
	Message    string
}

type HomeHandlerData struct {
	User          *User
	Categories    []Category
	TopPosts      []Post
	NewPosts      []Post
	Page          int
	TotalPages    int
	Notifications int
}

type NotificationsData struct {
	User          *User
	Notifications []Notification
	Unread        int
}

type OAuthData struct {
	Provider          string
	ProviderID        string
	Email             string
	Name              string
	SuggestedUsername string
	ValidationErr     string
	FormData          map[string]string
}

type PostHandlerData struct {
	User             *User
	SinglePost       Post
	ValidationErr    string
	Notifications    int
	CommentHighlight string
	UserHighlight    string
}

type ProfileData struct {
	User                *User
	UserPosts           []Post
	UserPostsPage       int
	UserPostsPages      int
	LikedPosts          []Post
	LikedPostsPage      int
	LikedPostsPages     int
	CommentedPosts      []Post
	CommentedPostsPage  int
	CommentedPostsPages int
	DislikedPosts       []Post
	DislikedPostsPage   int
	DislikedPostsPages  int
	UserHolder          string
	Notifications       int
}
