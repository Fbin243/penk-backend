package auth

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"tenkhours/pkg/db/coredb"
	"tenkhours/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Middleware struct {
	profilesRepo *coredb.ProfilesRepo
}

func NewMiddleware(profilesRepo *coredb.ProfilesRepo) *Middleware {
	return &Middleware{profilesRepo}
}

func (m *Middleware) CheckRequestBody(c *gin.Context) {
	var postData utils.GraphqlQueryData
	if err := json.NewDecoder(c.Request.Body).Decode(&postData); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	reqCtx := c.Request.Context()
	reqCtx = context.WithValue(reqCtx, PostDataKey, postData)
	c.Request = c.Request.WithContext(reqCtx)
	c.Next()
}

func (m *Middleware) CheckAuth(c *gin.Context) {
	reqCtx := c.Request.Context()
	authKey := c.Request.Header.Get("Authorization")
	if strings.HasPrefix(authKey, "Bearer ") {
		idToken := strings.Replace(authKey, "Bearer ", "", 1)
		firebaseProfile, err := GetProfileByIDToken(idToken)
		if err != nil {
			log.Printf("invalid id token: %v\n", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid id token"})
			return
		}

		profile, err := m.profilesRepo.GetProfileByFirebaseUID(firebaseProfile.UID)
		if err == mongo.ErrNoDocuments {
			log.Printf("user has not registered profile, so register it\n")
			newProfile := coredb.Profile{
				ID:                 primitive.NewObjectID(),
				Name:               firebaseProfile.Name,
				Email:              firebaseProfile.Email,
				FirebaseUID:        firebaseProfile.UID,
				ImageURL:           "",
				CreatedAt:          utils.Now(),
				UpdatedAt:          utils.Now(),
				AutoSnapshot:       true,
				AvailableSnapshots: 2,
			}

			createdProfile, err := m.profilesRepo.CreateNewProfile(&newProfile)
			if err != nil {
				log.Printf("failed to insert new profile: %v\n", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, "failed to register new profile")
			}

			c.Request = c.Request.WithContext(context.WithValue(reqCtx, ProfileKey, *createdProfile))
			c.Next()
			return
		}

		c.Request = c.Request.WithContext(context.WithValue(reqCtx, ProfileKey, *profile))
		c.Next()
		return
	}

	c.Next()
}
