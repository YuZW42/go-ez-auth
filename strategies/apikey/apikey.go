package apikey

import (
	"context"
	"net/http"

	"go-ez-auth/core"
)

// Config holds settings for the API key strategy.
type Config struct {
	HeaderName string         // HTTP header name for API key
	QueryParam string         // URL query parameter name for API key
	CredKey    string         // credential key used in UserStore lookup
	Store      core.UserStore // backend for user lookup
}

// Strategy implements core.Strategy for API key authentication.
type Strategy struct {
	config Config
}

// New creates an API key strategy with defaults.
func New(config Config) *Strategy {
	if config.HeaderName == "" {
		config.HeaderName = "X-API-Key"
	}
	if config.QueryParam == "" {
		config.QueryParam = "api_key"
	}
	if config.CredKey == "" {
		config.CredKey = "id"
	}
	return &Strategy{config: config}
}

// Name returns the strategy name.
func (s *Strategy) Name() string {
	return "apikey"
}

// Setup is a no-op for API key strategy.
func (s *Strategy) Setup() error {
	return nil
}

// Authenticate extracts the API key from header or query param and validates it.
func (s *Strategy) Authenticate(ctx context.Context, r *http.Request) (core.User, error) {
	// Try header
	key := r.Header.Get(s.config.HeaderName)
	// Fallback to query param
	if key == "" {
		key = r.URL.Query().Get(s.config.QueryParam)
	}
	if key == "" {
		return nil, core.ErrUnauthorized
	}
	// Lookup user by credential
	criteria := map[string]interface{}{s.config.CredKey: key}
	user, err := s.config.Store.FindUserByCredentials(ctx, criteria)
	if err != nil {
		return nil, core.ErrUnauthorized
	}
	return user, nil
}
