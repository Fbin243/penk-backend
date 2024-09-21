package main

import (
	"log"
	"os"

	"tenkhours/services/timetrackings/api"

	"github.com/joho/godotenv"
)

func main() {
	if godotenv.Load(".env."+os.Getenv("TENK_ENV")) != nil {
		log.Fatal("Error loading .env file")
	}

	app := api.NewApp()
	app.InitRouter()

	port, found := os.LookupEnv("TIME_TRACKINGS_PORT")
	if !found {
		port = "8081"
	}

	app.Engine.Run(":" + port)
}
