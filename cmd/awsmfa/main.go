package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/d-tsuji/awsmfa"
	"github.com/urfave/cli/v2"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	app := &cli.App{
		Name:  "awsmfa",
		Usage: "Refresh AWS MFA",
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
			&cli.StringFlag{
				Name:  "config-path",
				Usage: "AWS config path",
				Value: filepath.Join(homeDir, ".aws", "config"),
			},
			&cli.StringFlag{
				Name:  "credentials-path",
				Usage: "AWS credentials path",
				Value: filepath.Join(homeDir, ".aws", "credentials"),
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
				Name:  "token-code",
				Usage: "AWS MFA token",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
