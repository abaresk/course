package course

// Piece will be implemented by *Track, *fullNode, *halfNode, *CurveNode
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

type fullNode struct {
	radius int
	nexts  [4]Piece
}

type FullNodeClear struct {
	*fullNode
}

type FullNodeSet struct {
	flow Flow
	*fullNode
}

type FullNodeSwitched struct {
	sig Signal
	FullNodeSet
}

type halfNode struct {
	radius    int
	flowsLeft bool
	dir       Direction
	nexts     [3]Piece
}

type HalfNodeSet struct {
	*halfNode
}

type HalfNodeSwitched struct {
	sig Signal
	*halfNode
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

func (n *fullNode) getNext(dir Direction) (Piece, bool) {
	if index, ok := n.getIndex(dir); !ok {
		return nil, false
	} else {
		return n.nexts[index], true
	}
}

func (n *halfNode) getNext(dir Direction) (Piece, bool) {
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

func (n *fullNode) setNext(p Piece, dir Direction) bool {
	if index, ok := n.getIndex(dir); !ok {
		return false
	} else {
		n.nexts[index] = p
		return true
	}
}

func (n *halfNode) setNext(p Piece, dir Direction) bool {
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

func (n *fullNode) ports(point Point) []Point {
	out := []Point{}
	for dir := Up; dir <= Right; dir++ {
		out = append(out, point.Add(UnitVector[dir].Scale(n.radius+1)))
	}
	return out
}

func (n *halfNode) ports(point Point) []Point {
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

func (n *fullNode) getIndex(dir Direction) (int, bool) {
	return int(dir), true
}

func (n *halfNode) getIndex(dir Direction) (int, bool) {
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
func NewTrack(orient Orientation) Piece {
	return &Track{orient: orient}
}

func newFullNode() *fullNode {
	return &fullNode{radius: 1}
}

func newHalfNode(dir Direction, flowsLeft bool) *halfNode {
	return &halfNode{radius: 1, dir: dir, flowsLeft: flowsLeft}
}

func NewFullNodeClear() Piece {
	return FullNodeClear{newFullNode()}
}

func NewFullNodeSet(flow Flow) Piece {
	return FullNodeSet{flow, newFullNode()}
}

func NewFullNodeSwitched(flow Flow, sig Signal) Piece {
	return FullNodeSwitched{sig, FullNodeSet{flow, newFullNode()}}
}

func NewHalfNodeSet(dir Direction, flowsLeft bool) Piece {
	return HalfNodeSet{newHalfNode(dir, flowsLeft)}
}

func NewHalfNodeSwithced(dir Direction, flowsLeft bool, sig Signal) Piece {
	return HalfNodeSwitched{sig, newHalfNode(dir, flowsLeft)}
}

func NewCurveNode(quad Quadrant) Piece {
	return &CurveNode{radius: 1, quad: quad}
}
