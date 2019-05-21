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

	w.AddPiece(Point{0, 0}, NewFullNodeClear())
	w.AddPiece(Point{2, 0}, NewTrack(Horizontal))
	w.AddPiece(Point{4, 0}, NewFullNodeClear())
	w.AddPiece(Point{4, 2}, NewTrack(Vertical))
	w.AddPiece(Point{4, 4}, NewFullNodeClear())
	w.AddPiece(Point{2, 4}, NewTrack(Horizontal))
	w.AddPiece(Point{0, 4}, NewFullNodeClear())
	w.AddPiece(Point{0, 2}, NewTrack(Vertical))

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

	w.AddPiece(Point{0, 0}, NewFullNodeClear())   // Works!
	w.AddPiece(Point{1, 0}, NewTrack(Horizontal)) // invalid
	w.AddPiece(Point{2, 0}, NewFullNodeClear())   // invalid
	w.AddPiece(Point{2, 1}, NewTrack(Vertical))   // Works!
	w.AddPiece(Point{2, 2}, NewFullNodeClear())   // invalid
	w.AddPiece(Point{1, 2}, NewTrack(Horizontal)) // Works!
	w.AddPiece(Point{0, 2}, NewFullNodeClear())   // invalid
	w.AddPiece(Point{1, 1}, NewTrack(Vertical))   // invalid

	if len(w.pmap.objects) != 11 {
		t.Fatalf("Erroneous objects should be ignored")
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

	w.AddPiece(Point{0, 0}, NewFullNodeClear()) // 4
	w.AddPiece(Point{6, 3}, NewFullNodeClear()) // 3
	w.AddPiece(Point{3, 0}, NewFullNodeClear()) // 5
	w.AddPiece(Point{3, 3}, NewFullNodeClear()) // 2
	w.AddPiece(Point{0, 3}, NewFullNodeClear()) // 1
	w.AddPiece(Point{6, 0}, NewFullNodeClear()) // 6

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

	w.AddPiece(Point{0, 0}, NewFullNodeClear())
	w.AddPiece(Point{3, 0}, NewFullNodeClear())
	w.AddPiece(Point{1, -3}, NewFullNodeClear())

	n1 := w.Get(Point{0, 0}, 0).(Piece)

	s1, _ := n1.getNext(Right)
	s2, _ := s1.getNext(Left)

	if n1 != s2 {
		t.Fatalf("Nodes should have connected")
	}
	for _, n := range w.Get(Point{1, -3}, 0).(FullNodeClear).nexts {
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

	w.AddPiece(Point{0, 0}, NewTrack(Horizontal))
	w.AddPiece(Point{0, 0}, NewTrack(Vertical))

	if len(*w.pmap.get(Point{0, 0})) != 2 {
		t.Fatalf("Should be able to overlap 2 tracks w/ opposite orientation")
	}

	// Connecting vertically
	w.AddPiece(Point{0, 1}, NewTrack(Vertical))
	s1, _ := w.Get(Point{0, 1}, 0).(Piece).getNext(Down)

	if s1.(*Track).orient != Vertical {
		t.Fatalf("New track merged with wrong track")
	}

	w.AddPiece(Point{-1, 0}, NewTrack(Vertical))
	w.AddPiece(Point{-1, 0}, NewTrack(Horizontal))
	s2, _ := w.Get(Point{-1, 0}, 1).(Piece).getNext(Right)
	if s2.(*Track).orient != Horizontal {
		t.Fatalf("New track merged with wrong track")
	}
}

func TestAddRemovePieces(t *testing.T) {
	w := new(World)
	w.Init()

	w.AddPiece(Point{0, 0}, NewFullNodeClear())
	w.AddPiece(Point{2, 0}, NewTrack(Horizontal))
	w.AddPiece(Point{0, 2}, NewTrack(Vertical))

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

	if anyObjectsInNodeTerritory(w, Point{0, 0}, &fullNode{}) || len(w.pmap.objects) > 2 {
		t.Fatalf("Node and NodeBody's weren't all removed from pointmap")
	}
}

/*
Make 3/4 of circle. Add the 4th curve, but facing wrong way (shouldn't connect w/ neighbors).
Try to add new one on top of it (shouldn't work).
Remove old 4th curve and add the proper one.
*/
func TestCurveCircle(t *testing.T) {
	w := new(World)
	w.Init()

	// Make a circle
	w.AddPiece(Point{0, 0}, NewCurveNode(Two))
	w.AddPiece(Point{-3, 0}, NewCurveNode(One))
	w.AddPiece(Point{0, 3}, NewCurveNode(Three))

	s1 := w.Get(Point{-3, 0}, 0).(Piece)
	s2 := w.Get(Point{0, 0}, 0).(Piece)
	s3 := w.Get(Point{0, 3}, 0).(Piece)

	if u, _ := s1.getNext(Right); u != s2 {
		t.Fatalf("Curves didn't connect properly")
	}

	if u, _ := s2.getNext(Up); u != s3 {
		t.Fatalf("Curves didn't connect properly")
	}

	if u, _ := s2.getNext(Left); u != s1 {
		t.Fatalf("Curves didn't connect properly")
	}

	w.AddPiece(Point{-3, 3}, NewCurveNode(Two)) // invalid; it should not merge with either
	s4 := w.Get(Point{-3, 3}, 0).(Piece)
	if u, _ := s4.getNext(Up); u != nil {
		t.Fatalf("This curve should not have connected to anything")
	}

	w.AddPiece(Point{-3, 3}, NewCurveNode(Four)) // invalid; should not be placed (something's already there)
	if len(*w.pmap.get(Point{-3, 3})) != 1 {
		t.Fatalf("The above node should not have been allowed to be placed")
	}

	w.Delete(Point{-3, 3}, 0)
	w.AddPiece(Point{-3, 3}, NewCurveNode(Four)) // works!
	s5 := w.Get(Point{-3, 3}, 0).(Piece)
	if u, _ := s5.getNext(Down); u != s1 {
		t.Fatalf("Curves didn't connect properly")
	}

	if u, _ := s5.getNext(Right); u != s3 {
		t.Fatalf("Curves didn't connect properly")
	}

	if u, _ := s1.getNext(Up); u != s5 {
		t.Fatalf("Curves didn't connect properly")
	}
}

/*
Create the following:
	   t
	   n
	mmnnnmm
	m     m
	n     n
	nntttnn

Then delete the halfNode
*/
func TestHalfNodeLoop(t *testing.T) {
	w := new(World)
	w.Init()

	w.AddPiece(Point{3, 5}, NewTrack(Vertical))
	w.AddPiece(Point{3, 3}, NewHalfNodeSet(Up, true))
	w.AddPiece(Point{0, 3}, NewCurveNode(Four))
	w.AddPiece(Point{0, 0}, NewCurveNode(One))
	w.AddPiece(Point{2, 0}, NewTrack(Horizontal))
	w.AddPiece(Point{3, 0}, NewTrack(Horizontal))
	w.AddPiece(Point{4, 0}, NewTrack(Horizontal))
	w.AddPiece(Point{6, 0}, NewCurveNode(Two))
	w.AddPiece(Point{6, 3}, NewCurveNode(Three))

	s1 := w.Get(Point{3, 3}, 0).(Piece)
	s2 := w.Get(Point{3, 5}, 0).(Piece)
	s3 := w.Get(Point{0, 3}, 0).(Piece)
	s4 := w.Get(Point{6, 3}, 0).(Piece)
	if u, _ := s1.getNext(Up); u != s2 {
		t.Fatalf("halfNode didn't connect properly")
	}
	if u, _ := s2.getNext(Down); u != s1 {
		t.Fatalf("halfNode didn't connect properly")
	}
	if u, _ := s3.getNext(Right); u != s1 {
		t.Fatalf("halfNode didn't connect properly")
	}
	if u, _ := s1.getNext(Left); u != s3 {
		t.Fatalf("halfNode didn't connect properly")
	}
	if u, _ := s1.getNext(Right); u != s4 {
		t.Fatalf("halfNode didn't connect properly")
	}
	if u, _ := s4.getNext(Left); u != s1 {
		t.Fatalf("halfNode didn't connect properly")
	}

	// Now remove halfNode
	w.Delete(Point{3, 3}, 0)
	if u, _ := s2.getNext(Down); u != nil {
		t.Fatalf("halfNode didn't disconnect properly")
	}
	if len(w.pmap.objects) != 20 {
		t.Fatalf("halfNode wasn't removed properly")
	}
}

/*
First create some tracks, then add 1 track switch and 2 linked portals.
*/
func TestAddObjects(t *testing.T) {
	w := new(World)
	w.Init()

	w.AddPiece(Point{0, 0}, NewHalfNodeSet(Right, true))
	w.AddPiece(Point{2, 0}, NewTrack(Horizontal))
	w.AddPiece(Point{3, 0}, NewTrack(Horizontal))
	w.AddPiece(Point{4, 0}, NewTrack(Horizontal))
	w.AddPiece(Point{0, 1}, NewTrack(Vertical))
	w.AddPiece(Point{0, -1}, NewTrack(Vertical))

	w.AddObject(Point{5, 0}, NewTrackSwitch(SignalBlue)) // invalid, do nothing
	w.AddObject(Point{4, 0}, NewTrackSwitch(SignalBlue)) // valid

	w.AddObject(Point{1, 0}, NewPhaser()) // invalid, do nothing
	w.AddObject(Point{2, 0}, NewPhaser())

	// Portal -- needs a way to add linked items

	// Create AddLinkedItem(p1, p2 Point, arg)
}

/*
First create some tracks, then add 1 Christine car and a Zapper.
Then remove one directly, and remove another by deleting the track.
*/
func TestAddEnemies(t *testing.T) {
	w := new(World)
	w.Init()

	w.AddPiece(Point{0, 0}, NewFullNodeClear())
	w.AddPiece(Point{2, 0}, NewTrack(Horizontal))
	w.AddPiece(Point{3, 0}, NewTrack(Horizontal))
	w.AddPiece(Point{4, 0}, NewTrack(Horizontal))

	w.AddObject(Point{0, 0}, NewChristine(nil))
	w.AddObject(Point{4, 0}, NewZapper(nil))

	if _, ok := w.Get(Point{0, 0}, 1).(*ChristinePart); !ok {
		t.Fatalf("Christine car not properly added")
	}

	if _, ok := w.Get(Point{4, 0}, 1).(*ZapperPart); !ok {
		t.Fatalf("Zapper not properly added")
	}

	// Remove Items
	w.Delete(Point{0, 0}, 1)
	if l := w.pmap.points[Point{0, 0}]; l.Len() != 1 {
		t.Fatalf("Christine car not properly deleted")
	}

	w.Delete(Point{4, 0}, 0)
	if l := w.pmap.points[Point{4, 0}]; l.Len() != 0 {
		t.Fatalf("Deleted Track should have deleted any items on top")
	}
}

func TestLongEnemies(t *testing.T) {

}

func anyObjectsInNodeTerritory(w *World, point Point, n Node) bool {
	for _, t := range n.nodeTerritory(point) {
		if len(*w.pmap.get(t)) > 0 {
			return true
		}
	}
	return false
}
