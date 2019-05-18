package course

type PowerUpArg interface {
	isPowerUpArg()
}

type (
	PhaserArg     struct{}
	TimeSlowerArg struct{}
)

func (a PhaserArg) isPowerUpArg()     {}
func (a TimeSlowerArg) isPowerUpArg() {}

type PowerUp interface {
	isPowerUp()
	Object
}

type Phaser struct{}

type TimeSlower struct{}

func (i *Phaser) isPowerUp()     {}
func (i *TimeSlower) isPowerUp() {}

// Constructor
func NewPowerUp(arg PowerUpArg) PowerUp {
	switch arg.(type) {
	case PhaserArg:
		return &Phaser{}
	case TimeSlowerArg:
		return &TimeSlower{}
	}
	return nil
}
