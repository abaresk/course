package course

import "testing"

func TestNodeChains(t *testing.T) {
	n := NewNodePiece()
	n.AddTrack(Left)

	tr := n.getNext(Left)
	tr.AddNode(Right) // won't do anything, something's already there
	if tr.getNext(Right) != n {
		t.Fatalf("Improperly overwrote a node connection.")
	}

	tr.AddNode(Left)
	if tr.getNext(Left).getNext(Right) != tr {
		t.Fatalf("Improperly added connection to the new node")
	}

	tr.RemovePiece(Left)
	if tr.getNext(Left) != nil {
		t.Fatalf("Improper unlinking between track and node to the left")
	}
}

func TestDelete1(t *testing.T) {
	n := NewNodePiece()
	n.AddTrack(Up)

	tr1 := n.getNext(Up)
	tr1.AddTrack(Up)

	tr2 := tr1.getNext(Up)
	tr2.AddTrack(Up)

	tr3 := tr2.getNext(Up)
	tr3.AddNode(Up)

	tr2.Delete() // should leave a gap between tr1 and tr3
	if tr2.getNext(Up) != nil || tr2.getNext(Down) != nil || tr1.getNext(Up) != nil ||
		tr3.getNext(Down) != nil {
		t.Fatalf("Intermediate node was not properly deleted")
	}

	tr1.DeleteNeighbor(Down) // should leave n with only nil in nexts
	if n.getNext(Up) != nil || tr1.getNext(Down) != nil {
		t.Fatalf("Endpoint should point to nothing, and nothing should point to it")
	}
}

func TestLink1(t *testing.T) {
	// for this test, start with a track piece...
	// Create this:	nt
	// 				t t
	//				ntn
	// Add a node to complete the square. Then link it with the first track
	tr1 := NewTrackPiece(Horizontal)
	tr1.AddNode(Left)

	n1 := tr1.getNext(Left)
	n1.AddTrack(Down)

	tr2 := n1.getNext(Down)
	tr2.AddNode(Down)

	n2 := tr2.getNext(Down)
	n2.AddTrack(Right)

	tr3 := n2.getNext(Right)
	tr3.AddNode(Right)

	n3 := tr3.getNext(Right)
	n3.AddTrack(Up)

	tr4 := n3.getNext(Up)

	// Add the node and link it
	tr4.AddNode(Up)
	n4 := tr4.getNext(Up)
	n4.LinkPiece(tr1, Left)
	if tr1.getNext(Right) != n4 || n4.getNext(Left) != tr1 {
		t.Fatalf("Pieces didn't link properly")
	}
}

func TestLink2(t *testing.T) {
	// Create this:	ntttn
	// Delete the middle track. Add a new track in its place.
	// Link it to the track on the other side
	n1 := NewNodePiece()
	n1.AddTrack(Right)

	tr1 := n1.getNext(Right)
	tr1.AddTrack(Right)

	tr2 := tr1.getNext(Right)
	tr2.AddTrack(Right)

	tr3 := tr2.getNext(Right)
	tr3.AddNode(Right)

	// Delete middle track and add new track in its place
	tr2.Delete()
	tr3.AddTrack(Left)
	tr2New := tr3.getNext(Left)
	tr2New.LinkPiece(tr1, Left)
	if tr2.getNext(Left) != nil || tr2.getNext(Right) != nil {
		t.Fatalf("Deleted piece still references other pieces")
	}
	if tr1.getNext(Right) != tr2New || tr2New.getNext(Left) != tr1 ||
		tr3.getNext(Left) != tr2New || tr2New.getNext(Right) != tr3 {
		t.Fatalf("New piece didn't link properly with pieces around it")
	}
}

func TestInvalidInput(t *testing.T) {
	// try adding stuff to track in wrong direction
	// try merging two tracks going the wrong ways
}
