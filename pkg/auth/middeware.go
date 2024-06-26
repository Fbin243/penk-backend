package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"tenkhours/pkg/db/coredb"

	"github.com/gin-gonic/gin"
)

type Middleware struct {
	userRepo *coredb.UsersRepo
}

func NewMiddleware(usersRepo *coredb.UsersRepo) *Middleware {
	return &Middleware{usersRepo}
}

func (m *Middleware) Authentication(c *gin.Context) {
	authKey := c.Request.Header.Get("Authorization")
	if strings.HasPrefix(authKey, "Bearer ") {
		idToken := strings.Replace(authKey, "Bearer ", "", 1)
		profile, err := GetProfileByIDToken(idToken)
		if err != nil {
			fmt.Printf("failed to get profile: %v\n", err)
			c.String(http.StatusUnauthorized, "invalid id token")
			return
		}

		user, err := m.userRepo.GetUserByFirebaseUID(profile.UID)
		if err != nil {
			fmt.Printf("failed to get user: %v\n", err)
			c.String(http.StatusUnauthorized, "user not found")
			return
		}

		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), UserKey, user))
	}
	c.Next()
}
