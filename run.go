package awsmfa

import (
	"context"

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
	if err != nil {
		return err
	}

	// Create mfa-profile section in config, if not exists.
	config, err := ini.Load(c.configPath)
	if err != nil {
		return err
	}
	_ = config.Section(c.mfaProfileName)
	_, err = config.WriteTo(c.outConfigStream)
	if err != nil {
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

	_, err = credential.WriteTo(c.outCredentialsStream)
	if err != nil {
		return err
	}

	return nil
}
