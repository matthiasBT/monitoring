// Package logging provides utilities for setting up and configuring loggers.
package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

// SetupLogger initializes a new instance of logrus.Logger with predefined settings.
// It configures the logger to output to the standard output (os.Stdout) and sets the logging level to DebugLevel.
// This is useful for applications that require a basic logger setup without extensive configuration.
//
// Returns:
//
//	*logrus.Logger: A pointer to the newly created logrus.Logger instance, ready to use.
func SetupLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.DebugLevel)
	return logger
}
