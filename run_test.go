package awsmfa

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

func TestRun(t *testing.T) {
	type args struct {
		c *Config
	}
	tests := []struct {
		name            string
		args            args
		wantErr         error
		wantConfig      string
		wantCredentials string
	}{
		{
			name: "create mfa section",
			args: args{
				c: &Config{
					durationSeconds: 43200,
					serialNumber:    "arn:aws:iam::123456789012:mfa/test",
					mfaTokenCode:    "012345",
					mfaProfileName:  "mfa",
					configPath:      filepath.Join("testdata", ".aws", "config"),
					credentialsPath: filepath.Join("testdata", ".aws", "credentials"),
				},
			},
			wantErr:         nil,
			wantConfig:      filepath.Join("testdata", "want", "config"),
			wantCredentials: filepath.Join("testdata", "want", "credentials"),
		},
		{
			name: "update mfa section",
			args: args{
				c: &Config{
					durationSeconds: 43200,
					serialNumber:    "arn:aws:iam::123456789012:mfa/test",
					mfaTokenCode:    "012345",
					mfaProfileName:  "mfa",
					configPath:      filepath.Join("testdata", ".aws", "config_2"),
					credentialsPath: filepath.Join("testdata", ".aws", "credentials_2"),
				},
			},
			wantErr:         nil,
			wantConfig:      filepath.Join("testdata", "want", "config_2"),
			wantCredentials: filepath.Join("testdata", "want", "credentials_2"),
		},
		{
			name: "no config",
			args: args{
				c: &Config{
					durationSeconds: 43200,
					serialNumber:    "arn:aws:iam::123456789012:mfa/test",
					mfaTokenCode:    "012345",
					mfaProfileName:  "mfa",
					configPath:      filepath.Join("testdata", ".aws", "config_not_found"),
					credentialsPath: filepath.Join("testdata", ".aws", "credentials"),
				},
			},
			wantErr: &os.PathError{},
		},
		{
			name: "no credentials",
			args: args{
				c: &Config{
					durationSeconds: 43200,
					serialNumber:    "arn:aws:iam::123456789012:mfa/test",
					mfaTokenCode:    "012345",
					mfaProfileName:  "mfa",
					configPath:      filepath.Join("testdata", ".aws", "config"),
					credentialsPath: filepath.Join("testdata", ".aws", "credentials_not_found"),
				},
			},
			wantErr: &os.PathError{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.c.client = sts.New(session.Must(session.NewSessionWithOptions(session.Options{
				Profile: tt.args.c.profile,
				Config: aws.Config{
					Region:   aws.String("ap-northeast-1"),
					Endpoint: aws.String("http://localhost:4566"),
				},
			})))

			tmpDir := t.TempDir()
			tt.args.c.outConfigPath = filepath.Join(tmpDir, "config")
			tt.args.c.outCredentialsPath = filepath.Join(tmpDir, "credentials")

			err := Run(context.TODO(), tt.args.c)
			if err != nil {
				if tt.wantErr == nil || !errors.As(err, &tt.wantErr) {
					t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			compareFile(t, filepath.Join(tmpDir, "config"), tt.wantConfig)
			compareFile(t, filepath.Join(tmpDir, "credentials"), tt.wantCredentials)
		})
	}
}

func compareFile(t *testing.T, gotPath, wantPath string) {
	t.Helper()
	got, err := ioutil.ReadFile(gotPath)
	if err != nil {
		t.Error(err)
	}
	want, err := ioutil.ReadFile(wantPath)
	if err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff(string(want), string(got)); diff != "" {
		t.Errorf("File mismatch (-want +got):\n%s", diff)
	}
}
