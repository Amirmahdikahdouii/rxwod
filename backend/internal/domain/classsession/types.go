package classsession

type ClassSessionID string

type Capacity int

func (id ClassSessionID) String() string {
	return string(id)
}

func validateCapacity(capacity Capacity) error {
	if capacity <= 0 {
		return ErrInvalidCapacity
	}
	return nil
}
