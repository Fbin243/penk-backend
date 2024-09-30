package main

import (
	"fmt"
	"log"
	"os"

	"tenkhours/pkg/analytics"
	"tenkhours/pkg/auth"
	"tenkhours/pkg/db"
	"tenkhours/pkg/db/analyticsdb"
	"tenkhours/pkg/db/coredb"
	"tenkhours/services/analytics_v2/graph"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	env := os.Getenv("TENK_ENV")
	if env == "" {
		env = "development"
	}

	if godotenv.Load(".env."+env) != nil {
		log.Fatal("Error loading .env." + env + " file")
	}

	fmt.Println("------------------Running in environment:", env)

	app := gin.Default()
	app.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowHeaders:    []string{"Content-Type", "Authorization"},
	}))

	// Init dependencies and perform DI manually
	db := db.GetDBManager().DB
	profilesRepo := coredb.NewProfilesRepo(db)
	charactersRepo := coredb.NewCharactersRepo(db)
	snapshotsRepo := analyticsdb.NewSnapshotRepo(db)
	charactersHandler := analytics.NewCharactersHandler(snapshotsRepo, charactersRepo, profilesRepo)

	// Check authentication
	authMiddleware := auth.NewMiddleware(profilesRepo)
	app.Use(authMiddleware.CheckAuth)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{
			CharactersHandler: charactersHandler,
		},
	}))

	app.GET("/", func(c *gin.Context) {
		playgroundHandler := playground.Handler("GraphQL playground", "/graphql")
		playgroundHandler.ServeHTTP(c.Writer, c.Request)
	})

	app.POST("/graphql", func(c *gin.Context) {
		srv.ServeHTTP(c.Writer, c.Request)
	})

	port, found := os.LookupEnv("ANALYTICS_PORT")
	if !found {
		port = "8082"
	}
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	app.Run(":" + port)
}
