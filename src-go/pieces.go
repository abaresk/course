package course

// Piece will be implemented by *Track, *FullNode, *HalfNode, *CurveNode
type Piece interface {
	getNext(Direction) (Piece, bool)
	setNext(Piece, Direction) bool
	ports(Point) []Point
	Object
}

// Link p1 to p2 via dir
func linkPieces(p1, p2 Piece, dir Direction) bool {
	// fails if either are linked
	if l, ok := p1.getNext(dir); !ok || l != nil {
		return false
	}
	if l, ok := p2.getNext(dir.Plus(2)); !ok || l != nil {
		return false
	}
	p1.setNext(p2, dir)
	p2.setNext(p1, dir.Plus(2))
	return true
}

func unlinkNeighbor(p Piece, dir Direction) bool {
	if n, ok := p.getNext(dir); !ok {
		return false
	} else {
		return unlinkPieces(p, n, dir)
	}
}

// Unlink p1 from p2 via dir
func unlinkPieces(p1, p2 Piece, dir Direction) bool {
	// fails if either are unlinked
	if l, ok := p1.getNext(dir); !ok || l == nil {
		return false
	}
	if l, ok := p2.getNext(dir.Plus(2)); !ok || l == nil {
		return false
	}
	p1.setNext(nil, dir)
	p2.setNext(nil, dir.Plus(2))
	return true
}

func erase(p Piece) {
	for dir := Up; dir <= Right; dir++ {
		unlinkNeighbor(p, dir)
	}
}

type FullNode struct {
	radius int
	nexts  [4]Piece
}

type HalfNode struct {
	radius int
	dir    Direction
	nexts  [3]Piece
}

type CurveNode struct {
	radius int
	quad   Quadrant
	nexts  [2]Piece
}

type Track struct {
	orient Orientation
	nexts  [2]Piece
}

func (n *FullNode) getNext(dir Direction) (Piece, bool) {
	if index, ok := n.getIndex(dir); !ok {
		return nil, false
	} else {
		return n.nexts[index], true
	}
}

func (n *HalfNode) getNext(dir Direction) (Piece, bool) {
	if index, ok := n.getIndex(dir); !ok {
		return nil, false
	} else {
		return n.nexts[index], true
	}
}

func (n *CurveNode) getNext(dir Direction) (Piece, bool) {
	if index, ok := n.getIndex(dir); !ok {
		return nil, false
	} else {
		return n.nexts[index], true
	}
}

func (t *Track) getNext(dir Direction) (Piece, bool) {
	if index, ok := t.getIndex(dir); !ok {
		return nil, false
	} else {
		return t.nexts[index], true
	}
}

func (n *FullNode) setNext(p Piece, dir Direction) bool {
	if index, ok := n.getIndex(dir); !ok {
		return false
	} else {
		n.nexts[index] = p
		return true
	}
}

func (n *HalfNode) setNext(p Piece, dir Direction) bool {
	if index, ok := n.getIndex(dir); !ok {
		return false
	} else {
		n.nexts[index] = p
		return true
	}
}
func (n *CurveNode) setNext(p Piece, dir Direction) bool {
	if index, ok := n.getIndex(dir); !ok {
		return false
	} else {
		n.nexts[index] = p
		return true
	}
}
func (t *Track) setNext(p Piece, dir Direction) bool {
	if index, ok := t.getIndex(dir); !ok {
		return false
	} else {
		t.nexts[index] = p
		return true
	}
}

func (n *FullNode) ports(point Point) []Point {
	out := []Point{}
	for dir := Up; dir <= Right; dir++ {
		out = append(out, point.Add(UnitVector[dir].Scale(n.radius+1)))
	}
	return out
}

func (n *HalfNode) ports(point Point) []Point {
	out := []Point{}
	for i := 0; i < 3; i++ {
		out = append(out, point.Add(UnitVector[n.dir.Minus(1-i)].Scale(n.radius+1)))
	}
	return out
}

func (n *CurveNode) ports(point Point) []Point {
	out := []Point{}
	for i := 0; i < 2; i++ {
		out = append(out, point.Add(UnitVector[Direction(n.quad).Minus(i)].Scale(n.radius+1)))
	}
	return out
}

func (t *Track) ports(point Point) []Point {
	v := UnitVector[Port2Dir[Port{t.orient, Forward}]]
	return []Point{point.Add(v), point.Add(v.Scale(-1))}
}

func (n *FullNode) getIndex(dir Direction) (int, bool) {
	return int(dir), true
}

func (n *HalfNode) getIndex(dir Direction) (int, bool) {
	if n.dir == dir.Plus(2) {
		return 0, false
	}
	index := n.dir.Plus(1).Minus(int(dir))
	return int(index), true
}

func (n *CurveNode) getIndex(dir Direction) (int, bool) {
	switch Direction(n.quad) {
	case dir:
		return 0, true
	case dir.Plus(1):
		return 1, true
	default:
		return 0, false
	}
}

func (t *Track) getIndex(dir Direction) (int, bool) {
	switch {
	case t.orient == Vertical && dir == Up:
		return 1, true
	case t.orient == Vertical && dir == Down:
		return 0, true
	case t.orient == Horizontal && dir == Right:
		return 1, true
	case t.orient == Horizontal && dir == Left:
		return 0, true
	default:
		return 0, false
	}
}

// Constructors
func NewTrack(orient Orientation) *Track {
	return &Track{orient: orient}
}

func NewNode(arg NodeArg) Node {
	switch arg.(type) {
	case FullNodeArg:
		return &FullNode{radius: 1}
	case HalfNodeArg:
		return &HalfNode{radius: 1, dir: arg.(HalfNodeArg).dir}
	case CurveNodeArg:
		return &CurveNode{radius: 1, quad: arg.(CurveNodeArg).quad}
	}
	return nil
}
