package awsmfa

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

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
	outConfigPath      string
	outCredentialsPath string
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
		return nil, errors.New("--token-code or token arguments is required")
	}

	client := sts.New(session.Must(session.NewSessionWithOptions(session.Options{
		Profile: c.String("profile"),
		Config: aws.Config{
			Region: aws.String(os.Getenv("AWS_REGION")),
		},
	})))

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	return &Config{
		client:             client,
		profile:            c.String("profile"),
		mfaProfileName:     c.String("mfa-profile-name"),
		configPath:         filepath.Join(homeDir, ".aws", "config"),
		credentialsPath:    filepath.Join(homeDir, ".aws", "credentials"),
		durationSeconds:    c.Int64("duration-seconds"),
		serialNumber:       serialNumber,
		mfaTokenCode:       mfaTokenCode,
		outConfigPath:      filepath.Join(homeDir, ".aws", "config"),
		outCredentialsPath: filepath.Join(homeDir, ".aws", "credentials"),
	}, nil
}
