package anaconda

// Logger interface to allow us to log
type Logger interface {
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})

	Panic(args ...interface{})
	Panicf(format string, args ...interface{})

	// Log functions
	Critical(args ...interface{})
	Criticalf(format string, args ...interface{})

	Error(args ...interface{})
	Errorf(format string, args ...interface{})

	Warning(args ...interface{})
	Warningf(format string, args ...interface{})

	Notice(args ...interface{})
	Noticef(format string, args ...interface{})

	Info(args ...interface{})
	Infof(format string, args ...interface{})

	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
}

// SetLogger lets you give us the logger
// you'd like us to log in.
// Default logger is silent
func (c *TwitterApi) SetLogger(l Logger) {
	c.Log = l
}

type silentLogger struct {
}

func (_ silentLogger) Fatal(_ ...interface{})                 {}
func (_ silentLogger) Fatalf(_ string, _ ...interface{})      {}
func (_ silentLogger) Panic(_ ...interface{})                 {}
func (_ silentLogger) Panicf(_ string, _ ...interface{})      {}
func (_ silentLogger) Critical(_ ...interface{})              {}
func (_ silentLogger) Criticalf(_ string, _ ...interface{})   {}
func (_ silentLogger) Error(_ ...interface{})                 {}
func (_ silentLogger) Errorf(_ string, _ ...interface{})      {}
func (_ silentLogger) Warning(_ ...interface{})               {}
func (_ silentLogger) Warningf(_ string, _ ...interface{})    {}
func (_ silentLogger) Notice(_ ...interface{})                {}
func (_ silentLogger) Noticef(_ string, _ ...interface{})     {}
func (_ silentLogger) Info(_ ...interface{})                  {}
func (_ silentLogger) Infof(_ string, _ ...interface{})       {}
func (_ silentLogger) Debug(_ ...interface{})                 {}
func (_ silentLogger) Debugf(format string, _ ...interface{}) {}
