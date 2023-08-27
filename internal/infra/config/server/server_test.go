package server

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitServerConfig(t *testing.T) {
	tests := []struct {
		name    string
		cmdArgs []string
		envs    map[string]string
		want    Config
	}{
		{
			name:    "read from command line",
			cmdArgs: []string{"test", "-a", "0.0.0.0:8901"},
			envs:    map[string]string{},
			want:    Config{Addr: "0.0.0.0:8901", TemplatePath: templatePath},
		},
		{
			name:    "read from env",
			cmdArgs: []string{"test"},
			envs:    map[string]string{"ADDRESS": "localhost:8801"},
			want:    Config{Addr: "localhost:8801", TemplatePath: templatePath},
		},
		{
			name:    "env gets higher priority",
			cmdArgs: []string{"test", "-a", "localhost:8888"},
			envs:    map[string]string{"ADDRESS": "0.0.0.0:8080"},
			want:    Config{Addr: "0.0.0.0:8080", TemplatePath: templatePath},
		},
		{
			name:    "default value if no flag and no env",
			cmdArgs: []string{"test"},
			envs:    map[string]string{},
			want:    Config{Addr: DefAddr, TemplatePath: templatePath},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.cmdArgs
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			os.Environ()
			for name, val := range tt.envs {
				t.Setenv(name, val)
			}
			got, _ := InitConfig()
			assert.Equal(t, *got, tt.want, "Server configuration is different")
		})
	}
}
