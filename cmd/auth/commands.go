package auth

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

var GetJWTTokenCommand = cli.Command{
	Name:     "jwt",
	Usage:    "load user jwt by using email",
	Args:     true,
	Category: "auth",
	Before: func(ctx *cli.Context) error {
		env := os.Getenv("TENK_ENV")
		if env == "" {
			env = "development"
		}

		if godotenv.Load(".env."+env) != nil {
			return cli.Exit("Error loading .env."+env+" file", 1)
		}

		fmt.Println("------------------Running in environment:", env)
		return nil
	},
	Action: func(ctx *cli.Context) error {
		email := ctx.Args().First()
		if email == "" {
			return cli.Exit("email is required", 1)
		}

		token, err := GetCustomTokenByEmail(email)
		if err != nil {
			return cli.Exit(err, 1)
		}

		fmt.Println("Token:", token)

		return nil
	},
}
