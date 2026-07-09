package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/rxwod/backend/internal/domain/user"
)

type fakeUsers struct {
	byEmail map[user.Email]user.User
	byID    map[user.UserID]user.User
	saved   []user.User
}

func (f *fakeUsers) Save(_ context.Context, aggregate user.User) error {
	f.saved = append(f.saved, aggregate)
	f.byEmail[aggregate.Email()] = aggregate
	f.byID[aggregate.ID()] = aggregate
	return nil
}

func (f *fakeUsers) FindByEmail(_ context.Context, email user.Email) (user.User, error) {
	if aggregate, ok := f.byEmail[email]; ok {
		return aggregate, nil
	}
	return user.User{}, ErrUserNotFound
}

func (f *fakeUsers) FindByID(_ context.Context, id user.UserID) (user.User, error) {
	if aggregate, ok := f.byID[id]; ok {
		return aggregate, nil
	}
	return user.User{}, ErrUserNotFound
}

type fakeRefreshTokens struct {
	revokedUser user.UserID
}

func (fakeRefreshTokens) Save(context.Context, RefreshToken) error { return nil }
func (fakeRefreshTokens) FindByHash(context.Context, string) (RefreshToken, error) {
	return RefreshToken{}, errors.New("not found")
}

func (f *fakeRefreshTokens) RevokeAllForUser(_ context.Context, userID user.UserID, _ time.Time) error {
	f.revokedUser = userID
	return nil
}

type fakePasswordResets struct {
	tokens   map[string]PasswordResetToken
	markedID string
}

func (f *fakePasswordResets) Save(_ context.Context, token PasswordResetToken) error {
	if f.tokens == nil {
		f.tokens = make(map[string]PasswordResetToken)
	}
	f.tokens[token.TokenHash] = token
	return nil
}

func (f *fakePasswordResets) FindByHash(_ context.Context, tokenHash string) (PasswordResetToken, error) {
	token, ok := f.tokens[tokenHash]
	if !ok {
		return PasswordResetToken{}, errors.New("not found")
	}
	return token, nil
}

func (fakePasswordResets) InvalidateUnusedForUser(context.Context, user.UserID, time.Time) error {
	return nil
}

func (f *fakePasswordResets) MarkUsed(_ context.Context, id string, now time.Time) error {
	f.markedID = id
	if token, ok := f.tokensByID()[id]; ok {
		token.UsedAt = &now
		f.tokens[token.TokenHash] = token
	}
	return nil
}

func (f *fakePasswordResets) tokensByID() map[string]PasswordResetToken {
	byID := make(map[string]PasswordResetToken, len(f.tokens))
	for _, token := range f.tokens {
		byID[token.ID] = token
	}
	return byID
}

type fakeEmailSender struct {
	calls []emailCall
}

type emailCall struct {
	email    string
	resetURL string
}

func (f *fakeEmailSender) SendPasswordReset(_ context.Context, toEmail, resetURL string) error {
	f.calls = append(f.calls, emailCall{email: toEmail, resetURL: resetURL})
	return nil
}

type fakeHasher struct {
	lastHash string
}

func (f *fakeHasher) Hash(password string) (string, error) {
	f.lastHash = "hash:" + password
	return f.lastHash, nil
}

func (fakeHasher) Verify(string, string) error { return nil }

type fakeIssuer struct{}

func (fakeIssuer) Issue(user.UserID, time.Time) (string, time.Time, error) {
	return "access-token", time.Now().Add(time.Minute), nil
}
func (fakeIssuer) Verify(string) (user.UserID, error) { return "user-1", nil }

type fakeClock struct {
	now time.Time
}

func (f fakeClock) Now() time.Time { return f.now }

type sequentialIDGen struct {
	ids []string
	idx int
}

func (g *sequentialIDGen) NewID() string {
	if g.idx >= len(g.ids) {
		return "generated-id"
	}
	id := g.ids[g.idx]
	g.idx++
	return id
}

func testUser(t *testing.T, now time.Time) user.User {
	t.Helper()
	aggregate, err := user.NewUser(
		user.UserID("user-1"),
		user.Email("owner@example.com"),
		user.PasswordHash("old-hash"),
		user.DisplayName("Owner"),
		now,
	)
	if err != nil {
		t.Fatalf("new user: %v", err)
	}
	return aggregate
}

func newTestService(t *testing.T, users *fakeUsers, resets *fakePasswordResets, refresh *fakeRefreshTokens, email *fakeEmailSender, hasher *fakeHasher, now time.Time, ids ...string) *Service {
	t.Helper()
	if len(ids) == 0 {
		ids = []string{"reset-token", "reset-record"}
	}
	return NewService(
		users,
		refresh,
		resets,
		hasher,
		fakeIssuer{},
		nil,
		email,
		fakeClock{now: now},
		&sequentialIDGen{ids: ids},
		time.Hour,
		"http://localhost:5173",
		time.Hour,
	)
}

func hashTestToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func TestRegisterRejectsShortPassword(t *testing.T) {
	service := newTestService(t, &fakeUsers{byEmail: map[user.Email]user.User{}, byID: map[user.UserID]user.User{}}, &fakePasswordResets{}, &fakeRefreshTokens{}, &fakeEmailSender{}, &fakeHasher{}, time.Now().UTC())

	_, err := service.Register(context.Background(), RegisterCommand{
		Email:       "owner@example.com",
		Password:    "short",
		DisplayName: "Owner",
	})
	if !errors.Is(err, ErrPasswordTooShort) {
		t.Fatalf("expected ErrPasswordTooShort, got %v", err)
	}
}

func TestRequestPasswordResetUnknownEmail(t *testing.T) {
	emailSender := &fakeEmailSender{}
	service := newTestService(
		t,
		&fakeUsers{byEmail: map[user.Email]user.User{}, byID: map[user.UserID]user.User{}},
		&fakePasswordResets{},
		&fakeRefreshTokens{},
		emailSender,
		&fakeHasher{},
		time.Now().UTC(),
	)

	if err := service.RequestPasswordReset(context.Background(), "missing@example.com"); err != nil {
		t.Fatalf("RequestPasswordReset() error = %v, want nil", err)
	}
	if len(emailSender.calls) != 0 {
		t.Fatalf("expected no email calls, got %d", len(emailSender.calls))
	}
}

func TestRequestPasswordResetCreatesToken(t *testing.T) {
	now := time.Date(2026, 7, 9, 12, 0, 0, 0, time.UTC)
	aggregate := testUser(t, now)
	users := &fakeUsers{
		byEmail: map[user.Email]user.User{aggregate.Email(): aggregate},
		byID:    map[user.UserID]user.User{aggregate.ID(): aggregate},
	}
	resets := &fakePasswordResets{}
	emailSender := &fakeEmailSender{}
	service := newTestService(t, users, resets, &fakeRefreshTokens{}, emailSender, &fakeHasher{}, now)

	if err := service.RequestPasswordReset(context.Background(), string(aggregate.Email())); err != nil {
		t.Fatalf("RequestPasswordReset() error = %v", err)
	}
	if len(resets.tokens) != 1 {
		t.Fatalf("expected 1 reset token, got %d", len(resets.tokens))
	}
	if len(emailSender.calls) != 1 {
		t.Fatalf("expected 1 email call, got %d", len(emailSender.calls))
	}
	if emailSender.calls[0].email != string(aggregate.Email()) {
		t.Fatalf("email = %q, want %q", emailSender.calls[0].email, aggregate.Email())
	}
	if !strings.Contains(emailSender.calls[0].resetURL, "/reset-password/reset-token") {
		t.Fatalf("reset URL = %q, want token path", emailSender.calls[0].resetURL)
	}
}

func TestResetPasswordUpdatesHash(t *testing.T) {
	now := time.Date(2026, 7, 9, 12, 0, 0, 0, time.UTC)
	aggregate := testUser(t, now)
	users := &fakeUsers{
		byEmail: map[user.Email]user.User{aggregate.Email(): aggregate},
		byID:    map[user.UserID]user.User{aggregate.ID(): aggregate},
	}
	resets := &fakePasswordResets{
		tokens: map[string]PasswordResetToken{
			hashTestToken("reset-token"): {
				ID:        "reset-record",
				UserID:    aggregate.ID(),
				TokenHash: hashTestToken("reset-token"),
				ExpiresAt: now.Add(time.Hour),
				CreatedAt: now,
			},
		},
	}
	hasher := &fakeHasher{}
	service := newTestService(t, users, resets, &fakeRefreshTokens{}, &fakeEmailSender{}, hasher, now)

	if err := service.ResetPassword(context.Background(), "reset-token", "new-password"); err != nil {
		t.Fatalf("ResetPassword() error = %v", err)
	}
	if len(users.saved) != 1 {
		t.Fatalf("expected user save, got %d saves", len(users.saved))
	}
	if string(users.saved[0].PasswordHash()) != "hash:new-password" {
		t.Fatalf("password hash = %q, want %q", users.saved[0].PasswordHash(), "hash:new-password")
	}
	if resets.markedID != "reset-record" {
		t.Fatalf("marked token id = %q, want reset-record", resets.markedID)
	}
}

func TestResetPasswordRejectsInvalidToken(t *testing.T) {
	now := time.Now().UTC()
	service := newTestService(
		t,
		&fakeUsers{byEmail: map[user.Email]user.User{}, byID: map[user.UserID]user.User{}},
		&fakePasswordResets{},
		&fakeRefreshTokens{},
		&fakeEmailSender{},
		&fakeHasher{},
		now,
	)

	err := service.ResetPassword(context.Background(), "missing-token", "new-password")
	if !errors.Is(err, ErrResetTokenInvalid) {
		t.Fatalf("expected ErrResetTokenInvalid, got %v", err)
	}
}

func TestResetPasswordRejectsShortPassword(t *testing.T) {
	now := time.Now().UTC()
	service := newTestService(
		t,
		&fakeUsers{byEmail: map[user.Email]user.User{}, byID: map[user.UserID]user.User{}},
		&fakePasswordResets{},
		&fakeRefreshTokens{},
		&fakeEmailSender{},
		&fakeHasher{},
		now,
	)

	err := service.ResetPassword(context.Background(), "reset-token", "short")
	if !errors.Is(err, ErrPasswordTooShort) {
		t.Fatalf("expected ErrPasswordTooShort, got %v", err)
	}
}

func TestResetPasswordRevokesRefreshTokens(t *testing.T) {
	now := time.Date(2026, 7, 9, 12, 0, 0, 0, time.UTC)
	aggregate := testUser(t, now)
	users := &fakeUsers{
		byEmail: map[user.Email]user.User{aggregate.Email(): aggregate},
		byID:    map[user.UserID]user.User{aggregate.ID(): aggregate},
	}
	resets := &fakePasswordResets{
		tokens: map[string]PasswordResetToken{
			hashTestToken("reset-token"): {
				ID:        "reset-record",
				UserID:    aggregate.ID(),
				TokenHash: hashTestToken("reset-token"),
				ExpiresAt: now.Add(time.Hour),
				CreatedAt: now,
			},
		},
	}
	refresh := &fakeRefreshTokens{}
	service := newTestService(t, users, resets, refresh, &fakeEmailSender{}, &fakeHasher{}, now)

	if err := service.ResetPassword(context.Background(), "reset-token", "new-password"); err != nil {
		t.Fatalf("ResetPassword() error = %v", err)
	}
	if refresh.revokedUser != aggregate.ID() {
		t.Fatalf("revoked user = %q, want %q", refresh.revokedUser, aggregate.ID())
	}
}
