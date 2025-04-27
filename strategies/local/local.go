package local

import (
	"context"
	"net/http"

	"go-ez-auth/core"
)

// Config holds settings for the local username/password strategy.
// Users are authenticated via Basic Auth and validated against the UserStore.
type Config struct {
	UserStore core.UserStore
}

// Strategy implements core.Strategy for local auth.
type Strategy struct {
	config Config
}

// New creates a new local auth strategy.
func New(config Config) *Strategy {
	return &Strategy{config: config}
}

// Name returns the strategy name.
func (s *Strategy) Name() string {
	return "local"
}

// Setup is a no-op for the local strategy.
func (s *Strategy) Setup() error {
	return nil
}

// Authenticate parses Basic Auth credentials and delegates to UserStore.FindUserByCredentials.
func (s *Strategy) Authenticate(ctx context.Context, r *http.Request) (core.User, error) {
	username, password, ok := r.BasicAuth()
	if !ok {
		return nil, core.ErrUnauthorized
	}
	// Delegate credential lookup with criteria map
	user, err := s.config.UserStore.FindUserByCredentials(ctx, map[string]interface{}{"username": username, "password": password})
	if err != nil {
		return nil, core.ErrUnauthorized
	}
	return user, nil
}
