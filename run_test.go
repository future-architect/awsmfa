package awsmfa

import (
	"context"
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
		wantErr         bool
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
			wantErr:         false,
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
			wantErr:         false,
			wantConfig:      filepath.Join("testdata", "want", "config_2"),
			wantCredentials: filepath.Join("testdata", "want", "credentials_2"),
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
			configFile, err := os.Create(filepath.Join(tmpDir, "config"))
			if err != nil {
				t.Fatal(err)
			}
			defer configFile.Close()
			credentialsFile, err := os.Create(filepath.Join(tmpDir, "credentials"))
			if err != nil {
				t.Fatal(err)
			}
			defer credentialsFile.Close()
			tt.args.c.outConfigStream = configFile
			tt.args.c.outCredentialsStream = credentialsFile

			if err := Run(context.TODO(), tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
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
