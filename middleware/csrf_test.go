package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/csrf"
	"go-ez-auth/middleware"
)

func TestCSRFMiddleware_GeneratesToken(t *testing.T) {
	authKey := []byte("01234567890123456789012345678901")
	mw := middleware.CSRFMiddleware(authKey)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// CSRF token should be in context
		token := csrf.Token(r)
		if token == "" {
			t.Fatal("expected CSRF token in request context")
		}
		w.WriteHeader(http.StatusOK)
	})
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	mw(handler).ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", rr.Code)
	}
	// Should set cookie
	if len(rr.Result().Cookies()) == 0 {
		t.Error("expected CSRF cookie to be set")
	}
}
