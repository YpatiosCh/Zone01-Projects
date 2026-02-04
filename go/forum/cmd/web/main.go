package main

import (
	"fmt"
	"forum/internal/database"
	"forum/internal/routes"
	"net/http"

	"log"

	_ "github.com/mattn/go-sqlite3"
)

const (
	port = ":8080"
)

func main() {
	db, err := database.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// set up routes
	handler := routes.SetUpRoutes(db)

	// start server
	fmt.Printf("Starting server on http://localhost%v \n", port)
	if err := http.ListenAndServe(port, handler); err != nil {
		log.Fatal("Failed to start server:", err)
	}

}
