package main

import (
	"fmt"
	"log"
	"os"

	"tenkhours/pkg/db"
	"tenkhours/pkg/middlewares"
	analyticsRepo "tenkhours/services/analytics/repo"
	"tenkhours/services/core/business"
	"tenkhours/services/core/graph"
	"tenkhours/services/core/repo"

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
	mongodb := db.GetDBManager().DB
	redisClient := db.GetRedisClient()
	profilesRepo := repo.NewProfilesRepo(mongodb, redisClient)
	charactersRepo := repo.NewCharactersRepo(mongodb)
	goalsRepo := repo.NewGoalsRepo(mongodb)

	// TODO: Temporary inject analyticsRepos into profilesBiz for deleting related data
	capturedRepordsRepo := analyticsRepo.NewCapturedRecordsRepo(mongodb)
	snapshotsRepo := analyticsRepo.NewSnapshotsRepo(mongodb)

	profilesBiz := business.NewProfilesBusiness(profilesRepo, charactersRepo, capturedRepordsRepo, snapshotsRepo, redisClient)
	charactersBiz := business.NewCharactersBusiness(charactersRepo, profilesRepo)
	goalsBiz := business.NewGoalsBusiness(goalsRepo, charactersRepo)

	// Check authentication
	authMiddleware := middlewares.NewMiddleware(redisClient, profilesRepo)
	app.Use(authMiddleware.CheckAuth)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{
			ProfilesBusiness:   profilesBiz,
			CharactersBusiness: charactersBiz,
			GoalsBusiness:      goalsBiz,
		},
	}))

	app.POST("/graphql", func(c *gin.Context) {
		srv.ServeHTTP(c.Writer, c.Request)
	})

	port, found := os.LookupEnv("CORE_PORT")
	if !found {
		port = "8080"
	}
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	app.Run(":" + port)
}
