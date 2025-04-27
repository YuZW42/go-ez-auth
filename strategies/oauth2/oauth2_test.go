package oauth2_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-ez-auth/core"
	authoauth "go-ez-auth/strategies/oauth2"

	"golang.org/x/oauth2"
)

// dummyUser implements core.User for tests
type dummyUser struct{ id string }

func (d dummyUser) GetID() string                         { return d.id }
func (d dummyUser) GetAttributes() map[string]interface{} { return map[string]interface{}{"id": d.id} }

func TestOAuth2Strategy_Authenticate_Success(t *testing.T) {
	// Setup test server
	mux := http.NewServeMux()
	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		code := r.FormValue("code")
		if code != "testcode" {
			http.Error(w, "bad code", http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"access_token": "testtoken", "token_type": "Bearer"})
	})
	mux.HandleFunc("/userinfo", func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer testtoken" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"id": "u1"})
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	// Config
	oauthConfig := &oauth2.Config{
		ClientID:     "id",
		ClientSecret: "secret",
		Endpoint: oauth2.Endpoint{
			AuthURL:  server.URL + "/auth",
			TokenURL: server.URL + "/token",
		},
		RedirectURL: server.URL + "/callback",
	}
	cfg := authoauth.Config{
		OAuth2Config: oauthConfig,
		UserInfoURL:  server.URL + "/userinfo",
		ExtractUser: func(ctx context.Context, info map[string]interface{}) (core.User, error) {
			id, _ := info["id"].(string)
			return dummyUser{id}, nil
		},
	}
	strat := authoauth.New(cfg)

	// Simulate request with code
	req := httptest.NewRequest("GET", "/?code=testcode", nil)
	user, err := strat.Authenticate(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.GetID() != "u1" {
		t.Errorf("expected user u1, got %s", user.GetID())
	}
}

func TestOAuth2Strategy_Authenticate_NoCode(t *testing.T) {
	cfg := authoauth.Config{}
	strat := authoauth.New(cfg)
	req := httptest.NewRequest("GET", "/", nil)
	_, err := strat.Authenticate(context.Background(), req)
	if err != core.ErrUnauthorized {
		t.Errorf("expected ErrUnauthorized, got %v", err)
	}
}
