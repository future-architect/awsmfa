package awsmfa

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"gopkg.in/ini.v1"
)

// Run executes awsmfa.
func Run(ctx context.Context, c *Config) error {
	out, err := c.client.GetSessionToken(ctx, &sts.GetSessionTokenInput{
		DurationSeconds: aws.Int32(c.durationSeconds),
		SerialNumber:    aws.String(c.serialNumber),
		TokenCode:       aws.String(c.mfaTokenCode),
	})
	if err != nil {
		return err
	}
	if !c.quiet {
		log.Printf("Wrote session token for profile %s, expiration: %s\n", c.mfaProfileName, *out.Credentials.Expiration)
	}

	// Create mfa-profile section in config, if section not exists.
	config, err := ini.Load(c.configPath)
	if err != nil {
		return err
	}
	_ = config.Section(fmt.Sprintf("profile %s", c.mfaProfileName))

	tmpConfig, err := ioutil.TempFile(c.awsDir, "config.tmp.*")
	if err != nil {
		return err
	}
	_, err = config.WriteTo(tmpConfig)
	if err != nil {
		return err
	}
	_ = tmpConfig.Close()
	if err := os.Rename(tmpConfig.Name(), c.outConfigPath); err != nil {
		return err
	}

	// Update or create mfa-profile section in credentials.
	credential, err := ini.Load(c.credentialsPath)
	if err != nil {
		return err
	}
	section := credential.Section(c.mfaProfileName)

	section.Key("aws_access_key_id").SetValue(aws.ToString(out.Credentials.AccessKeyId))
	section.Key("aws_secret_access_key").SetValue(aws.ToString(out.Credentials.SecretAccessKey))
	section.Key("aws_session_token").SetValue(aws.ToString(out.Credentials.SessionToken))

	tmpCredentialsFile, err := ioutil.TempFile(c.awsDir, "credentials.tmp.*")
	if err != nil {
		return err
	}
	_, err = credential.WriteTo(tmpCredentialsFile)
	if err != nil {
		return err
	}
	_ = tmpCredentialsFile.Close()
	if err := os.Rename(tmpCredentialsFile.Name(), c.outCredentialsPath); err != nil {
		return err
	}

	return nil
}
