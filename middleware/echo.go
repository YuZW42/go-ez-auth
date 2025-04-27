package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go-ez-auth/core"
)

// EchoMiddleware returns an Echo middleware enforcing authentication via strategyNames.
func EchoMiddleware(strategyNames ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, err := AuthenticateRequest(strategyNames, c.Request())
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
			}
			c.Set(core.ContextUserKey, user)
			return next(c)
		}
	}
}
