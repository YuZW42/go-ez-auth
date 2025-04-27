package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go-ez-auth/core"
	"go-ez-auth/middleware"
	"go-ez-auth/stores"
	"go-ez-auth/strategies/apikey"
)

type dummyUserNet struct{ id string }

func (d dummyUserNet) GetID() string                         { return d.id }
func (d dummyUserNet) GetAttributes() map[string]interface{} { return nil }

func TestMiddleware_Unauthorized(t *testing.T) {
	// Register apikey strategy with an empty store
	core.RegisterStrategy(apikey.New(apikey.Config{Store: stores.NewAPIKeyStore(nil)}))
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mw := middleware.Middleware("apikey")
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	mw(handler).ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rr.Code)
	}
}

func TestMiddleware_Authorized(t *testing.T) {
	dummy := dummyUserNet{"u1"}
	store := stores.NewAPIKeyStore(map[string]core.User{"key": dummy})
	core.RegisterStrategy(apikey.New(apikey.Config{Store: store}))

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := core.UserFromContext(r.Context())
		if !ok {
			t.Fatal("user not found in context")
		}
		w.Write([]byte(user.GetID()))
	})
	mw := middleware.Middleware("apikey")

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-API-Key", "key")
	mw(handler).ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if rr.Body.String() != "u1" {
		t.Errorf("expected body 'u1', got '%s'", rr.Body.String())
	}
}
