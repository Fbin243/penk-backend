package time

import (
	"time"

	"github.com/urfave/cli/v2"
)

var GetTheCurrentTime = cli.Command{
	Name:     "get-current-time",
	Category: "time",
	Usage:    "Get the current time in ",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "format",
			Aliases: []string{"f"},
			Usage:   "The format to display the current time in",
			Value:   time.RFC3339Nano,
		},
	},
	Action: func(cCtx *cli.Context) error {
		format := cCtx.String("format")
		GetTheCurrentTimeInFormat(format)

		return nil
	},
}
