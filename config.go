package awsmfa

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/urfave/cli/v2"
)

// Config holds information about the configuration of awsmfa.
type Config struct {
	// sts client
	client *sts.Client

	// sts config
	profile         string
	mfaProfileName  string
	configPath      string
	credentialsPath string
	durationSeconds int32
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

	cfg, err := config.LoadDefaultConfig(c.Context,
		config.WithSharedConfigProfile(c.String("profile")),
		config.WithRegion(os.Getenv("AWS_REGION")),
	)
	if err != nil {
		return nil, err
	}
	client := sts.NewFromConfig(cfg)

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
		durationSeconds:    int32(c.Int("duration-seconds")),
		serialNumber:       serialNumber,
		mfaTokenCode:       mfaTokenCode,
		outConfigPath:      configPath,
		outCredentialsPath: credentialsPath,
		awsDir:             filepath.Join(homeDir, ".aws"),
		quiet:              c.Bool("quiet"),
	}, nil
}
