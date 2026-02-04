package models

import (
	"time"
)

type User struct {
	ID       string
	Email    string
	Username string
}

type Category struct {
	ID   string
	Name string
}

type Post struct {
	ID         string
	Username   string
	Title      string
	Content    string
	CreatedAt  time.Time
	Categories []Category
	Likes      int
	Dislikes   int
	IsLiked    bool
	IsDisliked bool
	Comments   []Comment
	HasImage   bool
	Image      string
}

// Comment represents a comment on a post with its associated data
type Comment struct {
	ID         string
	Username   string
	Content    string
	CreatedAt  time.Time
	Likes      int
	Dislikes   int
	IsLiked    bool
	IsDisliked bool
}

type Notification struct {
	ID           string
	UserID       string
	Type         string // "like", "dislike", "comment"
	SourceUserID string
	PostID       *string // pointer because it can be nil
	CommentID    *string // pointer because it can be nil
	Message      string
	Read         bool
	CreatedAt    time.Time
	SourceUser   string // Username of the user who triggered the notification
}
