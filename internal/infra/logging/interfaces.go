package logging

type ILogger interface {
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Infoln(args ...interface{})
	Fatal(args ...interface{})
}
