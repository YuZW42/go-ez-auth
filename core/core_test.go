package core_test

import (
	"context"
	"net/http"
	"testing"

	"go-ez-auth/core"
)

type dummyStrategy struct{}

func (d dummyStrategy) Name() string { return "dummy" }
func (d dummyStrategy) Setup() error { return nil }
func (d dummyStrategy) Authenticate(ctx context.Context, r *http.Request) (core.User, error) {
	return nil, nil
}

func TestStrategyRegistry(t *testing.T) {
	// Register dummy strategy
	core.RegisterStrategy(dummyStrategy{})

	// Retrieve and verify
	s, ok := core.GetStrategy("dummy")
	if !ok {
		t.Fatal("expected strategy 'dummy' to be registered")
	}
	if s.Name() != "dummy" {
		t.Errorf("expected name 'dummy', got %s", s.Name())
	}

	// List strategies
	names := core.ListStrategies()
	if len(names) != 1 || names[0] != "dummy" {
		t.Errorf("expected names ['dummy'], got %v", names)
	}
}
