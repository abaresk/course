package course

type PieceSlice []*Piece

// bi-directional mapping
type Pointmap struct {
	points map[Point]PieceSlice
	pieces map[*Piece]Point
}

func (p *Pointmap) Init() {
	p.points = make(map[Point]PieceSlice)
	p.pieces = make(map[*Piece]Point)
}

func (p *Pointmap) Add(point Point, piece *Piece) {
	p.points[point] = append(p.points[point], piece)
	p.pieces[piece] = point
}

func (p *Pointmap) Remove(piece *Piece) {
	removePieceFromSlice(p.Get(p.Find(piece)), piece)
	delete(p.pieces, piece)
}

func (p *Pointmap) Get(point Point) PieceSlice {
	return p.points[point]
}

func (p *Pointmap) Find(piece *Piece) Point {
	return p.pieces[piece]
}

func removePieceFromSlice(l PieceSlice, piece *Piece) {
	for i, v := range l {
		if piece == v {
			l[i] = l[len(l)-1]
			l[len(l)-1] = nil
			l = l[:len(l)-1]
			return
		}
	}
}
