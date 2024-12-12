package middlewares

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db"
	"tenkhours/pkg/utils"
	"tenkhours/services/core/repo"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

type Middleware struct {
	redisClient  *redis.Client
	profilesRepo *repo.ProfilesRepo
}

func NewMiddleware(redisClient *redis.Client, profilesRepo *repo.ProfilesRepo) *Middleware {
	return &Middleware{redisClient, profilesRepo}
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

		// Check if there is any active session in Redis
		keyFound, err := m.redisClient.Exists(context.Background(), firebaseProfile.UID).Result()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "key not found in redis"})
			return
		}

		// Cache hit, return the profile from redis
		if keyFound == 1 {
			profileJSON, err := m.redisClient.Get(context.Background(), firebaseProfile.UID).Result()
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to get profile from redis"})
				return
			}

			var profile repo.Profile
			err = json.Unmarshal([]byte(profileJSON), &profile)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to unmarshal profile"})
				return
			}

			reqCtx = context.WithValue(reqCtx, auth.ProfileKey, profile)
		} else {
			// Cache miss, create a new session from the profile in DB
			profile, err := m.profilesRepo.GetProfileByFirebaseUID(firebaseProfile.UID)
			if err == mongo.ErrNoDocuments {
				// profile not found, mean the new account
				newProfile := repo.Profile{
					BaseModel:              &db.BaseModel{},
					Name:                   firebaseProfile.Name,
					Email:                  firebaseProfile.Email,
					FirebaseUID:            firebaseProfile.UID,
					ImageURL:               firebaseProfile.Picture,
					AutoSnapshot:           true,
					AvailableSnapshots:     utils.DefaultSnapshotsNumber,
					LimitedCharacterNumber: utils.LimitedCharacterNumber,
				}

				// Create new profile for the new user in DB
				profile, err = m.profilesRepo.InsertOne(&newProfile)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to create new profile"})
					return
				}
			} else if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to get profile by firebase UID"})
				return
			}

			// Save profile in redis
			profileJSON, err := json.Marshal(profile)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to serialize profile"})
				return
			}

			err = m.redisClient.Set(context.Background(), firebaseProfile.UID, profileJSON, time.Hour).Err()
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to save profile in redis"})
				return
			}

			reqCtx = context.WithValue(reqCtx, auth.ProfileKey, *profile)
		}

		c.Request = c.Request.WithContext(reqCtx)
		c.Next()
		return
	}

	c.Next()
}
