/*
World object for a course. Everything runs through the World. It delegates
to the Piece and Pointmap objects.
*/
package course

// TODO: Think about what operations you want World to do.
// Then what state info would it need to hold to accomplish that?

/*
Stuff World should do:
	- Add Piece at a specified point (check for potential merges)
		* Check for validity. If invalid, do nothing
	- Delete piece at a point at a specified layer (default layer == 0)
*/
type World struct {
	pmap *Pointmap
}

func (w *World) Init() {
	w.pmap = new(Pointmap)
	w.pmap.Init()
}

func (w *World) AddNode(point Point) {
	if !w.validNodePoint(point) {
		return
	}
	n := NewNodePiece()
	w.pmap.Add(point, n)
	for _, p := range nodeTerritory(point) {
		if p != point {
			w.pmap.Add(p, &NodeBody{n})
		}
	}
	w.makeMerges(n, point, nodePorts(point))
}

func (w *World) AddTrack(point Point, orient Orientation) {
	if !w.validTrackPoint(point, orient) {
		return
	}
	t := NewTrackPiece(orient)
	w.pmap.Add(point, t)
	w.makeMerges(t, point, trackPorts(point, orient))
}

// There should be no Pieces in the node's territory
func (w *World) validNodePoint(point Point) bool {
	for _, p := range nodeTerritory(point) {
		if len(w.pmap.Get(p)) != 0 {
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
	l := w.pmap.Get(point)
	if len(l) == 0 {
		return true
	}
	if len(l) == 1 {
		if t, ok := extractTrack(l[0]); ok {
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
func (w *World) makeMerges(p *Piece, point Point, ports []Point) {
	for _, port := range ports {
		l := w.pmap.Get(port)
	Loop:
		for _, item := range l {
			switch item.(type) {
			case *NodeBody:
				if merged := w.mergeNodeBody(p, item.(*NodeBody), point); merged {
					break Loop
				}
			case *Piece:
				if merged := w.mergePiece(p, item.(*Piece), point); merged {
					break Loop
				}
			}
		}
	}
}

func (w *World) mergeNodeBody(p *Piece, b *NodeBody, center Point) bool {
	nodeCenter := w.pmap.Find(b.center)
	if dir, aligned := nodeCenter.DirTo(center); aligned {
		b.center.LinkPiece(p, dir)
		return true
	}
	return false
}

func (w *World) mergePiece(p, portPiece *Piece, center Point) bool {
	itemCenter := w.pmap.Find(portPiece)
	switch portPiece.GetData().(type) {
	case *Node:
		if dir, aligned := itemCenter.DirTo(center); aligned {
			portPiece.LinkPiece(p, dir)
			return true
		}
	case *Track:
		track := portPiece.GetData().(*Track)
		if dir, aligned := itemCenter.DirTo(center); aligned && track.orient == Dir2Orient[dir] {
			portPiece.LinkPiece(p, dir)
			return true
		}
	}
	return false
}

// Box of radius 1 surrounding the node
func nodeTerritory(point Point) []Point {
	radius := 1
	out := []Point{}
	for x := point.x - radius; x <= point.x+radius; x++ {
		for y := point.y - radius; y <= point.y+radius; y++ {
			out = append(out, Point{x, y})
		}
	}
	return out
}

func nodePorts(point Point) []Point {
	radius := 1
	out := []Point{}
	for dir := Up; dir <= Right; dir++ {
		out = append(out, point.Add(UnitVector[dir].Scale(radius+1)))
	}
	return out
}

func trackPorts(point Point, orient Orientation) []Point {
	v := UnitVector[Port2Dir[Port{orient, Forward}]]
	return []Point{point.Add(v), point.Add(v.Scale(-1))}
}
