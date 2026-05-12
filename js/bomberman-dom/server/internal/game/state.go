package game

import (
	"encoding/json"
	"fmt"
	"math"
	"time"
)

type matchState struct {
	Announcement *string `json:"announcement,omitempty"`
	ArenaWidth   *int    `json:"arena_width,omitempty"`
	ArenaHeight  *int    `json:"arena_height,omitempty"`
}

//TODO don't stop game immeidately on death
//TODO faster base movespeed, more upgrades, nerf speed upgrade

type publicState struct {
	Players    *map[ID]*publicPlayer    `json:"players,omitempty"`
	Walls      *map[ID]*publicWall      `json:"walls,omitempty"`
	Bombs      *map[ID]*publicBomb      `json:"bombs,omitempty"`
	Explosions *map[ID]*publicExplosion `json:"explosions,omitempty"`
	MatchState *matchState              `json:"match_state,omitempty"`
	Powerups   *map[ID]*publicPowerup   `json:"powerups,omitempty"`
}

var printIt = true

// returns a struct that can be used the marshal into a json to be send to the frontend
func (s *state) publicState(fullView bool) publicState {

	playerMap := make(map[ID]*publicPlayer)
	wallMap := make(map[ID]*publicWall)
	bombMap := make(map[ID]*publicBomb)
	explosionMap := make(map[ID]*publicExplosion)
	powerups := make(map[ID]*publicPowerup)

	publicState := publicState{
		Players:    &playerMap,
		Walls:      &wallMap,
		Bombs:      &bombMap,
		Explosions: &explosionMap,
		Powerups:   &powerups,
	}

	if s.announcement != "" {
		newString := s.announcement
		publicState.MatchState = &matchState{
			Announcement: &newString,
		}
		fmt.Println("ANNOUNCEMENT: ", s.announcement)
		s.announcement = ""

	}

	if s.sendArenaSize || fullView {
		if publicState.MatchState == nil {
			w := s.arenaWidth
			h := s.arenaHeight
			publicState.MatchState = &matchState{
				ArenaWidth:  &w,
				ArenaHeight: &h,
			}
		} else {
			w := s.arenaWidth
			h := s.arenaHeight
			publicState.MatchState.ArenaHeight = &h
			publicState.MatchState.ArenaWidth = &w
		}

		if !fullView {
			s.sendArenaSize = false
		}
	}

	for key, p := range *s.Players {
		pp := p.publicPlayer(fullView)
		if pp.changed || fullView {
			(*publicState.Players)[key] = &pp
		}
	}
	if len(*publicState.Players) == 0 {
		publicState.Players = nil
	}

	for key, w := range *s.Walls {
		pw := w.publicWall(fullView)
		if pw.changed || fullView {
			(*publicState.Walls)[key] = &pw
		}
	}
	if len(*publicState.Walls) == 0 {
		publicState.Walls = nil
	}

	for key, b := range *s.Bombs {
		pb := b.publicBomb(fullView)
		if pb.changed || fullView {
			(*publicState.Bombs)[key] = &pb
		}
	}
	if len(*publicState.Bombs) == 0 {
		publicState.Bombs = nil
	}

	for key, e := range *s.Explosions {
		pe := e.publicExplosion(fullView)
		if pe.changed || fullView {
			(*publicState.Explosions)[key] = &pe
		}
	}
	if len(*publicState.Explosions) == 0 {
		publicState.Explosions = nil
	}

	for key, e := range *s.Powerups {
		pu := e.publicPowerup(fullView)
		if pu.changed || fullView {
			(*publicState.Powerups)[key] = &pu
		}
	}
	if len(*publicState.Powerups) == 0 {
		publicState.Powerups = nil
	}

	return publicState
}

func (s *state) GetState() any {
	p := s.publicState(true)

	// b, _ := json.Marshal(p)
	// fmt.Println(string(b))
	// time.Sleep(time.Second * 2)

	return p
}

type state struct {
	//public
	Players      *map[ID]*player
	Walls        *map[ID]*wall
	Bombs        *map[ID]*bomb
	Explosions   *map[ID]*explosion
	Powerups     *map[ID]*powerup
	announcement string

	//private
	IDCounter         ID
	gameFinished      bool
	arena             *arena
	twoPlayerLayout   []pos
	threePlayerLayout []pos
	fourPlayerLayout  []pos
	fullPlayerLayout  []pos
	selectedLayout    []pos
	layouts           [5][]pos
	endOfGameTimer    int
	sendArenaSize     bool
	arenaHeight       int
	arenaWidth        int
}

type obstacle interface {
	getBody() *body
	destroy()
}

type arena [][]obstacle

func (a *arena) getObstacle(x, y int) obstacle {
	if y < 0 || x < 0 || len((*a)) <= y || len((*a)[y]) <= x {
		return nil
	}
	return (*a)[y][x]
}

func (a *arena) getObstacleSAroundHere(pos pos) []obstacle {
	obstacles := []obstacle{}
	pos = alignToGrid(pos)
	for _, dir := range cardinalDirections {
		newObstacle := a.getObstacle(int(pos.x+dir.x), int(pos.y+dir.y))
		if newObstacle != nil {
			obstacles = append(obstacles, newObstacle)
		}
	}
	for _, dir := range diagonalDirections {
		newObstacle := a.getObstacle(int(pos.x+dir.x), int(pos.y+dir.y))
		if newObstacle != nil {
			obstacles = append(obstacles, newObstacle)
		}
	}
	return obstacles
}

// layouts = [5][]pos{nil, nil, twoPlayerLayout, threePlayerLayout, fourPlayerLayout}

func NewState(playerNames []string) (*state, error) {
	playerCount := len(playerNames)
	if playerCount == 0 || playerCount > 4 {
		return nil, fmt.Errorf("naaaah, not the right number of players bub!")
	}

	playerMap := make(map[ID]*player)
	wallMap := make(map[ID]*wall)
	bombMap := make(map[ID]*bomb)
	explosionMap := make(map[ID]*explosion)
	powerups := make(map[ID]*powerup)

	s := &state{
		Players:        &playerMap,
		Walls:          &wallMap,
		Bombs:          &bombMap,
		Explosions:     &explosionMap,
		Powerups:       &powerups,
		IDCounter:      0,
		endOfGameTimer: EndOfGameTimer,
		sendArenaSize:  true,
	}

	P2TopLeft := pos{1, 1}
	P3TopLeft := pos{1, 1}
	P4TopLeft := pos{1, 1}

	P3TopRight := pos{P3ArenaWidth - 2, 1}
	P4TopRight := pos{P4ArenaWidth - 2, 1}

	P4BottomLeft := pos{4, P4ArenaHeight - 2}

	P2BottomRight := pos{P2ArenaWidth - 2, P2ArenaHeight - 2}
	P4BottomRight := pos{P4ArenaWidth - 2, P4ArenaHeight - 2}

	BottomMiddle := pos{P3ArenaWidth / 2, P3ArenaHeight - 2}
	s.twoPlayerLayout = []pos{P2TopLeft, P2BottomRight}
	s.threePlayerLayout = []pos{P3TopLeft, BottomMiddle, P3TopRight}
	s.fourPlayerLayout = []pos{P4TopLeft, P4BottomLeft, P4TopRight, P4BottomRight}
	s.fullPlayerLayout = append(s.fourPlayerLayout, s.threePlayerLayout[1])
	s.layouts = [5][]pos{nil, nil, s.twoPlayerLayout, s.threePlayerLayout, s.fourPlayerLayout}
	s.selectedLayout = s.layouts[playerCount]

	for i, name := range playerNames {
		s.newPlayer(s.selectedLayout[i].x*CellSize, s.selectedLayout[i].y*CellSize, "red", name)
	}

	fmt.Printf("newly created players: %#v\n", s.Players)

	arenHeight := 0
	arenWidth := 0
	switch playerCount {
	case 2:
		arenHeight = P2ArenaHeight
		arenWidth = P2ArenaWidth
	case 3:
		arenHeight = P3ArenaHeight
		arenWidth = P3ArenaWidth
	case 4:
		arenHeight = P4ArenaHeight
		arenWidth = P4ArenaWidth
	}

	s.arenaHeight = arenHeight
	s.arenaWidth = arenWidth

	s.createObstacles(playerCount, arenHeight, arenWidth)

	return s, nil
}

func Print(message string, a any) {
	fmt.Printf("%s: %#v\n\n", message, a)
}

// moves forward the state of the game
func (s *state) Advance(inputs map[int]map[string]bool) (any, bool) {
	// timeStart := time.Now().UnixMilli()
	//delete destroyed or removed objects
	// fmt.Print("")

	for id, wall := range *s.Walls {
		if wall.destroyed {
			intX := int(math.RoundToEven(wall.pos.x) + 0.001)
			intY := int(math.RoundToEven(wall.pos.y) + 0.001)
			delete((*s.Walls), id)
			(*s.arena)[intY][intX] = nil
		}
	}

	for id, bomb := range *s.Bombs {
		if bomb.destroyed {
			delete((*s.Bombs), id)
			x, y := posToCell(bomb.pos)
			(*s.arena)[y][x] = nil
		}
	}

	for id, exp := range *s.Explosions {
		if exp.destroyed {
			delete((*s.Explosions), id)
		}
	}

	for id, Player := range *s.Players {
		if Player.destroyed {
			if len(*s.Players) > 2 {
				s.announcement = Player.name + " is out!"
			}
			delete((*s.Players), id)
		}
	}

	for id, powerup := range *s.Powerups {
		if powerup.destroyed {
			delete((*s.Powerups), id)
		}
	}

	for playerKey, input := range inputs {
		for keyName, buttonState := range input {
			switch keyName {
			case "up", "w":
				p, ok := (*s.Players)[ID(playerKey)]
				if !ok {
					continue
				}
				p.input.up = buttonState
			case "down", "s":
				p, ok := (*s.Players)[ID(playerKey)]
				if !ok {
					continue
				}
				p.input.down = buttonState
			case "left", "a":
				p, ok := (*s.Players)[ID(playerKey)]
				if !ok {
					continue
				}
				p.input.left = buttonState
			case "right", "d":
				p, ok := (*s.Players)[ID(playerKey)]
				if !ok {
					continue
				}
				p.input.right = buttonState
			case "space", "bomb", "b", " ":
				p, ok := (*s.Players)[ID(playerKey)]
				if !ok {
					continue
				}
				p.input.bomb = buttonState
			}
		}
	}

	//step objects forward
	for _, player := range *s.Players {
		player.move()
	}

	//make objects act
	for _, player := range *s.Players {
		player.act()
	}

	//make objects act
	for _, bomb := range *s.Bombs {
		bomb.act()
	}

	//calculate bomb explosion
	//kill objects

	alivePlayers := 0
	potentialWinner := &player{}
	for _, player := range *s.Players {
		if player.destroyed == false {
			alivePlayers++
			potentialWinner = player
		}
	}
	switch alivePlayers {
	case 1:
		if s.gameFinished == false {
			s.gameFinished = true
			potentialWinner.invincible = true
			s.announcement = potentialWinner.name + " won!!"
		}
	case 0:
		if s.gameFinished == false {
			s.gameFinished = true
			s.announcement = "EVERYBODY LOSES!"
		}
	}

	if s.gameFinished {
		s.endOfGameTimer--
		// fmt.Println("decrementing")
	}

	visualizeWalls(s.arena)

	//TODO add panic catch, and terminate game

	p := s.publicState(false)

	if printIt {
		b, _ := json.Marshal(p)
		fmt.Println()
		fmt.Println()
		fmt.Println()
		fmt.Println("printing bytes below this!")
		fmt.Println()
		fmt.Println(string(b))
		fmt.Println()
		fmt.Println()
		fmt.Println()
		fmt.Println()
		time.Sleep(time.Second * 1)
		printIt = false
	}

	return p, s.endOfGameTimer > 0
}

func (s *state) IsGameOver() bool {

	return s.endOfGameTimer == 0
}
