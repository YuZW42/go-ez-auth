package local_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"go-ez-auth/core"
	"go-ez-auth/strategies/local"

	"golang.org/x/crypto/bcrypt"
)

// dummyUser implements core.User for tests
type dummyUser struct{ id string }

func (d dummyUser) GetID() string                         { return d.id }
func (d dummyUser) GetAttributes() map[string]interface{} { return nil }

// dummyStore implements core.UserStore for local strategy tests
type dummyStore struct {
	hash string
	user core.User
}

// FindUserByCredentials validates credentials from criteria map
func (d dummyStore) FindUserByCredentials(ctx context.Context, criteria map[string]interface{}) (core.User, error) {
	// Extract username & password
	pass, _ := criteria["password"].(string)
	if err := bcrypt.CompareHashAndPassword([]byte(d.hash), []byte(pass)); err != nil {
		return nil, core.ErrInvalidCredentials
	}
	return d.user, nil
}

func (d dummyStore) FindUserByID(ctx context.Context, id string) (core.User, error) {
	return nil, core.ErrUserNotFound
}

func TestLocalStrategy_Authenticate_Success(t *testing.T) {
	// prepare hashed password
	hash, _ := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.DefaultCost)
	user := dummyUser{"u1"}
	store := dummyStore{hash: string(hash), user: user}
	strat := local.New(local.Config{UserStore: store})

	req := httptest.NewRequest("GET", "/", nil)
	req.SetBasicAuth("user1", "secret123") // username ignored in store

	u, err := strat.Authenticate(context.Background(), req)
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if u.GetID() != "u1" {
		t.Errorf("expected id u1, got %s", u.GetID())
	}
}

func TestLocalStrategy_Authenticate_Unauthorized(t *testing.T) {
	store := dummyStore{hash: "", user: dummyUser{"u1"}}
	strat := local.New(local.Config{UserStore: store})

	req := httptest.NewRequest("GET", "/", nil)
	req.SetBasicAuth("user1", "wrongpass")

	_, err := strat.Authenticate(context.Background(), req)
	if err != core.ErrUnauthorized {
		t.Errorf("expected ErrUnauthorized, got %v", err)
	}
}
