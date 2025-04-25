package stores

import (
	"context"

	"go-ez-auth/core"
)

// APIKeyStore maps API keys to users.
type APIKeyStore struct {
	byKey map[string]core.User
}

// NewAPIKeyStore creates a store with the given key->User mapping.
func NewAPIKeyStore(mapping map[string]core.User) *APIKeyStore {
	m := make(map[string]core.User)
	for k, u := range mapping {
		m[k] = u
	}
	return &APIKeyStore{byKey: m}
}

// FindUserByID looks up a user by user.ID across all mapped values.
func (s *APIKeyStore) FindUserByID(ctx context.Context, id string) (core.User, error) {
	for _, u := range s.byKey {
		if u.GetID() == id {
			return u, nil
		}
	}
	return nil, core.ErrUserNotFound
}

// FindUserByCredentials looks for any string value in criteria matching a stored key.
func (s *APIKeyStore) FindUserByCredentials(ctx context.Context, criteria map[string]interface{}) (core.User, error) {
	for _, v := range criteria {
		key, ok := v.(string)
		if !ok {
			continue
		}
		if u, exists := s.byKey[key]; exists {
			return u, nil
		}
	}
	return nil, core.ErrInvalidCredentials
}
