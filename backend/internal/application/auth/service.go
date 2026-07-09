package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/rxwod/backend/internal/application/authz"
	"github.com/rxwod/backend/internal/domain/user"
	"github.com/rxwod/backend/internal/platform/clock"
	"github.com/rxwod/backend/internal/platform/idgen"
)

type Service struct {
	users              UserRepository
	refreshTokens      RefreshTokenRepository
	passwordResets     PasswordResetTokenRepository
	emailVerifications EmailVerificationTokenRepository
	hasher             PasswordHasher
	issuer             AccessTokenIssuer
	invites            InviteAccepter
	email              EmailSender
	clock              clock.Clock
	idgen              idgen.Generator
	refreshTTL         time.Duration
	frontendURL        string
	resetTTL           time.Duration
	verificationTTL    time.Duration
}

func NewService(
	users UserRepository,
	refreshTokens RefreshTokenRepository,
	passwordResets PasswordResetTokenRepository,
	emailVerifications EmailVerificationTokenRepository,
	hasher PasswordHasher,
	issuer AccessTokenIssuer,
	invites InviteAccepter,
	email EmailSender,
	clock clock.Clock,
	idgen idgen.Generator,
	refreshTTL time.Duration,
	frontendURL string,
	resetTTL time.Duration,
	verificationTTL time.Duration,
) *Service {
	return &Service{
		users:              users,
		refreshTokens:      refreshTokens,
		passwordResets:     passwordResets,
		emailVerifications: emailVerifications,
		hasher:             hasher,
		issuer:             issuer,
		invites:            invites,
		email:              email,
		clock:              clock,
		idgen:              idgen,
		refreshTTL:         refreshTTL,
		frontendURL:        strings.TrimRight(frontendURL, "/"),
		resetTTL:           resetTTL,
		verificationTTL:    verificationTTL,
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
	if err := s.sendVerificationEmail(ctx, aggregate); err != nil {
		return TokenDTO{}, err
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
		ID:            aggregate.ID().String(),
		Email:         string(aggregate.Email()),
		DisplayName:   string(aggregate.DisplayName()),
		EmailVerified: aggregate.IsEmailVerified(),
	}, nil
}

func (s *Service) AuthenticateAccessToken(token string) (user.UserID, error) {
	userID, err := s.issuer.Verify(token)
	if err != nil {
		return "", authz.ErrUnauthenticated
	}
	return userID, nil
}

func (s *Service) RequestPasswordReset(ctx context.Context, email string) error {
	normalized := user.NormalizeEmail(email)
	aggregate, err := s.users.FindByEmail(ctx, normalized)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil
		}
		return fmt.Errorf("find user: %w", err)
	}

	now := s.clock.Now()
	if err := s.passwordResets.InvalidateUnusedForUser(ctx, aggregate.ID(), now); err != nil {
		return fmt.Errorf("invalidate reset tokens: %w", err)
	}

	rawToken := s.idgen.NewID()
	record := PasswordResetToken{
		ID:        s.idgen.NewID(),
		UserID:    aggregate.ID(),
		TokenHash: hashToken(rawToken),
		ExpiresAt: now.Add(s.resetTTL),
		CreatedAt: now,
	}
	if err := s.passwordResets.Save(ctx, record); err != nil {
		return fmt.Errorf("save reset token: %w", err)
	}

	resetURL := s.frontendURL + "/reset-password/" + url.PathEscape(rawToken)
	if err := s.email.SendPasswordReset(ctx, string(aggregate.Email()), resetURL); err != nil {
		return fmt.Errorf("send password reset email: %w", err)
	}
	return nil
}

func (s *Service) ResetPassword(ctx context.Context, token, password string) error {
	if len(password) < 8 {
		return ErrPasswordTooShort
	}

	now := s.clock.Now()
	stored, err := s.passwordResets.FindByHash(ctx, hashToken(token))
	if err != nil {
		return ErrResetTokenInvalid
	}
	if stored.UsedAt != nil || !stored.ExpiresAt.After(now) {
		return ErrResetTokenInvalid
	}

	aggregate, err := s.users.FindByID(ctx, stored.UserID)
	if err != nil {
		return ErrResetTokenInvalid
	}

	passwordHash, err := s.hasher.Hash(password)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	updated, err := user.ReconstructUser(
		aggregate.ID(),
		aggregate.Email(),
		user.PasswordHash(passwordHash),
		aggregate.DisplayName(),
		aggregate.EmailVerifiedAt(),
		aggregate.CreatedAt(),
		now,
	)
	if err != nil {
		return err
	}
	if err := s.users.Save(ctx, updated); err != nil {
		return fmt.Errorf("save user: %w", err)
	}
	if err := s.passwordResets.MarkUsed(ctx, stored.ID, now); err != nil {
		return fmt.Errorf("mark reset token used: %w", err)
	}
	if err := s.refreshTokens.RevokeAllForUser(ctx, stored.UserID, now); err != nil {
		return fmt.Errorf("revoke refresh tokens: %w", err)
	}
	return nil
}

func (s *Service) VerifyEmail(ctx context.Context, token string) error {
	now := s.clock.Now()
	stored, err := s.emailVerifications.FindByHash(ctx, hashToken(token))
	if err != nil {
		return ErrVerificationTokenInvalid
	}
	if stored.UsedAt != nil || !stored.ExpiresAt.After(now) {
		return ErrVerificationTokenInvalid
	}

	aggregate, err := s.users.FindByID(ctx, stored.UserID)
	if err != nil {
		return ErrVerificationTokenInvalid
	}
	if aggregate.IsEmailVerified() {
		if err := s.emailVerifications.MarkUsed(ctx, stored.ID, now); err != nil {
			return fmt.Errorf("mark verification token used: %w", err)
		}
		return nil
	}

	updated := aggregate.MarkEmailVerified(now)
	if err := s.users.Save(ctx, updated); err != nil {
		return fmt.Errorf("save user: %w", err)
	}
	if err := s.emailVerifications.MarkUsed(ctx, stored.ID, now); err != nil {
		return fmt.Errorf("mark verification token used: %w", err)
	}
	return nil
}

func (s *Service) ResendVerificationEmail(ctx context.Context) error {
	userID, err := authz.CurrentUserID(ctx)
	if err != nil {
		return err
	}
	aggregate, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if aggregate.IsEmailVerified() {
		return nil
	}
	return s.sendVerificationEmail(ctx, aggregate)
}

func (s *Service) sendVerificationEmail(ctx context.Context, aggregate user.User) error {
	now := s.clock.Now()
	if err := s.emailVerifications.InvalidateUnusedForUser(ctx, aggregate.ID(), now); err != nil {
		return fmt.Errorf("invalidate verification tokens: %w", err)
	}

	rawToken := s.idgen.NewID()
	record := EmailVerificationToken{
		ID:        s.idgen.NewID(),
		UserID:    aggregate.ID(),
		TokenHash: hashToken(rawToken),
		ExpiresAt: now.Add(s.verificationTTL),
		CreatedAt: now,
	}
	if err := s.emailVerifications.Save(ctx, record); err != nil {
		return fmt.Errorf("save verification token: %w", err)
	}

	verifyURL := s.frontendURL + "/verify-email/" + url.PathEscape(rawToken)
	if err := s.email.SendEmailVerification(ctx, string(aggregate.Email()), verifyURL); err != nil {
		return fmt.Errorf("send verification email: %w", err)
	}
	return nil
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
		errors.Is(err, ErrRefreshTokenInvalid) ||
		errors.Is(err, ErrResetTokenInvalid) ||
		errors.Is(err, ErrVerificationTokenInvalid) ||
		errors.Is(err, ErrEmailNotVerified)
}
