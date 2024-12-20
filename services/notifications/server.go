package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"tenkhours/pkg/errors"
	"tenkhours/services/notifications/business"
	"tenkhours/services/notifications/graph"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	env := os.Getenv("TENK_ENV")
	if env == "" {
		env = "development"
	}

	if err := godotenv.Load(".env." + env); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	fmt.Println("Running in environment:", env)

	app := gin.Default()

	app.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowHeaders:    []string{"Content-Type", "Authorization"},
	}))

	notificationBiz := business.NewNotificationBusiness()

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{
			NotificationBusiness: notificationBiz,
		},
	}))
	srv.SetErrorPresenter(errors.DefaultPresenter)

	app.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})

	app.POST("/graphql", func(c *gin.Context) {
		log.Println("Received request on /graphql")
		srv.ServeHTTP(c.Writer, c.Request)
	})

	port := os.Getenv("NOTIFICATION_PORT")
	if port == "" {
		port = "8084"
	}

	if err := app.Run(":" + port); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
