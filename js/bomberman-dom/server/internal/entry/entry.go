package entry

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bomberman/server/internal/handlers"
	"bomberman/server/internal/manager"
)

var version = "v0.1"

func Run() {
	// Initialize Lobby Manager
	mgr := manager.NewManager()

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./front"))
	mux.Handle("/", fs)

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleWebSocket(mgr, w, r)
	})
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		handlers.Test(mgr, w, r)
	})

	// set server
	server := http.Server{
		// Addr: "0.0.0.0:8080",
		Addr:    "localhost:8080",
		Handler: mux,
	}

	go func() {
		log.Printf("Server %s running on http://%s\n", version, server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("ListenAndServe failed: %v", err)
		}
	}()

	// wait here for process termination signal to initiate graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Println("Shutting down server...")

	// Gracefully shutdown the lobby manager (disconnects players, stops loops)
	mgr.Shutdown()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Graceful server Shutdown Failed: %v", err)
	}
	log.Println("Server stopped")
}
