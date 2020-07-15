package log

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"
)

type ctxLoggerMarker struct{}

var (
	ctxLoggerKey = &ctxLoggerMarker{}

	NullLogger = NewLogger(Options{}, ioutil.Discard)

	DefaultLogger = NewLogger(Options{Level: "error"}, os.Stdout)

	debugOutput io.Writer
)

func SetDebugOutput(w io.Writer) {
	debugOutput = w
}

func Extract(ctx context.Context) *Entry {
	l, ok := ctx.Value(ctxLoggerKey).(*Entry)
	if !ok {
		return NewEntry(DefaultLogger)
	}
	return l
}

func ToContext(ctx context.Context, logger *Entry) context.Context {
	return context.WithValue(ctx, ctxLoggerKey, logger)
}

func Debugf(s string, args ...interface{}) {
	if debugOutput == nil {
		debugOutput = os.Stdout
	}

	logPrefix := "[" + time.Now().Format("2006/01/02 15:04:05") + "] "
	suffix := "\n"
	fmt.Fprintf(debugOutput, logPrefix+s+suffix, args...)
}
