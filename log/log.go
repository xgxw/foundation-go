package log

import (
	"io"

	"github.com/sirupsen/logrus"
)

var (
	logLevel = map[string]logrus.Level{
		"panic": logrus.PanicLevel,
		"fatal": logrus.FatalLevel,
		"error": logrus.ErrorLevel,
		"warn":  logrus.WarnLevel,
		"info":  logrus.InfoLevel,
		"debug": logrus.DebugLevel,
	}

	logFormatter = map[string]logrus.Formatter{
		"json": &logrus.JSONFormatter{},
		"text": &logrus.TextFormatter{},
	}
)

// Logger hold a logrus.Logger
type Logger struct {
	*logrus.Logger
}

// Options configuration for logrus
type Options struct {
	Formatter string `yaml:"formatter" mapstructure:"formatter"`
	Level     string `yaml:"level" mapstructure:"level"`
	Mode      string `yaml:"mode" mapstructure:"mode"`
}

// New create a logrus entry instance and return
func NewLogger(opts Options, out io.Writer) *Logger {
	l := logrus.New()
	l.Out = out

	if opts.Mode == "debug" {
		l.SetLevel(logrus.DebugLevel)
		l.Formatter = &logrus.TextFormatter{
			ForceColors:      true,
			FullTimestamp:    true,
			QuoteEmptyFields: true,
		}
		return &Logger{Logger: l}
	}

	if level, ok := logLevel[opts.Level]; ok {
		l.SetLevel(level)
	} else {
		l.SetLevel(logrus.ErrorLevel)
	}

	if formatter, ok := logFormatter[opts.Formatter]; ok {
		l.Formatter = formatter
	} else {
		l.Formatter = logFormatter["json"]
	}

	return &Logger{
		Logger: l,
	}
}

// NewEntry creates an Entry object based on current Logger.
func (l *Logger) NewEntry() *Entry {
	return NewEntry(l)
}

func (l *Logger) WithScope(value string) *Entry {
	return l.NewEntry().WithScope(value)
}

func (l *Logger) WithField(name string, value interface{}) *Entry {
	return l.NewEntry().WithField(name, value)
}

func (l *Logger) WithFields(fields Fields) *Entry {
	return l.NewEntry().WithFields(fields)
}
