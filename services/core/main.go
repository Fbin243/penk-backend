package main

import (
	"log"
	"net/http"
	"os"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/core"
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

	authMiddleware := auth.NewMiddleware()

	r.Use(authMiddleware.CheckRequestBody)

	r.Use(authMiddleware.CheckAuth)

	r.POST("/graphql", func(c *gin.Context) {
		postData, ok := c.Request.Context().Value(auth.PostDataKey).(utils.GraphqlQueryData)
		if !ok {
			c.String(http.StatusBadRequest, "invalid post data")
			return
		}

		result := graphql.Do(graphql.Params{
			Context:        c.Request.Context(),
			Schema:         core.InitSchema(),
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
