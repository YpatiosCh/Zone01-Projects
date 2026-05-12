package game

import (
	"math/rand/v2"
	"slices"
)

type publicWall struct {
	Permanent *bool    `json:"permanent,omitempty"`
	X         *float64 `json:"x,omitempty"`
	Y         *float64 `json:"y,omitempty"`

	Destroyed *bool `json:"destroy,omitempty"`

	changed bool
}

func (w *wall) publicWall(fullView bool) publicWall {
	pw := publicWall{}
	pw.changed = slices.Contains(w.changed[:], true)

	if w.changed[WPermanent] == true || fullView {
		pw.Permanent = &w.Permanent
	}
	if w.changed[WY] == true || fullView {
		pw.X = &w.pos.x
	}
	if w.changed[WX] == true || fullView {
		pw.Y = &w.pos.y
	}
	if w.changed[WX] == true || fullView {
		pw.Y = &w.pos.y
	}
	if w.changed[WDestroyed] == true || fullView {
		pw.Destroyed = &w.destroyed
	}

	if fullView == false {
		w.changed = [4]bool{}
	}

	return pw
}

const (
	WPermanent = iota
	WY
	WX
	WDestroyed
)

type wall struct {
	//public
	Permanent bool
	pos       pos
	destroyed bool

	//private
	changed        [4]bool
	ID             ID
	width          int
	height         int
	body           body
	representation string
	state          *state
}

func (s *state) createWall(x, y int, permanent bool) {

	potentialObstacle := (*s.arena)[y][x]
	if potentialObstacle != nil {
		return
	}

	rep := "#"
	if permanent {
		rep = ":"
	}

	xPos := float64(x) * CellSize
	yPos := float64(y) * CellSize
	w := &wall{
		ID:             s.newID(),
		pos:            pos{xPos, yPos},
		body:           newRect(pos{xPos, yPos}, CellSize, CellSize),
		Permanent:      permanent,
		representation: rep,
		state:          s,
	}

	(*s.arena)[y][x] = w
	(*s.Walls)[w.ID] = w
}

func (s *state) createObstacles(playerCount, height, width int) {
	//set up starting positions

	//init this thing
	arena := make(arena, height)
	for i := range height {
		arena[i] = make([]obstacle, width)
	}

	s.arena = &arena

	//horizontal perm walls
	for x := range width {
		s.createWall(x, 0, true)
		s.createWall(x, height-1, true)
	}

	//vertical perm walls
	for y := 1; y < height-1; y++ {
		s.createWall(0, y, true)
		s.createWall(width-1, y, true)
	}

	//inner walls
	ymin := 1
	ymax := height - 1
	xmin := 1
	xmax := width - 1
	for y := ymin; y < ymax; y++ {
	nextSpot:
		for x := xmin; x < xmax; x++ {
			if x%2 == 0 && y%2 == 0 && x > 1 && x < xmax-1 && y > 1 && y < ymax-1 {
				//permanent walls
				s.createWall(x, y, true)
			} else {
				//breakable walls
				wallPos := pos{float64(x) * CellSize, float64(y) * CellSize}
				for _, l := range s.selectedLayout {
					if mag(sub(wallPos, l)) < 1.1 {
						// fmt.Print(" ")
						continue nextSpot
					}
				}

				// fmt.Print("#")
				if rand.IntN(2) == 0 {
					s.createWall(x, y, false)
				}
			}
		}
	}
}

func (w *wall) destroy() {
	if w.Permanent {
		return
	}
	if rand.Float32() < 0.37 {
		w.state.newPowerup(w.pos)
	}
	w.destroyed = true
	w.changed[WDestroyed] = true
}

func (w *wall) represent() string {
	return w.representation
}

func (w *wall) getBody() *body {
	return &w.body
}
