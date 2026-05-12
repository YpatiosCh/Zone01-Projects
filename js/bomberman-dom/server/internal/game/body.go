package game

type pos struct {
	x float64
	y float64
}

func newCircle(pos pos, radius float64) body {
	return body{Circle, pos, radius, radius, radius}
}

func newRect(pos pos, width, height float64) body {
	return body{Rectangle, pos, width, height, -1}
}
