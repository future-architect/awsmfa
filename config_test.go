package awsmfa

import (
	"flag"
	"testing"

	"github.com/urfave/cli/v2"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name    string
		context *cli.Context
		wantErr bool
	}{
		{
			name: "normal cli config",
			context: func() *cli.Context {
				fs := flag.NewFlagSet("test", 0)
				sf := cli.StringFlag{
					Name:  "serial-number",
					Usage: "AWS serial-number",
				}
				_ = sf.Apply(fs)
				testCmd := []string{"--serial-number", "x", "123456"}
				_ = fs.Parse(testCmd)
				return cli.NewContext(nil, fs, nil)
			}(),
			wantErr: false,
		},
		{
			name: "missing serial-number",
			context: func() *cli.Context {
				fs := flag.NewFlagSet("test", 0)
				sf := cli.StringFlag{
					Name:  "serial-number",
					Usage: "AWS serial-number",
				}
				_ = sf.Apply(fs)
				testCmd := []string{"123456"}
				_ = fs.Parse(testCmd)
				return cli.NewContext(nil, fs, nil)
			}(),
			wantErr: true,
		},
		{
			name: "missing token-code",
			context: func() *cli.Context {
				fs := flag.NewFlagSet("test", 0)
				sf := cli.StringFlag{
					Name:  "serial-number",
					Usage: "AWS serial-number",
				}
				_ = sf.Apply(fs)
				testCmd := []string{"--serial-number", "x"}
				_ = fs.Parse(testCmd)
				return cli.NewContext(nil, fs, nil)
			}(),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewConfig(tt.context)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
