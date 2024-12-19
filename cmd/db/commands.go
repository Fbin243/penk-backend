package db

import "github.com/urfave/cli/v2"

var ImportTemplatesDataFromJSON = cli.Command{
	Name:     "import-templates",
	Category: "db",
	Usage:    "Import templates data from JSON",
	Action: func(cCtx *cli.Context) error {
		return readTemplatesFromJSON()
	},
}
