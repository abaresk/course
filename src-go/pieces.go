package course

// Tracks and nodes implement PieceData. Pieces are stored in their own
// structure so that they can be derefenced

type Piece struct {
	d PieceData
}

type PieceData interface {
	GetPiece(Direction) *Piece
	AddTrack(*Piece, Direction)
	AddNode(*Piece, Direction)
	LinkPiece(*Piece, *Piece, Direction)
	RemovePiece(Direction)
	reference(*Piece, Direction)
	unreference(Direction)
}

type Track struct {
	orient Orientation
	nexts  [2]*Piece
}

/*
Node: this represents any type of node (curve, 3way, 4way).
Tracks can be added or removed from ports on the go.
A "Nullnode" is when one of a piece's nexts is a nil pointer
*/
type Node struct {
	nexts [4]*Piece
}

/*
Pieces should be able to add a new node/track coming off a direction.
Pieces must also remove tracks or nodes in a direction.
*/

func (p *Piece) GetPiece(dir Direction) *Piece {
	return p.d.GetPiece(dir)
}

// Adds a track coming off the piece in Direction dir
func (p *Piece) AddTrack(dir Direction) {
	p.d.AddTrack(p, dir)
}

// Adds a node coming off the piece in Direction dir
func (p *Piece) AddNode(dir Direction) {
	p.d.AddNode(p, dir)
}

// Link this piece to another piece in Direction dir
// NOTE: should be the oppositie of RemovePiece
func (p *Piece) LinkPiece(p2 *Piece, dir Direction) {
	p.d.LinkPiece(p, p2, dir)
}

// Unlinks this piece with neighbor piece in Direction dir
func (p *Piece) RemovePiece(dir Direction) {
	p.d.RemovePiece(dir)
}

// Unlinks this piece with all its neighbors
func (p *Piece) Delete() {
	for i := Up; i < Right+1; i++ {
		p.d.RemovePiece(i)
	}
}

// Deletes the neighbor piece in Direction dir
func (p *Piece) DeleteNeighbor(dir Direction) {
	p2 := p.GetPiece(dir)
	p2.Delete()
}

func (n *Node) GetPiece(dir Direction) *Piece {
	return n.nexts[dir]
}

func (t *Track) GetPiece(dir Direction) *Piece {
	return t.nexts[Dir2Trackport[dir]]
}

func (p *Piece) reference(p2 *Piece, dir Direction) {
	p.d.reference(p2, dir)
}

func (p *Piece) unreference(dir Direction) {
	p.d.unreference(dir)
}

// Implementations
func (n *Node) AddTrack(p *Piece, dir Direction) {
	if n.nexts[dir] == nil {
		t := NewTrack(Dir2Orient[dir])
		t.nexts[Dir2Trackport[dir].Plus(1)] = p
		n.nexts[dir] = &Piece{t}
	}
}

func (t *Track) AddTrack(p *Piece, dir Direction) {
	port := Dir2Trackport[dir]
	if t.nexts[port] == nil {
		t2 := NewTrack(t.orient)
		t2.nexts[port.Plus(1)] = p
		t.nexts[port] = &Piece{t2}
	}
}

func (n *Node) AddNode(p *Piece, dir Direction) {
	if n.nexts[dir] == nil {
		n2 := new(Node)
		n2.nexts[dir.Plus(2)] = p
		n.nexts[dir] = &Piece{n2}
	}
}

func (t *Track) AddNode(p *Piece, dir Direction) {
	port := Dir2Trackport[dir]
	if t.nexts[port] == nil {
		n := new(Node)
		n.nexts[Port2Dir[Port{t.orient, port}].Plus(2)] = p
		t.nexts[port] = &Piece{n}
	}
}

func (n *Node) LinkPiece(p *Piece, p2 *Piece, dir Direction) {
	if n.nexts[dir] == nil {
		p2.reference(p, dir.Plus(2))
		n.reference(p2, dir)
	}
}

func (t *Track) LinkPiece(p *Piece, p2 *Piece, dir Direction) {
	port := Dir2Trackport[dir]
	if t.nexts[port] == nil {
		p2.reference(p, dir.Plus(2))
		t.reference(p2, dir)
	}
}

func (n *Node) RemovePiece(dir Direction) {
	if n.nexts[dir] != nil {
		n.nexts[dir].unreference(dir.Plus(2))
		n.unreference(dir)
	}
}

func (t *Track) RemovePiece(dir Direction) {
	port := Dir2Trackport[dir]
	if t.nexts[port] != nil {
		t.nexts[port].unreference(dir.Plus(2))
		t.unreference(dir)
	}
}

func (n *Node) reference(p2 *Piece, dir Direction) {
	n.nexts[dir] = p2
}

func (t *Track) reference(p2 *Piece, dir Direction) {
	port := Dir2Trackport[dir]
	t.nexts[port] = p2
}

func (n *Node) unreference(dir Direction) {
	n.nexts[dir] = nil
}

func (t *Track) unreference(dir Direction) {
	port := Dir2Trackport[dir]
	t.nexts[port] = nil
}

// Creators
func NewTrack(orient Orientation) *Track {
	t := new(Track)
	t.orient = orient
	return t
}

func NewTrackPiece(orient Orientation) *Piece {
	return &Piece{NewTrack(orient)}
}

func NewNodePiece() *Piece {
	n := new(Node)
	return &Piece{n}
}
