package interfaces

// todo: read about clean architecture in Go for web projects
type ILogger interface {
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Infoln(args ...interface{})
	Fatal(args ...interface{})
}
