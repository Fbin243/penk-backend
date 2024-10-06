package e2e

import (
	"fmt"
	"log"

	"github.com/urfave/cli/v2"
)

var TestUserFlowCommand = cli.Command{
	Name:     "test-user-flow",
	Category: "e2e",
	Usage:    "Test the user flow with the a user's UID",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "uid",
			Aliases:     []string{"u"},
			Usage:       "firebase `UID` of the user",
			Value:       "p4oqtfuvuUUmTX5yHqccIFqp89A2",
			DefaultText: "UID",
		},
	},
	Action: func(cCtx *cli.Context) error {
		uid := cCtx.String("uid")
		if uid == "" {
			return fmt.Errorf("user's UID is required")
		}

		err := TestUserFlow(uid)
		if err != nil {
			log.Println(err)
			return cli.Exit(err, 1)
		}

		return nil
	},
}
