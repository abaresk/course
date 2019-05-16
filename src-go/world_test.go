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

	w.AddNode(Point{0, 0}, FullNodeArg{})
	w.AddTrack(Point{2, 0}, Horizontal)
	w.AddNode(Point{4, 0}, FullNodeArg{})
	w.AddTrack(Point{4, 2}, Vertical)
	w.AddNode(Point{4, 4}, FullNodeArg{})
	w.AddTrack(Point{2, 4}, Horizontal)
	w.AddNode(Point{0, 4}, FullNodeArg{})
	w.AddTrack(Point{0, 2}, Vertical)

	tr1, _ := w.Get(Point{0, 0}, 0).(Piece).getNext(Up)
	tr2, _ := w.Get(Point{0, 4}, 0).(Piece).getNext(Down)
	if tr1 != tr2 {
		t.Fatalf("Improper merging at end of loop")
	}

	n1, _ := w.Get(Point{0, 2}, 0).(Piece).getNext(Down)
	n2, _ := w.Get(Point{2, 0}, 0).(Piece).getNext(Left)
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

	w.AddNode(Point{0, 0}, FullNodeArg{}) // Works!
	w.AddTrack(Point{1, 0}, Horizontal)   // invalid
	w.AddNode(Point{2, 0}, FullNodeArg{}) // invalid
	w.AddTrack(Point{2, 1}, Vertical)     // Works!
	w.AddNode(Point{2, 2}, FullNodeArg{}) // invalid
	w.AddTrack(Point{1, 2}, Horizontal)   // Works!
	w.AddNode(Point{0, 2}, FullNodeArg{}) // invalid
	w.AddTrack(Point{1, 1}, Vertical)     // invalid

	if len(w.pmap.pieces) != 11 {
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

	w.AddNode(Point{0, 0}, FullNodeArg{}) // 4
	w.AddNode(Point{6, 3}, FullNodeArg{}) // 3
	w.AddNode(Point{3, 0}, FullNodeArg{}) // 5
	w.AddNode(Point{3, 3}, FullNodeArg{}) // 2
	w.AddNode(Point{0, 3}, FullNodeArg{}) // 1
	w.AddNode(Point{6, 0}, FullNodeArg{}) // 6

	// try to get at 5 from each side
	n4 := w.Get(Point{0, 0}, 0).(Piece)
	n3 := w.Get(Point{6, 3}, 0).(Piece)

	s1, _ := n4.getNext(Up)
	s2, _ := s1.getNext(Right)
	s3, _ := s2.getNext(Down)

	u1, _ := n3.getNext(Down)
	u2, _ := u1.getNext(Left)
	if s3 != u2 {
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

	w.AddNode(Point{0, 0}, FullNodeArg{})
	w.AddNode(Point{3, 0}, FullNodeArg{})
	w.AddNode(Point{1, -3}, FullNodeArg{})

	n1 := w.Get(Point{0, 0}, 0).(Piece)

	s1, _ := n1.getNext(Right)
	s2, _ := s1.getNext(Left)

	if n1 != s2 {
		t.Fatalf("Nodes should have connected")
	}
	for _, n := range w.Get(Point{1, -3}, 0).(Piece).(*FullNode).nexts {
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

	if len(w.pmap.get(Point{0, 0})) != 2 {
		t.Fatalf("Should be able to overlap 2 tracks w/ opposite orientation")
	}

	// Connecting vertically
	w.AddTrack(Point{0, 1}, Vertical)
	s1, _ := w.Get(Point{0, 1}, 0).(Piece).getNext(Down)

	if s1.(*Track).orient != Vertical {
		t.Fatalf("New track merged with wrong track")
	}

	w.AddTrack(Point{-1, 0}, Vertical)
	w.AddTrack(Point{-1, 0}, Horizontal)
	s2, _ := w.Get(Point{-1, 0}, 1).(Piece).getNext(Right)
	if s2.(*Track).orient != Horizontal {
		t.Fatalf("New track merged with wrong track")
	}
}

func TestAddRemovePieces(t *testing.T) {
	w := new(World)
	w.Init()

	w.AddNode(Point{0, 0}, FullNodeArg{})
	w.AddTrack(Point{2, 0}, Horizontal)
	w.AddTrack(Point{0, 2}, Vertical)

	n := w.Get(Point{0, 0}, 0)
	tr1 := w.Get(Point{2, 0}, 0).(Piece)
	tr2 := w.Get(Point{0, 2}, 0).(Piece)

	s1, _ := tr1.getNext(Left)
	s2, _ := tr2.getNext(Down)

	if s1 != n || s2 != n {
		t.Fatalf("Tracks didn't bind to node")
	}

	w.Delete(Point{0, 0}, 0)
	s3, _ := tr1.getNext(Left)
	s4, _ := tr2.getNext(Down)
	if s3 != nil || s4 != nil {
		t.Fatalf("Tracks didn't unbind to node")
	}

	if anyItemsInNodeTerritory(w, Point{0, 0}, &FullNode{}) || len(w.pmap.pieces) > 2 {
		t.Fatalf("Node and NodeBody's weren't all removed from pointmap")
	}
}

func anyItemsInNodeTerritory(w *World, point Point, n Node) bool {
	for _, t := range n.nodeTerritory(point) {
		if len(w.pmap.get(t)) > 0 {
			return true
		}
	}
	return false
}
