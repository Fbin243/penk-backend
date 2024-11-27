package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

func ReadEnvFile() {
	// Read .env file based on the environment
	env := os.Getenv("TENK_ENV")
	if env == "" {
		env = "development"
	}

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error determining current working directory: %v", err)
	}

	repoName := "tenk-backend"
	rootDir := strings.Split(currentDir, repoName)[0] + repoName

	fmt.Println("Root directory:", rootDir)

	envFilePath := filepath.Join(rootDir, ".env."+env)

	if godotenv.Load(envFilePath) != nil {
		log.Fatal("Error loading .env." + env + " file")
	}

	fmt.Println("------------------Running in environment:", env)
}
