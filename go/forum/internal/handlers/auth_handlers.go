package handlers

import (
	"errors"
	"forum/internal/middleware"
	"forum/internal/services"
	"forum/internal/utils/validation"
	"forum/pkg/session"
	"html/template"
	"log"
	"net/http"
	"strings"
)

// RegisterUser handler handles user registration
func RegisterUser(authService *services.AuthService, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodGet {
			tmpl.ExecuteTemplate(w, "signup.html", nil)
			return
		}

		if r.Method == http.MethodPost {

			username := validation.SanitizeInput(r.FormValue("username"))
			email := validation.SanitizeInput(r.FormValue("email"))
			password := validation.SanitizeInput(r.FormValue("password"))
			confirmPassword := validation.SanitizeInput(r.FormValue("confirm-password"))

			errorStruct := validation.ValidateRegistration(username, email, password)

			if len(errorStruct.ErrorSlice) != 0 {
				data := struct{ ValidationErr string }{
					ValidationErr: strings.Join(errorStruct.Error(), " & "),
				}

				w.WriteHeader(400)
				tmpl.ExecuteTemplate(w, "signup.html", data)
				return
			}

			if password != confirmPassword {
				data := struct{ ValidationErr string }{
					ValidationErr: "Wrong password confirmation",
				}

				w.WriteHeader(400)
				tmpl.ExecuteTemplate(w, "signup.html", data)
				return
			}

			userID, err := authService.RegisterUser(username, email, password)
			if err.Code != 0 {
				w.WriteHeader(err.Code)
				data := struct{ ValidationErr string }{
					ValidationErr: err.Msg,
				}
				tmpl.ExecuteTemplate(w, "signup.html", data)
				return
			}

			log.Println("User registered with ID:", userID)

			// Redirect to login page after successful registration
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
	}
}

func LoginUser(authService *services.AuthService, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		data := struct{ ValidationErr string }{}
		user := middleware.GetUser(r)
		if r.Method == http.MethodGet {
			// only if user exists
			if user != nil {
				RenderError(tmpl, w, http.StatusMethodNotAllowed, "Method not allowed: You are already logged in")
			}
			tmpl.ExecuteTemplate(w, "login.html", nil)
			return
		}

		if r.Method == http.MethodPost {
			email := validation.SanitizeInput(r.FormValue("email"))
			password := validation.SanitizeInput(r.FormValue("password"))

			sessionID, err := authService.Login(email, password)
			if err.Code != 0 {
				data.ValidationErr = err.Msg
				w.WriteHeader(err.Code)
				tmpl.ExecuteTemplate(w, "login.html", data)
				return
			}

			// Set the session cookie
			http.SetCookie(w, &http.Cookie{
				Name:     "session_id",
				Value:    sessionID,
				Path:     "/",
				Expires:  session.CalculateSessionExpiry(),
				HttpOnly: true,
				Secure:   r.TLS != nil, // Set Secure flag if TLS is enabled
				SameSite: http.SameSiteStrictMode,
			})

			// render home page with user data
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}
}

func LogoutUser(authService *services.AuthService, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		data := struct{ ValidationErr string }{}

		// Only allow POST requests
		if r.Method != http.MethodPost {
			RenderError(tmpl, w, http.StatusMethodNotAllowed, "This method is not allowed")
			return
		}

		// Get the session ID from the cookie
		cookie, err := r.Cookie("session_id")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				RenderError(tmpl, w, http.StatusUnauthorized, "Unauthorized user: session missing")
				return
			} else {
				RenderError(tmpl, w, http.StatusInternalServerError, "Internal server error")
				return
			}
		}

		sessionID := cookie.Value

		// Delete the session from the database
		Error := authService.Logout(sessionID)
		if Error.Code != 0 {
			data.ValidationErr = Error.Msg
			w.WriteHeader(Error.Code)
			tmpl.ExecuteTemplate(w, "login.html", data)
			return
		}

		// Clear the session cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
