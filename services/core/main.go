package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/dummy"
	"tenkhours/pkg/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/joho/godotenv"
)

func main() {
	if godotenv.Load(".env."+os.Getenv("TENK_ENV")) != nil {
		log.Fatal("Error loading .env file")
	}

	r := gin.Default()

	log.Println("--> Hello from Core service")

	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowHeaders:    []string{"Content-Type", "Authorization"},
	}))

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Core service is running!")
	})

	r.POST("/graphql", func(c *gin.Context) {
		var postData utils.GraphqlQueryData
		if err := json.NewDecoder(c.Request.Body).Decode(&postData); err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		gqlCtx := context.TODO()

		authKey := c.Request.Header.Get("Authorization")
		if strings.HasPrefix(authKey, "Bearer ") {
			idToken := strings.Replace(authKey, "Bearer ", "", 1)

			profile, err := auth.GetProfileByIDToken(idToken)
			if err != nil {
				c.String(http.StatusUnauthorized, "invalid id token")
				return
			}

			gqlCtx = context.WithValue(gqlCtx, auth.ProfileContextKey, profile)
		}

		result := graphql.Do(graphql.Params{
			Context:        gqlCtx,
			Schema:         dummy.DummySchema,
			RequestString:  postData.Query,
			VariableValues: postData.Variables,
			OperationName:  postData.Operation,
		})

		c.JSON(http.StatusOK, result)
	})

	port, found := os.LookupEnv("CORE_PORT")
	if !found {
		port = "8080"
	}

	r.Run(":" + port)
}
