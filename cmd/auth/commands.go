package auth

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

var GetJWTTokenCommand = cli.Command{
	Name:     "jwt",
	Usage:    "load user jwt by using user's UID",
	Args:     true,
	Category: "auth",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "uid",
			Aliases:     []string{"u"},
			Usage:       "firebase `UID` of the user",
			Value:       "lstMYDOFoXWQ2s4TGyR4GTrGpKO2",
			DefaultText: "UID",
		},
	},
	Action: func(ctx *cli.Context) error {
		uid := ctx.String("uid")
		if uid == "" {
			return fmt.Errorf("user's UID is required")
		}

		token, err := GetIdTokenByUID(uid)
		if err != nil {
			return cli.Exit(err, 1)
		}

		fmt.Println("Token:", token)

		return nil
	},
}
