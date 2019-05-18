package course

// implemented by *ShockerPart,
type EnemyPart interface {
	getNext() EnemyPart
	setNext(EnemyPart) bool
	Object
}

type ShockerPart struct {
	next *ShockerPart
}

// gets the *ShockerPart that this one follows
func (s *ShockerPart) getNext() EnemyPart {
	return s.next
}

func (s *ShockerPart) setNext(p EnemyPart) bool {
	s2, ok := p.(*ShockerPart)
	if ok {
		s.next = s2
	}
	return ok
}

func (s *ShockerPart) isObject() {}
