# Bomberman WebSocket API for Frontend

## Connection
Connect to root `ws` endpoint with a `name`. You can optionally provide an `id` to reconnect to a dropped session.
```
wss://localhost:8080/ws?name=Mitsos&id=550e8400...
```
*(If no `id` is provided, the backend will automatically generate a UUID for the user).*

### 1. Assignment (Self)
Immediately upon connecting, the backend will send you your generated/confirmed ID and assigned map Position (1-4).
```json
{
  "type": "ASSIGN_ID",
  "data": {
    "id": "550e8400...",
    "name": "Mitsos",
    "position": 1
  }
}
```

### 2. Player Joined (Others)
Whenever anyone joins the lobby, the server broadcasts their details (including position 1-4) to you so you can render them on the lobby screen or the map.
```json
{
  "type": "PLAYER_JOINED",
  "data": {
    "name": "Mitsos",
    "position": 1
  }
}
```

### 3. State Auto-Recovery (Reconnections)
If a player accidentally drops their connection during an active game (PhasePreGame or PhaseRunning) and reconnects using their `id`, the server will instantly send them the full, current `GAME_STATE` payload (supplied by the Engine) directly to their connection before any 60Hz loop events so the frontend can draw the board instantly.
```json
{
  "type": "GAME_STATE",
  "data": { ... full engine state payload ... }
}
```

## Chat functionality

### 1. Sending Chat

To wire up chat, the frontend must send this exact message via WebSocket when a user presses enter:
```json
{
  "type": "SEND_MSG",
  "msg": "Hello everyone!"
}
```

### 2. Receiving Chat
The backend will immediately take that message and broadcast it to everyone in the lobby. The Backend has been specifically configured to output the exact format the frontend mock already expects:
```json
{
  "type": "GAME_CHAT",
  "message": {
    "user": "Mitsos",
    "text": "Hello everyone!"
  }
}
```

*Note: The server also uses this event type to broadcast connection events. For system messages, the `user` field will be  `"server"`.*
```json
{
  "type": "GAME_CHAT",
  "message": {
    "user": "server",
    "text": "Mitsos has disconnected"
  }
}
```

## Game Input functionality

### 1. Sending Key Presses
The server needs to persistently track when a player is holding down a key. To do this, the frontend must send a message when the key is **pressed** (`keydown`) and a matching message when the key is **released** (`keyup`). The payload must contain a `pressed` boolean field. 


Example Key Down:
```json
{
  "type": "INPUT",
  "key": "UP",
  "pressed": true
}
```

Example Key Up:
```json
{
  "type": "INPUT",
  "key": "UP",
  "pressed": false
}
```

## Game Timer functionality

### 1. Wait/Lobby Countdown (20s)
Wait till `LOBBY_TICK` is received from backend to tick down 20 seconds.
```json
{
  "type": "LOBBY_TICK",
  "data": {
    "countdown": 20
  }
}
```

### 2. Game Countdown (10s)
When 20 seconds reach 0 or 4 players joined, `PRE_GAME` is received including initial 60Hz loop state.
```json
{
  "type": "PRE_GAME",
  "data": {
    "countdown": 10,
    "state": { /* Initial map and players */ }
  }
}
```
And counts down:
```json
{
  "type": "PRE_GAME_TICK",
  "data": { "countdown": 9 }
}
```

### 3. Start Game
```json
{
  "type": "START_GAME"
}
```

### 4. Game Over
The server will wait 10 seconds after the primary Engine declares the game has concluded, then proactively send a `GAME_OVER` event and terminate the WebSocket connection.
```json
{
  "type": "GAME_OVER"
}
```
*Note for Frontend: Because the server instantly severs the TCP connection after sending this event, your application state must be designed to hold the "Victory/Game Over" screen open locally. When the user clears the screen, generate a new WebSocket connection.*
