package game

import (
	"fmt"
	"slices"
)

const (
	EUpLen = iota
	EDownLen
	ELeftLen
	ERightLen
	EX
	EY
	EDestroy
	ERemove
)

type publicExplosion struct {
	UpLen    *int     `json:"upLen,omitempty"`
	DownLen  *int     `json:"downLen,omitempty"`
	LeftLen  *int     `json:"leftLen,omitempty"`
	RightLen *int     `json:"rightLen,omitempty"`
	X        *float64 `json:"x,omitempty"`
	Y        *float64 `json:"y,omitempty"`

	Destroyed *bool `json:"destroy,omitempty"`

	Remove  *bool `json:"remove,omitempty"`
	changed bool
}

func (e *explosion) publicExplosion(fullView bool) publicExplosion {
	pe := publicExplosion{}
	pe.changed = slices.Contains(e.changed[:], true)
	if e.changed[EUpLen] || fullView {
		pe.UpLen = &e.upLen
	}
	if e.changed[EDownLen] || fullView {
		pe.DownLen = &e.downLen
	}
	if e.changed[ELeftLen] || fullView {
		pe.LeftLen = &e.leftLen
	}
	if e.changed[ERightLen] || fullView {
		pe.RightLen = &e.rightLen
	}
	if e.changed[EX] || fullView {
		pe.X = &e.p.x
	}
	if e.changed[EY] || fullView {
		pe.Y = &e.p.y
	}
	if e.changed[EDestroy] || fullView {
		pe.Destroyed = &e.destroyed
	}
	if e.changed[ERemove] || fullView {
		pe.Remove = &e.removed
	}
	if fullView == false {
		e.changed = [8]bool{}
	}
	return pe
}

type explosion struct {
	//public
	upLen     int
	downLen   int
	leftLen   int
	rightLen  int
	p         pos
	destroyed bool
	removed   bool

	changed [8]bool
	id      ID
}

func (s *state) newExplosion(p pos, size int) {
	//calculate length of each arm
	bombGridX, bombGridY := pos2GridPos(p)

	for _, player := range *s.Players {
		newPos := pos{float64(bombGridX), float64(bombGridY)}
		if dist(player.pos, newPos) < 0.7 && player.destroyed == false {
			player.takeDamage(p)
		}
	}

	legs := [4]int{}
	for legIndex, dir := range []pos{{0, -1}, {0, 1}, {-1, 0}, {1, 0}} {
		for i := range size {
			multiplier := 1 + i
			checkX, checkY := bombGridX+int(dir.x)*multiplier, bombGridY+int(dir.y)*multiplier

			for _, player := range *s.Players {
				newPos := pos{float64(checkX), float64(checkY)}
				if dist(player.pos, newPos) < 0.7 && player.destroyed == false {
					player.takeDamage(p)
				}
			}

			obst := (*s.arena)[checkY][checkX]
			if obst != nil {
				obst.destroy()
				break
			}

			legs[legIndex]++
		}
	}

	changed := [8]bool{}
	changeAll(changed[:])
	e := &explosion{
		upLen:     legs[0],
		downLen:   legs[1],
		leftLen:   legs[2],
		rightLen:  legs[3],
		p:         p,
		changed:   changed,
		id:        s.newID(),
		destroyed: true,
	}

	(*s.Explosions)[e.id] = e

	fmt.Printf("new explosion: %#v\n", e)
}
