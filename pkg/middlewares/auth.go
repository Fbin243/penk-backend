package middlewares

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"tenkhours/pkg/auth"
	"tenkhours/services/core/repo"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type Middleware struct {
	redisClient *redis.Client
}

func NewMiddleware(redisClient *redis.Client) *Middleware {
	return &Middleware{redisClient}
}

// Check if the request has a valid Authorization header with a Bearer token.
func (m *Middleware) CheckAuth(c *gin.Context) {
	reqCtx := c.Request.Context()
	authKey := c.Request.Header.Get("Authorization")
	if strings.HasPrefix(authKey, "Bearer ") {
		idToken := strings.Replace(authKey, "Bearer ", "", 1)
		firebaseProfile, err := auth.GetProfileByIDToken(idToken)
		if err != nil {
			log.Printf("invalid id token: %v\n", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid id token"})
			return
		}

		reqCtx = context.WithValue(reqCtx, auth.FirebaseProfileKey, *firebaseProfile)
		// Check profile in Redis if user is authorized
		profile := &repo.Profile{}
		profileJSON, err := m.redisClient.Get(reqCtx, firebaseProfile.UID).Result()
		if err != nil && err != redis.Nil {
			log.Printf("error getting profile from redis: %v\n", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error getting profile from redis"})
			return
		}

		if err == nil {
			err = json.Unmarshal([]byte(profileJSON), profile)
			if err != nil {
				return
			}

			// Set the profile in the context
			reqCtx = context.WithValue(reqCtx, auth.ProfileKey, *profile)
		}

		c.Request = c.Request.WithContext(reqCtx)
		c.Next()
		return
	}

	c.Next()
}
