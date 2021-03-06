package course

type Node interface {
	getRadius() int
	nodeTerritory(Point) []Point
	Piece
}

func (n *fullNode) getRadius() int {
	return n.radius
}

func (n *halfNode) getRadius() int {
	return n.radius
}

func (n *CurveNode) getRadius() int {
	return n.radius
}

func (n *fullNode) nodeTerritory(point Point) []Point {
	out := []Point{}
	for x := point.x - n.radius; x <= point.x+n.radius; x++ {
		for y := point.y - n.radius; y <= point.y+n.radius; y++ {
			out = append(out, Point{x, y})
		}
	}
	return out
}

func (n *halfNode) nodeTerritory(point Point) []Point {
	out := []Point{}
	for dx := -n.radius; dx <= n.radius; dx++ {
		for dy := -n.radius; dy <= n.radius; dy++ {
			if halfFilter(dx, dy, n.dir) {
				out = append(out, Point{point.x + dx, point.y + dy})
			}
		}
	}
	return out
}

func (n *CurveNode) nodeTerritory(point Point) []Point {
	out := []Point{}
	for dx := -n.radius; dx <= n.radius; dx++ {
		for dy := -n.radius; dy <= n.radius; dy++ {
			if quadrantFilter(dx, dy, n.quad) {
				out = append(out, Point{point.x + dx, point.y + dy})
			}
		}
	}
	return out
}

func halfFilter(dx, dy int, dir Direction) bool {
	switch dir {
	case Up:
		return dy >= 0
	case Left:
		return dx <= 0
	case Down:
		return dy <= 0
	case Right:
		return dx >= 0
	}
	return false
}

func quadrantFilter(dx, dy int, quad Quadrant) bool {
	switch quad {
	case One:
		return dx >= 0 && dy >= 0
	case Two:
		return dx <= 0 && dy >= 0
	case Three:
		return dx <= 0 && dy <= 0
	case Four:
		return dx >= 0 && dy <= 0
	}
	return false
}
