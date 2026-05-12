package game

type ID int

type input struct {
	left  bool
	right bool
	up    bool
	down  bool
	bomb  bool
}

func (s *state) newID() ID {
	s.IDCounter++
	return s.IDCounter
}

type Walls struct {
	walls map[ID]*wall
}
