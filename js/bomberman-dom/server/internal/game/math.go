package game

import "math"

func alignToGrid(pos pos) pos {
	pos.x = math.RoundToEven(pos.x)
	pos.y = math.RoundToEven(pos.y)
	return pos
}

func posToCell(p pos) (int, int) {
	p = alignToGrid(p)
	return int(p.x + 0.001), int(p.y + 0.001)
}

func closeEnough(a, b pos) bool {
	if dist(a, b) < 0.001 {
		return true
	}

	return false
}

func norm(p pos) pos {
	mag := mag(p)
	if mag == 0 {
		return pos{0, 0}
	}

	return pos{p.x / mag, p.y / mag}
}

func mag(pos pos) float64 {
	if pos.x == 0 && pos.y == 0 {
		return 0
	}

	return math.Sqrt(float64(pos.x*pos.x + pos.y*pos.y))
}

func sub(a, b pos) pos {
	a.x -= b.x
	a.y -= b.y

	return a
}

func add(a, b pos) pos {
	a.x += b.x
	a.y += b.y

	return a
}

func mult(a pos, multipler float64) pos {
	a.x *= multipler
	a.y *= multipler

	return a
}

func sameCell(a, b pos) bool {
	a = alignToGrid(a)
	b = alignToGrid(b)
	return closeEnough(a, b)
}

func dist(a, b pos) float64 {
	a = sub(a, b)
	return mag(a)
}

func diffSigns(a, b float64) bool {
	if (a < 0 && b < 0) || (a > 0 && b > 0) {
		return false
	}
	return true
}
