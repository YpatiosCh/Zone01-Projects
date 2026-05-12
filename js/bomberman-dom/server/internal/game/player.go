package game

import (
	"fmt"
	"slices"
)

type publicPlayer struct {
	Sprite       *string  `json:"sprite,omitempty"`
	Color        *string  `json:"color,omitempty"`
	Lives        *int     `json:"lives,omitempty"`
	MaxLives     *int     `json:"maxLives,omitempty"`
	CurrentBombs *int     `json:"currentBombs,omitempty"`
	MaxBombs     *int     `json:"maxBombs,omitempty"`
	Invincible   *bool    `json:"invincible,omitempty"`
	Alive        *bool    `json:"alive,omitempty"`
	Y            *float64 `json:"y,omitempty"`
	X            *float64 `json:"x,omitempty"`
	Animate      *string  `json:"animate,omitempty"`
	Direction    *string  `json:"direction,omitempty"`
	ItemPickedUp *string  `json:"item_picked_up,omitempty"`

	Destroyed *bool `json:"destroy,omitempty"`

	Remove *bool `json:"remove,omitempty"`

	changed bool
}

func (p *player) publicPlayer(fullView bool) publicPlayer {
	pp := publicPlayer{}
	pp.changed = slices.Contains(p.changed[:], true)

	if p.changed[PSprite] == true || fullView {
		pp.Sprite = &p.Sprite
	}
	if p.changed[PColor] == true || fullView {
		pp.Color = &p.Color
	}
	if p.changed[PLives] == true || fullView {
		pp.Lives = &p.Lives
	}
	if p.changed[PMaxLives] == true || fullView {
		pp.MaxLives = &p.MaxLives
	}
	if p.changed[PCurrentBombs] == true || fullView {
		pp.CurrentBombs = &p.currentBombs
	}
	if p.changed[PMaxBombs] == true || fullView {
		pp.MaxBombs = &p.maxBombs
	}
	if p.changed[PInvincible] == true || fullView {
		pp.Invincible = &p.invincible
	}
	if p.changed[PAlive] == true || fullView {
		pp.Alive = &p.Alive
	}
	if p.changed[PY] == true || fullView {
		pp.Y = &p.pos.y
	}
	if p.changed[PX] == true || fullView {
		pp.X = &p.pos.x
	}
	if p.changed[PDestroy] == true || fullView {
		pp.Destroyed = &p.destroyed
	}
	if p.changed[PRemove] == true || fullView {
		pp.Remove = &p.remove
	}
	if p.changed[PAnimate] == true || fullView {
		pp.Animate = &p.animate
	}
	if p.changed[PDirection] == true || fullView {
		pp.Direction = &p.direction
	}
	if p.changed[PItemPickedUp] == true || fullView {
		pp.ItemPickedUp = &p.itemPickedUp
	}

	if fullView == false {
		p.changed = [15]bool{}
	}

	return pp
}

const (
	PSprite = iota
	PColor
	PLives
	PMaxLives
	PCurrentBombs
	PMaxBombs
	PInvincible
	PAlive
	PY
	PX
	PDestroy
	PRemove
	PAnimate
	PDirection
	PItemPickedUp
)

type player struct {
	//public
	Sprite       string
	Color        string
	Lives        int
	MaxLives     int
	currentBombs int
	maxBombs     int
	invincible   bool
	Alive        bool
	pos          pos
	remove       bool
	destroyed    bool
	animate      string
	direction    string
	itemPickedUp string

	//private
	name               string
	changed            [15]bool
	explosionSize      int
	input              input
	ID                 ID
	speed              float64
	collisionRadius    float64
	bodies             []body
	state              *state
	bombCooldownMax    int
	bombCooldown       int
	impulse            pos
	invulTimer         int
	bombPlaceStopTimer int
}

var defaultBody = body{}

func (s *state) newPlayer(x, y float64, color, name string) {
	p := &player{
		ID:              s.newID(),
		Color:           color,
		Lives:           maxLives,
		MaxLives:        maxLives,
		maxBombs:        BombsMax,
		currentBombs:    BombsMax,
		Alive:           true,
		speed:           PlayerSpeed,
		bodies:          []body{newCircle(pos{x, y}, 0.45)},
		state:           s,
		explosionSize:   ExplosionStartSize,
		pos:             pos{x, y},
		bombCooldownMax: BombCooldownMax,
		invulTimer:      InvincibleDuration,
		name:            name,
	}

	changeAll(p.changed[:])

	(*s.Players)[p.ID] = p

	fmt.Printf("new player: %#v\n", p)
}

func (p *player) UpdatePosition(x, y float64) {
	p.pos.x = x
	p.pos.y = y
	p.bodies[0].pos.x = x
	p.bodies[0].pos.y = y
}

func (p *player) act() {

	if p.bombPlaceStopTimer > 0 {
		p.bombPlaceStopTimer--
	}

	if p.input.bomb == true && p.currentBombs > 0 && p.bombPlaceStopTimer == 0 {
		bombCreated := p.state.newBomb(p.pos, p.explosionSize, p.ID)
		if bombCreated {
			p.currentBombs -= 1
			p.bombPlaceStopTimer = BombPlaceStopTimer
		}
	}

	if p.invincible && p.state.endOfGameTimer == EndOfGameTimer {
		p.invulTimer--
		if p.invulTimer == 0 {
			p.invincible = false
			p.invulTimer = InvincibleDuration
		}
	}

	for _, powerup := range *p.state.Powerups {
		if bodiesTouch(p.bodies[0], powerup.bodies[0]) {

			p.itemPickedUp = powerup.name
			p.changed[PItemPickedUp] = true

			powerup.destroy()
			switch powerup.name {
			case "power-bomb":
				p.maxBombs++
				p.currentBombs++
			case "speed":
				p.speed += speedPowerupInc
				if p.speed > 0.09 {
					p.speed = 0.09
				}
			case "power":
				p.explosionSize++
			case "lives":
				p.Lives++
				p.changed[PLives] = true
			}
		}
	}
}

func (p *player) move() {
	if !p.Alive {
		return
	}

	//calc movement vector from inputs
	x := 0.0
	y := 0.0
	// if p.ID == ID(2) {
	// 	y += p.speed
	// }
	if p.input.up {
		y -= p.speed
	}
	if p.input.down {
		y += p.speed
	}
	if p.input.left {
		x -= p.speed
	}
	// if p.ID == ID(1) {
	// 	x += p.speed
	// }
	if p.input.right {
		x += p.speed
	}

	//apply movement
	position := norm(pos{x, y})

	yDelta := position.y*p.speed + p.impulse.y
	xDelta := position.x*p.speed + p.impulse.x

	if p.impulse.x > 0 {
		p.impulse.x *= impulseBleed
		if p.impulse.x < 0.01 {
			p.impulse.x = 0
		}
	} else if p.impulse.x < 0 {
		p.impulse.x *= impulseBleed
		if p.impulse.x > -0.01 {
			p.impulse.x = 0
		}
	}

	if p.impulse.y > 0 {
		p.impulse.y *= impulseBleed
		if p.impulse.y < 0.01 {
			p.impulse.y = 0
		}
	} else if p.impulse.y < 0 {
		p.impulse.y *= impulseBleed
		if p.impulse.y > -0.01 {
			p.impulse.y = 0
		}
	}

	resolution, _ := p.state.moveAndCollide(p.bodies[0], pos{xDelta, yDelta})
	// if p.ID == 1 {
	// 	fmt.Println("player pos: ", p.pos, ", resolution: ", resolution, ", inputs: ", p.input, str)
	// 	// time.Sleep(time.Millisecond * 30)
	// }

	if yDelta > 0 {
		p.direction = "down"
		p.changed[PDirection] = true
	} else if yDelta < 0 {
		p.direction = "up"
		p.changed[PDirection] = true
	}

	if xDelta > 0 {
		p.direction = "right"
		p.changed[PDirection] = true
	} else if xDelta < 0 {
		p.direction = "left"
		p.changed[PDirection] = true
	}

	if resolution.y < 0 && p.impulse.y > 0 {
		p.impulse.y = 0
	}
	if resolution.y > 0 && p.impulse.y < 0 {
		p.impulse.y = 0
	}
	if resolution.x < 0 && p.impulse.x > 0 {
		p.impulse.x = 0
	}
	if resolution.x > 0 && p.impulse.x < 0 {
		p.impulse.x = 0
	}

	if diffSigns(yDelta, resolution.y) {
		yDelta = resolution.y
	}
	if diffSigns(xDelta, resolution.x) {
		xDelta = resolution.x
	}

	if yDelta != 0 {
		p.changed[PY] = true
		p.pos.y += yDelta
		p.bodies[0].pos.y += yDelta
	}

	if xDelta != 0 {
		p.changed[PX] = true
		p.pos.x += xDelta
		p.bodies[0].pos.x += xDelta
	}

	if yDelta != 0 || xDelta != 0 {
		if p.animate != "walking" {
			p.animate = "walking"
			p.changed[PAnimate] = true
		}
	} else {
		if p.animate != "idle" {
			p.animate = "idle"
			p.changed[PAnimate] = true
		}
	}

	//return new state if it changed
	// if beforeX != p.X || beforeY != p.Y {
	// 	return fmt.Sprint("%s:")
	// }
}

func (p *player) takeDamage(attackerPos pos) {
	if p.invincible {
		return
	}

	p.impulse = mult(norm(sub(p.pos, attackerPos)), BombPushPower)

	p.Lives -= 1
	p.changed[PLives] = true
	if p.Lives == 0 {
		p.Alive = false
		p.changed[PAlive] = true

		p.destroyed = true
		p.changed[PDestroy] = true
	} else {
		p.invincible = true
	}
}
