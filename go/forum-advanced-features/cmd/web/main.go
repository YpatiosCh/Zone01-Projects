package main

import (
	"fmt"
	"forum/internal/config"
	"forum/internal/database"
	"forum/internal/middleware"
	"forum/internal/repository"
	"forum/internal/routes"
	"forum/internal/services"
	sessioncleanup "forum/internal/utils/sessionCleanup"
	"net/http"

	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := database.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// set up configuration
	config := config.LoadConfig()

	// set up database manager
	dbManager := repository.NewManager(db)

	// Clean expired sessions evwery 24 hours
	go sessioncleanup.CleanupExpiredSessionsCron(dbManager)

	// set up middleware
	middleware := middleware.NewMiddleware(dbManager)

	// set up services
	services := services.NewServiceContainer(dbManager, config)

	// set up routes
	handler := routes.SetUpRoutes(middleware, services, config)

	// start server
	fmt.Printf("Starting server on http://localhost%v \n", config.Port)
	if err := http.ListenAndServe(config.Port, handler); err != nil {
		log.Fatal("Failed to start server:", err)
	}

}
