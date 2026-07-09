package wodresult

import (
	"time"

	"github.com/rxwod/backend/internal/domain/gym"
	domainwod "github.com/rxwod/backend/internal/domain/wod"
)

type WODResult struct {
	id              WODResultID
	wodID           domainwod.WODID
	gymMembershipID gym.MembershipID
	scoreValue      ScoreValue
	isRx            bool
	notes           Notes
	createdAt       time.Time
	updatedAt       time.Time
}

func NewWODResult(
	id WODResultID,
	wodID domainwod.WODID,
	gymMembershipID gym.MembershipID,
	scoreValue ScoreValue,
	isRx bool,
	notes Notes,
	now time.Time,
) (WODResult, error) {
	if err := validateScoreValue(scoreValue); err != nil {
		return WODResult{}, err
	}

	return WODResult{
		id:              id,
		wodID:           wodID,
		gymMembershipID: gymMembershipID,
		scoreValue:      scoreValue,
		isRx:            isRx,
		notes:           notes,
		createdAt:       now,
		updatedAt:       now,
	}, nil
}

func ReconstructWODResult(
	id WODResultID,
	wodID domainwod.WODID,
	gymMembershipID gym.MembershipID,
	scoreValue ScoreValue,
	isRx bool,
	notes Notes,
	createdAt time.Time,
	updatedAt time.Time,
) (WODResult, error) {
	result, err := NewWODResult(id, wodID, gymMembershipID, scoreValue, isRx, notes, createdAt)
	if err != nil {
		return WODResult{}, err
	}
	result.updatedAt = updatedAt
	return result, nil
}

func (r WODResult) ID() WODResultID {
	return r.id
}

func (r WODResult) WODID() domainwod.WODID {
	return r.wodID
}

func (r WODResult) GymMembershipID() gym.MembershipID {
	return r.gymMembershipID
}

func (r WODResult) ScoreValue() ScoreValue {
	return r.scoreValue
}

func (r WODResult) IsRx() bool {
	return r.isRx
}

func (r WODResult) Notes() Notes {
	return r.notes
}

func (r WODResult) CreatedAt() time.Time {
	return r.createdAt
}

func (r WODResult) UpdatedAt() time.Time {
	return r.updatedAt
}
