package main

import (
	"fmt"
	"log"
	"os"

	"tenkhours/cmd/gql"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "tenk-cli",
		Usage: "say hi to tenk-cli",
		Action: func(*cli.Context) error {
			fmt.Println("hieu dep trai sieu cap vu tru")
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:     "gql-setup",
				Category: "gql",
				Usage:    "Setup a gql boilerplate",
				Action: func(cCtx *cli.Context) error {
					setupPath := cCtx.Args().First()
					if setupPath == "" {
						setupPath = "services/example"
					}

					return gql.SetupBoilerplate(setupPath)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
