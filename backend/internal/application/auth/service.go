package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/rxwod/backend/internal/application/authz"
	"github.com/rxwod/backend/internal/domain/user"
	"github.com/rxwod/backend/internal/platform/clock"
	"github.com/rxwod/backend/internal/platform/idgen"
)

type Service struct {
	users         UserRepository
	refreshTokens RefreshTokenRepository
	hasher        PasswordHasher
	issuer        AccessTokenIssuer
	invites       InviteAccepter
	clock         clock.Clock
	idgen         idgen.Generator
	refreshTTL    time.Duration
}

func NewService(
	users UserRepository,
	refreshTokens RefreshTokenRepository,
	hasher PasswordHasher,
	issuer AccessTokenIssuer,
	invites InviteAccepter,
	clock clock.Clock,
	idgen idgen.Generator,
	refreshTTL time.Duration,
) *Service {
	return &Service{
		users:         users,
		refreshTokens: refreshTokens,
		hasher:        hasher,
		issuer:        issuer,
		invites:       invites,
		clock:         clock,
		idgen:         idgen,
		refreshTTL:    refreshTTL,
	}
}

func (s *Service) Register(ctx context.Context, cmd RegisterCommand) (TokenDTO, error) {
	if len(cmd.Password) < 8 {
		return TokenDTO{}, ErrPasswordTooShort
	}
	email := user.NormalizeEmail(cmd.Email)
	passwordHash, err := s.hasher.Hash(cmd.Password)
	if err != nil {
		return TokenDTO{}, fmt.Errorf("hash password: %w", err)
	}

	now := s.clock.Now()
	aggregate, err := user.NewUser(
		user.UserID(s.idgen.NewID()),
		email,
		user.PasswordHash(passwordHash),
		user.DisplayName(strings.TrimSpace(cmd.DisplayName)),
		now,
	)
	if err != nil {
		return TokenDTO{}, err
	}
	if err := s.users.Save(ctx, aggregate); err != nil {
		return TokenDTO{}, fmt.Errorf("save user: %w", err)
	}
	if s.invites != nil {
		if err := s.invites.AcceptInvitesForEmail(ctx, email, aggregate.ID()); err != nil {
			return TokenDTO{}, fmt.Errorf("accept invitations: %w", err)
		}
	}
	return s.issueTokens(ctx, aggregate.ID(), now)
}

func (s *Service) Login(ctx context.Context, cmd LoginCommand) (TokenDTO, error) {
	email := user.NormalizeEmail(cmd.Email)
	aggregate, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		return TokenDTO{}, ErrInvalidCredentials
	}
	if err := s.hasher.Verify(cmd.Password, string(aggregate.PasswordHash())); err != nil {
		return TokenDTO{}, ErrInvalidCredentials
	}
	return s.issueTokens(ctx, aggregate.ID(), s.clock.Now())
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (AccessTokenDTO, error) {
	now := s.clock.Now()
	stored, err := s.refreshTokens.FindByHash(ctx, hashToken(refreshToken))
	if err != nil {
		return AccessTokenDTO{}, ErrRefreshTokenInvalid
	}
	if stored.RevokedAt != nil || !stored.ExpiresAt.After(now) {
		return AccessTokenDTO{}, ErrRefreshTokenInvalid
	}
	accessToken, expiresAt, err := s.issuer.Issue(stored.UserID, now)
	if err != nil {
		return AccessTokenDTO{}, fmt.Errorf("issue access token: %w", err)
	}
	return AccessTokenDTO{
		AccessToken: accessToken,
		ExpiresIn:   int64(expiresAt.Sub(now).Seconds()),
	}, nil
}

func (s *Service) CurrentUser(ctx context.Context) (UserDTO, error) {
	userID, err := authz.CurrentUserID(ctx)
	if err != nil {
		return UserDTO{}, err
	}
	aggregate, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return UserDTO{}, err
	}
	return UserDTO{
		ID:          aggregate.ID().String(),
		Email:       string(aggregate.Email()),
		DisplayName: string(aggregate.DisplayName()),
	}, nil
}

func (s *Service) AuthenticateAccessToken(token string) (user.UserID, error) {
	userID, err := s.issuer.Verify(token)
	if err != nil {
		return "", authz.ErrUnauthenticated
	}
	return userID, nil
}

func (s *Service) issueTokens(ctx context.Context, userID user.UserID, now time.Time) (TokenDTO, error) {
	accessToken, accessExpiresAt, err := s.issuer.Issue(userID, now)
	if err != nil {
		return TokenDTO{}, fmt.Errorf("issue access token: %w", err)
	}

	refreshToken := s.idgen.NewID()
	record := RefreshToken{
		ID:        s.idgen.NewID(),
		UserID:    userID,
		TokenHash: hashToken(refreshToken),
		ExpiresAt: now.Add(s.refreshTTL),
		CreatedAt: now,
	}
	if err := s.refreshTokens.Save(ctx, record); err != nil {
		return TokenDTO{}, fmt.Errorf("save refresh token: %w", err)
	}

	return TokenDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(accessExpiresAt.Sub(now).Seconds()),
	}, nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func IsAuthError(err error) bool {
	return errors.Is(err, ErrInvalidCredentials) ||
		errors.Is(err, ErrPasswordTooShort) ||
		errors.Is(err, ErrRefreshTokenInvalid)
}
