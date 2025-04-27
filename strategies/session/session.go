package session

import (
	"context"
	"net/http"

	"go-ez-auth/core"

	"github.com/gorilla/sessions"
)

// Config holds settings for the Session strategy.
type Config struct {
	Store       sessions.Store // Gorilla sessions store
	SessionName string         // name of the session (cookie)
	Key         string         // key in session.Values for user ID
	UserStore   core.UserStore // backend to lookup users
}

// Strategy implements core.Strategy for session-based auth.
type Strategy struct {
	config Config
}

// New returns a configured Session strategy.
func New(config Config) *Strategy {
	return &Strategy{config: config}
}

// Name returns the strategy name.
func (s *Strategy) Name() string {
	return "session"
}

// Setup configures cookie options for secure sessions if using CookieStore.
func (s *Strategy) Setup() error {
	// Apply secure defaults to CookieStore
	if cs, ok := s.config.Store.(*sessions.CookieStore); ok {
		cs.Options = &sessions.Options{
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		}
	}
	return nil
}

// Authenticate retrieves the session, extracts user ID, and looks up the user.
func (s *Strategy) Authenticate(ctx context.Context, r *http.Request) (core.User, error) {
	sess, err := s.config.Store.Get(r, s.config.SessionName)
	if err != nil {
		return nil, core.ErrUnauthorized
	}

	raw, ok := sess.Values[s.config.Key].(string)
	if !ok || raw == "" {
		return nil, core.ErrUnauthorized
	}

	user, err := s.config.UserStore.FindUserByID(ctx, raw)
	if err != nil {
		return nil, core.ErrUnauthorized
	}
	return user, nil
}

// Login saves the user's ID in a new session and protects against fixation by issuing a fresh cookie.
func (s *Strategy) Login(w http.ResponseWriter, r *http.Request, user core.User) error {
	sess, err := s.config.Store.Get(r, s.config.SessionName)
	if err != nil {
		return core.ErrUnauthorized
	}
	// Set user ID in session
	sess.Values[s.config.Key] = user.GetID()
	// Save writes a new cookie with updated data
	return sess.Save(r, w)
}
