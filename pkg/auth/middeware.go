package auth

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"tenkhours/pkg/db"
	"tenkhours/pkg/db/coredb"
	"tenkhours/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Middleware struct {
	userRepo *coredb.UsersRepo
}

func NewMiddleware() *Middleware {
	return &Middleware{userRepo: coredb.NewUsersRepo(db.GetDB())}
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
		profile, err := GetProfileByIDToken(idToken)
		if err != nil {
			log.Printf("invalid id token: %v\n", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid id token"})
			return
		}

		user, err := m.userRepo.GetUserByFirebaseUID(profile.UID)
		if err != nil {
			log.Printf("user has not registered, so register it\n")
			newUser := coredb.User{
				ID:          primitive.NewObjectID(),
				Name:        profile.Name,
				Email:       profile.Email,
				FirebaseUID: profile.UID,
				ImageURL:    "",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}

			createdUser, err := m.userRepo.CreateNewUser(&newUser)
			if err != nil {
				log.Printf("failed to insert user: %v\n", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, "failed to register new user")
			}

			c.Request = c.Request.WithContext(context.WithValue(reqCtx, UserKey, *createdUser))
			c.Next()
			return
		}

		c.Request = c.Request.WithContext(context.WithValue(reqCtx, UserKey, *user))
		c.Next()
		return
	}

	c.Next()
}
