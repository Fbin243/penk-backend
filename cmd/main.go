package main

import (
	"fmt"
	"log"
	"os"

	"tenkhours/cmd/auth"
	"tenkhours/cmd/gql"
	"tenkhours/cmd/test"
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
			&test.TestAPICommand,
			&time.GetTheCurrentTime,
		},
		Before: func(ctx *cli.Context) error {
			if godotenv.Load(".env.test") != nil {
				return cli.Exit("Error loading .env.test"+" file", 1)
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
