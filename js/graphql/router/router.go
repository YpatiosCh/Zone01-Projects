package router

import (
	"net/http"
	"zone01-dashboard/handlers"
)

func New() http.Handler {
	fs := http.FileServer(http.Dir("static"))

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/", handlers.HandleHomePage)

	mux.HandleFunc("/api/auth/login", handlers.HandleLogin)
	mux.HandleFunc("/api/auth/logout", handlers.HandleLogout)
	mux.HandleFunc("/api/profile", handlers.HandleProfile)

	return mux
}
