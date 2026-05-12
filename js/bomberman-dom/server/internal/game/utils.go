package game

func changeAll(slice []bool) {
	for i := range slice {
		slice[i] = true
	}
}

func gridPos2Pos(x, y int) pos {
	return pos{float64(x) * CellSize, float64(y) * CellSize}
}

func pos2GridPos(pos pos) (int, int) {
	return int(pos.x / CellSize), int(pos.y / CellSize)
}

func visualizeWalls(wallGrid *arena) {
	return
}
