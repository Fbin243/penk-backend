package main

import (
	"log"
	"net"
	"os"

	"tenkhours/pkg/errors"
	"tenkhours/proto/pb/notification"
	"tenkhours/services/notification/composer"
	"tenkhours/services/notification/transport/graph"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	env := os.Getenv("TENK_ENV")
	if env == "" {
		env = "development"
	}

	if err := godotenv.Load(".env." + env); err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	app := gin.Default()

	app.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowHeaders:    []string{"Content-Type", "Authorization", "X-Device-Id", "X-User-Id"},
	}))

	app.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})

	// Init GraphQL server
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: composer.ComposeGraphQLResolver(),
	}))
	srv.SetErrorPresenter(errors.DefaultPresenter)
	app.POST("/graphql", func(c *gin.Context) {
		srv.ServeHTTP(c.Writer, c.Request)
	})

	// Start RPC server
	go startRPCServer()

	port, found := os.LookupEnv("NOTIFICATION_PORT")
	if !found {
		port = "8084"
	}

	if err := app.Run(":" + port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

func startRPCServer() {
	port, found := os.LookupEnv("NOTIFICATION_GRPC_PORT")
	if !found {
		port = "50054"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Can not create gRPC: %v", err)
	}

	s := grpc.NewServer()
	notification.RegisterNotificationServer(s, composer.ComposeRPCHandler())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("error running gRPC server: %v", err)
	}
}
