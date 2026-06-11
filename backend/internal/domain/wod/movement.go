package wod

type Movement struct {
	id        MovementID
	position  int
	name      string
	reps      *RepCount
	loadValue *LoadValue
	loadUnit  *LoadUnit
	notes     string
}

func NewMovement(
	id MovementID,
	position int,
	name string,
	reps *RepCount,
	loadValue *LoadValue,
	loadUnit *LoadUnit,
	notes string,
) (Movement, error) {
	m := Movement{
		id:        id,
		position:  position,
		name:      name,
		reps:      reps,
		loadValue: loadValue,
		loadUnit:  loadUnit,
		notes:     notes,
	}
	return m, m.Validate()
}

func (m Movement) ID() MovementID {
	return m.id
}

func (m Movement) Position() int {
	return m.position
}

func (m Movement) Name() string {
	return m.name
}

func (m Movement) Reps() *RepCount {
	if m.reps == nil {
		return nil
	}
	value := *m.reps
	return &value
}

func (m Movement) LoadValue() *LoadValue {
	if m.loadValue == nil {
		return nil
	}
	value := *m.loadValue
	return &value
}

func (m Movement) LoadUnit() *LoadUnit {
	if m.loadUnit == nil {
		return nil
	}
	value := *m.loadUnit
	return &value
}

func (m Movement) Notes() string {
	return m.notes
}

func (m Movement) Validate() error {
	if m.position <= 0 {
		return ErrInvalidPosition
	}
	if len(m.name) == 0 {
		return ErrInvalidMovement
	}
	if m.reps != nil && *m.reps <= 0 {
		return ErrInvalidReps
	}
	if m.loadValue != nil && m.loadUnit == nil {
		return ErrInvalidMovement
	}
	if m.loadUnit != nil {
		if err := validateLoadUnit(*m.loadUnit); err != nil {
			return err
		}
	}
	return nil
}

func cloneMovements(movements []Movement) []Movement {
	cloned := make([]Movement, len(movements))
	copy(cloned, movements)
	return cloned
}
