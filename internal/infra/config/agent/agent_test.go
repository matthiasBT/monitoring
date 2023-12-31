package agent

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitConfig(t *testing.T) {
	tests := []struct {
		name    string
		cmdArgs []string
		envs    map[string]string
		want    Config
	}{
		{
			name:    "read from env",
			cmdArgs: []string{"test"},
			envs:    map[string]string{"ADDRESS": "0.0.0.0:8765", "REPORT_INTERVAL": "4", "POLL_INTERVAL": "1"},
			want: Config{
				Addr:                 "0.0.0.0:8765",
				ReportInterval:       4,
				PollInterval:         1,
				UpdateURL:            updateURL,
				RetryAttempts:        DefRetryAttempts,
				RetryIntervalBackoff: DefRetryIntervalBackoff,
				RetryIntervalInitial: DefRetryIntervalInitial,
				RateLimit:            DefRateLimit,
			},
		},
		{
			name:    "read from command line",
			cmdArgs: []string{"test", "-a", "0.0.0.0:8901", "-p", "7", "-r", "25"},
			envs:    map[string]string{},
			want: Config{
				Addr:                 "0.0.0.0:8901",
				ReportInterval:       25,
				PollInterval:         7,
				UpdateURL:            updateURL,
				RetryAttempts:        DefRetryAttempts,
				RetryIntervalBackoff: DefRetryIntervalBackoff,
				RetryIntervalInitial: DefRetryIntervalInitial,
				RateLimit:            DefRateLimit,
			},
		},
		{
			name:    "env gets higher priority",
			cmdArgs: []string{"test", "-a", "0.0.0.0:8901", "-p", "7", "-r", "25"},
			envs:    map[string]string{"ADDRESS": "0.0.0.0:8765", "REPORT_INTERVAL": "4", "POLL_INTERVAL": "1"},
			want: Config{
				Addr:                 "0.0.0.0:8765",
				ReportInterval:       4,
				PollInterval:         1,
				UpdateURL:            updateURL,
				RetryAttempts:        DefRetryAttempts,
				RetryIntervalBackoff: DefRetryIntervalBackoff,
				RetryIntervalInitial: DefRetryIntervalInitial,
				RateLimit:            DefRateLimit,
			},
		},
		{
			name:    "default values if no flag and no env",
			cmdArgs: []string{"test"},
			envs:    map[string]string{},
			want: Config{
				Addr:                 DefAddr,
				ReportInterval:       DefReportInterval,
				PollInterval:         DefPollInterval,
				UpdateURL:            updateURL,
				RetryAttempts:        DefRetryAttempts,
				RetryIntervalBackoff: DefRetryIntervalBackoff,
				RetryIntervalInitial: DefRetryIntervalInitial,
				RateLimit:            DefRateLimit,
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
			got, _ := InitConfig()
			assert.Equal(t, *got, tt.want, "Agent configuration is different")
		})
	}
}

func TestInitConfigWithJSON(t *testing.T) {
	tests := []struct {
		name     string
		cmdArgs  []string
		envs     map[string]string
		jsonArgs string
		want     Config
	}{
		{
			name: "priority: env > cmd > json > default",
			cmdArgs: []string{
				"test",
				"-a", "0.0.0.0:8901",
				"-c", "/tmp/test-monitoring-agent-config-priority-cmd.json",
			},
			envs: map[string]string{
				"CONFIG": "/tmp/test-monitoring-agent-config-priority-env.json",
			},
			jsonArgs: `{"address":"0.0.0.0:9001","crypto_key":"max.key"}`,
			want: Config{
				ConfigPath:           "/tmp/test-monitoring-agent-config-priority-env.json", // from env
				Addr:                 "0.0.0.0:8901",                                        // from cmd
				CryptoKey:            "max.key",                                             // from JSON
				ReportInterval:       DefReportInterval,                                     // defaults since no values
				PollInterval:         DefPollInterval,
				RateLimit:            DefRateLimit,
				UpdateURL:            updateURL,
				RetryAttempts:        DefRetryAttempts,
				RetryIntervalBackoff: DefRetryIntervalBackoff,
				RetryIntervalInitial: DefRetryIntervalInitial,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile, err := os.Create("/tmp/test-monitoring-agent-config-priority-env.json")
			if err != nil {
				t.Fatalf("Failed to create temporary file: %v", err)
			}
			defer os.Remove(tmpFile.Name()) // Clean up after the test
			if _, err := tmpFile.Write([]byte(tt.jsonArgs)); err != nil {
				t.Fatalf("Failed to write to temporary file: %v", err)
			}
			if err := tmpFile.Close(); err != nil {
				t.Fatalf("Failed to close temporary file: %v", err)
			}

			os.Args = tt.cmdArgs
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			for name, val := range tt.envs {
				t.Setenv(name, val)
			}
			got, _ := InitConfig()
			assert.Equal(t, *got, tt.want, "Agent configuration is different")
		})
	}
}
