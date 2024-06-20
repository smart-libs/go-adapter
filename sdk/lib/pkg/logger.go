package sdk

import (
	"context"
	"fmt"
	"os"
)

type (
	Logger interface {
		Debug(string)
		Error(string, error)
	}

	noLogger       struct{}
	fmtPrintLogger struct{}

	ctxLoggerKey int
)

func (n noLogger) Debug(string)        {}
func (n noLogger) Error(string, error) {}

func (f fmtPrintLogger) Debug(s string)          { fmt.Println(s) }
func (f fmtPrintLogger) Error(s string, e error) { _, _ = fmt.Fprintf(os.Stderr, "%s: %s", s, e) }

var (
	NoLogger              Logger = noLogger{}
	PrintLogger           Logger = fmtPrintLogger{}
	singletonCtxLoggerKey        = ctxLoggerKey(1)
	LoggerFrom                   = func(ctx context.Context) Logger {
		logger, found := ctx.Value(singletonCtxLoggerKey).(Logger)
		if found {
			return logger
		}
		return NoLogger
	}
)

func NewContextWithLogger(parent context.Context, l Logger) context.Context {
	return context.WithValue(parent, singletonCtxLoggerKey, l)
}
