/*
Supplementary to World and Pointmap
Pointmap stores Objects, which are either *Piece (see: pieces.go) or *NodeBody

*NodeBody is essentially a dummy piece that refers back to the Node it belongs to
*/

package course

type ObjectType int

// Listed from lowest priority value to highest
const (
	ObjectPiece ObjectType = iota
	ObjectEnemy
	ObjectItem
)

func getObjectType(ob Object) ObjectType {
	switch ob.(type) {
	case Piece, *NodeBody:
		return ObjectPiece
	case EnemyPart:
		return ObjectEnemy
	case Item:
		return ObjectItem
	}
	return ObjectType(MaxInt)
}

// Piece and *Nodebody implement this interface
type Object interface {
	isObject()
}

type NodeBody struct {
	center Piece
}

func (n *FullNode) isObject()  {}
func (n *HalfNode) isObject()  {}
func (n *CurveNode) isObject() {}
func (t *Track) isObject()     {}
func (n *NodeBody) isObject()  {}
