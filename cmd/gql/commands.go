package gql

import (
	"github.com/urfave/cli/v2"
)

var SetupBoilerplateCommand = cli.Command{
	Name:     "gql-setup",
	Category: "gql",
	Usage:    "Setup a gql boilerplate",
	Action: func(cCtx *cli.Context) error {
		setupPath := cCtx.Args().First()
		if setupPath == "" {
			setupPath = "services/example"
		}

		return SetupBoilerplate(setupPath)
	},
}
