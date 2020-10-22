package awsmfa

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"

	"github.com/urfave/cli/v2"
)

type Config struct {
	// sts client
	client *sts.STS

	// sts config
	profile         string
	mfaProfileName  string
	configPath      string
	credentialsPath string
	durationSeconds int64
	serialNumber    string
	mfaTokenCode    string

	// output
	outConfigStream      io.Writer
	outCredentialsStream io.Writer
	closeFunc            func() error
}

func NewConfig(c *cli.Context) (*Config, error) {
	serialNumber := c.String("serial-number")
	if serialNumber == "" {
		accountID := c.String("account-id")
		role := c.String("role")
		serialNumber = fmt.Sprintf("arn:aws:iam::%s:mfa/%s", accountID, role)
	}

	mfaTokenCode := c.String("token-code")
	if mfaTokenCode == "" {
		mfaTokenCode = c.Args().First()
	}
	if mfaTokenCode == "" {
		return nil, errors.New("awsmfa: --token-code or token arguments is required")
	}

	client := sts.New(session.Must(session.NewSessionWithOptions(session.Options{
		Profile: c.String("profile"),
		Config: aws.Config{
			Region: aws.String(os.Getenv("AWS_REGION")),
		},
	})))

	configFile, err := os.Create(c.String("config-path"))
	if err != nil {
		return nil, err
	}
	credentialsFile, err := os.Create(c.String("credentials-path"))
	if err != nil {
		return nil, err
	}
	f := func() error {
		if err := configFile.Close(); err != nil {
			return err
		}
		if err := credentialsFile.Close(); err != nil {
			return err
		}
		return nil
	}

	return &Config{
		profile:              c.String("profile"),
		mfaProfileName:       c.String("mfa-profile-name"),
		configPath:           c.String("config-path"),
		credentialsPath:      c.String("credentials-path"),
		durationSeconds:      c.Int64("duration-seconds"),
		serialNumber:         serialNumber,
		mfaTokenCode:         mfaTokenCode,
		client:               client,
		outConfigStream:      configFile,
		outCredentialsStream: credentialsFile,
		closeFunc:            f,
	}, nil
}
