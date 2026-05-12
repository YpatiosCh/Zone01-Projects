package protocol

// Data sent from Client to Server
type ClientMessage struct {
	Type    string `json:"type"`
	Key     string `json:"key,omitempty"`     // For INPUT: UP, DOWN, LEFT, RIGHT, BOMB
	Pressed bool   `json:"pressed,omitempty"` // True on keydown, false on keyup
	Msg     string `json:"msg,omitempty"`     // For SEND_MSG
}

// Data sent from Server to Client
type ServerMessage struct {
	Type    string            `json:"type"`
	Data    any               `json:"data,omitempty"`
	Message map[string]string `json:"message,omitempty"` // For GAME_CHAT {user, text}
}
