package longresolver

var dfLogger Logger = &NilLogger{}

func SetLogger(l Logger) {
	dfLogger = l
}

func GetDfLogger() Logger {
	return dfLogger
}

type Logger interface {
	Infof(format string, msgs ...interface{})
	Errorf(format string, msgs ...interface{})
	Panicf(format string, msgs ...interface{})
}

type NilLogger struct{}

func (l *NilLogger) Infof(format string, msgs ...interface{})  {}
func (l *NilLogger) Errorf(format string, msgs ...interface{}) {}
func (l *NilLogger) Panicf(format string, msgs ...interface{}) {}
