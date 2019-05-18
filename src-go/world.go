/*
World object for a course. Everything runs through the World. It delegates
to the Piece and Pointmap objects.
*/
package course

/*
World should:
	- Add Piece at a specified point (check for potential merges)
		* Check for validity. If invalid, do nothing
	- Get Piece at a point at a specific layer (default layer == 0)
	- Delete Piece at a point at a specified layer (default layer == 0)
*/
type World struct {
	pmap *Pointmap
}

func (w *World) Init() {
	w.pmap = new(Pointmap)
	w.pmap.init()
}

func (w *World) AddNode(point Point, arg NodeArg) {
	var n Node
	if n = NewNode(arg); n == nil || !w.validNodePoint(n, point) {
		return
	}
	w.pmap.add(point, n)
	for _, p := range n.nodeTerritory(point) {
		if p != point {
			w.pmap.add(p, &NodeBody{n})
		}
	}
	w.makeMerges(n, point)
}

func (w *World) AddTrack(point Point, orient Orientation) {
	if !w.validTrackPoint(point, orient) {
		return
	}
	t := NewTrack(orient)
	w.pmap.add(point, t)
	w.makeMerges(t, point)
}

func (w *World) Get(point Point, layer int) Object {
	var out Object
	l := *w.pmap.get(point)
	if layer >= 0 && layer < len(l) {
		out = l[layer]
	}
	return out
}

// Used to get a piece at a point and layer. Returns nil
// if there is no *Piece located there
func (w *World) GetPiece(point Point, layer int) Piece {
	obj := w.Get(point, layer)
	piece, _ := obj.(Piece)
	return piece
}

// You can only delete Pieces (not NodeBody's)
func (w *World) Delete(point Point, layer int) {
	var piece Piece
	if piece = w.GetPiece(point, layer); piece == nil {
		return
	}

	// Delete all NodeBody's
	if _, ok := piece.(Node); ok {
		for _, t := range piece.(Node).nodeTerritory(point) {
			w.pmap.remove(w.Get(t, 0))
		}
	}
	w.pmap.remove(piece)
	erase(piece)
}

// There should be no Pieces in the node's territory
func (w *World) validNodePoint(n Node, point Point) bool {
	for _, p := range n.nodeTerritory(point) {
		if len(w.pmap.getObjectPieces(p)) != 0 {
			return false
		}
	}
	return true
}

// At point: 1 track w/ opposite orient allowed, no other Piece
// At adjacent points: anything allowed
// The following setup is valid:
//	  |
// 	──|──
//	  |
func (w *World) validTrackPoint(point Point, orient Orientation) bool {
	l := w.pmap.getObjectPieces(point)
	if len(l) == 0 {
		return true
	}
	if len(l) == 1 {
		if t, ok := l[0].(*Track); ok {
			return orient != t.orient
		}
	}
	return false
}

/*
Merge checking:
	- Node: 	Check each port. Can merge with a track that faces the node. Can
				also merge with a node whose center you are aligned with.
	- Track: 	Check both ports. Can merge with a track with same orientation.
				Can also merge with a node whose center you are aligned with.
	* You should merge with at most 1 Piece for each port

*/

// p == piece you're trying to merge
// point == center point of Track or Node
// ports == boundary of Track or Node
func (w *World) makeMerges(p Piece, point Point) {
	ports := p.ports(point)
	for _, port := range ports {
		l := *w.pmap.get(port)
	Loop:
		for _, obj := range l {
			switch obj.(type) {
			case *NodeBody:
				if merged := w.mergeNodeBody(p, obj.(*NodeBody), point); merged {
					break Loop
				}
			case Piece:
				if merged := w.mergePiece(p, obj.(Piece), point); merged {
					break Loop
				}
			}
		}
	}
}

func (w *World) mergeNodeBody(p Piece, b *NodeBody, center Point) bool {
	nodeCenter := w.pmap.find(b.center)
	if dir, aligned := nodeCenter.DirTo(center); aligned {
		linkPieces(b.center, p, dir)
		return true
	}
	return false
}

func (w *World) mergePiece(p, portPiece Piece, center Point) bool {
	objCenter := w.pmap.find(portPiece)
	track, isTrack := portPiece.(*Track)
	if dir, aligned := objCenter.DirTo(center); aligned && (!isTrack || track.orient == Dir2Orient[dir]) {
		linkPieces(portPiece, p, dir)
		return true
	}
	return false
}
