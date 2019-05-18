package course

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

// implemented by *ZapperPart, *ChristinePart, etc.
type EnemyPart interface {
	getNext() EnemyPart
	setNext(EnemyPart) bool
	Object
}

type ZapperPart struct {
	next *ZapperPart
}

type ChristinePart struct {
	next *ChristinePart
}

// gets the *ZapperPart that this one follows
func (s *ZapperPart) getNext() EnemyPart {
	return s.next
}

func (e *ZapperPart) setNext(p EnemyPart) bool {
	s2, ok := p.(*ZapperPart)
	if ok {
		e.next = s2
	}
	return ok
}

func (e *ChristinePart) getNext() EnemyPart {
	return e.next
}

func (e *ChristinePart) setNext(p EnemyPart) bool {
	e2, ok := p.(*ChristinePart)
	if ok {
		e.next = e2
	}
	return ok
}

// Constructor
func NewEnemy(arg EnemyArg) EnemyPart {
	switch arg.(type) {
	case ZapperArg:
		return &ZapperPart{arg.(ZapperArg).next}
	case ChristineArg:
		return &ChristinePart{arg.(ChristineArg).next}
	}
	return nil
}
