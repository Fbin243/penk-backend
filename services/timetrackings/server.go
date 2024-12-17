package main

import (
	"fmt"
	"log"
	"os"

	"tenkhours/pkg/db"
	"tenkhours/pkg/middlewares"
	"tenkhours/services/core/repo"
	fishBiz "tenkhours/services/currency/business"
	fishRepo "tenkhours/services/currency/repo"
	"tenkhours/services/timetrackings/business"
	"tenkhours/services/timetrackings/graph"
	timetrackingsRepo "tenkhours/services/timetrackings/repo"

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
	timetrackingsRepo := timetrackingsRepo.NewTimeTrackingsRepo(mongodb)
	fishRepo := fishRepo.NewFishRepo(mongodb)
	fishBusiness := fishBiz.NewFishBusiness(fishRepo, charactersRepo, profilesRepo, redisClient)
	timetrackingsBiz := business.NewTimeTrackingsBusiness(timetrackingsRepo, charactersRepo, fishRepo, fishBusiness, profilesRepo, redisClient)

	// Check authentication
	authMiddleware := middlewares.NewMiddleware(redisClient, profilesRepo)
	app.Use(authMiddleware.CheckAuth)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{TimeTrackingsBusiness: timetrackingsBiz},
	}))

	app.POST("/graphql", func(c *gin.Context) {
		srv.ServeHTTP(c.Writer, c.Request)
	})

	port, found := os.LookupEnv("TIME_TRACKINGS_PORT")
	if !found {
		port = "8082"
	}
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	app.Run(":" + port)
}
