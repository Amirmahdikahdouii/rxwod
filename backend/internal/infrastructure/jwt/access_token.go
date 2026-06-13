package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/rxwod/backend/internal/domain/user"
)

var ErrInvalidToken = errors.New("access token is invalid")

type AccessTokenIssuer struct {
	secret []byte
	ttl    time.Duration
	clock  func() time.Time
}

func NewAccessTokenIssuer(secret string, ttl time.Duration) *AccessTokenIssuer {
	return &AccessTokenIssuer{
		secret: []byte(secret),
		ttl:    ttl,
		clock:  time.Now,
	}
}

type tokenHeader struct {
	Algorithm string `json:"alg"`
	Type      string `json:"typ"`
}

type tokenClaims struct {
	Subject   string `json:"sub"`
	ExpiresAt int64  `json:"exp"`
	IssuedAt  int64  `json:"iat"`
}

func (i *AccessTokenIssuer) Issue(userID user.UserID, now time.Time) (string, time.Time, error) {
	expiresAt := now.Add(i.ttl)
	header := tokenHeader{Algorithm: "HS256", Type: "JWT"}
	claims := tokenClaims{Subject: userID.String(), ExpiresAt: expiresAt.Unix(), IssuedAt: now.Unix()}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("marshal jwt header: %w", err)
	}
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("marshal jwt claims: %w", err)
	}

	encodedHeader := base64.RawURLEncoding.EncodeToString(headerJSON)
	encodedClaims := base64.RawURLEncoding.EncodeToString(claimsJSON)
	signingInput := encodedHeader + "." + encodedClaims
	signature := i.sign(signingInput)
	return signingInput + "." + signature, expiresAt, nil
}

func (i *AccessTokenIssuer) Verify(token string) (user.UserID, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", ErrInvalidToken
	}
	signingInput := parts[0] + "." + parts[1]
	if !hmac.Equal([]byte(parts[2]), []byte(i.sign(signingInput))) {
		return "", ErrInvalidToken
	}

	var header tokenHeader
	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return "", ErrInvalidToken
	}
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return "", ErrInvalidToken
	}
	if header.Algorithm != "HS256" || header.Type != "JWT" {
		return "", ErrInvalidToken
	}

	var claims tokenClaims
	claimsBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", ErrInvalidToken
	}
	if err := json.Unmarshal(claimsBytes, &claims); err != nil {
		return "", ErrInvalidToken
	}
	if claims.Subject == "" || !time.Unix(claims.ExpiresAt, 0).After(i.clock().UTC()) {
		return "", ErrInvalidToken
	}
	return user.UserID(claims.Subject), nil
}

func (i *AccessTokenIssuer) sign(input string) string {
	mac := hmac.New(sha256.New, i.secret)
	mac.Write([]byte(input))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
