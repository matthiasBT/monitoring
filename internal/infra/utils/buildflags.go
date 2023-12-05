package utils

import "fmt"

// PrintBuildFlags prints build flags of version, date and commit
// these flags are set in compile-time like this:
// 1. jump to project directory
// 2. execute the corresponding command to build a binary with non-default flag values:
// go build -ldflags "-X 'main.buildVersion=1.0.0' -X 'main.buildDate=$(date)' -X 'main.buildCommit=$(git rev-parse HEAD)'" ./cmd/agent
// go build -ldflags "-X 'main.buildVersion=1.0.0' -X 'main.buildDate=$(date)' -X 'main.buildCommit=$(git rev-parse HEAD)'" ./cmd/server
func PrintBuildFlags(version, date, commit string) {
	fmt.Printf("Build version: %s\n", getValueOrDefault(version, "N/A"))
	fmt.Printf("Build date: %s\n", getValueOrDefault(date, "N/A"))
	fmt.Printf("Build commit: %s\n", getValueOrDefault(commit, "N/A"))
}

// getValueOrDefault returns the value if it's not empty, otherwise returns the default value.
func getValueOrDefault(value, defaultValue string) string {
	if value != "" {
		return value
	}
	return defaultValue
}
