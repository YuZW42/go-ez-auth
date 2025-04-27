package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-ez-auth/core"
)

// GinMiddleware returns a gin.HandlerFunc that enforces authentication using given strategies.
func GinMiddleware(strategyNames ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := AuthenticateRequest(strategyNames, c.Request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.Set(core.ContextUserKey, user)
		c.Next()
	}
}
