package wod

import "time"

type Variant interface {
	wodVariant()
	ID() WODID
	Name() WODName
	Description() WODDescription
	Type() WODType
	Status() WODStatus
	ScoringKind() ScoringKind
	Movements() []Movement
	CreatedAt() time.Time
	UpdatedAt() time.Time
}

type AMRAPWOD = WOD[AMRAPConfig]
type ForTimeWOD = WOD[ForTimeConfig]
type TabataWOD = WOD[TabataConfig]
type EMOMWOD = WOD[EMOMConfig]

type SavedAMRAP struct{ wod AMRAPWOD }
type SavedForTime struct{ wod ForTimeWOD }
type SavedTabata struct{ wod TabataWOD }
type SavedEMOM struct{ wod EMOMWOD }

func NewSavedAMRAP(wod AMRAPWOD) SavedAMRAP { return SavedAMRAP{wod: wod} }
func NewSavedForTime(wod ForTimeWOD) SavedForTime { return SavedForTime{wod: wod} }
func NewSavedTabata(wod TabataWOD) SavedTabata { return SavedTabata{wod: wod} }
func NewSavedEMOM(wod EMOMWOD) SavedEMOM { return SavedEMOM{wod: wod} }

func (SavedAMRAP) wodVariant() {}
func (SavedForTime) wodVariant() {}
func (SavedTabata) wodVariant() {}
func (SavedEMOM) wodVariant() {}

func (s SavedAMRAP) ID() WODID { return s.wod.ID() }
func (s SavedForTime) ID() WODID { return s.wod.ID() }
func (s SavedTabata) ID() WODID { return s.wod.ID() }
func (s SavedEMOM) ID() WODID { return s.wod.ID() }

func (s SavedAMRAP) Name() WODName { return s.wod.Name() }
func (s SavedForTime) Name() WODName { return s.wod.Name() }
func (s SavedTabata) Name() WODName { return s.wod.Name() }
func (s SavedEMOM) Name() WODName { return s.wod.Name() }

func (s SavedAMRAP) Description() WODDescription { return s.wod.Description() }
func (s SavedForTime) Description() WODDescription { return s.wod.Description() }
func (s SavedTabata) Description() WODDescription { return s.wod.Description() }
func (s SavedEMOM) Description() WODDescription { return s.wod.Description() }

func (s SavedAMRAP) Type() WODType { return WODTypeAMRAP }
func (s SavedForTime) Type() WODType { return WODTypeForTime }
func (s SavedTabata) Type() WODType { return WODTypeTabata }
func (s SavedEMOM) Type() WODType { return WODTypeEMOM }

func (s SavedAMRAP) Status() WODStatus { return s.wod.Status() }
func (s SavedForTime) Status() WODStatus { return s.wod.Status() }
func (s SavedTabata) Status() WODStatus { return s.wod.Status() }
func (s SavedEMOM) Status() WODStatus { return s.wod.Status() }

func (s SavedAMRAP) ScoringKind() ScoringKind { return s.wod.Scoring().Kind() }
func (s SavedForTime) ScoringKind() ScoringKind { return s.wod.Scoring().Kind() }
func (s SavedTabata) ScoringKind() ScoringKind { return s.wod.Scoring().Kind() }
func (s SavedEMOM) ScoringKind() ScoringKind { return s.wod.Scoring().Kind() }

func (s SavedAMRAP) Movements() []Movement { return s.wod.Movements() }
func (s SavedForTime) Movements() []Movement { return s.wod.Movements() }
func (s SavedTabata) Movements() []Movement { return s.wod.Movements() }
func (s SavedEMOM) Movements() []Movement { return s.wod.Movements() }

func (s SavedAMRAP) CreatedAt() time.Time { return s.wod.CreatedAt() }
func (s SavedForTime) CreatedAt() time.Time { return s.wod.CreatedAt() }
func (s SavedTabata) CreatedAt() time.Time { return s.wod.CreatedAt() }
func (s SavedEMOM) CreatedAt() time.Time { return s.wod.CreatedAt() }

func (s SavedAMRAP) UpdatedAt() time.Time { return s.wod.UpdatedAt() }
func (s SavedForTime) UpdatedAt() time.Time { return s.wod.UpdatedAt() }
func (s SavedTabata) UpdatedAt() time.Time { return s.wod.UpdatedAt() }
func (s SavedEMOM) UpdatedAt() time.Time { return s.wod.UpdatedAt() }

func (s SavedAMRAP) WOD() AMRAPWOD { return s.wod }
func (s SavedForTime) WOD() ForTimeWOD { return s.wod }
func (s SavedTabata) WOD() TabataWOD { return s.wod }
func (s SavedEMOM) WOD() EMOMWOD { return s.wod }

func MatchVariant[T any](variant Variant, handlers map[WODType]func(Variant) (T, error)) (T, error) {
	var zero T
	handler, ok := handlers[variant.Type()]
	if !ok {
		return zero, ErrUnknownWODType
	}
	return handler(variant)
}
