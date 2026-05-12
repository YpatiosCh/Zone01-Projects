# Multiplayer Game Server

This directory contains the backend server for the multiplayer game. The server is responsible for handling HTTP requests, upgrading connections to WebSockets, managing game lobbies, executing the core 60Hz game loop, and synchronizing game state across all connected clients.

The server is designed to be completely agnostic of the specific game rules, interacting with the game logic strictly through an `Engine` interface. 

## Architecture & Flow

The server follows a pipeline architecture from HTTP connection down to the localized game loop. 

1. **HTTP Entry & WebSocket Upgrade**  
   The application starts at `cmd/main.go`, which invokes `internal/entry.Run()`. The HTTP server serves static frontend files from `./front` and listens on the `/ws` endpoint. When a client hits `/ws`, `handlers.HandleWebSocket` upgrades the HTTP request to a persistent WebSocket connection and extracts the player's `id` and `name` from the query parameters.
2. **Player Abstraction**  
   A new `player.Player` object is instantiated for the connection. This object manages concurrent Read and Write pumps (goroutines) to safely handle incoming and outgoing JSON messages over the WebSocket.
3. **Lobby Management (Matchmaking)**  
   The `player.Player` is handed off to the `manager.Manager`. The Manager acts as the global matchmaker with a strict thread-safe mutex. It attempts to:
   - Reconnect the player to an active lobby if their `id` matches a disconnected user.
   - Place the player in an existing open lobby (waiting for players).
   - Create an entirely new `lobby.Lobby` if no open lobbies exist.
4. **Lobby & Game Loop**  
   Inside the `lobby.Lobby`, players are assigned map positions (1-4). The lobby runs an isolated state machine with several phases (`Waiting`, `Lobby`, `PreGame`, `Running`, `GameOver`). Once the `Running` phase is reached, a 60Hz `runGameLoop` takes over. The loop aggregates asynchronous keyboard inputs (e.g., `UP`, `DOWN`, `LEFT`, `RIGHT`, `ACTION_BUTTON`) from the frontend, feeds them to the game `Engine.Advance()` method, and broadcasts the serialized graphical output state back to all players in real-time.



## Lobby Lifecycle & Phases

The `Lobby` state machine progresses strictly through the following phases (`int32` atomic state):

1. **`PhaseWaiting` (0)**: The default state. The lobby is waiting for players to join. It stays in this state until at least `MinPlayers` (2) are connected.
2. **`PhaseLobby` (1)**: Triggered when 2 to 3 players are connected. A 20-second countdown ticker begins. If players drop below 2, the countdown cancels and reverts to `PhaseWaiting`. If the lobby hits `MaxPlayers` (4), the countdown immediately fast-forwards.
3. **`PhasePreGame` (2)**: Initializes the game `Engine` passing only the `playerCount`. The full board state is generated and a 10-second final countdown begins. Reconnections during this phase instantly receive the populated board state.
4. **`PhaseRunning` (3)**: The 60Hz physics ticker starts processing inputs and driving the `Engine.Advance()` method.
5. **`PhaseGameOver` (4)**: Set when `Engine.Advance()` returns `false` (meaning the game is complete). The lobby lingers for a few seconds so players can see the victory screen before sockets are forcibly closed.

## Core Loop Mechanics (`runGameLoop`)

Once a game transitions to `PhaseRunning`, the `Lobby.runGameLoop()` goroutine takes over. This loop is decoupled from WebSocket network reading (`run()`), which prevents slow clients from lagging the server physics.

- **Accumulator Pattern**: It uses a standard fixed-timestep accumulator targeting `60 FPS` (`16.66ms` per frame). If the server hiccups or thread sleeps, the accumulator catches up by processing multiple logic ticks back-to-back without sleeping.
- **Input Gathering**: `lobby.go` listens asynchronously for key events (keydown/keyup) from clients. When the physics tick executes `handleGameTick()`, it safely locks and clones the current snapshot of all held keys across all 4 players (`gatheredInputs` map).
- **Engine Advancement**: The server passes those boolean input maps directly to the abstract `Engine`. The server itself calculates no collision, movement, or game-specific feature logic.
- **Broadcast**: The engine returns a state delta (or full state). The `Lobby` packages this into a JSON `FRAME_UPDATE` and non-blockingly enqueues it to all connected sockets.

## Component Breakdown

- **`cmd/main.go`**  
  The main entry point. Simply invokes the server startup.

- **`internal/entry/entry.go`**  
  Sets up the `http.ServeMux`, the `Manager`, and starts the HTTP server. It also intercepts OS termination signals to gracefully shutdown all lobbies and connections before exiting.

- **`internal/handlers/websocket.go`**  
  Contains the WebSocket Upgrader logic. It parses the incoming request for query parameters and bootstraps the `Player` instance before passing it to the `Manager`.

- **`internal/manager/manager.go`**  
  The global registry of all active lobbies. It features strict locking to safely allocate players to lobbies and cleans up dead lobbies when games finish or become permanently empty.

- **`internal/lobby/lobby.go`**  
  The most complex component in the server. It manages:
  - Multi-phase countdowns (Waiting for 2 players -> 20s Lobby Countdown -> 10s Pre-Game Countdown -> Running).
  - Clean soft-disconnects and active reconnections, allowing players to briefly drop their connection without ruining the match.
  - A mathematically precise 60 FPS game loop (`runGameLoop`) that prevents drift, enforces latency ceilings, and interacts with the decoupled game `Engine`.

- **`internal/player/player.go`**  
  Wraps the raw Gorilla WebSocket connection. It contains `ReadPump` and `WritePump` to handle ping/pong heartbeats, parse inbound `protocol.ClientMessage`s, and securely enqueue outbound byte payloads without blocking the main game loop.

- **`internal/protocol/protocol.go`**  
  Defines the strongly typed JSON structures used for bidirectional communication:
  - `ClientMessage`: Sent from Frontend to Backend (Inputs, Chat, Leave).
  - `ServerMessage`: Sent from Backend to Frontend (State broadcasts, Chat, Lobby updates, Engine frames).

- **`internal/WEBSOCKET_API.md`**  
  The comprehensive reference documentation detailing the JSON schemas, events, and flow expectations for the Frontend application integrating with this Backend.

## Graceful Shutdowns

The server is built with resilience in mind. If the server process is killed (e.g., `SIGINT`), the `entry` package traps the signal and invokes `manager.Shutdown()`. The Manager synchronously broadcasts closure messages to all running lobbies, which in turn close all player WebSocket connections immediately, ensuring no ghost connections are left hanging. Lobbies automatically destroy themselves if they sit completely empty for 30 seconds.
