package course

type (
	Direction   int
	Trackport   int
	Orientation int
	PieceType   int
)

const (
	Up Direction = iota
	Left
	Down
	Right
)

func (dir Direction) Plus(c int) Direction {
	return Direction(modulo(int(dir)+c, 4))
}

const (
	Forward Trackport = iota
	Backward
)

func (port Trackport) Plus(c int) Trackport {
	return Trackport(modulo(int(port)+c, 2))
}

const (
	Vertical Orientation = iota
	Horizontal
)

var Dir2Orient = map[Direction]Orientation{
	Up:    Vertical,
	Left:  Horizontal,
	Down:  Vertical,
	Right: Horizontal,
}

var Dir2Trackport = map[Direction]Trackport{
	Up:    Forward,
	Left:  Backward,
	Down:  Backward,
	Right: Forward,
}

type Port struct {
	orient Orientation
	port   Trackport
}

var Port2Dir = map[Port]Direction{
	Port{Vertical, Forward}:    Up,
	Port{Vertical, Backward}:   Down,
	Port{Horizontal, Forward}:  Right,
	Port{Horizontal, Backward}: Left,
}

var UnitVector = map[Direction]Point{
	Up:    Point{0, 1},
	Left:  Point{-1, 0},
	Down:  Point{0, -1},
	Right: Point{1, 0},
}

type Point struct {
	x, y int
}

func (p Point) Add(other Point) Point {
	return Point{p.x + other.x, p.y + other.y}
}

func (p Point) Sub(other Point) Point {
	return Point{p.x - other.x, p.y - other.y}
}

func (p Point) Scale(c int) Point {
	return Point{p.x * c, p.y * c}
}

func (p Point) DirTo(other Point) (Direction, bool) {
	diff := other.Sub(p)
	if diff.x == 0 && diff.y > 0 {
		return Up, true
	} else if diff.x < 0 && diff.y == 0 {
		return Left, true
	} else if diff.x == 0 && diff.y < 0 {
		return Down, true
	} else if diff.x > 0 && diff.y == 0 {
		return Right, true
	}
	return Up, false
}
