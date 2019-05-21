/*
Supplementary to World and Pointmap
Pointmap stores Objects, which are either ObjectPieces (Piece and *NodeBody),
ObjectEnemy's or ObjectItem's.

*NodeBody is essentially a dummy piece that refers back to the Node it belongs to
*/

package course

type ObjectType int

// Listed from lowest priority value to highest
const (
	ObjectPiece ObjectType = iota
	ObjectTrackItem
	ObjectEnemy
	ObjectPowerUp
)

func getObjectType(ob Object) ObjectType {
	switch ob.(type) {
	case Piece, *NodeBody:
		return ObjectPiece
	case TrackItem:
		return ObjectTrackItem
	case EnemyPart:
		return ObjectEnemy
	case PowerUp:
		return ObjectPowerUp
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

func (n *fullNode) isObject()  {}
func (n *halfNode) isObject()  {}
func (n *CurveNode) isObject() {}
func (t *Track) isObject()     {}
func (n *NodeBody) isObject()  {}

func (i *Portal) isObject()      {}
func (i *TrackSwitch) isObject() {}
func (i *Sludge) isObject()      {}

func (e *ZapperPart) isObject()    {}
func (e *ChristinePart) isObject() {}

func (i *Phaser) isObject()     {}
func (i *TimeSlower) isObject() {}
