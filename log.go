package anaconda

import (
	"log"
	"os"
)

// The Logger interface provides optional logging ability for the streaming API.
// It can also be used to log the rate limiting headers if desired.
type Logger interface {
	Err(m string)
	Info(m string)
	Debug(m string)
}

// SetLogger sets the Logger used by the API client.
// The default logger is silent. BasicLogger will log to STDERR
// using the log package from the standard library.
func (c *TwitterApi) SetLogger(l Logger) {
	c.Log = l
}

type silentLogger struct {
}

func (_ silentLogger) Err(_ string)   {}
func (_ silentLogger) Info(_ string)  {}
func (_ silentLogger) Debug(_ string) {}

// BasicLogger is the equivalent of using log from the standard
// library to print to STDERR.
var BasicLogger Logger

type basicLogger struct {
	log *log.Logger //func New(out io.Writer, prefix string, flag int) *Logger
}

func init() {
	BasicLogger = &basicLogger{log: log.New(os.Stderr, log.Prefix(), log.LstdFlags)}
}

func (l basicLogger) Err(m string)   { l.log.Print("E ", m) }
func (l basicLogger) Info(m string)  { l.log.Print("I ", m) }
func (l basicLogger) Debug(m string) { l.log.Print("D ", m) }
