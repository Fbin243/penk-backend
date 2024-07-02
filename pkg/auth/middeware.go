package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"tenkhours/pkg/db/coredb"
	"tenkhours/pkg/utils"

	"github.com/gin-gonic/gin"
)

type Middleware struct {
	userRepo *coredb.UsersRepo
}

func NewMiddleware() *Middleware {
	return &Middleware{userRepo: coredb.NewUsersRepo()}
}

func (m *Middleware) CheckRequestBody(c *gin.Context) {
	var postData utils.GraphqlQueryData
	if err := json.NewDecoder(c.Request.Body).Decode(&postData); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	reqCtx := c.Request.Context()
	if strings.Contains(postData.Query, "mutation") && strings.Contains(postData.Query, "registerAccount") {
		reqCtx = context.WithValue(reqCtx, RegisterKey, true)
	} else {
		reqCtx = context.WithValue(reqCtx, RegisterKey, false)
	}

	reqCtx = context.WithValue(reqCtx, PostDataKey, postData)

	c.Request = c.Request.WithContext(reqCtx)
	c.Next()
}

func (m *Middleware) CheckAuth(c *gin.Context) {
	reqCtx := c.Request.Context()
	wantToRegister := reqCtx.Value(RegisterKey).(bool)

	authKey := c.Request.Header.Get("Authorization")
	if strings.HasPrefix(authKey, "Bearer ") {
		idToken := strings.Replace(authKey, "Bearer ", "", 1)
		profile, err := GetProfileByIDToken(idToken)
		if err != nil {
			fmt.Printf("failed to get profile: %v\n", err)
			c.String(http.StatusUnauthorized, "invalid id token")
			return
		}

		if wantToRegister {
			c.Request = c.Request.WithContext(context.WithValue(reqCtx, ProfileKey, profile))
			c.Next()
			return
		}

		user, err := m.userRepo.GetUserByFirebaseUID(profile.UID)
		if err != nil {
			fmt.Printf("failed to get user: %v\n", err)
			c.String(http.StatusUnauthorized, "user not found")
			return
		}

		c.Request = c.Request.WithContext(context.WithValue(reqCtx, UserKey, user))
	}
	c.Next()
}
