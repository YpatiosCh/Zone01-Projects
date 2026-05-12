package main

import (
	"net/http"
	"path/filepath"
	"strings"
)

// CustomFileServer creates a file server with proper MIME type handling
func CustomFileServer(root string) http.Handler {
	fs := http.FileServer(http.Dir(root))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set proper MIME types based on file extension
		ext := strings.ToLower(filepath.Ext(r.URL.Path))

		switch ext {
		case ".css":
			w.Header().Set("Content-Type", "text/css; charset=utf-8")
		case ".js":
			w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		case ".html", ".htm":
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
		case ".json":
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
		case ".png":
			w.Header().Set("Content-Type", "image/png")
		case ".jpg", ".jpeg":
			w.Header().Set("Content-Type", "image/jpeg")
		case ".gif":
			w.Header().Set("Content-Type", "image/gif")
		case ".svg":
			w.Header().Set("Content-Type", "image/svg+xml")
		case ".ico":
			w.Header().Set("Content-Type", "image/x-icon")
		case ".woff":
			w.Header().Set("Content-Type", "font/woff")
		case ".woff2":
			w.Header().Set("Content-Type", "font/woff2")
		case ".ttf":
			w.Header().Set("Content-Type", "font/ttf")
		}

		// Add caching headers for static assets
		if ext == ".css" || ext == ".js" || strings.HasPrefix(ext, ".woff") {
			w.Header().Set("Cache-Control", "public, max-age=3600") // 1 hour cache
		}

		// Serve the file
		fs.ServeHTTP(w, r)
	})
}

// LogRequest logs HTTP requests for debugging
func LogRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// You can uncomment this for debugging
		// fmt.Printf("🌐 %s %s %s\n", r.Method, r.URL.Path, r.RemoteAddr)
		handler.ServeHTTP(w, r)
	})
}
