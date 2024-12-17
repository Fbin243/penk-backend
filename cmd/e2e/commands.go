package e2e

import (
	"fmt"
	"log"
	"tenkhours/test/common"

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
			Value:       "lstMYDOFoXWQ2s4TGyR4GTrGpKO2",
			DefaultText: "UID",
		},
		&cli.StringSliceFlag{
			Name:    "flows",
			Aliases: []string{"f"},
			Usage:   "The flow to test",
			Value:   cli.NewStringSlice(string(common.ProfilesFlowKey)),
		},
	},
	Action: func(cCtx *cli.Context) error {
		uid := cCtx.String("uid")
		if uid == "" {
			return fmt.Errorf("user's UID is required")
		}

		flows := cCtx.StringSlice("flows")
		if len(flows) == 0 {
			return fmt.Errorf("flows are required")
		}

		err := TestUserFlow(uid, flows)
		if err != nil {
			log.Println(err)
			return cli.Exit(err, 1)
		}

		return nil
	},
}
