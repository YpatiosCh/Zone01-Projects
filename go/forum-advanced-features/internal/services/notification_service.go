package services

import (
	"fmt"
	"forum/internal/models"
	"forum/internal/repository"
	"time"
)

type NotificationService struct {
	repo *repository.Manager
}

func NewNotificationService(repo *repository.Manager) *NotificationService {
	return &NotificationService{
		repo: repo,
	}
}

// CreateNotification creates a new notification
func (s *NotificationService) CreateNotification(userID, notificationType, sourceUserID, message string, postID, commentID *string) (string, error) {
	// Prevent self-notifications
	if userID == sourceUserID {
		return "", nil // Don't create notification, but don't error either
	}

	// Check if the recipient user exists
	users, err := s.repo.Get("user", "id", userID)
	if err != nil {
		return "", fmt.Errorf("failed to check user existence: %w", err)
	}
	if len(users) == 0 {
		return "", fmt.Errorf("recipient user not found")
	}

	data := map[string]interface{}{
		"user_id":        userID,
		"type":           notificationType,
		"source_user_id": sourceUserID,
		"message":        message,
		"read":           false,
	}

	// Add postID if provided
	if postID != nil {
		data["post_id"] = *postID
	}

	// Add commentID if provided
	if commentID != nil {
		data["comment_id"] = *commentID
	}

	// Use existing repository CreateRecord method
	return s.repo.CreateRecord("notification", data)
}

// GetUserNotifications gets all notifications for a user
func (s *NotificationService) GetUserNotifications(userID string) ([]models.Notification, error) {
	// Get all notifications for the user
	allNotifications, err := s.repo.Get("notification", "user_id", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}

	if len(allNotifications) == 0 {
		return []models.Notification{}, nil
	}

	var result []models.Notification
	for _, notif := range allNotifications {
		// // Skip if notification is already read
		// if notif["read"].(bool) == true {
		// 	continue
		// }

		// Get source user's username
		sourceUsers, err := s.repo.Get("user", "id", notif["source_user_id"].(string))
		if err != nil {
			continue // Skip this notification if source user not found
		}
		if len(sourceUsers) == 0 {
			continue // Skip this notification if source user not found
		}

		sourceUsername := sourceUsers[0]["username"].(string)

		notification := models.Notification{
			ID:           notif["id"].(string),
			UserID:       notif["user_id"].(string),
			Type:         notif["type"].(string),
			SourceUserID: notif["source_user_id"].(string),
			Message:      notif["message"].(string),
			Read:         notif["read"].(bool),
			CreatedAt:    notif["created_at"].(time.Time),
			SourceUser:   sourceUsername,
		}

		// Handle nullable post_id
		if postID, ok := notif["post_id"].(string); ok && postID != "" {
			notification.PostID = &postID
		}

		// Handle nullable comment_id
		if commentID, ok := notif["comment_id"].(string); ok && commentID != "" {
			notification.CommentID = &commentID
		}

		result = append(result, notification)
	}

	// Sort notifications by created_at in descending order
	for i := 0; i < len(result)-1; i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i].CreatedAt.Before(result[j].CreatedAt) {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	return result, nil
}

// createReactionNotification creates a notification when someone reacts to a post or comment
func (s *NotificationService) CreateReactionNotification(sourceUserID, table, ID, reactionType string) {
	var recipientUserID string
	var message string
	var postID, commentID *string

	if table == "post" {
		// Get the post to find the author
		posts, err := s.repo.Get("post", "id", ID)
		if err != nil || len(posts) == 0 {
			return
		}
		recipientUserID = posts[0]["user_id"].(string)
		postTitle := posts[0]["title"].(string)

		// Create message
		if reactionType == "like" {
			message = fmt.Sprintf("liked your post: %s", postTitle)
		} else {
			message = fmt.Sprintf("disliked your post: %s", postTitle)
		}

		postID = &ID
	} else { // table == "comment"
		// Get the comment to find the author
		comments, err := s.repo.Get("comment", "id", ID)
		if err != nil || len(comments) == 0 {
			return
		}
		recipientUserID = comments[0]["user_id"].(string)
		postid := comments[0]["post_id"].(string)

		// Create message
		if reactionType == "like" {
			message = "liked your comment"
		} else {
			message = "disliked your comment"
		}
		postID = &postid
		commentID = &ID
	}

	s.CreateNotification(recipientUserID, reactionType, sourceUserID, message, postID, commentID)
}

// CreateCommentNotification creates a notification when someone comments on a post
func (c *NotificationService) CreateCommentNotification(sourceUserID, postID, commentID string) {
	// Get the post to find the author
	posts, err := c.repo.Get("post", "id", postID)
	if err != nil || len(posts) == 0 {
		return
	}

	recipientUserID := posts[0]["user_id"].(string)
	postTitle := posts[0]["title"].(string)

	// Create message
	message := fmt.Sprintf("commented on your post: %s", postTitle)

	c.CreateNotification(recipientUserID, "comment", sourceUserID, message, &postID, &commentID)
}

// GetUnreadNotificationCount gets the count of unread notifications for a user
func (s *NotificationService) GetUnreadNotificationCount(userID string) (int, error) {
	// Get all notifications for the user
	notifications, err := s.repo.Get("notification", "user_id", userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get notifications: %w", err)
	}

	// Count unread notifications
	unreadCount := 0
	for _, notif := range notifications {
		if notif["read"].(bool) == false {
			unreadCount++
		}
	}

	return unreadCount, nil
}

// MarkNotificationAsRead marks a specific notification as read
func (s *NotificationService) MarkNotificationAsRead(notificationID, userID string) error {
	// First check if the notification exists and belongs to the user
	notifications, err := s.repo.Get("notification", "id", notificationID)
	if err != nil {
		return fmt.Errorf("failed to get notification: %w", err)
	}
	if len(notifications) == 0 {
		return fmt.Errorf("notification not found")
	}

	// Check if notification belongs to the user
	if notifications[0]["user_id"].(string) != userID {
		return fmt.Errorf("unauthorized: notification doesn't belong to user")
	}

	// Update the notification to mark as read
	updateData := map[string]interface{}{
		"read": true,
	}

	return s.repo.UpdateRecord("notification", notificationID, updateData)
}

// DeleteNotification deletes a specific notification
func (s *NotificationService) DeleteNotification(notificationID, userID string) error {
	// First check if the notification exists and belongs to the user
	notifications, err := s.repo.Get("notification", "id", notificationID)
	if err != nil {
		return fmt.Errorf("failed to get notification: %w", err)
	}
	if len(notifications) == 0 {
		return fmt.Errorf("notification not found")
	}

	// Check if notification belongs to the user (security)
	if notifications[0]["user_id"].(string) != userID {
		return fmt.Errorf("unauthorized: notification doesn't belong to user")
	}

	// Delete the notification
	return s.repo.DeleteRecord("notification", notificationID)
}
