package handlers

import (
	"fmt"
	"log"
	"net/http"

	"bomberman/server/internal/manager"
	"bomberman/server/internal/player"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 4096,
	// Allocating CheckOrigin to allow all connections for development
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Test(m *manager.Manager, w http.ResponseWriter, r *http.Request) {
	fmt.Println("THIS IS A TEST")
}

// HandleWebSocket handles websocket requests
func HandleWebSocket(m *manager.Manager, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Get player ID and name from query params or generate defaults
	id := r.URL.Query().Get("id") //get uuid from front in case of reconnection
	if id == "" {
		id = uuid.NewString() //generate new uuid if new connection
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		name = "Player" // fallback
		return
	}

	log.Printf("New connection from %s (ID: %s)", name, id)

	// Create new player with nil channel initially (updated New signature)
	p := player.New(id, name, conn)

	// Join lobby and get reference
	l := m.Join(p)

	// Connect player to lobby input
	p.ToLobby = l.Input
	p.Unregister = l.GetUnregisterChan() // Needs new method in lobby

	// Start pumps
	go p.WritePump()
	go p.ReadPump()
}
