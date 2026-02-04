package middleware

import (
	"context"
	"forum/internal/models"
	"forum/internal/repository"
	"net/http"
	"time"
)

// Define a key type for the context
type contextKey string

// User context key
const UserContextKey contextKey = "user"

type Middleware struct {
	repo *repository.Manager
}

// NewMiddleware creates a new Middleware instance
func NewMiddleware(repo *repository.Manager) *Middleware {
	return &Middleware{
		repo: repo,
	}
}

func (m *Middleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the session cookie
		cookie, err := r.Cookie("session_id")
		if err != nil {
			// no cookie means not authenticated, proceed doing nothing
			next.ServeHTTP(w, r)
			return
		}

		// validate the session
		sessionID := cookie.Value
		sessions, err := m.repo.Get("session", "id", sessionID)
		if err != nil || sessions == nil || len(sessions) == 0 {
			// session not found, clear cookie and proceed doing nothing
			http.SetCookie(w, &http.Cookie{
				Name:     "session_id",
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
			})
			// continue without authentication
			next.ServeHTTP(w, r)
			return
		}

		session := sessions[0]

		// Check if session is expired
		if expiresAt, ok := session["expires_at"].(string); ok {
			expiry, err := time.Parse(time.RFC3339, expiresAt)
			if err == nil && time.Now().After(expiry) {
				// Session expired, clear cookie
				http.SetCookie(w, &http.Cookie{
					Name:     "session_id",
					Value:    "",
					Path:     "/",
					MaxAge:   -1,
					HttpOnly: true,
				})
				// Continue without authentication
				next.ServeHTTP(w, r)
				return
			}
		}

		// Session is valid, get the user
		userID, ok := session["user_id"].(string)
		if !ok {
			// User ID not found in session, clear cookie
			http.SetCookie(w, &http.Cookie{
				Name:     "session_id",
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
			})
			// Continue without authentication
			next.ServeHTTP(w, r)
			return
		}

		// Get user data
		users, err := m.repo.Get("user", "id", userID)
		if err != nil || users == nil || len(users) == 0 {
			// User not found, clear cookie
			http.SetCookie(w, &http.Cookie{
				Name:     "session_id",
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
			})
			// Continue without authentication
			next.ServeHTTP(w, r)
			return
		}

		user := users[0]

		// Create a new context with the user
		ctx := context.WithValue(r.Context(), UserContextKey, user)

		// Call the next handler with the new context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAuth is middleware that checks if a user is passed in the context (authenticated)
// If not authenticated, it redirects to the login page
func (m *Middleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUser(r)
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// GetUser gets the authenticated user from the request context
// Returns nil if no user is authenticated
func GetUser(r *http.Request) *models.User {
	user, ok := r.Context().Value(UserContextKey).(map[string]interface{})
	if !ok {
		return nil
	}

	// Also return nil for empty maps
	if len(user) == 0 {
		return nil
	}

	// Convert the map to a User struct
	userModel := &models.User{
		ID:       user["id"].(string),
		Email:    user["email"].(string),
		Username: user["username"].(string),
	}

	return userModel
}
