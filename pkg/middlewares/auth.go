package middlewares

import (
	"context"
	"net/http"
	"strings"

	"tenkhours/pkg/auth"

	"github.com/gin-gonic/gin"
)

// Check if the request has a valid Authorization header with a Bearer token.
func RequireAuth(ac *authClient) func(c *gin.Context) {
	return func(c *gin.Context) {
		reqCtx := c.Request.Context()
		authKey := c.Request.Header.Get("Authorization")
		if strings.HasPrefix(authKey, "Bearer ") {
			idToken := strings.Replace(authKey, "Bearer ", "", 1)

			// Instropect the token to get or make an auth session
			authSession, err := ac.IntrospectToken(reqCtx, idToken)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				return
			}

			if authSession != nil {
				reqCtx = context.WithValue(reqCtx, auth.AuthSessionKey, *authSession)
			}

			c.Request = c.Request.WithContext(reqCtx)
			c.Next()
			return
		}

		c.Next()
	}
}
