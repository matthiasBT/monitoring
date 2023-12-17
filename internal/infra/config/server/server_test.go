package server

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
			envs: map[string]string{
				"ADDRESS":           "0.0.0.0:8765",
				"STORE_INTERVAL":    "4",
				"FILE_STORAGE_PATH": "/foo/bar.json",
				"RESTORE":           "false",
			},
			want: Config{
				Addr:                 "0.0.0.0:8765",
				StoreInterval:        4,
				FileStoragePath:      "/foo/bar.json",
				Restore:              false,
				TemplatePath:         templatePath,
				RetryAttempts:        DefRetryAttempts,
				RetryIntervalBackoff: DefRetryIntervalBackoff,
				RetryIntervalInitial: DefRetryIntervalInitial,
			},
		},
		{
			name:    "read from command line",
			cmdArgs: []string{"test", "-i", "27", "-a", "1.2.3.4:5678", "-f", "/lol/kek.txt", "-r"},
			envs:    make(map[string]string),
			want: Config{
				Addr:                 "1.2.3.4:5678",
				StoreInterval:        27,
				FileStoragePath:      "/lol/kek.txt",
				Restore:              true,
				TemplatePath:         templatePath,
				RetryAttempts:        DefRetryAttempts,
				RetryIntervalBackoff: DefRetryIntervalBackoff,
				RetryIntervalInitial: DefRetryIntervalInitial,
			},
		},
		{
			name:    "env gets higher priority",
			cmdArgs: []string{"test", "-i", "27", "-a", "1.2.3.4:5678", "-f", "/lol/kek.txt", "-r"},
			envs: map[string]string{
				"ADDRESS":           "0.0.0.0:8765",
				"STORE_INTERVAL":    "4",
				"FILE_STORAGE_PATH": "/foo/bar.json",
				"RESTORE":           "false",
			},
			want: Config{
				Addr:                 "0.0.0.0:8765",
				StoreInterval:        4,
				FileStoragePath:      "/foo/bar.json",
				Restore:              false,
				TemplatePath:         templatePath,
				RetryAttempts:        DefRetryAttempts,
				RetryIntervalBackoff: DefRetryIntervalBackoff,
				RetryIntervalInitial: DefRetryIntervalInitial,
			},
		},
		{
			name:    "default values if no flag and no env",
			cmdArgs: []string{"test"},
			envs:    make(map[string]string),
			want: Config{
				Addr:                 DefAddr,
				StoreInterval:        DefStoreInterval,
				FileStoragePath:      DefFileStoragePath,
				Restore:              DefRestore,
				TemplatePath:         templatePath,
				RetryAttempts:        DefRetryAttempts,
				RetryIntervalBackoff: DefRetryIntervalBackoff,
				RetryIntervalInitial: DefRetryIntervalInitial,
			},
		},
		{
			name:    "empty storage path is allowed",
			cmdArgs: []string{"test", "-f", ""},
			envs:    make(map[string]string),
			want: Config{
				Addr:                 DefAddr,
				StoreInterval:        DefStoreInterval,
				FileStoragePath:      "",
				Restore:              DefRestore,
				TemplatePath:         templatePath,
				RetryAttempts:        DefRetryAttempts,
				RetryIntervalBackoff: DefRetryIntervalBackoff,
				RetryIntervalInitial: DefRetryIntervalInitial,
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
			assert.Equal(t, *got, tt.want, "Server configuration is different")
		})
	}
}

func TestFlushesSync(t *testing.T) {
	tests := []struct {
		name   string
		config Config
		want   bool
	}{
		{
			name: "stores",
			config: Config{
				Addr:            "0.0.0.0:8765",
				StoreInterval:   0,
				FileStoragePath: "/foo/bar.json",
				Restore:         false,
				TemplatePath:    templatePath,
			},
			want: true,
		},
		{
			name: "doesn't store",
			config: Config{
				Addr:            "0.0.0.0:8765",
				StoreInterval:   1,
				FileStoragePath: "/foo/bar.json",
				Restore:         false,
				TemplatePath:    templatePath,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.config.FlushesSync()
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestFlushes(t *testing.T) {
	tests := []struct {
		name   string
		config Config
		want   bool
	}{
		{
			name: "flushes",
			config: Config{
				Addr:            DefAddr,
				StoreInterval:   DefStoreInterval,
				FileStoragePath: "/foo/bar.json",
				Restore:         DefRestore,
				TemplatePath:    templatePath,
			},
			want: true,
		},
		{
			name: "doesn't flush",
			config: Config{
				Addr:            DefAddr,
				StoreInterval:   DefStoreInterval,
				FileStoragePath: "",
				Restore:         DefRestore,
				TemplatePath:    templatePath,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.config.Flushes()
			assert.Equal(t, got, tt.want)
		})
	}
}

func ptruint(val uint) *uint {
	return &val
}

func ptrbool(val bool) *bool {
	return &val
}
