package awsmfa

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"

	"github.com/urfave/cli/v2"
)

// Config holds information about the configuration of awsmfa.
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
	awsDir             string

	// print log config
	quiet bool
}

// NewConfig generates the Config.
func NewConfig(c *cli.Context) (*Config, error) {
	serialNumber := c.String("serial-number")
	if serialNumber == "" {
		return nil, errors.New("--serial-number is required")
	}

	mfaTokenCode := c.Args().First()
	if mfaTokenCode == "" {
		return nil, errors.New("[token-code] arguments is required")
	}

	client := sts.New(session.Must(session.NewSessionWithOptions(session.Options{
		Profile: c.String("profile"),
		Config: aws.Config{
			Region: aws.String(os.Getenv("AWS_REGION")),
		},
	})))

	var (
		configPath      string
		credentialsPath string
	)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	if v := os.Getenv("AWS_CONFIG_FILE"); v != "" {
		configPath = v
	} else {
		configPath = filepath.Join(homeDir, ".aws", "config")
	}
	if v := os.Getenv("AWS_SHARED_CREDENTIALS_FILE"); v != "" {
		credentialsPath = v
	} else {
		credentialsPath = filepath.Join(homeDir, ".aws", "credentials")
	}

	return &Config{
		client:             client,
		profile:            c.String("profile"),
		mfaProfileName:     c.String("mfa-profile-name"),
		configPath:         configPath,
		credentialsPath:    credentialsPath,
		durationSeconds:    c.Int64("duration-seconds"),
		serialNumber:       serialNumber,
		mfaTokenCode:       mfaTokenCode,
		outConfigPath:      configPath,
		outCredentialsPath: credentialsPath,
		awsDir:             filepath.Join(homeDir, ".aws"),
		quiet:              c.Bool("quiet"),
	}, nil
}
