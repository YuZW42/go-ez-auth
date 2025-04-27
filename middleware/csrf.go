package middleware

import (
	"net/http"

	"github.com/gorilla/csrf"
)

// CSRFMiddleware returns a net/http middleware that protects against CSRF attacks.
// authKey must be 32-bytes. Pass additional csrf.Options as needed (e.g., CookieName, Secure).
func CSRFMiddleware(authKey []byte, opts ...csrf.Option) func(http.Handler) http.Handler {
	// default options: HttpOnly cookie
	defaultOpts := []csrf.Option{
		csrf.HttpOnly(true),
	}
	allOpts := append(defaultOpts, opts...)
	return csrf.Protect(authKey, allOpts...)
}
