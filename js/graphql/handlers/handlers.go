package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"text/template"
	"time"
	"zone01-dashboard/graphql"
	"zone01-dashboard/response"
	"zone01-dashboard/session"
)

type LoginRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type UserReq struct {
	UserID int `json:"user_id"`
}

const sessionCookieName = "zone01_session"

var Templates = template.Must(template.ParseGlob("templates/*.html"))

func HandleHomePage(w http.ResponseWriter, r *http.Request) {
	Templates.ExecuteTemplate(w, "index.html", nil)
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.JSON(w, http.StatusMethodNotAllowed, false, "method not allowed", nil)
		return
	}

	var payload LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		response.JSON(w, http.StatusBadRequest, false, "invalid payload", nil)
		return
	}

	payload.Identifier = strings.TrimSpace(payload.Identifier)
	if payload.Identifier == "" || payload.Password == "" {
		response.JSON(w, http.StatusBadRequest, false, "identifier and password are required", nil)
		return
	}

	// send request to zone01 :)
	token, claims, code := graphql.SignIn(payload.Identifier, payload.Password)
	if code != 200 {
		response.JSON(w, code, false, "invalid credentials", nil)
		return
	}
	// save JWT token in env
	userID, err := session.ExtractUserID(claims)
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, false, "unable to parse user id", nil)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(claims.Exp, 0),
	})

	response.JSON(w, http.StatusOK, true, "successful login", map[string]interface{}{
		"user_id": userID,
	})
}

func HandleProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.JSON(w, http.StatusMethodNotAllowed, false, "method not allowed", nil)
		return
	}

	// decode userID from body
	var userReq UserReq
	if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
		response.JSON(w, http.StatusBadRequest, false, "invalid payload", nil)
		return
	}

	cookie, err := r.Cookie(sessionCookieName)
	if err != nil || strings.TrimSpace(cookie.Value) == "" {
		response.JSON(w, http.StatusUnauthorized, false, "missing or invalid session", nil)
		return
	}

	claims, err := session.DecodeClaims(cookie.Value)
	if err != nil {
		response.JSON(w, http.StatusUnauthorized, false, "invalid session token", nil)
		return
	}

	claimedUserID, err := session.ExtractUserID(claims)
	if err != nil {
		response.JSON(w, http.StatusUnauthorized, false, "unable to extract session user id", nil)
		return
	}

	if userReq.UserID != 0 && userReq.UserID != claimedUserID {
		response.JSON(w, http.StatusForbidden, false, "forbidden: mismatched user id", nil)
		return
	}

	userID := claimedUserID
	token := cookie.Value

	// create arguments to query
	variables := make(map[string]interface{})
	variables["userId"] = userID

	// Run queries
	user, code, err := graphql.Query(graphql.USER_INFO, variables, token)
	if err != nil {
		response.JSON(w, code, false, err.Error(), nil)
		return
	}
	xp, code, err := graphql.Query(graphql.USER_XP, variables, token)
	if err != nil {
		response.JSON(w, code, false, err.Error(), nil)
		return
	}
	progress, code, err := graphql.Query(graphql.USER_PROGRESS, variables, token)
	if err != nil {
		response.JSON(w, code, false, err.Error(), nil)
		return
	}
	projects, code, err := graphql.Query(graphql.USER_PROJECTS, variables, token)
	if err != nil {
		response.JSON(w, code, false, err.Error(), nil)
		return
	}
	collabs, code, err := graphql.GetCollabs(userID, token)
	if err != nil {
		response.JSON(w, code, false, err.Error(), nil)
		return
	}

	// prepare data to respond
	data := map[string]interface{}{
		"user":          user,
		"xp":            xp,
		"progress":      progress,
		"projects":      projects,
		"collaborators": collabs,
	}

	response.JSON(w, http.StatusOK, true, "Successful retrieval of profile data", data)
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.JSON(w, http.StatusMethodNotAllowed, false, "method not allowed", nil)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	response.JSON(w, http.StatusOK, true, "logged out", nil)
}
