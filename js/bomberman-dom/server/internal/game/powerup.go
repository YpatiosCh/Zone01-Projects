package game

import (
	"math/rand/v2"
	"slices"
)

type publicPowerup struct {
	Y         *float64 `json:"y,omitempty"`
	X         *float64 `json:"x,omitempty"`
	Name      *string  `json:"name,omitempty"`
	Remove    *bool    `json:"remove,omitempty"`
	Destroyed *bool    `json:"destroy,omitempty"`
	changed   bool
}

const (
	PUY = iota
	PUX
	PUName
	PURemove
	PUDestroyed
)

type powerup struct {
	//public
	pos       pos
	name      string
	remove    bool
	destroyed bool

	//private
	id      ID
	changed [5]bool
	bodies  []body
	state   *state
}

func (p *powerup) publicPowerup(fullView bool) publicPowerup {
	pp := publicPowerup{}
	pp.changed = slices.Contains(p.changed[:], true)
	if p.changed[PUY] || fullView {
		pp.Y = &p.pos.y
	}
	if p.changed[PUX] || fullView {
		pp.X = &p.pos.x
	}
	if p.changed[PUName] || fullView {
		pp.Name = &p.name
	}
	if p.changed[PURemove] || fullView {
		pp.Remove = &p.remove
	}
	if p.changed[PUDestroyed] || fullView {
		pp.Destroyed = &p.destroyed
	}
	if fullView == false {
		p.changed = [5]bool{}
	}
	return pp
}

var powerupTypes = []string{"power-bomb", "speed", "power", "power", "power-bomb", "speed", "power", "power-bomb", "lives"}

func (s *state) newPowerup(pos pos) {

	allChanged := [5]bool{}
	changeAll(allChanged[:])

	newPowerup := &powerup{
		pos:     pos,
		name:    powerupTypes[rand.IntN(len(powerupTypes))],
		id:      s.newID(),
		changed: allChanged,
		state:   s,
		bodies:  []body{newCircle(pos, CellSize*0.45)},
	}

	(*s.Powerups)[newPowerup.id] = newPowerup
}

func (p *powerup) destroy() {
	p.destroyed = true
	p.changed[PUDestroyed] = true
	p.remove = true
	p.changed[PURemove] = true
}
