package course

type (
	Direction   int
	Trackport   int
	Orientation int
	Quadrant    int
	PieceType   int
	Signal      int
	Flow        int
)

const (
	One Quadrant = iota
	Two
	Three
	Four
)

const (
	Up Direction = iota
	Left
	Down
	Right
)

const (
	SignalRed Signal = iota
	SignalBlue
	SignalGreen
	SignalYellow
)

const (
	FlowI        Flow = iota // first/third Quadrant
	FlowStraight             // criss-cross
	FlowII                   // second/fourth Quadrant
)

func (dir Direction) Plus(c int) Direction {
	return Direction(modulo(int(dir)+c, 4))
}

func (dir Direction) Minus(c int) Direction {
	return Direction(modulo(int(dir)-c, 4))
}

const (
	Backward Trackport = iota
	Forward
)

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
