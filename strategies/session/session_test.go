package session_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"go-ez-auth/core"
	"go-ez-auth/stores"
	"go-ez-auth/strategies/session"

	"github.com/gorilla/sessions"
)

type dummyUser struct{ id string }

func (d dummyUser) GetID() string                         { return d.id }
func (d dummyUser) GetAttributes() map[string]interface{} { return nil }

func TestAuthenticate_NoSession(t *testing.T) {
	store := sessions.NewCookieStore([]byte("secret"))
	s := session.New(session.Config{Store: store, SessionName: "sess", Key: "user_id", UserStore: stores.NewInMemoryUserStore()})
	req := httptest.NewRequest("GET", "/", nil)
	_, err := s.Authenticate(context.Background(), req)
	if err != core.ErrUnauthorized {
		t.Errorf("expected ErrUnauthorized, got %v", err)
	}
}

func TestAuthenticate_InvalidUserID(t *testing.T) {
	store := sessions.NewCookieStore([]byte("secret"))
	s := session.New(session.Config{Store: store, SessionName: "sess", Key: "user_id", UserStore: stores.NewInMemoryUserStore()})
	// create initial session
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	sess, _ := store.Get(req, "sess")
	sess.Values["user_id"] = "nope"
	sess.Save(req, w)

	// attach cookie
	req2 := httptest.NewRequest("GET", "/", nil)
	for _, c := range w.Result().Cookies() {
		req2.AddCookie(c)
	}
	_, err := s.Authenticate(context.Background(), req2)
	if err != core.ErrUnauthorized {
		t.Errorf("expected ErrUnauthorized for unknown user, got %v", err)
	}
}

func TestAuthenticate_ValidSession(t *testing.T) {
	store := sessions.NewCookieStore([]byte("secret"))
	dummy := dummyUser{"u1"}
	us := stores.NewInMemoryUserStore(dummy)
	s := session.New(session.Config{Store: store, SessionName: "sess", Key: "user_id", UserStore: us})

	// set session with valid user
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	sess, _ := store.Get(req, "sess")
	sess.Values["user_id"] = "u1"
	sess.Save(req, w)

	req2 := httptest.NewRequest("GET", "/", nil)
	for _, c := range w.Result().Cookies() {
		req2.AddCookie(c)
	}
	user, err := s.Authenticate(context.Background(), req2)
	if err != nil || user.GetID() != "u1" {
		t.Fatalf("expected user u1, got %v %v", user, err)
	}
}
