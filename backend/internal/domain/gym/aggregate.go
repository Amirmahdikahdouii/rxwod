package gym

import (
	"time"

	"github.com/rxwod/backend/internal/domain/user"
)

type Gym struct {
	id        GymID
	name      GymName
	ownerID   user.UserID
	createdAt time.Time
	updatedAt time.Time
}

func NewGym(id GymID, name GymName, ownerID user.UserID, now time.Time) (Gym, error) {
	if err := validateGymName(name); err != nil {
		return Gym{}, err
	}
	return Gym{
		id:        id,
		name:      name,
		ownerID:   ownerID,
		createdAt: now,
		updatedAt: now,
	}, nil
}

func ReconstructGym(id GymID, name GymName, ownerID user.UserID, createdAt time.Time, updatedAt time.Time) (Gym, error) {
	gym, err := NewGym(id, name, ownerID, createdAt)
	if err != nil {
		return Gym{}, err
	}
	gym.updatedAt = updatedAt
	return gym, nil
}

func (g Gym) ID() GymID {
	return g.id
}

func (g Gym) Name() GymName {
	return g.name
}

func (g Gym) OwnerID() user.UserID {
	return g.ownerID
}

func (g Gym) CreatedAt() time.Time {
	return g.createdAt
}

func (g Gym) UpdatedAt() time.Time {
	return g.updatedAt
}
