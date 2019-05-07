package course

import "testing"

func TestNodeChains(t *testing.T) {
	n := NewNodePiece()
	n.AddTrack(Left)

	tr := n.GetPiece(Left)
	tr.AddNode(Right) // won't do anything, something's already there
	if tr.GetPiece(Right) != n {
		t.Fatalf("Improperly overwrote a node connection.")
	}

	tr.AddNode(Left)
	if tr.GetPiece(Left).GetPiece(Right) != tr {
		t.Fatalf("Improperly added connection to the new node")
	}

	tr.RemovePiece(Left)
	if tr.GetPiece(Left) != nil {
		t.Fatalf("Improper unlinking between track and node to the left")
	}
}

func TestDelete1(t *testing.T) {
	n := NewNodePiece()
	n.AddTrack(Up)

	tr1 := n.GetPiece(Up)
	tr1.AddTrack(Up)

	tr2 := tr1.GetPiece(Up)
	tr2.AddTrack(Up)

	tr3 := tr2.GetPiece(Up)
	tr3.AddNode(Up)

	tr2.Delete() // should leave a gap between tr1 and tr3
	if tr2.GetPiece(Up) != nil || tr2.GetPiece(Down) != nil || tr1.GetPiece(Up) != nil ||
		tr3.GetPiece(Down) != nil {
		t.Fatalf("Intermediate node was not properly deleted")
	}

	tr1.DeleteNeighbor(Down) // should leave n with only nil in nexts
	if n.GetPiece(Up) != nil || tr1.GetPiece(Down) != nil {
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

	n1 := tr1.GetPiece(Left)
	n1.AddTrack(Down)

	tr2 := n1.GetPiece(Down)
	tr2.AddNode(Down)

	n2 := tr2.GetPiece(Down)
	n2.AddTrack(Right)

	tr3 := n2.GetPiece(Right)
	tr3.AddNode(Right)

	n3 := tr3.GetPiece(Right)
	n3.AddTrack(Up)

	t4 := n3.GetPiece(Up)
}

func TestLink2(t *testing.T) {
	// Create this:	ntttn
	// Delete the middle track. Add a new track in its place.
	// Link it to the track on the other side
}

func TestInvalidInput(t *testing.T) {
	// try adding stuff to track in wrong direction
	// try merging two tracks going the wrong ways
}
