package main

import (
	"fmt"
	"log"
	"os"

	"tenkhours/services/core/api"

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

	app := api.NewApp()
	app.InitRouter()

	port, found := os.LookupEnv("CORE_PORT")
	if !found {
		port = "8080"
	}

	app.Engine.Run(":" + port)
}
