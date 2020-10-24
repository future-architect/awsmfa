package awsmfa

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
	"gopkg.in/ini.v1"
)

func Run(ctx context.Context, c *Config) error {
	out, err := c.client.GetSessionTokenWithContext(ctx, &sts.GetSessionTokenInput{
		DurationSeconds: aws.Int64(c.durationSeconds),
		SerialNumber:    aws.String(c.serialNumber),
		TokenCode:       aws.String(c.mfaTokenCode),
	})
	if !c.quiet {
		log.Println(out)
	}
	if err != nil {
		return err
	}

	// Create mfa-profile section in config, if not exists.
	config, err := ini.Load(c.configPath)
	if err != nil {
		return err
	}
	_ = config.Section(fmt.Sprintf("profile %s", c.mfaProfileName))

	outConfigFile, err := os.Create(fmt.Sprintf("%s.tmp", c.outConfigPath))
	if err != nil {
		return err
	}
	_, err = config.WriteTo(outConfigFile)
	if err != nil {
		return err
	}
	outConfigFile.Close()
	if err := os.Rename(fmt.Sprintf("%s.tmp", c.outConfigPath), c.outConfigPath); err != nil {
		return err
	}

	// Update or create mfa-profile section in credentials.
	credential, err := ini.Load(c.credentialsPath)
	if err != nil {
		return err
	}
	section := credential.Section(c.mfaProfileName)

	section.Key("aws_access_key_id").SetValue(aws.StringValue(out.Credentials.AccessKeyId))
	section.Key("aws_secret_access_key").SetValue(aws.StringValue(out.Credentials.SecretAccessKey))
	section.Key("aws_session_token").SetValue(aws.StringValue(out.Credentials.SessionToken))

	outCredentialsFile, err := os.Create(fmt.Sprintf("%s.tmp", c.outCredentialsPath))
	if err != nil {
		return err
	}
	_, err = credential.WriteTo(outCredentialsFile)
	if err != nil {
		return err
	}
	outCredentialsFile.Close()
	if err := os.Rename(fmt.Sprintf("%s.tmp", c.outCredentialsPath), c.outCredentialsPath); err != nil {
		return err
	}

	return nil
}
