package game

import (
	"slices"
)

type publicBomb struct {
	Y *float64 `json:"y,omitempty"`
	X *float64 `json:"x,omitempty"`

	Destroyed *bool `json:"destroy,omitempty"`

	Remove  *bool `json:"remove,omitempty"`
	changed bool
}

const (
	BY = iota
	BX
	BDestroy
	BRemove
)

type bomb struct {
	//public
	pos          pos
	destroyed    bool
	remove       bool
	durationLeft int

	//private
	id           ID
	changed      [5]bool
	bodies       []body
	size         int
	state        *state
	ownerId      ID
	explodeDelay int
	exploding    bool
}

func (b *bomb) publicBomb(fullView bool) publicBomb {
	pb := publicBomb{}
	pb.changed = slices.Contains(b.changed[:], true)
	if b.changed[BY] || fullView {
		pb.Y = &b.pos.y
	}
	if b.changed[BX] || fullView {
		pb.X = &b.pos.x
	}
	if b.changed[BDestroy] || fullView {
		pb.Destroyed = &b.destroyed
	}
	if b.changed[BRemove] || fullView {
		pb.Remove = &b.remove
	}
	if fullView == false {
		b.changed = [5]bool{}
	}
	return pb
}

func (s *state) newBomb(pos pos, size int, ownerId ID) bool {
	pos = alignToGrid(pos)

	allChanged := [5]bool{}
	changeAll(allChanged[:])

	cellX, cellY := posToCell(pos)
	if s.isSpaceFree(cellX, cellY) {
		return false
	}

	newBomb := &bomb{
		pos:          pos,
		durationLeft: BombDelayFrames,
		id:           s.newID(),
		changed:      allChanged,
		state:        s,
		size:         size,
		bodies:       []body{newRect(pos, CellSize, CellSize)},
		ownerId:      ownerId,
		explodeDelay: ExplodeDelay,
	}

	(*s.Bombs)[newBomb.id] = newBomb
	(*s.arena)[cellY][cellX] = newBomb

	return true
}

func (b *bomb) getBody() *body {
	return &b.bodies[0]
}

func (s *state) isSpaceFree(x, y int) bool {

	return s.arena.getObstacle(x, y) != nil
}

func (b *bomb) act() {
	b.durationLeft -= 1
	if b.durationLeft == 0 && !b.exploding {
		b.startExplosion()
	}

	if b.exploding {
		b.explodeDelay--
		if b.explodeDelay == 0 {
			b.explode()
		}
	}
}

func (b *bomb) startExplosion() {
	b.exploding = true
}

func (b *bomb) explode() {
	player, ok := (*b.state.Players)[b.ownerId]
	if ok {
		player.currentBombs++
	}
	b.destroyed = true
	b.remove = true
	b.changed[BDestroy] = true
	b.changed[BRemove] = true
	b.state.newExplosion(b.pos, b.size)
}

func (b *bomb) destroy() {
	if b.destroyed {
		return
	}
	b.startExplosion()
}
