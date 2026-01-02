package httpadpt

import (
	"context"
	"fmt"
	"log/slog"
	"sync/atomic"
	"time"
)

type (
	handleLogMiddleware struct {
		provider  LoggerProvider
		decorated Handler
	}
)

func getPath(input Request) string {
	if input == nil {
		return ""
	}
	if input.URL() == nil {
		return ""
	}
	return input.URL().Path
}

func getStatus(output *Response) int {
	if output == nil {
		return 500
	}
	if output.StatusCode == nil {
		return 500
	}
	return *output.StatusCode
}

var next atomic.Uint64

func getRequestID(start time.Time) string {
	localID := next.Add(1)
	return fmt.Sprintf("%s.%d", start.Format("20060102150405.000000"), localID)
}

func getDuration(start time.Time) string {
	return time.Since(start).String()
}

func getMethod(input Request) string {
	if input == nil {
		return "GET"
	}
	if input.Method() == "" {
		return "GET"
	}
	return input.Method()
}

func (h handleLogMiddleware) Invoke(ctx context.Context, input Request, output *Response) error {
	start := time.Now()
	xid := slog.String("xid", getRequestID(start))
	logger := h.provider(ctx)
	if logger != nil {
		ctx = DefaultLoggerToContext(ctx, logger.With(xid))
	}
	defer func() {
		if logger != nil {
			pathLabel := slog.String("path", getPath(input))
			statusCodeLabel := slog.Int("status", getStatus(output))
			durationLabel := slog.String("duration", getDuration(start))
			methodLabel := slog.String("method", getMethod(input))
			logger.Info("HTTP.Request", pathLabel, methodLabel, statusCodeLabel, xid, durationLabel)
		}
	}()

	return h.decorated.Invoke(ctx, input, output)
}

func NewHandleWithLoggerProviderMiddleware(provider LoggerProvider) Middleware {
	return func(next Handler) Handler {
		return handleLogMiddleware{provider: provider, decorated: next}
	}
}

func NewHandleWithSLogMiddleware(logger *slog.Logger) Middleware {
	return func(next Handler) Handler {
		return handleLogMiddleware{provider: func(_ context.Context) *slog.Logger {
			return logger
		}, decorated: next}
	}
}
