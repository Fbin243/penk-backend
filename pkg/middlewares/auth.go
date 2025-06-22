package middlewares

import (
	"context"
	"net/http"
	"strings"

	"tenkhours/pkg/auth"

	"github.com/gin-gonic/gin"
)

// Check if the request has a valid Authorization header with a Bearer token.
func RequireAuth(ac *AuthClient) func(c *gin.Context) {
	return func(c *gin.Context) {
		reqCtx := c.Request.Context()
		authorization := c.Request.Header.Get("Authorization")
		deviceID := c.Request.Header.Get("X-Device-Id")
		userID := c.Request.Header.Get("X-User-Id")
		if userID != "" {
			// Instropect the token to get or make an auth session
			var idToken string
			if strings.HasPrefix(authorization, "Bearer ") {
				idToken = strings.Split(authorization, "Bearer ")[1]
			}

			authSession, err := ac.IntrospectUser(reqCtx, idToken, userID, deviceID)
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
