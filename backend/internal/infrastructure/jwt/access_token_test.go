package jwt

import (
	"errors"
	"testing"
	"time"

	"github.com/rxwod/backend/internal/domain/user"
)

func TestAccessTokenIssueAndVerify(t *testing.T) {
	issuer := NewAccessTokenIssuer("secret", time.Minute)

	token, expiresAt, err := issuer.Issue(user.UserID("user-1"), time.Now().UTC())
	if err != nil {
		t.Fatalf("issue token: %v", err)
	}
	if !expiresAt.After(time.Now().UTC()) {
		t.Fatalf("expected future expiry, got %s", expiresAt)
	}

	userID, err := issuer.Verify(token)
	if err != nil {
		t.Fatalf("verify token: %v", err)
	}
	if userID != "user-1" {
		t.Fatalf("expected user-1, got %s", userID)
	}
}

func TestAccessTokenRejectsTamperedToken(t *testing.T) {
	issuer := NewAccessTokenIssuer("secret", time.Minute)
	token, _, err := issuer.Issue(user.UserID("user-1"), time.Now().UTC())
	if err != nil {
		t.Fatalf("issue token: %v", err)
	}

	_, err = issuer.Verify(token + "tampered")
	if !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("expected ErrInvalidToken, got %v", err)
	}
}
