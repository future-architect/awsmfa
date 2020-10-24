package main

import (
	"context"
	"log"
	"os"

	"github.com/d-tsuji/awsmfa"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:      "awsmfa",
		Usage:     "Refresh AWS MFA",
		UsageText: "awsmfa [global options] [token-code]",
		Action: func(c *cli.Context) error {
			config, err := awsmfa.NewConfig(c)
			if err != nil {
				return err
			}
			return awsmfa.Run(context.Background(), config)
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "profile",
				Usage: "profile name",
				Value: "default",
			},
			&cli.StringFlag{
				Name:  "mfa-profile-name",
				Usage: "MFA profile name",
				Value: "mfa",
			},
			&cli.Int64Flag{
				Name:  "duration-seconds",
				Usage: "AWS IAM user session valid seconds",
				Value: 43200,
			},
			&cli.StringFlag{
				Name:  "serial-number",
				Usage: "AWS serial-number",
			},
			&cli.StringFlag{
				Name:  "account-id",
				Usage: "AWS account id",
			},
			&cli.StringFlag{
				Name:  "role",
				Usage: "AWS role",
			},
			&cli.StringFlag{
				Name:  "token-code",
				Usage: "AWS MFA token",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
