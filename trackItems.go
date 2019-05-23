package course

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

func NewPortal(dest *Portal) *Portal {
	return &Portal{dest}
}

func NewTrackSwitch(sig Signal) *TrackSwitch {
	return &TrackSwitch{sig}
}

func NewSludge() *Sludge { return &Sludge{} }
