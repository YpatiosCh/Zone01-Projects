package main

import (
	"ascii-art-web-stylize/tools"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// Parse the HTML templates used for rendering
var templates = template.Must(template.ParseFiles("templates/index.html", "templates/result.html", "templates/errors.html"))

func main() {
	// Serve static files from the 'static' directory (e.g., CSS, JS, images)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Serve files from the 'api' directory for API requests
	http.Handle("/api/", http.StripPrefix("/api/", http.FileServer(http.Dir("api"))))

	// Handle the root path and form submission routes
	http.HandleFunc("/", tools.HomeHandler)                     // Handles GET requests to the home page
	http.HandleFunc("/ascii-art", tools.AsciiArtHandler)        // Handles POST requests for ASCII art generation
	http.HandleFunc("/api/ascii-art", tools.ApiAsciiArtHandler) // Handles POST requests for ASCII art via API

	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil)) // Starts the server on port 8080
}
