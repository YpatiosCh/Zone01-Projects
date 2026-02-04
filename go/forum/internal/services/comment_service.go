package services

import "forum/internal/repository"

type CommentService struct {
	repo *repository.Manager
}

func NewCommentService(repo *repository.Manager) *CommentService {
	return &CommentService{
		repo: repo,
	}
}

func (c *CommentService) CreateComment(userID, postID, content string) (string, error) {

	// start transaction
	tx, err := c.repo.Db.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	// prepare data to create the comment
	commentData := map[string]interface{}{
		"user_id": userID,
		"post_id": postID,
		"content": content,
	}

	// Create the comment in the database
	commentID, err := c.repo.CreateRecord("comment", commentData)
	if err != nil {
		return "", err
	}

	// commit transaction
	if err = tx.Commit(); err != nil {
		return "", err
	}

	return commentID, nil
}
