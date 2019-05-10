/*
Supplementary to World and Pointmap
Pointmap stores mapPieces, which are either *Piece (see: pieces.go) or *NodeBody

*NodeBody is essentially a dummy piece that refers back to the Node it belongs to
*/

package course

// *Piece and *Nodebody implement this interface
type mapPiece interface {
	isMapPiece()
}

type NodeBody struct {
	center *Piece
}

func (n *NodeBody) isMapPiece() {}
func (p *Piece) isMapPiece()    {}

func extractPiece(m mapPiece) (*Piece, bool) {
	if piece, isPiece := m.(*Piece); isPiece {
		return piece, true
	}
	return nil, false
}

func extractTrack(m mapPiece) (*Track, bool) {
	if piece, isPiece := extractPiece(m); isPiece {
		if track, isTrack := piece.GetData().(*Track); isTrack {
			return track, true
		}
	}
	return nil, false
}

func extractNode(m mapPiece) (*Node, bool) {
	if piece, isPiece := extractPiece(m); isPiece {
		if node, isNode := piece.GetData().(*Node); isNode {
			return node, true
		}
	}
	return nil, false
}
