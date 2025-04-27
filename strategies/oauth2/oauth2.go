package oauth2

import (
	"context"
	"encoding/json"
	"net/http"

	"golang.org/x/oauth2"
	"go-ez-auth/core"
)

// Config holds settings for the OAuth2/OIDC strategy.
// OAuth2Config: configured client, redirect URL, endpoints, scopes.
// UserInfoURL: endpoint to fetch user profile with Bearer token.
// ExtractUser: maps userinfo JSON to core.User.
type Config struct {
	OAuth2Config *oauth2.Config
	UserInfoURL  string
	ExtractUser  func(ctx context.Context, info map[string]interface{}) (core.User, error)
}

// Strategy implements core.Strategy using OAuth2/OIDC.
type Strategy struct {
	config Config
}

// New creates an OAuth2 strategy from Config.
func New(config Config) *Strategy {
	return &Strategy{config: config}
}

// Name returns the strategy name.
func (s *Strategy) Name() string {
	return "oauth2"
}

// Setup is a no-op for OAuth2 strategy.
func (s *Strategy) Setup() error {
	return nil
}

// Authenticate handles OAuth2 callback: exchanges code, fetches userinfo, and returns core.User.
func (s *Strategy) Authenticate(ctx context.Context, r *http.Request) (core.User, error) {
	code := r.URL.Query().Get("code")
	if code == "" {
		return nil, core.ErrUnauthorized
	}

	// Exchange code for token
	tok, err := s.config.OAuth2Config.Exchange(ctx, code)
	if err != nil {
		return nil, core.ErrUnauthorized
	}

	// Fetch user info
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.config.UserInfoURL, nil)
	if err != nil {
		return nil, core.ErrUnauthorized
	}
	req.Header.Set("Authorization", "Bearer "+tok.AccessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, core.ErrUnauthorized
	}
	defer resp.Body.Close()

	var info map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, core.ErrUnauthorized
	}

	// Extract core.User
	user, err := s.config.ExtractUser(ctx, info)
	if err != nil {
		return nil, core.ErrUnauthorized
	}
	return user, nil
}
