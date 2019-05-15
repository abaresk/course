package course

type ItemSlice []Item

// bi-directional mapping
type Pointmap struct {
	points map[Point]ItemSlice
	pieces map[Item]Point
}

func (p *Pointmap) init() {
	p.points = make(map[Point]ItemSlice)
	p.pieces = make(map[Item]Point)
}

func (p *Pointmap) add(point Point, piece Item) {
	p.points[point] = append(p.points[point], piece)
	p.pieces[piece] = point
}

func (p *Pointmap) remove(item Item) {
	point := p.find(item)
	p.points[point] = removePieceFromSlice(p.get(point), item)
	delete(p.pieces, item)
}

func (p *Pointmap) get(point Point) ItemSlice {
	return p.points[point]
}

func (p *Pointmap) find(piece Item) Point {
	return p.pieces[piece]
}

func removePieceFromSlice(l ItemSlice, piece Item) ItemSlice {
	for i, v := range l {
		if piece == v {
			l[i] = l[len(l)-1]
			l[len(l)-1] = nil
			l = l[:len(l)-1]
			return l
		}
	}
	return l
}
