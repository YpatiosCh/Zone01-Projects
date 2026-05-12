package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

func runClient(name string, leaveOnStart bool) {
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
	log.Printf("[%s] connecting to %s", name, u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalf("[%s] dial error: %v", name, err)
	}
	defer c.Close()

	// Wait a moment before registering
	time.Sleep(100 * time.Millisecond)

	// Send registration
	regMsg := map[string]interface{}{
		"type": "REGISTER",
		"name": name,
	}
	if err := c.WriteJSON(regMsg); err != nil {
		log.Fatalf("[%s] write error: %v", name, err)
	}

	gameStarted := false
	frameCount := 0

	// Message reader loop
	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Printf("[%s] read error: %v", name, err)
				return
			}

			var data map[string]interface{}
			json.Unmarshal(message, &data)

			msgType, ok := data["type"].(string)
			if !ok {
				log.Printf("[%s] warning: message missing type field: %s", name, string(message))
				continue
			}

			// Print specific important messages
			if msgType == "PLAYER_JOINED" || msgType == "PRE_GAME" || msgType == "LOBBY_TICK" || msgType == "PRE_GAME_TICK" || msgType == "START_GAME" || msgType == "GAME_STATE" || msgType == "GAME_CHAT" {
				if msgType == "GAME_STATE" {
					fmt.Printf("\n[%s] <= %s (Master State Sync Received!)\n", name, msgType)
				} else if msgType == "GAME_CHAT" {
					msgData, _ := data["message"].(map[string]interface{})
					fmt.Printf("\n[%s] CHAT: %s: %s\n", name, msgData["user"], msgData["text"])
				} else {
					fmt.Printf("\n[%s] <= %s\n", name, string(message))
				}
			}

			if msgType == "START_GAME" {
				gameStarted = true
				if leaveOnStart {
					go func() {
						time.Sleep(2 * time.Second)
						log.Printf("\n[%s] => LEAVING THE LOBBY EARLY!\n", name)
						// Send a standard close message to guarantee the server's ReadPump notices the disconnect
						c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
						c.Close()
					}()
				}
			}

			// Only print FRAME_UPDATE occasionally so we don't spam terminal
			if msgType == "FRAME_UPDATE" {
				frameCount++
				if frameCount%60 == 0 {
					fmt.Printf("\n[%s] Received 60 FRAME_UPDATE physics ticks...\n", name)
				}
			}
		}
	}()

	// Simulating gameplay inputs
	for {
		if gameStarted {
			// Press RIGHT
			fmt.Printf("\n[%s] => Pressing Right Key\n", name)
			c.WriteJSON(map[string]interface{}{
				"type":    "INPUT",
				"key":     "right",
				"pressed": true,
			})
			time.Sleep(1 * time.Second)

			// Release RIGHT
			fmt.Printf("\n[%s] => Releasing Right Key\n", name)
			c.WriteJSON(map[string]interface{}{
				"type":    "INPUT",
				"key":     "right",
				"pressed": false,
			})

			// Press BOMB
			fmt.Printf("\n[%s] => Pressing BOMB Key\n", name)
			c.WriteJSON(map[string]interface{}{
				"type":    "INPUT",
				"key":     "bomb",
				"pressed": true,
			})
			time.Sleep(100 * time.Millisecond)
			c.WriteJSON(map[string]interface{}{
				"type":    "INPUT",
				"key":     "bomb",
				"pressed": false,
			})

			// Simulate Lag by manually requesting the full State
			if name == "Alice" {
				fmt.Printf("\n[%s] => Requesting Full GAME_STATE (Lag Simulation)\n", name)
				c.WriteJSON(map[string]interface{}{
					"type": "REQUEST_STATE",
				})
			}

			time.Sleep(2 * time.Second)
		} else {
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func main() {
	go runClient("Alice", true)
	time.Sleep(500 * time.Millisecond)
	go runClient("Bob", false)

	// Block forever
	select {}
}
