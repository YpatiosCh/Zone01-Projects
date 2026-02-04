package services

import (
	"fmt"
	"forum/internal/repository"
)

type ReactionService struct {
	repo *repository.Manager
}

func NewReactionService(repo *repository.Manager) *ReactionService {
	return &ReactionService{
		repo: repo,
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

	return nil
}
