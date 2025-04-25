package stores

import (
	"context"
	"go-ez-auth/core"
)

// InMemoryUserStore is a simple UserStore backed by an in-memory map.
type InMemoryUserStore struct {
	users map[string]core.User
}

// NewInMemoryUserStore creates a new store with optional initial users.
func NewInMemoryUserStore(initialUsers ...core.User) *InMemoryUserStore {
	m := make(map[string]core.User)
	for _, u := range initialUsers {
		m[u.GetID()] = u
	}
	return &InMemoryUserStore{users: m}
}

// FindUserByID retrieves a user by ID.
func (s *InMemoryUserStore) FindUserByID(ctx context.Context, id string) (core.User, error) {
	if u, ok := s.users[id]; ok {
		return u, nil
	}
	return nil, core.ErrUserNotFound
}

// FindUserByCredentials supports lookup by "id" field in criteria, else returns ErrInvalidCredentials.
func (s *InMemoryUserStore) FindUserByCredentials(ctx context.Context, criteria map[string]interface{}) (core.User, error) {
	if idVal, ok := criteria["id"].(string); ok {
		return s.FindUserByID(ctx, idVal)
	}
	return nil, core.ErrInvalidCredentials
}
