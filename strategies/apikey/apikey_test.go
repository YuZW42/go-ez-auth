package apikey_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-ez-auth/core"
	"go-ez-auth/stores"
	"go-ez-auth/strategies/apikey"
)

type dummyUser struct{ id, key string }
func (d dummyUser) GetID() string                 { return d.id }
func (d dummyUser) GetAttributes() map[string]interface{} { return nil }

func TestAuthenticate_NoKey(t *testing.T) {
	s := apikey.New(apikey.Config{Store: stores.NewInMemoryUserStore()})
	req := httptest.NewRequest("GET", "/", nil)
	_, err := s.Authenticate(context.Background(), req)
	if err != core.ErrUnauthorized {
		t.Errorf("expected ErrUnauthorized, got %v", err)
	}
}

func TestAuthenticate_Header(t *testing.T) {
	dummy := dummyUser{"u1", "key123"}
	s := apikey.New(apikey.Config{Store: stores.NewInMemoryUserStore(dummy), CredKey: "id"})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-API-Key", "key123")
	user, err := s.Authenticate(context.Background(), req)
	if err != nil || user.GetID() != "u1" {
		t.Fatalf("expected u1, got %v %v", user, err)
	}
}

func TestAuthenticate_QueryParam(t *testing.T) {
	dummy := dummyUser{"u2", "key456"}
	s := apikey.New(apikey.Config{Store: stores.NewInMemoryUserStore(dummy), CredKey: "id", QueryParam: "api_key"})

	req := httptest.NewRequest("GET", "/?api_key=key456", nil)
	user, err := s.Authenticate(context.Background(), req)
	if err != nil || user.GetID() != "u2" {
		t.Fatalf("expected u2, got %v %v", user, err)
	}
}
