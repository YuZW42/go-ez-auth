package jwt

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"go-ez-auth/core"
	jwtLib "github.com/golang-jwt/jwt/v5"
)

// Config holds settings for the JWT strategy.
type Config struct {
	SigningKey    []byte
	SigningMethod string
	Issuer        string
	Audience      string
	Store         core.UserStore
}

// Strategy implements the core.Strategy interface for JWT.
type Strategy struct {
	config Config
}

// New creates a JWT strategy with the given config, defaulting to HS256 if unspecified.
func New(config Config) *Strategy {
	if config.SigningMethod == "" {
		config.SigningMethod = jwtLib.SigningMethodHS256.Alg()
	}
	return &Strategy{config: config}
}

// Name returns the strategy name.
func (s *Strategy) Name() string {
	return "jwt"
}

// Setup is a no-op for JWT strategy.
func (s *Strategy) Setup() error {
	return nil
}

// Authenticate extracts and validates a JWT from the Authorization header.
func (s *Strategy) Authenticate(ctx context.Context, r *http.Request) (core.User, error) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return nil, core.ErrUnauthorized
	}
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, core.ErrUnauthorized
	}
	tokenString := parts[1]

	claims := &jwtLib.RegisteredClaims{}
	token, err := jwtLib.ParseWithClaims(tokenString, claims, func(token *jwtLib.Token) (interface{}, error) {
		if token.Method.Alg() != s.config.SigningMethod {
			return nil, errors.New("unexpected signing method")
		}
		return s.config.SigningKey, nil
	})
	if err != nil || !token.Valid {
		return nil, core.ErrUnauthorized
	}

	// Validate issuer
	if s.config.Issuer != "" && claims.Issuer != s.config.Issuer {
		return nil, core.ErrUnauthorized
	}
	// Validate audience
	if s.config.Audience != "" {
		found := false
		for _, a := range claims.Audience {
			if a == s.config.Audience {
				found = true
				break
			}
		}
		if !found {
			return nil, core.ErrUnauthorized
		}
	}

	userID := claims.Subject
	if userID == "" {
		return nil, core.ErrUnauthorized
	}
	// If a store is provided, lookup the user
	if s.config.Store != nil {
		user, err := s.config.Store.FindUserByID(ctx, userID)
		if err != nil {
			return nil, core.ErrUnauthorized
		}
		return user, nil
	}
	// Otherwise return a simple user with claims as attributes
	attrs := map[string]interface{}{ // include basic claims
		"issuer":   claims.Issuer,
		"audience": claims.Audience,
		"expires":  claims.ExpiresAt,
	}
	return &jwtUser{id: userID, attributes: attrs}, nil
}

// jwtUser is a simple User implementation for JWTStrategy.
type jwtUser struct {
	id         string
	attributes map[string]interface{}
}

func (u *jwtUser) GetID() string {
	return u.id
}

func (u *jwtUser) GetAttributes() map[string]interface{} {
	return u.attributes
}
