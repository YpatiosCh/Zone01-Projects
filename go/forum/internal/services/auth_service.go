package services

import (
	"fmt"
	"forum/internal/repository"
	"forum/pkg/crypt"
	"forum/pkg/session"
)

type Error struct {
	Msg  string
	Code int
}

// Common errors
var (
	ErrInvalidCredentials = Error{Msg: "invalid email or password", Code: 401}
	ErrUserNotFound       = Error{Msg: "user not found", Code: 404}
	ErrEmailExists        = Error{Msg: "email already exists", Code: 409}
	ErrUsernameExists     = Error{Msg: "username already exists", Code: 409}
	ErrInvalidSessionID   = Error{Msg: "invalid or expired session", Code: 401}
)

// AuthService handles authentication-related business logic
type AuthService struct {
	repo *repository.Manager
}

// NewAuthService creates a new AuthService
func NewAuthService(repo *repository.Manager) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

// RegisterUser registers a new user
func (s *AuthService) RegisterUser(username, email, password string) (string, Error) {
	// Check if email already exists
	existingUsers, err := s.repo.Get("user", "email", email)
	if err != nil {
		Error := Error{Msg: err.Error(), Code: 500}
		return "", Error
	}
	if existingUsers != nil && len(existingUsers) > 0 {
		return "", ErrEmailExists
	}

	// Check if username already exists
	existingUsers, err = s.repo.Get("user", "username", username)
	if err != nil {
		Error := Error{Msg: err.Error(), Code: 500}
		return "", Error
	}
	if existingUsers != nil && len(existingUsers) > 0 {
		return "", ErrUsernameExists
	}

	// Hash password
	hashedPassword, err := crypt.GenerateHash(password)
	if err != nil {
		Error := Error{Msg: err.Error(), Code: 500}
		return "", Error
	}

	// Create user
	userID, err := s.repo.CreateRecord("user", map[string]interface{}{
		"username": username,
		"email":    email,
		"password": hashedPassword,
	})
	if err != nil {
		Error := Error{Msg: err.Error(), Code: 500}
		return "", Error
	}

	return userID, Error{}
}

// Login authenticates a user and creates a single session (deleting any existing sessions)
func (s *AuthService) Login(email, password string) (string, Error) {
	// Get user by username
	users, err := s.repo.Get("user", "email", email)
	if err != nil {
		Error := Error{Msg: err.Error(), Code: 500}
		return "", Error
	}
	if users == nil || len(users) == 0 {
		return "", ErrInvalidCredentials
	}

	user := users[0]
	userID, ok := user["id"].(string)
	if !ok {
		Error := Error{Msg: "User Not found", Code: 500}
		return "", Error
	}

	// Verify password
	storedPassword, ok := user["password"].(string)
	if !ok {
		Error := Error{Msg: "Invalid password", Code: 500}
		return "", Error
	}

	if !crypt.VerifyPassword(password, storedPassword) {
		return "", ErrInvalidCredentials
	}

	// Delete any existing sessions for this user
	existingSessions, err := s.repo.Get("session", "user_id", userID)
	if err != nil {
		Error := Error{Msg: err.Error(), Code: 500}
		return "", Error
	}

	if existingSessions != nil && len(existingSessions) > 0 {
		for _, session := range existingSessions {
			sessionID, ok := session["id"].(string)
			if ok {
				err = s.repo.DeleteRecord("session", sessionID)
				if err != nil {
					// Log the error but continue - don't block login if cleanup fails
					// Consider adding proper logging here
					fmt.Printf("Error deleting existing session %s: %v\n", sessionID, err)
				}
			}
		}
	}

	// Generate session token
	sessionID := session.GenerateSessionToken()
	expiry := session.CalculateSessionExpiry()

	// Store session in database
	_, err = s.repo.CreateRecord("session", map[string]interface{}{
		"id":         sessionID,
		"user_id":    userID,
		"expires_at": expiry,
	})
	if err != nil {
		Error := Error{Msg: err.Error(), Code: 500}
		return "", Error
	}

	return sessionID, Error{}
}

// Logout invalidates a user's session
func (s *AuthService) Logout(sessionID string) Error {
	// Find session
	sessions, err := s.repo.Get("session", "id", sessionID)
	if err != nil {
		Error := Error{Msg: err.Error(), Code: 500}
		return Error
	}
	if sessions == nil || len(sessions) == 0 {
		return ErrInvalidSessionID
	}
	if err := s.repo.DeleteRecord("session", sessionID); err != nil {
		Error := Error{Msg: err.Error(), Code: 500}
		return Error
	}
	// Delete session
	return Error{}
}
