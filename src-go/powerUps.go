package course

type PowerUp interface {
	isPowerUp()
	Object
}

type Phaser struct{}

type TimeSlower struct{}

func (i *Phaser) isPowerUp()     {}
func (i *TimeSlower) isPowerUp() {}

func NewPhaser() *Phaser         { return &Phaser{} }
func NewTimeSlower() *TimeSlower { return &TimeSlower{} }
