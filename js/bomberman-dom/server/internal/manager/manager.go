package manager

import (
	"log"
	"sync"

	"bomberman/server/internal/lobby"
	"bomberman/server/internal/player"
)

type Manager struct {
	lobbies []*lobby.Lobby
	mu      sync.Mutex
}

func NewManager() *Manager {
	return &Manager{
		lobbies: make([]*lobby.Lobby, 0),
	}
}

// Join adds a player to a lobby and returns the lobby they joined.
func (m *Manager) Join(p *player.Player) *lobby.Lobby {
	for {
		m.mu.Lock()
		var target *lobby.Lobby

		// 1. Check if player belongs to an existing lobby (active or disconnected)
		for _, l := range m.lobbies {
			if l.HasPlayer(p.ID) {
				log.Printf("Player %s (ID: %s) reconnecting to previous lobby", p.Name, p.ID)
				target = l
				break
			}
		}

		// 2. Check for an open lobby
		if target == nil {
			for _, l := range m.lobbies {
				if l.CanJoin() {
					log.Printf("Player %s joining open lobby", p.Name)
					target = l
					break
				}
			}
		}

		// 3. Create new lobby
		if target == nil {
			log.Printf("Creating new lobby for player %s", p.Name)
			target = lobby.NewLobby()
			go func(l *lobby.Lobby) {
				l.Run()
				// When Run returns, the lobby has closed. Trigger cleanup.
				m.RemoveLobby(l)
			}(target)
			m.lobbies = append(m.lobbies, target)
		}
		m.mu.Unlock()

		if target.RegisterPlayer(p) {
			return target
		}
		// The lobby closed between selection and registration. Retry.
	}
}

// RemoveLobby safely deletes a lobby from the manager's tracking list when it closes.
func (m *Manager) RemoveLobby(target *lobby.Lobby) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, l := range m.lobbies {
		if l == target {
			// Fast removal by swapping with the last element
			m.lobbies[i] = m.lobbies[len(m.lobbies)-1]
			m.lobbies = m.lobbies[:len(m.lobbies)-1]
			log.Println("Lobby cleanly removed from Manager.")
			return
		}
	}
}

func (m *Manager) Shutdown() {
	m.mu.Lock()
	log.Printf("Manager Shutdown initiated: Closing %d active lobbies.", len(m.lobbies))

	// Copy the slice and release the lock immediately to prevent deadlocks
	// when the closing lobby's goroutine inevitably attempts to call m.RemoveLobby()
	lobbiesToClose := append([]*lobby.Lobby(nil), m.lobbies...)
	m.mu.Unlock()

	for _, l := range lobbiesToClose {
		l.Close()
	}
}
