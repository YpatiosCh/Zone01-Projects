package services

import (
	"encoding/json"
	"fmt"
	"forum/internal/config"
	"forum/internal/models"
	"forum/internal/repository"
	uuid "forum/pkg/UUID"
	"forum/pkg/session"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type OAuthService struct {
	repo   *repository.Manager
	config *config.AppConfig
}

func NewOAuthService(repo *repository.Manager, config *config.AppConfig) *OAuthService {
	return &OAuthService{
		repo:   repo,
		config: config,
	}
}

// Google OAuth URLs
func (s *OAuthService) GetGoogleAuthURL(state string) string {
	redirectURL := s.config.AppURL + "/auth/google/callback"
	return fmt.Sprintf(
		"https://accounts.google.com/o/oauth2/auth?client_id=%s&redirect_uri=%s&scope=%s&response_type=code&state=%s",
		s.config.GoogleClientID,
		redirectURL,
		"openid email profile",
		state,
	)
}

// GitHub OAuth URLs
func (s *OAuthService) GetGitHubAuthURL(state string) string {
	redirectURL := s.config.AppURL + "/auth/github/callback"
	return fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=%s&state=%s",
		s.config.GitHubClientID,
		redirectURL,
		"user:email",
		state,
	)
}

// Handle Google OAuth callback
func (s *OAuthService) HandleGoogleCallback(code string) (string, Error) {
	// Exchange code for access token
	token, err := s.ExchangeGoogleCode(code)
	if err != nil {
		return "", Error{Msg: "Failed to authenticate with Google", Code: 500}
	}

	// Get user info from Google
	userInfo, err := s.GetGoogleUserInfo(token)
	if err != nil {
		return "", Error{Msg: "Failed to get user information", Code: 500}
	}

	// Handle user creation/login
	return s.handleOAuthUser("google", userInfo)
}

// Handle GitHub OAuth callback
func (s *OAuthService) HandleGitHubCallback(code string) (string, Error) {
	// Exchange code for access token
	token, err := s.ExchangeGitHubCode(code)
	if err != nil {
		return "", Error{Msg: "Failed to authenticate with GitHub", Code: 500}
	}

	// Get user info from GitHub
	userInfo, err := s.GetGitHubUserInfo(token)
	if err != nil {
		return "", Error{Msg: "Failed to get user information", Code: 500}
	}

	// Handle user creation/login
	return s.handleOAuthUser("github", userInfo)
}

func (s *OAuthService) handleOAuthUser(provider string, userInfo *models.User) (string, Error) {
	// STEP 1: Check if OAuth user already exists (returning user)
	existingOAuthUsers, err := s.repo.GetByTwoColumns("user", "oauth_provider", provider, "oauth_provider_id", userInfo.ID)
	if err != nil {
		return "", Error{Msg: err.Error(), Code: 500}
	}

	var userID string

	if len(existingOAuthUsers) > 0 {
		// RETURNING OAUTH USER - just login
		userID = existingOAuthUsers[0]["id"].(string)
	} else {
		// NEW OAUTH USER - check if email exists
		existingEmailUsers, err := s.repo.Get("user", "email", userInfo.Email)
		if err != nil {
			return "", Error{Msg: err.Error(), Code: 500}
		}

		if len(existingEmailUsers) > 0 {
			// EMAIL EXISTS - link OAuth to existing regular user
			userID = existingEmailUsers[0]["id"].(string)

			// Update existing user with OAuth info
			err = s.repo.UpdateRecord("user", userID, map[string]interface{}{
				"oauth_provider":    provider,
				"oauth_provider_id": userInfo.ID,
			})
			if err != nil {
				return "", Error{Msg: "Failed to link OAuth account", Code: 500}
			}
		} else {
			// COMPLETELY NEW USER - return special error with user info encoded
			userInfoString := fmt.Sprintf("USERNAME_NEEDED|%s|%s|%s|%s",
				provider, userInfo.ID, userInfo.Email, userInfo.Username)
			return "", Error{Msg: userInfoString, Code: 299}
		}
	}

	// Create session for existing user
	sessionID, loginerr := s.LoginoAuthUser(userID)
	if loginerr.Code != 0 {
		return "", loginerr
	}

	return sessionID, Error{}
}

// Check if OAuth user already exists (returning user)
func (s *OAuthService) CheckExistingOAuthUser(provider, providerID string) (*models.User, error) {
	existingUsers, err := s.repo.GetByTwoColumns("user", "oauth_provider", provider, "oauth_provider_id", providerID)
	if err != nil {
		return nil, err
	}

	if len(existingUsers) == 0 {
		return nil, nil
	}

	user := existingUsers[0]
	return &models.User{
		ID:       user["id"].(string),
		Email:    user["email"].(string),
		Username: user["username"].(string),
	}, nil
}

// Create session for user
func (s *OAuthService) CreateSessionForUser(userID string) (string, error) {
	sessionID := session.GenerateSessionToken()
	expiry := session.CalculateSessionExpiry()

	_, err := s.repo.CreateRecord("session", map[string]interface{}{
		"id":         sessionID,
		"user_id":    userID,
		"expires_at": expiry,
	})
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

// Exchange Google authorization code for access token
func (s *OAuthService) ExchangeGoogleCode(code string) (string, error) {
	redirectURL := s.config.AppURL + "/auth/google/callback"

	data := url.Values{}
	data.Set("client_id", s.config.GoogleClientID)
	data.Set("client_secret", s.config.GoogleClientSecret)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", redirectURL)

	resp, err := http.Post(
		"https://oauth2.googleapis.com/token",
		"application/x-www-form-urlencoded",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	accessToken, ok := result["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("no access token in response")
	}

	return accessToken, nil
}

// Get Google user information using access token
func (s *OAuthService) GetGoogleUserInfo(token string) (*models.User, error) {
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo map[string]interface{}
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	// Convert to our User model
	return &models.User{
		ID:       userInfo["id"].(string), // Google user ID
		Email:    userInfo["email"].(string),
		Username: userInfo["name"].(string), // Display name
	}, nil
}

// Exchange GitHub authorization code for access token
func (s *OAuthService) ExchangeGitHubCode(code string) (string, error) {
	data := url.Values{}
	data.Set("client_id", s.config.GitHubClientID)
	data.Set("client_secret", s.config.GitHubClientSecret)
	data.Set("code", code)

	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	accessToken, ok := result["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("no access token in response")
	}

	return accessToken, nil
}

// Get GitHub user information using access token
func (s *OAuthService) GetGitHubUserInfo(token string) (*models.User, error) {
	// Step 1: Get basic GitHub user info
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo map[string]interface{}
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	// Step 2: Get fallback email if missing
	email := fmt.Sprintf("%v", userInfo["email"])
	if email == "" || email == "<nil>" {
		emailReq, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
		if err != nil {
			return nil, err
		}
		emailReq.Header.Set("Authorization", "Bearer "+token)

		emailResp, err := http.DefaultClient.Do(emailReq)
		if err != nil {
			return nil, err
		}
		defer emailResp.Body.Close()

		emailBody, err := io.ReadAll(emailResp.Body)
		if err != nil {
			return nil, err
		}

		var emails []map[string]interface{}
		if err := json.Unmarshal(emailBody, &emails); err != nil {
			return nil, err
		}

		for _, e := range emails {
			if primary, ok := e["primary"].(bool); ok && primary {
				if verified, ok := e["verified"].(bool); ok && verified {
					email = fmt.Sprintf("%v", e["email"])
					break
				}
			}
		}
	}

	// Step 3: Fallback for username
	username := userInfo["name"]
	if username == nil || username == "" {
		username = userInfo["login"]
	}

	return &models.User{
		ID:       fmt.Sprintf("%v", userInfo["id"]),
		Email:    email,
		Username: fmt.Sprintf("%v", username),
	}, nil
}

// CreateOAuthUser creates a new user account with OAuth info and chosen username
func (s *OAuthService) CreateOAuthUser(provider, providerID, email, name, username string) (string, Error) {
	// Check if username already exists
	existingUsers, err := s.repo.Get("user", "username", username)
	if err != nil {
		return "", Error{Msg: err.Error(), Code: 500}
	}
	if len(existingUsers) > 0 {
		return "", Error{Msg: "Username already exists. Please choose a different one.", Code: 409}
	}

	// Check if email already exists (shouldn't happen, but safety check)
	existingEmailUsers, err := s.repo.Get("user", "email", email)
	if err != nil {
		return "", Error{Msg: err.Error(), Code: 500}
	}
	if len(existingEmailUsers) > 0 {
		return "", Error{Msg: "Email already registered", Code: 409}
	}

	// Create new user with OAuth info
	userID := uuid.GenerateUUID()
	userID, err = s.repo.CreateRecord("user", map[string]interface{}{
		"id":                userID,
		"email":             email,
		"username":          username,
		"password":          nil, // OAuth users don't have passwords
		"oauth_provider":    provider,
		"oauth_provider_id": providerID,
	})
	if err != nil {
		return "", Error{Msg: "Failed to create user account", Code: 500}
	}

	return userID, Error{}
}

func (s *OAuthService) LoginoAuthUser(userID string) (string, Error) {

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

	// Create new session for the user
	sessionID, err := s.CreateSessionForUser(userID)
	if err != nil {
		Error := Error{Msg: err.Error(), Code: 500}
		return "", Error
	}
	return sessionID, Error{}
}
