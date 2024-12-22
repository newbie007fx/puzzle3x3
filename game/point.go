package game

type Point struct {
	PointX int
	PointY int
}

func (p Point) Move(action ACTION) Point {
	x := p.PointX
	y := p.PointY

	switch action {
	case MOVE_UP:
		x = x - 1
	case MOVE_DOWN:
		x = x + 1
	case MOVE_LEFT:
		y = y - 1
	case MOVE_RIGHT:
		y = y + 1
	}

	return Point{
		PointX: x,
		PointY: y,
	}
}
