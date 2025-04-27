package middleware

import (
	"context"
	"net/http"

	"go-ez-auth/core"
)

// AuthenticateRequest tries each named strategy in order and returns the first successful user.
func AuthenticateRequest(strategyNames []string, r *http.Request) (core.User, error) {
	for _, name := range strategyNames {
		strat, ok := core.GetStrategy(name)
		if !ok {
			continue
		}
		user, err := strat.Authenticate(r.Context(), r)
		if err == nil {
			return user, nil
		}
	}
	return nil, core.ErrUnauthorized
}

// Middleware returns a net/http middleware that enforces authentication.
func Middleware(strategyNames ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := AuthenticateRequest(strategyNames, r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), core.ContextUserKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
