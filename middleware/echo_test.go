package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"go-ez-auth/core"
	"go-ez-auth/middleware"
	"go-ez-auth/stores"
	"go-ez-auth/strategies/apikey"
)

// dummyUserEcho implements core.User for tests.
type dummyUserEcho struct{ id string }

func (d dummyUserEcho) GetID() string                         { return d.id }
func (d dummyUserEcho) GetAttributes() map[string]interface{} { return nil }

func TestEchoMiddleware_Unauthorized(t *testing.T) {
	e := echo.New()
	e.Use(middleware.EchoMiddleware("apikey"))
	e.GET("/", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

func TestEchoMiddleware_Authorized(t *testing.T) {
	dummy := dummyUserEcho{"u1"}
	store := stores.NewAPIKeyStore(map[string]core.User{"key": dummy})
	core.RegisterStrategy(apikey.New(apikey.Config{Store: store}))

	e := echo.New()
	e.Use(middleware.EchoMiddleware("apikey"))
	e.GET("/", func(c echo.Context) error {
		u := c.Get(core.ContextUserKey).(core.User)
		return c.String(http.StatusOK, u.GetID())
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-API-Key", "key")
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if rec.Body.String() != "u1" {
		t.Errorf("expected body 'u1', got '%s'", rec.Body.String())
	}
}
