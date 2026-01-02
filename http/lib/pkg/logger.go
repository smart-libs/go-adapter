package httpadpt

import (
	"context"
	"log/slog"
)

type (
	// LoggerKey is used by the default implementation here to add the logger in the context
	LoggerKey struct{}

	LoggerProvider = func(ctx context.Context) *slog.Logger

	// LoggerToContext is used by the adapter to add to the context a new logger with labels to be
	// printed in all log lines, like XID
	LoggerToContext = func(ctx context.Context, current *slog.Logger) context.Context
)

var (
	DefaultLoggerProvider  = defaultLoggerProvider
	DefaultLoggerToContext = defaultLoggerToContext
)

func defaultLoggerProvider(ctx context.Context) *slog.Logger {
	ctxLogger, ok := ctx.Value(LoggerKey{}).(*slog.Logger)
	if ok {
		return ctxLogger
	}
	return slog.Default()
}

func defaultLoggerToContext(ctx context.Context, current *slog.Logger) context.Context {
	return context.WithValue(ctx, LoggerKey{}, current)
}
