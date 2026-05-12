package main

import (
	"errors"
	"log"
	"net/http"
	"zone01-dashboard/router"
)

const (
	ListenAddr = ":8080"
)

func main() {
	server := &http.Server{
		Addr:    ListenAddr,
		Handler: router.New(),
	}

	log.Printf("Server listening on %s", ListenAddr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server failed: %v", err)
	}
}
