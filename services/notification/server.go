package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"tenkhours/pkg/cron"
	"tenkhours/pkg/errors"
	"tenkhours/pkg/middlewares"
	"tenkhours/proto/pb/notification"
	"tenkhours/services/notification/composer"
	"tenkhours/services/notification/transport/graph"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/samber/lo"
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

	// Initialize cron scheduler
	c := cron.NewCron()
	c.RunEverySeconds(func() {
		log.Printf("Running daily notification job at %v", time.Now())
		// Sync today's reminders
		if err := composer.GetComposer().NotificationBiz.SyncTodayReminders(context.Background()); err != nil {
			log.Printf("Error syncing today's reminders: %v", err)
			return
		}
		log.Printf("Successfully synced reminders for today")
	}, 10)

	c.RunOnce(func() {
		log.Printf("Running daily notification job at %v", time.Now())
	}, lo.ToPtr(time.Now().Unix()))

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
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutting down...")
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
