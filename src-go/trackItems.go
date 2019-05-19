package course

type Signal int

const (
	SignalRed Signal = iota
	SignalBlue
	SignalGreen
	SignalYellow
)

type TrackItem interface {
	isTrackItem()
	Object
}

type Portal struct {
	dest *Portal
}

type TrackSwitch struct {
	sig Signal
}

type Sludge struct{}

func (i *Portal) isTrackItem()      {}
func (i *TrackSwitch) isTrackItem() {}
func (i *Sludge) isTrackItem()      {}

// Constructor
func NewTrackItem(arg TrackItemArg) TrackItem {
	switch arg.(type) {
	case PortalArg:
		return &Portal{arg.(PortalArg).dest}
	case TrackSwitchArg:
		return &TrackSwitch{arg.(TrackSwitchArg).sig}
	case SludgeArg:
		return &Sludge{}
	}
	return nil
}
