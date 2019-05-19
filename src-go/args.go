package course

type PieceArg interface {
	isPieceArg()
}

type ItemArg interface {
	isItemArg()
}

/*
Pieces
*/
type FullNodeArg struct{}

type HalfNodeArg struct {
	dir Direction
}

type CurveNodeArg struct {
	quad Quadrant
}

type TrackArg struct {
	orient Orientation
}

func (a FullNodeArg) isPieceArg()  {}
func (a HalfNodeArg) isPieceArg()  {}
func (a CurveNodeArg) isPieceArg() {}
func (a TrackArg) isPieceArg()     {}

/*
Items
*/
// Enemies
type EnemyArg interface {
	isEnemyArg()
}

type ZapperArg struct {
	next *ZapperPart
}

type ChristineArg struct {
	next *ChristinePart
}

func (e ZapperArg) isEnemyArg()    {}
func (e ChristineArg) isEnemyArg() {}

// PowerUps
type PowerUpArg interface {
	isPowerUpArg()
}

type (
	PhaserArg     struct{}
	TimeSlowerArg struct{}
)

func (a PhaserArg) isPowerUpArg()     {}
func (a TimeSlowerArg) isPowerUpArg() {}

// TrackItems
type TrackItemArg interface {
	isTrackItemArg()
}

type PortalArg struct {
	dest *Portal
}
type TrackSwitchArg struct {
	sig Signal
}
type SludgeArg struct{}

func (a PortalArg) isTrackItemArg()      {}
func (a TrackSwitchArg) isTrackItemArg() {}
func (a SludgeArg) isTrackItemArg()      {}

func (e ZapperArg) isItemArg()      {}
func (e ChristineArg) isItemArg()   {}
func (a PhaserArg) isItemArg()      {}
func (a TimeSlowerArg) isItemArg()  {}
func (a PortalArg) isItemArg()      {}
func (a TrackSwitchArg) isItemArg() {}
func (a SludgeArg) isItemArg()      {}
