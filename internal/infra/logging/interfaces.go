// Package logging defines interfaces and utilities for structured logging.
package logging

// ILogger is an interface that defines standard logging methods.
// It provides a generic interface for logging mechanisms and can be implemented
// by different logging frameworks. This interface includes methods for various
// logging levels such as Info, Fatal, Error, Warning, and Debug. Each level has
// a standard method, a formatted string method (suffixed with 'f'), and a line method
// for logging data without formatting (suffixed with 'ln').
//
// The methods are:
//   - Info, Infof, Infoln: Log informational messages that highlight the progress of the application.
//   - Fatal: Log critical messages after which the program should terminate.
//   - Errorf: Log error messages that indicate a failure in a specific operation.
//   - Warningf: Log potentially harmful situations.
//   - Debugf: Log detailed information useful for debugging during development.
//
// Implementations of ILogger should ensure that these methods behave according to
// the semantics of the underlying logging framework.
type ILogger interface {
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Infoln(args ...interface{})
	Fatal(args ...interface{})
	Errorf(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
}
