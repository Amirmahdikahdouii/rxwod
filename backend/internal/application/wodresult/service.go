package wodresult

import (
	"context"
	"fmt"
	"sort"

	appauthz "github.com/rxwod/backend/internal/application/authz"
	appgym "github.com/rxwod/backend/internal/application/gym"
	appwod "github.com/rxwod/backend/internal/application/wod"
	domainauthz "github.com/rxwod/backend/internal/domain/authz"
	domainwod "github.com/rxwod/backend/internal/domain/wod"
	domainwodresult "github.com/rxwod/backend/internal/domain/wodresult"
	"github.com/rxwod/backend/internal/platform/clock"
	"github.com/rxwod/backend/internal/platform/idgen"
)

type Service struct {
	repo    Repository
	wodRepo appwod.Repository
	gymRepo appgym.Repository
	clock   clock.Clock
	idgen   idgen.Generator
}

func NewService(
	repo Repository,
	wodRepo appwod.Repository,
	gymRepo appgym.Repository,
	clock clock.Clock,
	idgen idgen.Generator,
) *Service {
	return &Service{
		repo:    repo,
		wodRepo: wodRepo,
		gymRepo: gymRepo,
		clock:   clock,
		idgen:   idgen,
	}
}

func (s *Service) SubmitResult(ctx context.Context, cmd SubmitResultCommand) (ResultDTO, error) {
	principal, err := appauthz.Require(ctx, domainauthz.PermissionWODResultSubmit)
	if err != nil {
		return ResultDTO{}, err
	}

	if _, err := s.wodRepo.FindByID(ctx, principal.GymID, domainwod.WODID(cmd.WODID)); err != nil {
		return ResultDTO{}, err
	}

	membership, err := s.gymRepo.FindActiveMembership(ctx, principal.GymID, principal.UserID)
	if err != nil {
		return ResultDTO{}, err
	}

	now := s.clock.Now()
	result, err := domainwodresult.NewWODResult(
		domainwodresult.WODResultID(s.idgen.NewID()),
		domainwod.WODID(cmd.WODID),
		membership.ID(),
		domainwodresult.ScoreValue(cmd.ScoreValue),
		cmd.IsRx,
		domainwodresult.Notes(cmd.Notes),
		now,
	)
	if err != nil {
		return ResultDTO{}, err
	}

	if err := s.repo.Save(ctx, result); err != nil {
		return ResultDTO{}, fmt.Errorf("save wod result: %w", err)
	}

	saved, err := s.repo.FindByWODAndMembership(ctx, domainwod.WODID(cmd.WODID), membership.ID())
	if err != nil {
		return ResultDTO{}, fmt.Errorf("load saved wod result: %w", err)
	}

	return toResultDTO(saved), nil
}

func (s *Service) GetLeaderboard(ctx context.Context, wodID string) (LeaderboardDTO, error) {
	principal, err := appauthz.Require(ctx, domainauthz.PermissionWODResultRead)
	if err != nil {
		return LeaderboardDTO{}, err
	}

	w, err := s.wodRepo.FindByID(ctx, principal.GymID, domainwod.WODID(wodID))
	if err != nil {
		return LeaderboardDTO{}, err
	}

	scoringKind := metconScoringKind(w)

	rows, err := s.repo.ListLeaderboardByWOD(ctx, domainwod.WODID(wodID))
	if err != nil {
		return LeaderboardDTO{}, fmt.Errorf("list leaderboard: %w", err)
	}

	sortLeaderboardRows(rows, scoringKind)

	return toLeaderboardDTO(wodID, rows), nil
}

func sortLeaderboardRows(rows []LeaderboardRow, scoringKind domainwod.ScoringKind) {
	ascending := scoringKind == domainwod.ScoringTimeToComplete
	sort.SliceStable(rows, func(i, j int) bool {
		a, b := rows[i].Result, rows[j].Result
		if a.IsRx() != b.IsRx() {
			return a.IsRx()
		}
		if ascending {
			return a.ScoreValue() < b.ScoreValue()
		}
		return a.ScoreValue() > b.ScoreValue()
	})
}

func metconScoringKind(w domainwod.WOD) domainwod.ScoringKind {
	for _, stage := range w.Stages() {
		if stage.Kind() == domainwod.StageMetcon {
			return stage.ScoringKind()
		}
	}
	return domainwod.ScoringRoundsReps
}

func toResultDTO(result domainwodresult.WODResult) ResultDTO {
	return ResultDTO{
		ID:              result.ID().String(),
		WODID:           result.WODID().String(),
		GymMembershipID: result.GymMembershipID().String(),
		ScoreValue:      int(result.ScoreValue()),
		IsRx:            result.IsRx(),
		Notes:           string(result.Notes()),
		CreatedAt:       result.CreatedAt(),
		UpdatedAt:       result.UpdatedAt(),
	}
}

func toLeaderboardDTO(wodID string, rows []LeaderboardRow) LeaderboardDTO {
	entries := make([]LeaderboardEntryDTO, 0, len(rows))
	for i, row := range rows {
		result := row.Result
		entries = append(entries, LeaderboardEntryDTO{
			Rank:            i + 1,
			GymMembershipID: result.GymMembershipID().String(),
			DisplayName:     row.DisplayName,
			ScoreValue:      int(result.ScoreValue()),
			IsRx:            result.IsRx(),
			Notes:           string(result.Notes()),
			UpdatedAt:       result.UpdatedAt(),
		})
	}
	return LeaderboardDTO{
		WODID:   wodID,
		Entries: entries,
	}
}
