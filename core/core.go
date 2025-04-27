package core

import (
	"context"
	"errors"
	"net/http"
)

// Strategy defines the methods for authentication strategies.
type Strategy interface {
	Name() string
	Setup() error
	Authenticate(ctx context.Context, r *http.Request) (User, error)
}

// User represents an authenticated user.
type User interface {
	GetID() string
	GetAttributes() map[string]interface{}
}

// UserStore defines methods for user lookup.
type UserStore interface {
	FindUserByID(ctx context.Context, id string) (User, error)
	// FindUserByCredentials can be used for credential-based lookups (e.g., username/password, API keys).
	FindUserByCredentials(ctx context.Context, criteria map[string]interface{}) (User, error)
}

// strategyRegistry holds registered authentication strategies.
var strategyRegistry = make(map[string]Strategy)

// RegisterStrategy registers a new authentication strategy by name.
func RegisterStrategy(s Strategy) {
	strategyRegistry[s.Name()] = s
}

// GetStrategy retrieves a registered strategy by name.
func GetStrategy(name string) (Strategy, bool) {
	s, ok := strategyRegistry[name]
	return s, ok
}

// ListStrategies returns the names of all registered strategies.
func ListStrategies() []string {
	names := make([]string, 0, len(strategyRegistry))
	for name := range strategyRegistry {
		names = append(names, name)
	}
	return names
}

// ContextUserKey is the context key for storing the authenticated User.
const ContextUserKey = "go-ez-auth-user"

// UserFromContext retrieves the authenticated User from context.
func UserFromContext(ctx context.Context) (User, bool) {
	user, ok := ctx.Value(ContextUserKey).(User)
	return user, ok
}

// Standard error variables.
var (
	ErrUnauthorized       = errors.New("unauthorized")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
)
