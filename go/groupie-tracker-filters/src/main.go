package main

import (
	"fmt"
	"groupie-tracker/handlers"
	"log"
	"net/http"
)

func main() {

	if err := handlers.InitializeArtists(); err != nil {
		log.Fatalf("Failed to initialize artists: %v", err)
	}
	// Serve static files such as CSS, JavaScript, and images from the "static" directory.
	// Using http.FileServer ensures efficient handling of static resources.
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Define the route handlers for various endpoints.
	// Each handler function is responsible for managing a specific route.

	// Home route handler - serves the main page of the application.
	http.HandleFunc("/", handlers.HomeHandler)

	// Filter route handler - serves filter requests from users.
	http.HandleFunc("/filter", handlers.FilterArtistsHandler)

	// Artist route handler - serves details of a specific artist based on the provided ID.
	http.HandleFunc("/artist/", handlers.HandleArtist)

	// Search route handler - handles search requests from users.
	http.HandleFunc("/search", handlers.SearchHandler)

	// Suggest route handler - provides suggestions for user input (e.g., autocomplete).
	http.HandleFunc("/suggest", handlers.SuggestHandler)

	// Start the HTTP server on port 8080 and log any fatal errors.
	// It's good practice to use log.Fatal to ensure errors are logged if the server fails to start.
	fmt.Println("Server started on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
