package lobby

import (
	"bomberman/server/internal/game"
	"bomberman/server/internal/player"
	"bomberman/server/internal/protocol"
	"encoding/json"
	"fmt"
	"log"
	"maps"
	"sync"
	"sync/atomic"
	"time"
)

const (
	MaxPlayers               = 4
	MinPlayers               = 2
	LobbyCountdownSeconds    = 15 // change to 20 before production
	PreGameCountdownSeconds  = 4  // change to 10 before production
	GameOverDelaySeconds     = 1
	GameTickRate             = 60
	EmptyLobbyTimeoutSeconds = 30
)

const (
	PhaseWaiting = iota
	PhaseLobby
	PhasePreGame
	PhaseRunning
	PhaseGameOver
)

// Engine defines the interface for the game logic.
// This allows server/internal to remain completely agnostic of game rules.
type Engine interface {
	//IsGameOver() bool // currently not in use
	GetState() any

	// Advance is called 60 times a second with all collected inputs.
	// It returns an object to be broadcast to all players,
	// or nil if nothing should be broadcast.
	Advance(inputs map[int]map[string]bool) (any, bool)
}

// Lobby maintains the set of active players and broadcasts messages.
type Lobby struct {
	players         map[*player.Player]bool   // Registered players.
	Input           chan player.PlayerMessage // Inbound messages from the players.
	broadcast       chan []byte               // Outbound messages to players (byte slice for raw sending).
	register        chan *player.Player       // Register requests from the players.
	unregister      chan *player.Player       // Unregister requests from players.
	quit            chan struct{}             // Channel to signal the lobby to shut down
	game            Engine                    // The game instance managed by this lobby.
	playerPositions map[string]int            // Tracks assigned position 1-4 for actively connected or disconnected UUIDs.
	phase           int32
	countdown       int                     // Remaining seconds before the game starts (used in PhaseLobby and PhasePreGame)
	countdownTicker *time.Ticker            // 1-second timer triggering countdown broadcasts
	mu              sync.Mutex              // Mutex to protect players map during concurrent access
	gameMu          sync.Mutex              // Mutex to serialize game engine access across goroutines
	inputsMu        sync.Mutex              // Mutex to protect gatheredInputs
	gatheredInputs  map[int]map[string]bool // Accumulated inputs for the current tick
	tickCounter     int                     // Debugging variable profiling the number of 60Hz physics frames executed
	emptyTimer      *time.Timer             // Timer triggering lobby shutdown when completely empty
	closeOnce       sync.Once
}

// NewLobby creates a new Lobby.
func NewLobby() *Lobby {
	return &Lobby{
		Input:           make(chan player.PlayerMessage, 1024),
		broadcast:       make(chan []byte, 256),
		register:        make(chan *player.Player, 10),
		unregister:      make(chan *player.Player, 10),
		quit:            make(chan struct{}),
		players:         make(map[*player.Player]bool),
		playerPositions: make(map[string]int),
		phase:           int32(PhaseWaiting),
		game:            nil,
		gatheredInputs:  make(map[int]map[string]bool),
	}
}

// Run pumps events to the lobby.
func (l *Lobby) Run() {
	go l.runBroadcastLoop()

	for {
		select {
		case p := <-l.register:
			l.handleRegister(p)

		case p := <-l.unregister:
			if l.handleUnregister(p) {
				return
			}

		case msg := <-l.Input:
			// fmt.Println("case input", msg)
			l.handlePlayerMessage(msg)

		case <-l.countdownChan():
			// fmt.Println("case countdown")
			l.handleCountdownTick()

		case <-l.emptyTimeoutChan():
			log.Println("Lobby empty timeout reached. Closing.")
			l.Close()
			return // Exit the Run loop

		case <-l.quit:
			log.Println("Lobby received quit signal. Shutting down loops.")
			return // Exit the Run loop
		}
	}
}

func (l *Lobby) handleRegister(p *player.Player) {
	l.mu.Lock()
	pendingBroadcasts := make([]protocol.ServerMessage, 0, 4)
	shouldStartGame := false

	// Enforce capacity at the point of commit to avoid CanJoin/Register TOCTOU overflow.
	if _, known := l.playerPositions[p.ID]; !known && len(l.players) >= MaxPlayers {
		l.mu.Unlock()
		close(p.Send)
		log.Printf("Rejecting player %s (%s): lobby full at registration time.", p.Name, p.ID)
		return
	}

	// Assign or restore position.
	if pos, exists := l.playerPositions[p.ID]; exists {
		p.Position = pos
	} else {
		p.Position = l.getAvailablePosition()
		l.playerPositions[p.ID] = p.Position
	}

	// Always drop an old ghost socket if it exists to prevent duplicate inputs from the same user ID.
	for existing := range l.players {
		if existing.ID == p.ID {
			delete(l.players, existing)
			close(existing.Send)
			log.Printf("Player %s (%s) new connection replaced ghost socket.", p.Name, p.ID)
			break
		}
	}

	// Stop the empty shutdown timer if players are joining
	if l.emptyTimer != nil {
		l.emptyTimer.Stop()
		l.emptyTimer = nil
		log.Println("Player joined/reconnected. Cancelled empty lobby shutdown timer.")
	}

	// Snapshot players already in the lobby before the new player is added.
	players := make(map[int]string, len(l.players))
	for existingPlayer := range l.players {
		players[existingPlayer.Position] = existingPlayer.Name
	}

	// Send assignment directly to the connecting player.
	assignMsg, _ := json.Marshal(protocol.ServerMessage{
		Type: "ASSIGN_ID",
		Data: map[string]any{
			"id":           p.ID,
			"position":     p.Position,
			"name":         p.Name,
			"playersCount": len(l.players) + 1,
			"players":      players,
		},
	})
	p.Send <- assignMsg

	// Reconnecting to an active game
	if phase := l.getPhase(); phase == PhasePreGame || phase == PhaseRunning {
		l.players[p] = true
		log.Printf("Player %s (%s) reconnected to lobby at position %d.", p.Name, p.ID, p.Position)

		pendingBroadcasts = append(pendingBroadcasts, protocol.ServerMessage{
			Type: "GAME_CHAT",
			Message: map[string]string{
				"user": "server",
				"text": fmt.Sprintf("%s has reconnected", p.Name),
			},
		})

		// If the game has already started, immediately push full state.
		l.gameMu.Lock()
		if l.game != nil {
			stateMsg, _ := json.Marshal(protocol.ServerMessage{
				Type: "GAME_STATE",
				Data: l.game.GetState(),
			})
			select {
			case p.Send <- stateMsg:
			default:
			}
		}
		l.gameMu.Unlock()
		l.mu.Unlock()
		for _, msg := range pendingBroadcasts {
			l.broadcastJSON(msg)
		}
		return
	}

	l.players[p] = true
	log.Printf("Player %s (%s) joined lobby. Total players: %d, Position: %d", p.Name, p.ID, len(l.players), p.Position)

	pendingBroadcasts = append(pendingBroadcasts, protocol.ServerMessage{
		Type: "PLAYER_JOINED",
		Data: map[string]any{
			"name":         p.Name,
			"position":     p.Position,
			"playersCount": len(l.players),
		},
	})

	// Check if we should start the timer or fast-forward.
	if l.getPhase() == PhaseWaiting && len(l.players) >= MinPlayers {
		log.Printf("Two players connected. Starting %ds lobby countdown...", LobbyCountdownSeconds)
		l.setPhase(PhaseLobby)
		l.countdown = LobbyCountdownSeconds
		if l.countdownTicker == nil {
			l.countdownTicker = time.NewTicker(time.Second)
		}
		pendingBroadcasts = append(pendingBroadcasts, protocol.ServerMessage{
			Type: "LOBBY_TICK",
			Data: map[string]int{"countdown": l.countdown},
		})
	} else if l.getPhase() == PhaseLobby && len(l.players) >= MaxPlayers {
		log.Println("Lobby full. Fast-forwarding to pre-game...")
		playerNames := getPlayerNames(l.players)
		preGameMsg, stateMsg, ok := l.startPreGameLocked(playerNames)
		if ok {
			if l.countdown <= 0 {
				shouldStartGame = true
				pendingBroadcasts = append(pendingBroadcasts, stateMsg)
			} else {
				pendingBroadcasts = append(pendingBroadcasts, preGameMsg, stateMsg)
			}
		}
	}
	l.mu.Unlock()

	for _, msg := range pendingBroadcasts {
		l.broadcastJSON(msg)
	}
	if shouldStartGame {
		l.startGame()
	}
}

func (l *Lobby) handleUnregister(p *player.Player) bool {
	l.mu.Lock()
	pendingBroadcasts := make([]protocol.ServerMessage, 0, 4)
	if _, ok := l.players[p]; !ok {
		l.mu.Unlock()
		return false
	}

	if phase := l.getPhase(); phase == PhaseRunning || phase == PhasePreGame {
		// Soft disconnect: keep in game.
		log.Printf("Player %s (%s) soft-disconnected.", p.Name, p.ID)

		pendingBroadcasts = append(pendingBroadcasts, protocol.ServerMessage{
			Type: "GAME_CHAT",
			Message: map[string]string{
				"user": "server",
				"text": fmt.Sprintf("%s has disconnected", p.Name),
			},
		})

		delete(l.players, p)
		close(p.Send)
		l.inputsMu.Lock()
		delete(l.gatheredInputs, p.Position)
		l.inputsMu.Unlock()

		// Check if lobby is completely empty
		if len(l.players) == 0 {
			log.Println("All players disconnected (soft). Starting 30 seconds empty timer...")
			if l.emptyTimer == nil {
				l.emptyTimer = time.NewTimer(EmptyLobbyTimeoutSeconds * time.Second)
			} else {
				l.emptyTimer.Reset(EmptyLobbyTimeoutSeconds * time.Second)
			}
		}

		l.mu.Unlock()
		for _, msg := range pendingBroadcasts {
			l.broadcastJSON(msg)
		}
		return false
	}

	// Hard disconnect if game hasn't started.
	// NOTE: We only delete out of playerPositions on hard disconnect so they lose their slot indefinitely.
	delete(l.playerPositions, p.ID)
	log.Printf("Player %s (%s) left lobby.", p.Name, p.ID)

	pendingBroadcasts = append(pendingBroadcasts, protocol.ServerMessage{
		Type: "GAME_CHAT",
		Message: map[string]string{
			"user": "server",
			"text": fmt.Sprintf("%s has disconnected", p.Name),
		},
	})

	delete(l.players, p)
	close(p.Send)
	l.inputsMu.Lock()
	delete(l.gatheredInputs, p.Position)
	l.inputsMu.Unlock()

	// Since frontend has no PLAYER_LEFT message handler, we use a two-step broadcast to clear the visual slot
	pendingBroadcasts = append(pendingBroadcasts, protocol.ServerMessage{
		Type: "PLAYER_JOINED",
		Data: map[string]any{
			"position":     p.Position,
			"name":         p.Name,
			"playersCount": len(l.players),
		},
	}, protocol.ServerMessage{
		Type: "GAME_STATE",
		Data: map[string]any{
			"players": map[int]any{
				p.Position: nil,
			},
			"go_to_game": false,
		},
	})

	// Cancel timer if players drop below 2.
	if l.getPhase() == PhaseLobby && len(l.players) < MinPlayers {
		log.Println("Not enough players left. Canceling countdown.")
		l.setPhase(PhaseWaiting)
		if l.countdownTicker != nil {
			l.countdownTicker.Stop()
			l.countdownTicker = nil
		}
		pendingBroadcasts = append(pendingBroadcasts, protocol.ServerMessage{
			Type: "LOBBY_CANCEL",
			Data: map[string]int{"countdown": -1},
		})
	}

	// If the lobby is completely empty and hasn't started, initiate cleanup.
	shouldClose := len(l.players) == 0 && l.getPhase() == PhaseWaiting
	l.mu.Unlock()
	for _, msg := range pendingBroadcasts {
		l.broadcastJSON(msg)
	}
	if shouldClose {
		log.Println("Lobby is completely empty. Shutting down lobby.")
		l.Close()
		return true
	}
	return false
}

func (l *Lobby) handlePlayerMessage(msg player.PlayerMessage) {
	// Explicit leave request from frontend.
	if msg.Msg.Type == "LEAVE_GAME" {
		fmt.Println("LEAVE_GAME")
		l.handleLeaveGame(msg.Player)
		return
	}

	// fmt.Println("received player input, handling it")
	if msg.Msg.Type == "SEND_MSG" {
		chatMsg := protocol.ServerMessage{
			Type: "GAME_CHAT",
			Message: map[string]string{
				"user": msg.Player.Name,
				"text": msg.Msg.Msg,
			},
		}
		l.broadcastJSON(chatMsg)
		return
	}

	if msg.Msg.Type == "INPUT" {
		if msg.Msg.Key == "leave_game" {
			return
		}
		l.inputsMu.Lock()
		defer l.inputsMu.Unlock()
		// Keep per-player key states so held keys persist across ticks.
		inputState, exists := l.gatheredInputs[msg.Player.Position]
		if !exists {
			inputState = make(map[string]bool)
		}
		inputState[msg.Msg.Key] = msg.Msg.Pressed
		l.gatheredInputs[msg.Player.Position] = inputState
		return
	}

	if msg.Msg.Type == "REQUEST_STATE" {
		// Give the requesting player a full master state dump immediately.
		phase := l.getPhase()
		if phase == PhasePreGame || phase == PhaseRunning {
			l.gameMu.Lock()
			defer l.gameMu.Unlock()
		}
		if (phase == PhasePreGame || phase == PhaseRunning) && l.game != nil {
			stateMsg := protocol.ServerMessage{
				Type: "GAME_STATE",
				Data: l.game.GetState(),
			}
			b, err := json.Marshal(stateMsg)
			if err == nil {
				// Serialize send-vs-close of player channels with l.mu.
				l.mu.Lock()
				if _, present := l.players[msg.Player]; present {
					select {
					case msg.Player.Send <- b:
					default:
						// Drop if buffer is full.
					}
				}
				l.mu.Unlock()
			}
		} else {
			noGameMsg := protocol.ServerMessage{Type: "NO_GAME"}
			b, err := json.Marshal(noGameMsg)
			if err == nil {
				l.mu.Lock()
				if _, present := l.players[msg.Player]; present {
					select {
					case msg.Player.Send <- b:
					default:
						// Drop if buffer is full.
					}
				}
				l.mu.Unlock()
			}
		}
	}
}

func (l *Lobby) handleLeaveGame(p *player.Player) {
	l.mu.Lock()
	pendingBroadcasts := make([]protocol.ServerMessage, 0, 2)

	// Leaving a game means the ID should no longer be reserved in this lobby.
	delete(l.playerPositions, p.ID)

	if _, ok := l.players[p]; ok {
		delete(l.players, p)
		close(p.Send)
	}

	l.inputsMu.Lock()
	delete(l.gatheredInputs, p.Position)
	l.inputsMu.Unlock()

	pendingBroadcasts = append(pendingBroadcasts, protocol.ServerMessage{
		Type: "GAME_CHAT",
		Message: map[string]string{
			"user": "server",
			"text": fmt.Sprintf("%s left the game", p.Name),
		},
	})

	if l.getPhase() == PhaseLobby && len(l.players) < MinPlayers {
		l.setPhase(PhaseWaiting)
		if l.countdownTicker != nil {
			l.countdownTicker.Stop()
			l.countdownTicker = nil
		}
		pendingBroadcasts = append(pendingBroadcasts, protocol.ServerMessage{
			Type: "LOBBY_CANCEL",
			Data: map[string]int{"countdown": -1},
		})
	}

	shouldClose := false
	if len(l.players) == 0 {
		phase := l.getPhase()
		if phase == PhasePreGame || phase == PhaseRunning || phase == PhaseGameOver {
			log.Println("All players left. Starting empty lobby timeout...")
			if l.emptyTimer == nil {
				l.emptyTimer = time.NewTimer(EmptyLobbyTimeoutSeconds * time.Second)
			} else {
				l.emptyTimer.Reset(EmptyLobbyTimeoutSeconds * time.Second)
			}
		} else {
			shouldClose = true
		}
	}

	l.mu.Unlock()

	for _, msg := range pendingBroadcasts {
		l.broadcastJSON(msg)
	}

	if shouldClose {
		l.Close()
	}
}

func (l *Lobby) handleBroadcast(message []byte) {
	l.mu.Lock()
	defer l.mu.Unlock()
	for p := range l.players {
		select {
		case p.Send <- message:
		default:
			close(p.Send)
			delete(l.players, p)
		}
	}
}

func (l *Lobby) handleCountdownTick() {
	l.mu.Lock()
	pendingBroadcasts := make([]protocol.ServerMessage, 0, 2)
	shouldStartGame := false

	l.countdown--
	if l.getPhase() == PhaseLobby {
		if l.countdown <= 0 {
			playerNames := getPlayerNames(l.players)
			preGameMsg, stateMsg, ok := l.startPreGameLocked(playerNames)
			if ok {
				if l.countdown <= 0 {
					shouldStartGame = true
					pendingBroadcasts = append(pendingBroadcasts, stateMsg)
				} else {
					pendingBroadcasts = append(pendingBroadcasts, preGameMsg, stateMsg)
				}
			} else {
				// Avoid negative countdown spiral if pregame setup fails.
				l.countdown = 1
			}
		} else {
			fmt.Println("Lobby countdown: ", l.countdown)
			pendingBroadcasts = append(pendingBroadcasts, protocol.ServerMessage{
				Type: "LOBBY_TICK",
				Data: map[string]int{"countdown": l.countdown},
			})
		}
		l.mu.Unlock()
		for _, msg := range pendingBroadcasts {
			l.broadcastJSON(msg)
		}
		if shouldStartGame {
			l.startGame()
		}
		return
	}

	if l.getPhase() == PhasePreGame {
		if l.countdown <= 0 {
			shouldStartGame = true
		} else {
			fmt.Println("Pregame countdown: ", l.countdown)
			pendingBroadcasts = append(pendingBroadcasts, protocol.ServerMessage{
				Type: "PRE_GAME_TICK",
				Data: map[string]int{"countdown": l.countdown},
			})
		}
	}
	l.mu.Unlock()
	for _, msg := range pendingBroadcasts {
		l.broadcastJSON(msg)
	}
	if shouldStartGame {
		l.startGame()
	}
}

func (l *Lobby) handleGameTick() bool {
	l.tickCounter++
	if l.tickCounter%200 == 0 {
		// fmt.Println("game tick: ", l.tickCounter)
	}
	phase := l.getPhase()
	if phase != PhaseRunning && phase != PhaseGameOver {
		return false
	}

	// If the game indicates it's over, handle cleanup.
	// l.gameMu.Lock()
	// if l.game.IsGameOver() {
	// 	l.gameMu.Unlock()
	// 	log.Println("Game Over detected. Shutting down lobby.")
	// 	l.setPhase(PhaseGameOver)
	// 	l.endGameWithDelay()
	// 	return false
	// }
	// l.gameMu.Unlock()

	// Hand all accumulated inputs to the engine.
	// Inputs are keyed by gameplay position, not websocket presence, so transient
	// connection issues do not freeze a player's controls mid-match.
	l.inputsMu.Lock()
	inputsCopy := make(map[int]map[string]bool)
	for position, keys := range l.gatheredInputs {
		keyCopy := make(map[string]bool)
		maps.Copy(keyCopy, keys)
		inputsCopy[position] = keyCopy
	}
	l.inputsMu.Unlock()

	t1 := time.Now()
	l.gameMu.Lock()
	response, gameInProgress := l.game.Advance(inputsCopy)
	l.gameMu.Unlock()
	engineDur := time.Since(t1)

	t2 := time.Now()
	if response != nil {
		msg := protocol.ServerMessage{
			Type: "FRAME_UPDATE",
			Data: response,
		}
		l.broadcastJSON(msg)
	}
	broadcastDur := time.Since(t2)

	if l.tickCounter%60 == 0 {
		log.Printf("Tick %d timing - Engine Advance: %v, Broadcast: %v", l.tickCounter, engineDur, broadcastDur)
	}

	if !gameInProgress && phase != PhaseGameOver {
		log.Println("THE GAME HAS FINISHED")
		l.setPhase(PhaseGameOver)
		l.endGameWithDelay()
	}

	return false
}

// runGameLoop executes concurrently to run the engine at exactly 60 FPS.
func (l *Lobby) runGameLoop() {
	tickDuration := time.Second / GameTickRate
	lastTime := time.Now()
	accumulator := time.Duration(0)
	lastLogTime := time.Now()
	frames := 0
	const maxCatchUpSteps = 5

	for {
		select {
		case <-l.quit:
			return
		default:
			now := time.Now()
			frameDelta := now.Sub(lastTime)
			lastTime = now

			// Cap delta so long stalls do not trigger a giant catch-up burst.
			if frameDelta > 250*time.Millisecond {
				frameDelta = 250 * time.Millisecond
			}
			accumulator += frameDelta

			steps := 0
			for accumulator >= tickDuration && steps < maxCatchUpSteps {
				if l.handleGameTick() {
					return // Game finished
				}
				accumulator -= tickDuration
				steps++
				frames++
			}

			if frames > 0 && frames%60 == 0 {
				elapsed := time.Since(lastLogTime)
				log.Printf("Game Loop [%p] - 60 frames processed in %v (Target: 1s). Effective FPS: %.2f", l, elapsed, 60.0/elapsed.Seconds())
				lastLogTime = time.Now()
			}

			if steps == maxCatchUpSteps && accumulator >= tickDuration {
				// Drop excess backlog to keep real-time behavior stable under heavy load.
				accumulator = 0
			}

			waitDuration := tickDuration - accumulator
			if waitDuration <= 0 {
				continue
			}

			timer := time.NewTimer(waitDuration)
			select {
			case <-l.quit:
				if !timer.Stop() {
					<-timer.C
				}
				return
			case <-timer.C:
			}
		}
	}
}

func (l *Lobby) countdownChan() <-chan time.Time {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.countdownTicker != nil {
		return l.countdownTicker.C
	}
	return nil
}

func (l *Lobby) emptyTimeoutChan() <-chan time.Time {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.emptyTimer != nil {
		return l.emptyTimer.C
	}
	return nil
}

func (l *Lobby) getAvailablePosition() int {
	taken := make(map[int]bool)
	for _, pos := range l.playerPositions {
		taken[pos] = true
	}
	for i := 1; i <= MaxPlayers; i++ {
		if !taken[i] {
			return i
		}
	}
	return 0 // Should not be reached due to CanJoin restriction
}

func (l *Lobby) broadcastJSON(v interface{}) {
	// if l.tickCounter%200 == 0 {
	// fmt.Println("broadcasting: ", l.tickCounter)
	// }
	data, err := json.Marshal(v)
	if err != nil {
		log.Printf("Json marshal error: %v", err)
		return
	}
	select {
	case <-l.quit:
		return
	case l.broadcast <- data:
	default:
		log.Printf("Broadcast queue full, dropping message type %T", v)
	}
}

func (l *Lobby) startLobbyCountdown() {
	l.setPhase(PhaseLobby)
	l.countdown = LobbyCountdownSeconds
	if l.countdownTicker == nil {
		l.countdownTicker = time.NewTicker(time.Second)
	}
	l.broadcastJSON(protocol.ServerMessage{Type: "LOBBY_TICK", Data: map[string]int{"countdown": l.countdown}})
}

func (l *Lobby) startPreGameLocked(playerNames []string) (protocol.ServerMessage, protocol.ServerMessage, bool) {
	engine, err := game.NewState(playerNames)
	if err != nil {
		log.Printf("Failed to create game state: %v", err)
		return protocol.ServerMessage{}, protocol.ServerMessage{}, false
	}

	l.setPhase(PhasePreGame)
	l.countdown = PreGameCountdownSeconds
	if l.countdownTicker == nil {
		l.countdownTicker = time.NewTicker(time.Second)
	}

	l.gameMu.Lock()
	l.game = engine
	gameState := l.game.GetState()
	l.gameMu.Unlock()

	preGameMsg := protocol.ServerMessage{
		Type: "PRE_GAME_TICK",
		Data: map[string]any{
			"countdown": l.countdown,
		},
	}

	stateMsg := protocol.ServerMessage{
		Type: "INITIAL_STATE",
		Data: gameState,
	}
	return preGameMsg, stateMsg, true
}

func (l *Lobby) startGame() {
	l.mu.Lock()
	phase := l.getPhase()
	if phase != PhasePreGame && phase != PhaseGameOver {
		l.mu.Unlock()
		return
	}
	if phase == PhasePreGame {
		l.gameMu.Lock()
		hasGame := l.game != nil
		l.gameMu.Unlock()
		if !hasGame {
			l.mu.Unlock()
			log.Println("Cannot start game: engine is nil.")
			return
		}
	}

	l.setPhase(PhaseRunning)
	if l.countdownTicker != nil {
		l.countdownTicker.Stop()
		l.countdownTicker = nil
	}
	l.mu.Unlock()

	log.Println("Game starting!")

	// Send START_GAME message
	l.broadcastJSON(protocol.ServerMessage{Type: "START_GAME"})

	// Start dedicated game loop goroutine
	go l.runGameLoop()
}

func (l *Lobby) cancelCountdown() {
	l.setPhase(PhaseWaiting)
	if l.countdownTicker != nil {
		l.countdownTicker.Stop()
		l.countdownTicker = nil
	}
	l.broadcastJSON(protocol.ServerMessage{Type: "LOBBY_CANCEL", Data: map[string]int{"countdown": -1}})
}

// CanJoin checks if a player can join this lobby.
func (l *Lobby) CanJoin() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if len(l.players) >= MaxPlayers {
		return false
	}
	phase := l.getPhase()
	if phase == PhasePreGame || phase == PhaseRunning || phase == PhaseGameOver {
		return false
	}
	return true
}

func (l *Lobby) RegisterPlayer(p *player.Player) bool {
	select {
	case <-l.quit:
		return false
	case l.register <- p:
		return true
	}
}

// GetUnregisterChan returns the channel used to signal player disconnection to the lobby.
func (l *Lobby) GetUnregisterChan() chan<- *player.Player {
	return l.unregister
}

func (l *Lobby) GetPlayerCount() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.players)
}

// HasPlayer checks if the player ID is currently tracked in this lobby's assigned positions.
func (l *Lobby) HasPlayer(id string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	_, exists := l.playerPositions[id]
	return exists
}

// endGameWithDelay handles end of game sequence.
// waits 10s, then sends GAME_OVER and closes.
func (l *Lobby) endGameWithDelay() {

	go func() {
		time.Sleep(GameOverDelaySeconds * time.Second)
		msg, _ := json.Marshal(protocol.ServerMessage{Type: "GAME_OVER"})
		l.mu.Lock()
		for p := range l.players {
			select {
			case p.Send <- msg:
			default:
			}
		}
		l.mu.Unlock()
		l.Close()
	}()
}

// Close gracefully shuts down the lobby and disconnects all players
func (l *Lobby) Close() {
	l.closeOnce.Do(func() {
		close(l.quit)

		l.mu.Lock()
		defer l.mu.Unlock()

		// Stop tickers
		if l.countdownTicker != nil {
			l.countdownTicker.Stop()
		}
		if l.emptyTimer != nil {
			l.emptyTimer.Stop()
		}

		// Close all player connections
		for p := range l.players {
			close(p.Send)
			delete(l.players, p)
		}
	})
}

func (l *Lobby) getPhase() int {
	return int(atomic.LoadInt32(&l.phase))
}

func (l *Lobby) setPhase(phase int) {
	atomic.StoreInt32(&l.phase, int32(phase))
}

func (l *Lobby) runBroadcastLoop() {
	for {
		select {
		case <-l.quit:
			return
		case message := <-l.broadcast:
			l.handleBroadcast(message)
		}
	}
}

func getPlayerNames(players map[*player.Player]bool) []string {
	playerNames := make([]string, len(players))
	for p, b := range players {
		if b {
			playerNames[p.Position-1] = p.Name
		}
	}
	return playerNames
}
