package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if godotenv.Load(".env.development") != nil {
		log.Fatal("Error loading .env file")
	}

	r := gin.Default()

	log.Println("--> Hello from Analytics service")

	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowHeaders:    []string{"Content-Type", "Authorization"},
	}))

	port, found := os.LookupEnv("ANALYTICS_PORT")
	if !found {
		port = "8081"
	}

	r.Run(":" + port)
}
