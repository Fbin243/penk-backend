package test

import (
	"fmt"
	"log"

	"tenkhours/test/common"

	"github.com/urfave/cli/v2"
)

var TestAPICommand = cli.Command{
	Name:     "api-test",
	Category: "test",
	Usage:    "Test the user flow with the a user's UID",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "uid",
			Aliases:     []string{"u"},
			Usage:       "firebase `UID` of the user",
			Value:       "OSXiAfFWsGcRATGzspuHqDjrdSf2",
			DefaultText: "UID",
		},
		&cli.StringSliceFlag{
			Name:    "flows",
			Aliases: []string{"f"},
			Usage:   "The flow to test",
			Value:   cli.NewStringSlice(string(common.ProfileFlowKey)),
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

		err := TestAPI(uid, flows)
		if err != nil {
			log.Println(err)
			return cli.Exit(err, 1)
		}

		return nil
	},
}
