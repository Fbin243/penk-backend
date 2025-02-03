package main

import (
	"log"
	"net"
	"os"

	"tenkhours/pkg/errors"
	"tenkhours/pkg/middlewares"
	"tenkhours/pkg/pb"
	"tenkhours/services/core/composer"
	"tenkhours/services/core/transport/graph"

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

	if godotenv.Load(".env."+env) != nil {
		log.Fatal("Error loading .env." + env + " file")
	}

	app := gin.Default()
	app.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowHeaders:    []string{"Content-Type", "Authorization"},
	}))

	app.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})

	// Check authentication
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

	defer composer.GetComposer().CurrencyConn.Close()
	defer composer.GetComposer().AnalyticConn.Close()

	// Start RPC server
	go startRPCServer()

	port, found := os.LookupEnv("CORE_PORT")
	if !found {
		port = "8080"
	}

	app.Run(":" + port)
}

func startRPCServer() {
	// Create the server for gRPC API
	s := grpc.NewServer()
	pb.RegisterCoreServer(s, composer.ComposeRPCHandler())

	port, found := os.LookupEnv("CORE_GRPC_PORT")
	if !found {
		port = "50051"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
