package jwt_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	jwtLib "github.com/golang-jwt/jwt/v5"
	"go-ez-auth/core"
	"go-ez-auth/stores"
	"go-ez-auth/strategies/jwt"
)

type dummyUser struct{ id string }
func (d dummyUser) GetID() string                 { return d.id }
func (d dummyUser) GetAttributes() map[string]interface{} { return nil }

func TestAuthenticate_NoHeader(t *testing.T) {
	s := jwt.New(jwt.Config{SigningKey: []byte("secret")})
	req, _ := http.NewRequest("GET", "/", nil)
	_, err := s.Authenticate(context.Background(), req)
	if err != core.ErrUnauthorized {
		t.Errorf("expected ErrUnauthorized, got %v", err)
	}
}

func TestAuthenticate_InvalidToken(t *testing.T) {
	s := jwt.New(jwt.Config{SigningKey: []byte("secret")})
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")
	_, err := s.Authenticate(context.Background(), req)
	if err != core.ErrUnauthorized {
		t.Errorf("expected ErrUnauthorized for invalid token, got %v", err)
	}
}

func TestAuthenticate_ValidToken_NoStore(t *testing.T) {
	// Create a signed token with custom issuer/audience
	key := []byte("secret")
	claims := jwtLib.RegisteredClaims{
		Subject:   "user123",
		Issuer:    "issuer",
		Audience:  jwtLib.ClaimStrings{"aud1"},
		ExpiresAt: jwtLib.NewNumericDate(time.Now().Add(time.Hour)),
	}
	token := jwtLib.NewWithClaims(jwtLib.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(key)
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	s := jwt.New(jwt.Config{
		SigningKey:    key,
		SigningMethod: jwtLib.SigningMethodHS256.Alg(),
		Issuer:        "issuer",
		Audience:      "aud1",
	})

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	user, err := s.Authenticate(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.GetID() != "user123" {
		t.Errorf("expected user123, got %s", user.GetID())
	}
	attrs := user.GetAttributes()
	if attrs["issuer"] != "issuer" {
		t.Errorf("expected issuer 'issuer', got %v", attrs["issuer"])
	}
	aud, ok := attrs["audience"].(jwtLib.ClaimStrings)
	if !ok {
		t.Errorf("expected jwtLib.ClaimStrings, got %T", attrs["audience"])
	}
	found := false
	for _, a := range aud {
		if a == "aud1" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected audience aud1, got %v", aud)
	}
}

func TestAuthenticate_ValidToken_WithStore(t *testing.T) {
	// Create a signed token with subject only
	key := []byte("secret")
	claims := jwtLib.RegisteredClaims{
		Subject:   "user456",
		ExpiresAt: jwtLib.NewNumericDate(time.Now().Add(time.Hour)),
	}
	token := jwtLib.NewWithClaims(jwtLib.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(key)
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	dummy := dummyUser{"user456"}
	store := stores.NewInMemoryUserStore(dummy)
	s := jwt.New(jwt.Config{
		SigningKey: key,
		Store:      store,
	})

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	user, err := s.Authenticate(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.GetID() != "user456" {
		t.Errorf("expected user456, got %s", user.GetID())
	}
}
