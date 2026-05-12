package game

import (
	"fmt"
	"math"
)

func bodiesTouch(A, B body) bool {

	switch {
	case A.shapeType == Circle && B.shapeType == Circle:
		pos := sub(A.pos, B.pos)
		return mag(pos) < A.radius+B.radius
	case A.shapeType == Circle && B.shapeType == Rectangle || B.shapeType == Circle && A.shapeType == Rectangle:
		circle := A
		rect := B
		if rect.shapeType == Circle && circle.shapeType == Rectangle {
			circle, rect = rect, circle
		}

		dist := pos{}
		dist.x = math.Abs(circle.pos.x - rect.pos.x)
		dist.y = math.Abs(circle.pos.y - rect.pos.y)

		if dist.x > (rect.width/2 + circle.radius) {
			return false
		}
		if dist.y > (rect.height/2 + circle.radius) {
			return false
		}

		if dist.x <= (rect.width / 2) {
			return true
		}
		if dist.y <= (rect.height / 2) {
			return true
		}

		cornerDist := math.Pow((dist.x-rect.width/2), 2) + math.Pow((dist.y-rect.height/2), 2)
		return (cornerDist <= (circle.radius * circle.radius))

	case A.shapeType == Rectangle && B.shapeType == Rectangle:
		xDist := math.Abs(A.pos.x - B.pos.x)
		yDist := math.Abs(A.pos.y - B.pos.y)
		if xDist < A.width+B.width && yDist < A.height+B.height {
			return true
		}
	}
	return false
}

func (s *state) moveAndCollide(mover body, vector pos) (pos, string) {
	str := ""
	str += fmt.Sprint("  mover.pos: ", mover.pos)
	cellPos := alignToGrid(mover.pos)
	str += fmt.Sprint("  cellpos: ", cellPos)
	// middleCell := s.arena.getObstacle(int(cellPos.x+0.0001), int(cellPos.y+0.0001))
	// if middleCell != nil {
	// 	middle := middleCell.getBody()
	// 	//if already inside an obstacle, you can move freely inside it, but only away from its center
	// 	if mover.pos.x+mover.width > middle.pos.x-middle.width/2 &&
	// 		mover.pos.x-mover.width < middle.pos.x+middle.width/2 &&
	// 		mover.pos.y+mover.width > middle.pos.y-middle.height/2 &&
	// 		mover.pos.y-mover.width < middle.pos.y+middle.height/2 {

	// 		// //following checks stop bomberman from moving towards the middle of the object it's already inside of
	// 		if dist(mover.pos, middle.pos) > 0.35 {
	// 			switch {
	// 			case mover.pos.x > middle.pos.x:
	// 				vector.x = max(0.0, vector.x)
	// 			case mover.pos.x < middle.pos.x:
	// 				vector.x = min(0.0, vector.x)
	// 			}
	// 			switch {
	// 			case mover.pos.y > middle.pos.y:
	// 				vector.y = max(0.0, vector.y)
	// 			case mover.pos.y < middle.pos.y:
	// 				vector.y = min(0.0, vector.y)
	// 			}
	// 		}
	// 		str += fmt.Sprint("  YE in obs ")
	// 	} else {
	// 		str += fmt.Sprint("  NO in obs ")
	// 	}
	// }

	mover.pos = add(mover.pos, vector)

	obstacles := s.arena.getObstacleSAroundHere(mover.pos)

	var closestObstacle *body
	smallestDist := 10000000000.0
	for _, ob := range obstacles {
		dist := dist(mover.pos, ob.getBody().pos)
		if smallestDist > dist {
			smallestDist = dist
			closestObstacle = ob.getBody()
		}
	}

	resolution, collision := resolve(mover, closestObstacle)

	resolution = mult(resolution, 1.001)
	str += fmt.Sprint("  collision: ", collision)

	return add(vector, resolution), str
}

func resolve(mover body, obstacle *body) (pos, bool) {
	if obstacle == nil {
		return pos{}, false
	}
	if mover.shapeType == Circle && obstacle.shapeType == Rectangle {
		closestRectPoint := pos{}
		closestRectPoint.y = clamp(mover.pos.y, obstacle.pos.y-obstacle.height/2, obstacle.pos.y+obstacle.height/2)
		closestRectPoint.x = clamp(mover.pos.x, obstacle.pos.x-obstacle.width/2, obstacle.pos.x+obstacle.width/2)
		distance := dist(closestRectPoint, mover.pos)
		if distance < mover.width {
			return mult(norm(sub(mover.pos, closestRectPoint)), math.Abs(distance-mover.width)), true
		}
	}

	return pos{}, false
}

func clamp(n, mn, mx float64) float64 {
	return max(mn, min(mx, n))
}
