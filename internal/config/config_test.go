package config

import (
	"flag"
	"os"
	"testing"

	"github.com/matthiasBT/monitoring/internal/adapters"
	"github.com/stretchr/testify/assert"
)

func TestInitServerConfig(t *testing.T) {
	logger := adapters.SetupLogger()
	tests := []struct {
		name    string
		cmdArgs []string
		envs    map[string]string
		want    ServerConfig
	}{
		{
			name:    "read from command line",
			cmdArgs: []string{"test", "-a", "0.0.0.0:8901"},
			envs:    map[string]string{},
			want:    ServerConfig{Addr: "0.0.0.0:8901", TemplatePath: templatePath},
		},
		{
			name:    "read from env",
			cmdArgs: []string{"test"},
			envs:    map[string]string{"ADDRESS": "localhost:8801"},
			want:    ServerConfig{Addr: "localhost:8801", TemplatePath: templatePath},
		},
		{
			name:    "env gets higher priority",
			cmdArgs: []string{"test", "-a", "localhost:8888"},
			envs:    map[string]string{"ADDRESS": "0.0.0.0:8080"},
			want:    ServerConfig{Addr: "0.0.0.0:8080", TemplatePath: templatePath},
		},
		{
			name:    "default value if no flag and no env",
			cmdArgs: []string{"test"},
			envs:    map[string]string{},
			want:    ServerConfig{Addr: ServerDefAddr, TemplatePath: templatePath},
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
			got := InitServerConfig(logger)
			assert.Equal(t, *got, tt.want, "Server configuration is different")
		})
	}
}

func TestInitAgentConfig(t *testing.T) {
	logger := adapters.SetupLogger()
	tests := []struct {
		name    string
		cmdArgs []string
		envs    map[string]string
		want    AgentConfig
	}{
		{
			name:    "read from env",
			cmdArgs: []string{"test"},
			envs:    map[string]string{"ADDRESS": "0.0.0.0:8765", "REPORT_INTERVAL": "4", "POLL_INTERVAL": "1"},
			want:    AgentConfig{Addr: "0.0.0.0:8765", ReportInterval: 4, PollInterval: 1},
		},
		{
			name:    "read from command line",
			cmdArgs: []string{"test", "-a", "0.0.0.0:8901", "-p", "7", "-r", "25"},
			envs:    map[string]string{},
			want:    AgentConfig{Addr: "0.0.0.0:8901", ReportInterval: 25, PollInterval: 7},
		},
		{
			name:    "env gets higher priority",
			cmdArgs: []string{"test", "-a", "0.0.0.0:8901", "-p", "7", "-r", "25"},
			envs:    map[string]string{"ADDRESS": "0.0.0.0:8765", "REPORT_INTERVAL": "4", "POLL_INTERVAL": "1"},
			want:    AgentConfig{Addr: "0.0.0.0:8765", ReportInterval: 4, PollInterval: 1},
		},
		{
			name:    "default values if no flag and no env",
			cmdArgs: []string{"test"},
			envs:    map[string]string{},
			want: AgentConfig{
				Addr: AgentDefAddr, ReportInterval: AgentDefReportInterval, PollInterval: AgentDefPollInterval,
			},
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
			got := InitAgentConfig(logger)
			assert.Equal(t, *got, tt.want, "Agent configuration is different")
		})
	}
}
