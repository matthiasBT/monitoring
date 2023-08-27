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
				Addr:            "0.0.0.0:8765",
				StoreInterval:   ptruint(4),
				FileStoragePath: "/foo/bar.json",
				Restore:         ptrbool(false),
				TemplatePath:    templatePath,
			},
		},
		{
			name:    "read from command line",
			cmdArgs: []string{"test", "-i", "27", "-a", "1.2.3.4:5678", "-f", "/lol/kek.txt", "-r"},
			envs:    make(map[string]string),
			want: Config{
				Addr:            "1.2.3.4:5678",
				StoreInterval:   ptruint(27),
				FileStoragePath: "/lol/kek.txt",
				Restore:         ptrbool(true),
				TemplatePath:    templatePath,
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
				Addr:            "0.0.0.0:8765",
				StoreInterval:   ptruint(4),
				FileStoragePath: "/foo/bar.json",
				Restore:         ptrbool(false),
				TemplatePath:    templatePath,
			},
		},
		{
			name:    "default values if no flag and no env",
			cmdArgs: []string{"test"},
			envs:    make(map[string]string),
			want: Config{
				Addr:            DefAddr,
				StoreInterval:   ptruint(DefStoreInterval),
				FileStoragePath: DefFileStoragePath,
				Restore:         ptrbool(DefRestore),
				TemplatePath:    templatePath,
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

func TestStoresSync(t *testing.T) {
	tests := []struct {
		name   string
		config Config
		want   bool
	}{
		{
			name: "stores",
			config: Config{
				Addr:            "0.0.0.0:8765",
				StoreInterval:   ptruint(0),
				FileStoragePath: "/foo/bar.json",
				Restore:         ptrbool(false),
				TemplatePath:    templatePath,
			},
			want: true,
		},
		{
			name: "doesn't store",
			config: Config{
				Addr:            "0.0.0.0:8765",
				StoreInterval:   ptruint(1),
				FileStoragePath: "/foo/bar.json",
				Restore:         ptrbool(false),
				TemplatePath:    templatePath,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.config.StoresSync()
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
