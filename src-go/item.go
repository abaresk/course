/*
Supplementary to World and Pointmap
Pointmap stores Items, which are either *Piece (see: pieces.go) or *NodeBody

*NodeBody is essentially a dummy piece that refers back to the Node it belongs to
*/

package course

// Piece and *Nodebody implement this interface
type Item interface {
	isItem()
}

type NodeBody struct {
	center Piece
}

func (n *FullNode) isItem()  {}
func (n *HalfNode) isItem()  {}
func (n *CurveNode) isItem() {}
func (t *Track) isItem()     {}
func (n *NodeBody) isItem()  {}
