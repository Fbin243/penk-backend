package main

import (
	"log"
	"tenkhours/pkg/db"

	"github.com/joho/godotenv"
)

func main() {
	if godotenv.Load(".env.development") != nil {
		log.Fatal("Error loading .env file")
	}

	log.Println("--> Hello from Core service")

	log.Println("Connected to DB", db.GetDB().Name())
}
