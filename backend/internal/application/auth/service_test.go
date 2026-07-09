package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/rxwod/backend/internal/application/authz"
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
	resetCalls       []emailCall
	verificationCalls []verificationCall
}

type emailCall struct {
	email    string
	resetURL string
}

type verificationCall struct {
	email     string
	verifyURL string
}

func (f *fakeEmailSender) SendPasswordReset(_ context.Context, toEmail, resetURL string) error {
	f.resetCalls = append(f.resetCalls, emailCall{email: toEmail, resetURL: resetURL})
	return nil
}

func (f *fakeEmailSender) SendEmailVerification(_ context.Context, toEmail, verifyURL string) error {
	f.verificationCalls = append(f.verificationCalls, verificationCall{email: toEmail, verifyURL: verifyURL})
	return nil
}

type fakeEmailVerifications struct {
	tokens   map[string]EmailVerificationToken
	markedID string
}

func (f *fakeEmailVerifications) Save(_ context.Context, token EmailVerificationToken) error {
	if f.tokens == nil {
		f.tokens = make(map[string]EmailVerificationToken)
	}
	f.tokens[token.TokenHash] = token
	return nil
}

func (f *fakeEmailVerifications) FindByHash(_ context.Context, tokenHash string) (EmailVerificationToken, error) {
	token, ok := f.tokens[tokenHash]
	if !ok {
		return EmailVerificationToken{}, errors.New("not found")
	}
	return token, nil
}

func (fakeEmailVerifications) InvalidateUnusedForUser(context.Context, user.UserID, time.Time) error {
	return nil
}

func (f *fakeEmailVerifications) MarkUsed(_ context.Context, id string, now time.Time) error {
	f.markedID = id
	for hash, token := range f.tokens {
		if token.ID == id {
			token.UsedAt = &now
			f.tokens[hash] = token
			break
		}
	}
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

func newTestService(t *testing.T, users *fakeUsers, resets *fakePasswordResets, verifications *fakeEmailVerifications, refresh *fakeRefreshTokens, email *fakeEmailSender, hasher *fakeHasher, now time.Time, ids ...string) *Service {
	t.Helper()
	if len(ids) == 0 {
		ids = []string{"reset-token", "reset-record"}
	}
	if verifications == nil {
		verifications = &fakeEmailVerifications{}
	}
	return NewService(
		users,
		refresh,
		resets,
		verifications,
		hasher,
		fakeIssuer{},
		nil,
		email,
		fakeClock{now: now},
		&sequentialIDGen{ids: ids},
		time.Hour,
		"http://localhost:5173",
		time.Hour,
		time.Hour,
	)
}

func hashTestToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func TestRegisterRejectsShortPassword(t *testing.T) {
	service := newTestService(t, &fakeUsers{byEmail: map[user.Email]user.User{}, byID: map[user.UserID]user.User{}}, &fakePasswordResets{}, &fakeEmailVerifications{}, &fakeRefreshTokens{}, &fakeEmailSender{}, &fakeHasher{}, time.Now().UTC())

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
		&fakeEmailVerifications{},
		&fakeRefreshTokens{},
		emailSender,
		&fakeHasher{},
		time.Now().UTC(),
	)

	if err := service.RequestPasswordReset(context.Background(), "missing@example.com"); err != nil {
		t.Fatalf("RequestPasswordReset() error = %v, want nil", err)
	}
	if len(emailSender.resetCalls) != 0 {
		t.Fatalf("expected no email calls, got %d", len(emailSender.resetCalls))
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
	service := newTestService(t, users, resets, &fakeEmailVerifications{}, &fakeRefreshTokens{}, emailSender, &fakeHasher{}, now)

	if err := service.RequestPasswordReset(context.Background(), string(aggregate.Email())); err != nil {
		t.Fatalf("RequestPasswordReset() error = %v", err)
	}
	if len(resets.tokens) != 1 {
		t.Fatalf("expected 1 reset token, got %d", len(resets.tokens))
	}
	if len(emailSender.resetCalls) != 1 {
		t.Fatalf("expected 1 email call, got %d", len(emailSender.resetCalls))
	}
	if emailSender.resetCalls[0].email != string(aggregate.Email()) {
		t.Fatalf("email = %q, want %q", emailSender.resetCalls[0].email, aggregate.Email())
	}
	if !strings.Contains(emailSender.resetCalls[0].resetURL, "/reset-password/reset-token") {
		t.Fatalf("reset URL = %q, want token path", emailSender.resetCalls[0].resetURL)
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
	service := newTestService(t, users, resets, &fakeEmailVerifications{}, &fakeRefreshTokens{}, &fakeEmailSender{}, hasher, now)

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
		&fakeEmailVerifications{},
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
		&fakeEmailVerifications{},
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
	service := newTestService(t, users, resets, &fakeEmailVerifications{}, refresh, &fakeEmailSender{}, &fakeHasher{}, now)

	if err := service.ResetPassword(context.Background(), "reset-token", "new-password"); err != nil {
		t.Fatalf("ResetPassword() error = %v", err)
	}
	if refresh.revokedUser != aggregate.ID() {
		t.Fatalf("revoked user = %q, want %q", refresh.revokedUser, aggregate.ID())
	}
}

func TestRegisterSendsVerificationEmail(t *testing.T) {
	now := time.Date(2026, 7, 9, 12, 0, 0, 0, time.UTC)
	users := &fakeUsers{byEmail: map[user.Email]user.User{}, byID: map[user.UserID]user.User{}}
	verifications := &fakeEmailVerifications{}
	emailSender := &fakeEmailSender{}
	service := newTestService(
		t,
		users,
		&fakePasswordResets{},
		verifications,
		&fakeRefreshTokens{},
		emailSender,
		&fakeHasher{},
		now,
		"user-id", "verify-token", "verify-record", "refresh-token", "refresh-record",
	)

	_, err := service.Register(context.Background(), RegisterCommand{
		Email:       "owner@example.com",
		Password:    "password123",
		DisplayName: "Owner",
	})
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	if len(verifications.tokens) != 1 {
		t.Fatalf("expected 1 verification token, got %d", len(verifications.tokens))
	}
	if len(emailSender.verificationCalls) != 1 {
		t.Fatalf("expected 1 verification email, got %d", len(emailSender.verificationCalls))
	}
	if !strings.Contains(emailSender.verificationCalls[0].verifyURL, "/verify-email/verify-token") {
		t.Fatalf("verify URL = %q, want token path", emailSender.verificationCalls[0].verifyURL)
	}
}

func TestVerifyEmailMarksUserVerified(t *testing.T) {
	now := time.Date(2026, 7, 9, 12, 0, 0, 0, time.UTC)
	aggregate := testUser(t, now)
	users := &fakeUsers{
		byEmail: map[user.Email]user.User{aggregate.Email(): aggregate},
		byID:    map[user.UserID]user.User{aggregate.ID(): aggregate},
	}
	verifications := &fakeEmailVerifications{
		tokens: map[string]EmailVerificationToken{
			hashTestToken("verify-token"): {
				ID:        "verify-record",
				UserID:    aggregate.ID(),
				TokenHash: hashTestToken("verify-token"),
				ExpiresAt: now.Add(time.Hour),
				CreatedAt: now,
			},
		},
	}
	service := newTestService(t, users, &fakePasswordResets{}, verifications, &fakeRefreshTokens{}, &fakeEmailSender{}, &fakeHasher{}, now)

	if err := service.VerifyEmail(context.Background(), "verify-token"); err != nil {
		t.Fatalf("VerifyEmail() error = %v", err)
	}
	if len(users.saved) != 1 || !users.saved[0].IsEmailVerified() {
		t.Fatalf("expected verified user save")
	}
	if verifications.markedID != "verify-record" {
		t.Fatalf("marked token id = %q, want verify-record", verifications.markedID)
	}
}

func TestVerifyEmailRejectsInvalidToken(t *testing.T) {
	now := time.Now().UTC()
	service := newTestService(
		t,
		&fakeUsers{byEmail: map[user.Email]user.User{}, byID: map[user.UserID]user.User{}},
		&fakePasswordResets{},
		&fakeEmailVerifications{},
		&fakeRefreshTokens{},
		&fakeEmailSender{},
		&fakeHasher{},
		now,
	)

	err := service.VerifyEmail(context.Background(), "missing-token")
	if !errors.Is(err, ErrVerificationTokenInvalid) {
		t.Fatalf("expected ErrVerificationTokenInvalid, got %v", err)
	}
}

func TestResendVerificationSkipsVerifiedUser(t *testing.T) {
	now := time.Date(2026, 7, 9, 12, 0, 0, 0, time.UTC)
	aggregate := testUser(t, now).MarkEmailVerified(now)
	users := &fakeUsers{
		byEmail: map[user.Email]user.User{aggregate.Email(): aggregate},
		byID:    map[user.UserID]user.User{aggregate.ID(): aggregate},
	}
	emailSender := &fakeEmailSender{}
	service := newTestService(t, users, &fakePasswordResets{}, &fakeEmailVerifications{}, &fakeRefreshTokens{}, emailSender, &fakeHasher{}, now)

	ctx := authz.WithPrincipal(context.Background(), authz.Principal{UserID: aggregate.ID()})
	if err := service.ResendVerificationEmail(ctx); err != nil {
		t.Fatalf("ResendVerificationEmail() error = %v", err)
	}
	if len(emailSender.verificationCalls) != 0 {
		t.Fatalf("expected no verification email, got %d", len(emailSender.verificationCalls))
	}
}

func TestResendVerificationSendsForUnverified(t *testing.T) {
	now := time.Date(2026, 7, 9, 12, 0, 0, 0, time.UTC)
	aggregate := testUser(t, now)
	users := &fakeUsers{
		byEmail: map[user.Email]user.User{aggregate.Email(): aggregate},
		byID:    map[user.UserID]user.User{aggregate.ID(): aggregate},
	}
	verifications := &fakeEmailVerifications{}
	emailSender := &fakeEmailSender{}
	service := newTestService(
		t,
		users,
		&fakePasswordResets{},
		verifications,
		&fakeRefreshTokens{},
		emailSender,
		&fakeHasher{},
		now,
		"verify-token", "verify-record",
	)

	ctx := authz.WithPrincipal(context.Background(), authz.Principal{UserID: aggregate.ID()})
	if err := service.ResendVerificationEmail(ctx); err != nil {
		t.Fatalf("ResendVerificationEmail() error = %v", err)
	}
	if len(verifications.tokens) != 1 {
		t.Fatalf("expected 1 verification token, got %d", len(verifications.tokens))
	}
	if len(emailSender.verificationCalls) != 1 {
		t.Fatalf("expected 1 verification email, got %d", len(emailSender.verificationCalls))
	}
}
