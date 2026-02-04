package models

import "time"

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
	Comments   int
}

// SinglePostResponse represents a complete post view with all comments
type SinglePost struct {
	ID         string
	Username   string
	Categories []Category
	Title      string
	Content    string
	CreatedAt  time.Time
	Likes      int
	Dislikes   int
	Comments   []Comment
}

// Comment represents a comment on a post with its associated data
type Comment struct {
	ID        string
	Username  string
	Content   string
	CreatedAt time.Time
	Likes     int
	Dislikes  int
}
