package main

import (
	"bomberman/server/internal/game"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"time"
)

func main() {
	players := []string{"Alex", "Kat", "Vag", "Ypat"}
	state, err := game.NewState(players)
	if err != nil {
		fmt.Println("error1: ", err.Error())
	}
	b, _ := json.Marshal(state.GetState())
	fmt.Println(string(b))

	for f := range 10000 {
		inputMap := make(map[string]bool)
		// inputMap["up"] = true
		if (f+30)%40 == 0 {
			inputMap["bomb"] = true
			fmt.Println("DROPPING BOMB")
			time.Sleep(500 * time.Millisecond)
		}
		playerInputMap := make(map[int]map[string]bool)
		playerInputMap[1] = inputMap
		state, gameInProgress := state.Advance(playerInputMap)
		if !gameInProgress {
			bytes, _ := json.Marshal(state)
			fmt.Println(string(bytes))
			return
		}
		bytes, _ := json.Marshal(state)
		fmt.Println(string(bytes))
		time.Sleep(time.Millisecond * 17)
	}
}

func randBool() bool {
	return rand.IntN(2) == 0
}
