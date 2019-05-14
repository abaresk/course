/*
Supplementary to World and Pointmap
Pointmap stores Items, which are either *Piece (see: pieces.go) or *NodeBody

*NodeBody is essentially a dummy piece that refers back to the Node it belongs to
*/

package course

// *Piece and *Nodebody implement this interface
type Item interface {
	isItem()
}

type NodeBody struct {
	center *Piece
}

func (n *NodeBody) isItem() {}
func (p *Piece) isItem()    {}

func extractPiece(m Item) (*Piece, bool) {
	if piece, isPiece := m.(*Piece); isPiece {
		return piece, true
	}
	return nil, false
}

func extractTrack(m Item) (*Track, bool) {
	if piece, isPiece := extractPiece(m); isPiece {
		if track, isTrack := piece.GetData().(*Track); isTrack {
			return track, true
		}
	}
	return nil, false
}

func extractNode(m Item) (*Node, bool) {
	if piece, isPiece := extractPiece(m); isPiece {
		if node, isNode := piece.GetData().(*Node); isNode {
			return node, true
		}
	}
	return nil, false
}
