package course

type TrackItemArg interface {
	isTrackItemArg()
}

type PortalArg struct {
	dest *Portal
}
type TrackSwitchArg struct {
	pair Node
}
type SludgeArg struct{}

func (a PortalArg) isTrackItemArg()      {}
func (a TrackSwitchArg) isTrackItemArg() {}
func (a SludgeArg) isTrackItemArg()      {}

type TrackItem interface {
	isTrackItem()
	Object
}

type Portal struct {
	dest *Portal
}

type TrackSwitch struct {
	pair Node
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
		return &TrackSwitch{arg.(TrackSwitchArg).pair}
	case SludgeArg:
		return &Sludge{}
	}
	return nil
}
