package handlers

import (
	"fmt"
	"forum/internal/models"
	"forum/internal/services"
	uuid "forum/pkg/UUID"
	"forum/pkg/session"
	"html/template"
	"net/http"
	"net/url"
	"strings"
)

// GitHub Auth Handler - redirects user to GitHub
func AuthHandler(service services.Services, GetUrl func(string) string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state := uuid.GenerateUUID()

		// Store state in cookie for security
		http.SetCookie(w, &http.Cookie{
			Name:     "oauth_state",
			Value:    state,
			Path:     "/",
			HttpOnly: true,
			Secure:   r.TLS != nil,
			MaxAge:   300, // 5 minutes
		})

		authURL := GetUrl(state)
		http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
	}
}

// Google Callback Handler - handles return from Google
func GoogleCallbackHandler(service services.Services, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Verify state parameter for security
		stateCookie, err := r.Cookie("oauth_state")
		if err != nil || stateCookie.Value != r.URL.Query().Get("state") {
			RenderError(tmpl, w, http.StatusBadRequest, "Invalid state parameter") // ASK!!
			return
		}

		// Clear state cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "oauth_state",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
		})

		code := r.URL.Query().Get("code")
		if code == "" {
			RenderError(tmpl, w, http.StatusBadRequest, "No authorization code received")
			return
		}

		// Handle Google OAuth callback
		sessionID, authErr := service.OAuth().HandleGoogleCallback(code)

		if authErr.Code != 0 {
			if strings.HasPrefix(authErr.Msg, "USERNAME_NEEDED|") {

				// Parse user info from error message
				parts := strings.Split(authErr.Msg, "|")
				if len(parts) != 5 {
					RenderError(tmpl, w, http.StatusInternalServerError, "Invalid user info")
					return
				}

				provider := parts[1]
				providerID := parts[2]
				email := parts[3]
				name := parts[4]

				// Redirect to username selection page with user data
				redirectURL := fmt.Sprintf("/auth/username?provider=%s&provider_id=%s&email=%s&name=%s",
					url.QueryEscape(provider),
					url.QueryEscape(providerID),
					url.QueryEscape(email),
					url.QueryEscape(name))

				http.Redirect(w, r, redirectURL, http.StatusSeeOther)
				return
			}

			RenderError(tmpl, w, authErr.Code, authErr.Msg)
			return
		}

		// Set session cookie and redirect to homepage
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Path:     "/",
			Expires:  session.CalculateSessionExpiry(),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode, // CRITICAL: Changed from Strict to Lax
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// GitHub Callback Handler - handles return from GitHub
func GitHubCallbackHandler(service services.Services, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verify state parameter for security
		stateCookie, err := r.Cookie("oauth_state")
		if err != nil || stateCookie.Value != r.URL.Query().Get("state") {
			RenderError(tmpl, w, http.StatusBadRequest, "Invalid state parameter")
			return
		}

		// Clear state cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "oauth_state",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
		})

		code := r.URL.Query().Get("code")
		if code == "" {
			RenderError(tmpl, w, http.StatusBadRequest, "No authorization code received")
			return
		}

		// Handle GitHub OAuth callback
		sessionID, authErr := service.OAuth().HandleGitHubCallback(code)
		if authErr.Code != 0 {
			if strings.HasPrefix(authErr.Msg, "USERNAME_NEEDED|") {
				// Parse user info from error message
				parts := strings.Split(authErr.Msg, "|")
				if len(parts) != 5 {
					RenderError(tmpl, w, http.StatusInternalServerError, "Invalid user info")
					return
				}

				provider := parts[1]
				providerID := parts[2]
				email := parts[3]
				name := parts[4]

				// Redirect to username selection page with user data
				redirectURL := fmt.Sprintf("/auth/username?provider=%s&provider_id=%s&email=%s&name=%s",
					url.QueryEscape(provider),
					url.QueryEscape(providerID),
					url.QueryEscape(email),
					url.QueryEscape(name))

				http.Redirect(w, r, redirectURL, http.StatusSeeOther)
				return
			}
			RenderError(tmpl, w, authErr.Code, authErr.Msg)
			return
		}

		// Set the session cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Path:     "/",
			Expires:  session.CalculateSessionExpiry(),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode, // CRITICAL: Changed from Strict to Lax

		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// OAuth Username Selection Handler - shows username selection page
func OAuthUsernameHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			RenderError(tmpl, w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		// Get OAuth data from URL parameters (temporary storage)
		provider := r.URL.Query().Get("provider")
		providerID := r.URL.Query().Get("provider_id")
		email := r.URL.Query().Get("email")
		name := r.URL.Query().Get("name")

		// if provider == "" || providerID == "" || email == "" {
		// 	RenderError(tmpl, w, http.StatusBadRequest, "Missing OAuth information")
		// 	return
		// }

		data := models.OAuthData{
			Provider:   provider,
			ProviderID: providerID,
			Email:      email,
			Name:       name,
			FormData:   make(map[string]string),
		}

		tmpl.ExecuteTemplate(w, "oauth-username.html", data)
	}
}

// Set Username Handler - processes username selection form
func SetUsernameHandler(service services.Services, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			RenderError(tmpl, w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		// Get form data
		username := strings.TrimSpace(r.FormValue("username"))
		provider := r.FormValue("provider")
		providerID := r.FormValue("provider_id")
		email := r.FormValue("email")
		name := r.FormValue("name")

		// Validate required fields
		if username == "" || provider == "" || providerID == "" || email == "" {
			renderUsernameError(tmpl, w, "All fields are required", provider, providerID, email, name, username)
			return
		}

		// Basic username validation
		if len(username) < 3 || len(username) > 20 {
			renderUsernameError(tmpl, w, "Username must be between 3 and 20 characters", provider, providerID, email, name, username)
			return
		}

		// Create OAuth user with chosen username
		userID, err := service.OAuth().CreateOAuthUser(provider, providerID, email, name, username)
		if err.Code != 0 {
			renderUsernameError(tmpl, w, err.Msg, provider, providerID, email, name, username)
			return
		}

		//
		sessionID, err := service.OAuth().LoginoAuthUser(userID)
		if err.Code != 0 {
			RenderError(tmpl, w, err.Code, err.Msg)
			return
		}

		// Set the session cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Path:     "/",
			Expires:  session.CalculateSessionExpiry(),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode, // CRITICAL: Changed from Strict to Lax

		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// Helper function to render username selection with error
func renderUsernameError(tmpl *template.Template, w http.ResponseWriter, errorMsg, provider, providerID, email, name, username string) {
	data := models.OAuthData{
		Provider:      provider,
		ProviderID:    providerID,
		Email:         email,
		Name:          name,
		ValidationErr: errorMsg,
		FormData:      map[string]string{"username": username},
	}

	w.WriteHeader(400)
	tmpl.ExecuteTemplate(w, "oauth-username.html", data)
}
