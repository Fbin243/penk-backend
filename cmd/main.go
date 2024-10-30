package main

import (
	"fmt"
	"log"
	"os"

	"tenkhours/cmd/auth"
	"tenkhours/cmd/e2e"
	"tenkhours/cmd/gql"
	"tenkhours/cmd/time"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "tenk-cli",
		Usage: "say hi to tenk-cli",
		Action: func(*cli.Context) error {
			fmt.Println("Welcome to Tenk Hours CLI")
			return nil
		},
		Commands: []*cli.Command{
			&gql.SetupBoilerplateCommand,
			&auth.GetJWTTokenCommand,
			&e2e.TestUserFlowCommand,
			&time.GetTheCurrentTime,
		},
		Before: func(ctx *cli.Context) error {
			env := os.Getenv("TENK_ENV")
			if env == "" {
				env = "development"
			}

			if godotenv.Load(".env."+env) != nil {
				return cli.Exit("Error loading .env."+env+" file", 1)
			}

			fmt.Println("------------------Running in environment:", env)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
