package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"tenkhours/pkg/errors"
	"tenkhours/pkg/middlewares"
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

	// Create a context that will be canceled on shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize the Composer to get all dependencies
	comp := composer.GetComposer()

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

	authClient, conn := middlewares.ComposeAuthClient()
	defer conn.Close()
	app.Use(middlewares.RequireAuth(authClient))

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

	// Start the Reminder Cron job in a separate goroutine
	go comp.ReminderCron.Start()

	// Start the Kafka Consumer in a separate goroutine
	if err := comp.NotificationQueue.Start(ctx); err != nil {
		log.Printf("Failed to start notification queue: %v", err)
	}

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start the server in a goroutine
	go func() {
		port, found := os.LookupEnv("NOTIFICATION_PORT")
		if !found {
			port = "8084"
		}

		if err := app.Run(":" + port); err != nil {
			log.Printf("Failed to run server: %v", err)
			cancel()
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutting down...")

	// Cancel context to stop all goroutines
	cancel()

	// Close Kafka connections
	if err := comp.NotificationQueue.Close(); err != nil {
		log.Printf("Error closing notification queue: %v", err)
	}

	log.Println("Shutdown complete")
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
