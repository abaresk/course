package course

import "container/heap"

/*
ObjectHeap stores objects according to their priority levels
Priorities:
	0	-> course pieces: tracks, nodes
	1 	-> enemies: EnemyParts
	2	-> items: 	power-ups, portals
*/
type ObjectHeap []Object

func (h ObjectHeap) Len() int           { return len(h) }
func (h ObjectHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h ObjectHeap) Less(i, j int) bool { return getObjectType(h[i]) < getObjectType(h[j]) }

func (h *ObjectHeap) Push(x interface{}) {
	*h = append(*h, x.(Object))
}

func (h *ObjectHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// bi-directional mapping
type Pointmap struct {
	points map[Point]*ObjectHeap
	pieces map[Object]Point
}

func (p *Pointmap) init() {
	p.points = make(map[Point]*ObjectHeap)
	p.pieces = make(map[Object]Point)
}

func (p *Pointmap) add(point Point, piece Object) {
	// ObjectHeap needs to be initialized
	if l, _ := p.points[point]; l == nil {
		h := &ObjectHeap{}
		heap.Init(h)
		p.points[point] = h
	}
	heap.Push(p.points[point], piece)
	p.pieces[piece] = point
}

func (p *Pointmap) remove(piece Object) {
	point := p.find(piece)
	removeObjFromHeap(p.get(point), piece)
	delete(p.pieces, piece)
}

func (p *Pointmap) get(point Point) *ObjectHeap {
	if out := p.points[point]; out != nil {
		return out
	}
	return &ObjectHeap{}
}

func (p *Pointmap) find(piece Object) Point {
	return p.pieces[piece]
}

// Pop until you find the object, then push everything back on
func removeObjFromHeap(h *ObjectHeap, obj Object) {
	tmp := []Object{}
	for h.Len() > 0 {
		x := heap.Pop(h).(Object)
		if x == obj {
			break
		}
		tmp = append(tmp, x)
	}
	for _, v := range tmp {
		heap.Push(h, v)
	}
}

// Pieces are always first to get popped
func (p *Pointmap) getObjectPieces(point Point) []Object {
	out := []Object{}
	if h := p.get(point); h != nil {
		for h.Len() > 0 {
			x := heap.Pop(h).(Object)
			if getObjectType(x) != ObjectPiece {
				heap.Push(h, x)
				break
			}
			out = append(out, x)
		}
		for _, v := range out {
			heap.Push(h, v)
		}
	}
	return out
}
