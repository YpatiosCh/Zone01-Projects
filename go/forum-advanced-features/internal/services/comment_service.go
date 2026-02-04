package services

import (
	"errors"
	"forum/internal/models"
	"forum/internal/repository"
	"forum/internal/utils/validation"
	"net/http"
	"strings"
	"time"
)

type CommentService struct {
	repo         *repository.Manager
	reaction     *ReactionService
	notification *NotificationService // Add this line
}

func NewCommentService(repo *repository.Manager, reaction *ReactionService, notification *NotificationService) *CommentService {
	return &CommentService{
		repo:         repo,
		reaction:     reaction,
		notification: notification, // Add this line
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

	// Create a notification for the comment
	c.notification.CreateCommentNotification(userID, postID, commentID)
	return commentID, nil
}

// ConvertToCommentStruct - keep as helper
func (s *CommentService) ConvertToCommentStruct(comments []map[string]interface{}) ([]models.Comment, error) {
	var result []models.Comment
	for _, comment := range comments {
		// Get user
		user, err := s.repo.Get("user", "id", comment["user_id"].(string))
		if err != nil {
			return nil, err
		}
		username := user[0]["username"].(string)

		// Get reactions
		likes, dislikes, err := s.reaction.GetReactionsForComment(comment["id"].(string))
		if err != nil {
			return nil, err
		}

		result = append(result, models.Comment{
			ID:        comment["id"].(string),
			Username:  username,
			Content:   comment["content"].(string),
			CreatedAt: comment["created_at"].(time.Time),
			Likes:     likes,
			Dislikes:  dislikes,
		})
	}

	if len(result) == 0 {
		return nil, nil
	}

	return result, nil
}

func (s *CommentService) UpdateComment(userID, commentID, content string) (int, error) {
	commentInfo, err := s.repo.Get("comment", "id", commentID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	errStruct := validation.ValidateComment(content)
	if errStruct.ErrorSlice != nil {
		errmsg := strings.Join(errStruct.ErrorSlice, " & ")
		return http.StatusBadRequest, errors.New(errmsg)
	}

	comment := commentInfo[0]
	user := comment["user_id"]
	if user != userID {
		return http.StatusUnauthorized, errors.New("you can only edit your own comments")
	}
	data := map[string]interface{}{
		"content": content,
	}

	err = s.repo.UpdateRecord("comment", commentID, data)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (s *CommentService) GetSingleComment(commentId string) (string, models.Comment, error) {
	var result models.Comment
	// get comment by id
	comments, err := s.repo.Get("comment", "id", commentId)
	if err != nil {
		return "", result, err
	}
	comment := comments[0]

	// get the user who created the post
	user, err := s.repo.Get("user", "id", comment["user_id"].(string))
	if err != nil {
		return "", result, err
	}
	username := user[0]["username"].(string)

	// get reactions
	likes, dislikes, err := s.reaction.GetReactionsForComment(comment["id"].(string))
	if err != nil {
		return "", result, err
	}

	result = models.Comment{
		ID:        comment["id"].(string),
		Username:  username,
		Content:   comment["content"].(string),
		CreatedAt: comment["created_at"].(time.Time),
		Likes:     likes,
		Dislikes:  dislikes,
	}

	return comment["post_id"].(string), result, nil
}

func (s *CommentService) DeleteComment(userID, commentID string) (int, error) {
	commentInfo, err := s.repo.Get("comment", "id", commentID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	comment := commentInfo[0]
	user := comment["user_id"]
	if user != userID {
		return http.StatusUnauthorized, errors.New("you can only edit your own comments")
	}
	commentReactions, err := s.repo.Get("reaction", "comment_id", commentID)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if len(commentReactions) > 0 {
		for _, commentReaction := range commentReactions {
			commentReactionID := commentReaction["id"].(string)
			err = s.repo.DeleteRecord("reaction", commentReactionID)
			if err != nil {
				return http.StatusInternalServerError, err
			}
		}
	}

	err = s.repo.DeleteRecord("comment", commentID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
