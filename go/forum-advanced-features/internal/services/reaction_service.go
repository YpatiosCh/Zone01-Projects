package services

import (
	"fmt"
	"forum/internal/repository"
)

type ReactionService struct {
	repo         *repository.Manager
	notification *NotificationService // Add this line
}

func NewReactionService(repo *repository.Manager, notification *NotificationService) *ReactionService {
	return &ReactionService{
		repo:         repo,
		notification: notification, // Add this line
	}
}

func (s *ReactionService) ToggleReaction(userID, table, ID, reactionType string) error {
	// Validate reaction type
	if reactionType != "like" && reactionType != "dislike" {
		return fmt.Errorf("invalid reaction type: %s", reactionType)
	}

	// Verify table exists
	tables, err := s.repo.Get(table, "id", ID)
	if err != nil {
		return fmt.Errorf("failed to verify post: %w", err)
	}
	if len(tables) == 0 {
		return fmt.Errorf("table row not found")
	}

	// Check if user already has a reaction on this post
	existingReactions, err := s.repo.Get("reaction", "user_id", userID)
	if err != nil {
		return fmt.Errorf("failed to check existing reactions: %w", err)
	}

	// Find if user has already reacted to this post
	var existingReaction map[string]interface{}
	for _, reaction := range existingReactions {
		// Check if this reaction is for the current post
		if reaction["post_id"] == ID || reaction["comment_id"] == ID {
			existingReaction = reaction
			break
		}
	}

	// If user already has the same reaction, delete it (toggle off)
	if existingReaction != nil && existingReaction["type"] == reactionType {
		reactionID := existingReaction["id"].(string)
		err = s.repo.DeleteRecord("reaction", reactionID)
		if err != nil {
			return fmt.Errorf("failed to remove reaction: %w", err)
		}
		return nil
	}

	// If user has the opposite reaction, update it
	if existingReaction != nil {
		reactionID := existingReaction["id"].(string)
		err = s.repo.UpdateRecord("reaction", reactionID, map[string]interface{}{
			"type": reactionType,
		})
		if err != nil {
			return fmt.Errorf("failed to update reaction: %w", err)
		}
		// create notification for the user
		s.notification.CreateReactionNotification(userID, table, ID, reactionType)
		return nil
	}

	if table == "post" {
		// If no existing reaction, create a new one
		_, err = s.repo.CreateRecord("reaction", map[string]interface{}{
			"user_id": userID,
			"post_id": ID,
			"type":    reactionType,
		})
		if err != nil {
			return fmt.Errorf("failed to create reaction: %w", err)
		}
	} else { // else table is comment
		// If no existing reaction, create a new one
		_, err = s.repo.CreateRecord("reaction", map[string]interface{}{
			"user_id":    userID,
			"comment_id": ID,
			"type":       reactionType,
		})
		if err != nil {
			return fmt.Errorf("failed to create reaction: %w", err)
		}
	}

	// create notification for the user
	s.notification.CreateReactionNotification(userID, table, ID, reactionType)
	return nil
}

// UserLikedDislikedComment checks if a user has liked or disliked a comment
func (s *ReactionService) UserLikedDislikedComment(userID, commentID string) (bool, bool, error) {
	reactions, err := s.repo.GetByTwoColumns("reaction", "user_id", userID, "comment_id", commentID)
	if err != nil {
		return false, false, err
	}

	if reactions == nil || len(reactions) == 0 {
		return false, false, nil
	}

	var isLiked, isDisliked bool
	for _, reaction := range reactions {
		if reactionType, ok := reaction["type"].(string); ok {
			if reactionType == "like" {
				isLiked = true
			} else if reactionType == "dislike" {
				isDisliked = true
			}
		}
	}

	return isLiked, isDisliked, nil
}

// UserLikedDislikedPost checks if a user has liked or disliked a post
func (s *ReactionService) UserLikedDislikedPost(userID, postID string) (bool, bool, error) {
	reactions, err := s.repo.GetByTwoColumns("reaction", "user_id", userID, "post_id", postID)
	if err != nil {
		return false, false, err
	}

	if reactions == nil || len(reactions) == 0 {
		return false, false, nil
	}

	var isLiked, isDisliked bool
	for _, reaction := range reactions {
		if reactionType, ok := reaction["type"].(string); ok {
			if reactionType == "like" {
				isLiked = true
			} else if reactionType == "dislike" {
				isDisliked = true
			}
		}
	}

	return isLiked, isDisliked, nil
}

// GetReactionsForComment - keep as helper
func (s *ReactionService) GetReactionsForComment(commentID string) (int, int, error) {
	reactions, err := s.repo.Get("reaction", "comment_id", commentID)
	if err != nil {
		return 0, 0, err
	}

	if reactions == nil {
		return 0, 0, nil
	}

	var likes, dislikes int
	for _, reaction := range reactions {
		if reactionType, ok := reaction["type"].(string); ok {
			if reactionType == "like" {
				likes++
			} else if reactionType == "dislike" {
				dislikes++
			}
		}
	}

	return likes, dislikes, nil
}

// GetReactionsForPost - keep as helper
func (s *ReactionService) GetReactionsForPost(postID string) (int, int, error) {
	reactions, err := s.repo.Get("reaction", "post_id", postID)
	if err != nil {
		return 0, 0, err
	}

	if reactions == nil {
		return 0, 0, nil
	}

	var likes, dislikes int
	for _, reaction := range reactions {
		if reactionType, ok := reaction["type"].(string); ok {
			if reactionType == "like" {
				likes++
			} else if reactionType == "dislike" {
				dislikes++
			}
		}
	}

	return likes, dislikes, nil
}
