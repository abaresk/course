package course

import (
	"testing"
)

/*
Create the following:
	bbb bbb
	b*b-b*b
	bbb bbb
	 |   |
	bbb bbb
	b*b-b*b
	bbb bbb
*/
func TestMakeLoop(t *testing.T) {
	w := new(World)
	w.Init()

	w.AddNode(Point{0, 0})
	w.AddTrack(Point{2, 0}, Horizontal)
	w.AddNode(Point{4, 0})
	w.AddTrack(Point{4, 2}, Vertical)
	w.AddNode(Point{4, 4})
	w.AddTrack(Point{2, 4}, Horizontal)
	w.AddNode(Point{0, 4})
	w.AddTrack(Point{0, 2}, Vertical)

	tr1, tr2 := w.pmap.Get(Point{0, 0})[0].(*Piece).getNext(Up), w.pmap.Get(Point{0, 4})[0].(*Piece).getNext(Down)
	if tr1 != tr2 {
		t.Fatalf("Improper merging at end of loop")
	}

	n1, n2 := w.pmap.Get(Point{0, 2})[0].(*Piece).getNext(Down), w.pmap.Get(Point{2, 0})[0].(*Piece).getNext(Left)
	if n1 != n2 {
		t.Fatalf("Improper merging at end of loop")
	}
}

/*
Create the following:
	ntn
	t t
	ntn
*/
func TestMakeLoopBad(t *testing.T) {
	w := new(World)
	w.Init()

	w.AddNode(Point{0, 0})              // Works!
	w.AddTrack(Point{1, 0}, Horizontal) // invalid
	w.AddNode(Point{2, 0})              // invalid
	w.AddTrack(Point{2, 1}, Vertical)   // Works!
	w.AddNode(Point{2, 2})              // invalid
	w.AddTrack(Point{1, 2}, Horizontal) // Works!
	w.AddNode(Point{0, 2})              // invalid
	w.AddTrack(Point{1, 1}, Vertical)   // invalid

	if len(w.pmap.pieces) > 11 {
		t.Fatalf("Erroneous pieces should be ignored")
	}
}

/*
Create the following:
	bbb bbb bbb
	b1b b2b b3b
	bbb bbb bbb

	bbb bbb bbb
	b4b b5b b6b
	bbb bbb bbb

Result: Each node should be linked to its neighbors
*/
func TestNodeChain(t *testing.T) {
	w := new(World)
	w.Init()

	w.AddNode(Point{0, 0}) // 4
	w.AddNode(Point{6, 3}) // 3
	w.AddNode(Point{3, 0}) // 5
	w.AddNode(Point{3, 3}) // 2
	w.AddNode(Point{0, 3}) // 1
	w.AddNode(Point{6, 0}) // 6

	// try to get at 5 from each side
	n4, n3 := w.pmap.Get(Point{0, 0})[0].(*Piece), w.pmap.Get(Point{6, 3})[0].(*Piece)
	if n4.getNext(Up).getNext(Right).getNext(Down) != n3.getNext(Down).getNext(Left) {
		t.Fatalf("Node chain isn't connected properly")
	}

}

/*
Create the following:
	bbb bbb
	b*b b*b
	bbb bbb
	 bbb
	 b*b
	 bbb

Expected: top two should link, but bottom should not link to either
*/
func TestMisalignedNodes(t *testing.T) {
	w := new(World)
	w.Init()

	w.AddNode(Point{0, 0})
	w.AddNode(Point{3, 0})
	w.AddNode(Point{1, -3})

	if n1 := w.pmap.Get(Point{0, 0})[0].(*Piece); n1 != n1.getNext(Right).getNext(Left) {
		t.Fatalf("Nodes should have connected")
	}
	for _, n := range w.pmap.Get(Point{1, -3})[0].(*Piece).GetData().(*Node).nexts {
		if n != nil {
			t.Fatalf("Misaligned node should not be connected to other nodes")
		}
	}
}

/*
Create a track overlap.
Add a connecting track (make sure it merges to the correct one)
Add parallel track. Then add crossover to parallel track (make sure it merges with the other)

Final layout:
	 |
	++
*/
func TestTrackOverlap(t *testing.T) {
	w := new(World)
	w.Init()

	w.AddTrack(Point{0, 0}, Horizontal)
	w.AddTrack(Point{0, 0}, Vertical)

	if len(w.pmap.Get(Point{0, 0})) != 2 {
		t.Fatalf("Should be able to overlap 2 tracks w/ opposite orientation")
	}

	// Connecting vertically
	w.AddTrack(Point{0, 1}, Vertical)
	if w.pmap.Get(Point{0, 1})[0].(*Piece).getNext(Down).GetData().(*Track).orient != Vertical {
		t.Fatalf("New track merged with wrong track")
	}

	w.AddTrack(Point{-1, 0}, Vertical)
	w.AddTrack(Point{-1, 0}, Horizontal)
	if w.pmap.Get(Point{-1, 0})[1].(*Piece).getNext(Right).GetData().(*Track).orient != Horizontal {
		t.Fatalf("New track merged with wrong track")
	}
}

func TestAddRemovePieces(t *testing.T) {
	w := new(World)
	w.Init()

	w.AddNode(Point{0, 0})
	w.AddTrack(Point{2, 0}, Horizontal)
	w.AddTrack(Point{0, 2}, Vertical)

	n := w.getNext(Point{0, 0}, 0)
	tr1 := w.getNext(Point{2, 0}, 0)
	tr2 := w.getNext(Point{0, 2}, 0)

	if tr1.getNext(Left) != n || tr2.getNext(Down) != n {
		t.Fatalf("Tracks didn't bind to node")
	}

	w.Delete(Point{0, 0}, 0)
	if tr1.getNext(Left) != nil || tr2.getNext(Down) != nil {
		t.Fatalf("Tracks didn't unbind to node")
	}

	if anyItemsInNodeTerritory(w, Point{0, 0}) || len(w.pmap.pieces) > 2 {
		t.Fatalf("Node and NodeBody's weren't all removed from pointmap")
	}
}

func anyItemsInNodeTerritory(w *World, point Point) bool {
	for _, t := range nodeTerritory(point) {
		if len(w.pmap.Get(t)) > 0 {
			return true
		}
	}
	return false
}
