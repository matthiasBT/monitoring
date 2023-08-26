package agent

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitAgentConfig(t *testing.T) {
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
			got, _ := InitAgentConfig()
			assert.Equal(t, *got, tt.want, "Agent configuration is different")
		})
	}
}
