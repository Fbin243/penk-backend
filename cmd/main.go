package main

import (
	"fmt"
	"log"
	"os"

	"tenkhours/cmd/auth"
	"tenkhours/cmd/gql"

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
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
