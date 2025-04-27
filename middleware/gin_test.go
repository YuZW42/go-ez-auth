package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go-ez-auth/core"
	"go-ez-auth/middleware"
	"go-ez-auth/stores"
	"go-ez-auth/strategies/apikey"
)

type dummyUserGin struct{ id string }

func (d dummyUserGin) GetID() string { return d.id }
func (d dummyUserGin) GetAttributes() map[string]interface{} { return nil }

func TestGinMiddleware_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	e := gin.New()
	e.Use(middleware.GinMiddleware("apikey"))
	e.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

func TestGinMiddleware_Authorized(t *testing.T) {
	dummy := dummyUserGin{"u1"}
	store := stores.NewAPIKeyStore(map[string]core.User{"key": dummy})
	core.RegisterStrategy(apikey.New(apikey.Config{Store: store}))

	gin.SetMode(gin.TestMode)
	e := gin.New()
	e.Use(middleware.GinMiddleware("apikey"))
	e.GET("/", func(c *gin.Context) {
		u, exists := c.Get(core.ContextUserKey)
		if !exists {
			t.Fatal("user not found in context")
		}
		user := u.(core.User)
		c.String(http.StatusOK, user.GetID())
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
