package stores_test

import (
  "context"
  "testing"

  "go-ez-auth/core"
  "go-ez-auth/stores"
)

type dummyUser struct{ id string }
func (d dummyUser) GetID() string { return d.id }
func (d dummyUser) GetAttributes() map[string]interface{} { return nil }

func TestInMemoryUserStore(t *testing.T) {
  u1 := dummyUser{"u1"}
  s := stores.NewInMemoryUserStore(u1)

  // Find existing
  got, err := s.FindUserByID(context.Background(), "u1")
  if err != nil || got.GetID() != "u1" {
    t.Fatalf("expected u1, got %v %v", got, err)
  }

  // Not found
  if _, err := s.FindUserByID(context.Background(), "nope"); err != core.ErrUserNotFound {
    t.Errorf("expected ErrUserNotFound, got %v", err)
  }

  // Credentials lookup (valid)
  got2, err := s.FindUserByCredentials(context.Background(), map[string]interface{}{"id": "u1"})
  if err != nil || got2.GetID() != "u1" {
    t.Errorf("credentials lookup failed: %v %v", got2, err)
  }

  // Credentials lookup (invalid)
  if _, err := s.FindUserByCredentials(context.Background(), map[string]interface{}{"foo": "bar"}); err != core.ErrInvalidCredentials {
    t.Errorf("expected ErrInvalidCredentials, got %v", err)
  }
}