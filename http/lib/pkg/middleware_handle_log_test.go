package httpadpt

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/url"
	"strings"
	"testing"
	"time"
)

func Test_getPath(t *testing.T) {
	tests := []struct {
		name     string
		input    Request
		expected string
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: "",
		},
		{
			name: "nil URL",
			input: &mockRequest{
				url: nil,
			},
			expected: "",
		},
		{
			name: "valid URL with path",
			input: &mockRequest{
				url: &url.URL{Path: "/api/v1/users"},
			},
			expected: "/api/v1/users",
		},
		{
			name: "empty path",
			input: &mockRequest{
				url: &url.URL{Path: ""},
			},
			expected: "",
		},
		{
			name: "root path",
			input: &mockRequest{
				url: &url.URL{Path: "/"},
			},
			expected: "/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getPath(tt.input)
			if result != tt.expected {
				t.Errorf("getPath() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func Test_getStatus(t *testing.T) {
	tests := []struct {
		name     string
		output   *Response
		expected int
	}{
		{
			name:     "nil output",
			output:   nil,
			expected: 500,
		},
		{
			name: "nil StatusCode",
			output: &Response{
				StatusCode: nil,
			},
			expected: 500,
		},
		{
			name: "valid status code 200",
			output: &Response{
				StatusCode: intPtr(200),
			},
			expected: 200,
		},
		{
			name: "valid status code 404",
			output: &Response{
				StatusCode: intPtr(404),
			},
			expected: 404,
		},
		{
			name: "valid status code 500",
			output: &Response{
				StatusCode: intPtr(500),
			},
			expected: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getStatus(tt.output)
			if result != tt.expected {
				t.Errorf("getStatus() = %d, want %d", result, tt.expected)
			}
		})
	}
}

func Test_getMethod(t *testing.T) {
	tests := []struct {
		name     string
		input    Request
		expected string
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: "GET",
		},
		{
			name: "empty method",
			input: &mockRequest{
				method: "",
			},
			expected: "GET",
		},
		{
			name: "GET method",
			input: &mockRequest{
				method: "GET",
			},
			expected: "GET",
		},
		{
			name: "POST method",
			input: &mockRequest{
				method: "POST",
			},
			expected: "POST",
		},
		{
			name: "PUT method",
			input: &mockRequest{
				method: "PUT",
			},
			expected: "PUT",
		},
		{
			name: "DELETE method",
			input: &mockRequest{
				method: "DELETE",
			},
			expected: "DELETE",
		},
		{
			name: "PATCH method",
			input: &mockRequest{
				method: "PATCH",
			},
			expected: "PATCH",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getMethod(tt.input)
			if result != tt.expected {
				t.Errorf("getMethod() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func Test_getRequestID(t *testing.T) {
	start := time.Now()
	id1 := getRequestID(start)
	id2 := getRequestID(start)

	// IDs should be different
	if id1 == id2 {
		t.Errorf("getRequestID() returned same ID twice: %q", id1)
	}

	// ID should contain timestamp and counter
	if !strings.Contains(id1, ".") {
		t.Errorf("getRequestID() should contain a dot separator, got %q", id1)
	}

	// Verify format: timestamp.counter (format is YYYYMMDDHHmmss.ffffff.counter)
	parts := strings.Split(id1, ".")
	if len(parts) < 2 {
		t.Errorf("getRequestID() should have format timestamp.counter, got %q", id1)
	}

	// Verify timestamp format (should be YYYYMMDDHHmmss.ffffff = 20 characters)
	// The format is: YYYYMMDDHHmmss.ffffff
	timestampPart := parts[0] + "." + parts[1]
	if len(timestampPart) != 21 { // 14 + 1 (dot) + 6 (microseconds)
		t.Errorf("getRequestID() timestamp part should be 21 characters (YYYYMMDDHHmmss.ffffff), got %d: %q", len(timestampPart), timestampPart)
	}

	// Verify counter is present and numeric
	if len(parts) < 3 || parts[2] == "" {
		t.Errorf("getRequestID() counter part should not be empty, got parts: %v", parts)
	}
}

func Test_getDuration(t *testing.T) {
	start := time.Now()
	time.Sleep(10 * time.Millisecond) // Small delay to ensure duration > 0
	duration := getDuration(start)

	// Duration should be a valid time.Duration string
	if duration == "" {
		t.Error("getDuration() returned empty string")
	}

	// Should contain time unit (ms, s, etc.)
	if !strings.Contains(duration, "ms") && !strings.Contains(duration, "s") && !strings.Contains(duration, "ns") {
		t.Errorf("getDuration() should contain time unit, got %q", duration)
	}

	// Test with zero duration (should still return valid string)
	zeroStart := time.Now()
	zeroDuration := getDuration(zeroStart)
	if zeroDuration == "" {
		t.Error("getDuration() returned empty string for zero duration")
	}
}

func Test_handleLogMiddleware_Invoke(t *testing.T) {
	tests := []struct {
		name           string
		provider       LoggerProvider
		input          Request
		output         *Response
		handlerError   error
		expectLog      bool
		expectError    bool
		expectedStatus int
	}{
		{
			name: "successful request",
			provider: func(ctx context.Context) *slog.Logger {
				return slog.Default()
			},
			input: &mockRequest{
				method: "GET",
				url:    &url.URL{Path: "/api/test"},
			},
			output: &Response{
				StatusCode: intPtr(200),
			},
			handlerError: nil,
			expectLog:    true,
			expectError:  false,
		},
		{
			name: "request with error",
			provider: func(ctx context.Context) *slog.Logger {
				return slog.Default()
			},
			input: &mockRequest{
				method: "POST",
				url:    &url.URL{Path: "/api/create"},
			},
			output: &Response{
				StatusCode: intPtr(500),
			},
			handlerError: nil,
			expectLog:    true,
			expectError:  false,
		},
		{
			name: "nil input",
			provider: func(ctx context.Context) *slog.Logger {
				return slog.Default()
			},
			input:        nil,
			output:       &Response{StatusCode: intPtr(200)},
			handlerError: nil,
			expectLog:    true,
			expectError:  false,
		},
		{
			name: "nil output",
			provider: func(ctx context.Context) *slog.Logger {
				return slog.Default()
			},
			input: &mockRequest{
				method: "GET",
				url:    &url.URL{Path: "/api/test"},
			},
			output:       nil,
			handlerError: nil,
			expectLog:    true,
			expectError:  false,
		},
		{
			name: "nil logger provider",
			provider: func(ctx context.Context) *slog.Logger {
				return nil
			},
			input: &mockRequest{
				method: "GET",
				url:    &url.URL{Path: "/api/test"},
			},
			output:       &Response{StatusCode: intPtr(200)},
			handlerError: nil,
			expectLog:    false,
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var logBuffer bytes.Buffer
			logger := slog.New(slog.NewJSONHandler(&logBuffer, nil))

			// Override provider if needed
			provider := tt.provider
			if tt.provider == nil {
				provider = func(ctx context.Context) *slog.Logger {
					return logger
				}
			} else if tt.name == "nil logger provider" {
				// Keep the nil provider
			} else {
				provider = func(ctx context.Context) *slog.Logger {
					return logger
				}
			}

			testHandler := &testLogHandler{
				err: tt.handlerError,
			}

			middleware := handleLogMiddleware{
				provider:  provider,
				decorated: testHandler,
			}

			ctx := context.Background()
			err := middleware.Invoke(ctx, tt.input, tt.output)

			if (err != nil) != tt.expectError {
				t.Errorf("Invoke() error = %v, expectError = %v", err, tt.expectError)
			}

			if tt.expectLog && provider != nil {
				logOutput := logBuffer.String()
				if logOutput == "" {
					t.Error("Expected log output but got empty")
				} else {
					// Verify log contains expected fields
					if !strings.Contains(logOutput, "HTTP.Request") {
						t.Errorf("Log should contain 'HTTP.Request', got: %s", logOutput)
					}
				}
			}
		})
	}
}

func Test_handleLogMiddleware_Invoke_LogsCorrectFields(t *testing.T) {
	var logBuffer bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logBuffer, nil))

	testHandler := &testLogHandler{err: nil}
	middleware := handleLogMiddleware{
		provider: func(ctx context.Context) *slog.Logger {
			return logger
		},
		decorated: testHandler,
	}

	input := &mockRequest{
		method: "POST",
		url:    &url.URL{Path: "/api/v1/users"},
	}
	output := &Response{
		StatusCode: intPtr(201),
	}

	ctx := context.Background()
	err := middleware.Invoke(ctx, input, output)

	if err != nil {
		t.Errorf("Invoke() error = %v, want nil", err)
	}

	logOutput := logBuffer.String()
	if logOutput == "" {
		t.Fatal("Expected log output but got empty")
	}

	// Parse JSON log
	var logEntry map[string]interface{}
	if err := json.Unmarshal([]byte(logOutput), &logEntry); err != nil {
		t.Fatalf("Failed to parse log JSON: %v, log: %s", err, logOutput)
	}

	// Verify required fields
	requiredFields := []string{"msg", "xid", "path", "method", "status", "duration"}
	for _, field := range requiredFields {
		if _, ok := logEntry[field]; !ok {
			t.Errorf("Log entry missing required field: %s, log: %s", field, logOutput)
		}
	}

	// Verify field values
	if logEntry["msg"] != "HTTP.Request" {
		t.Errorf("Log msg = %v, want 'HTTP.Request'", logEntry["msg"])
	}

	if logEntry["path"] != "/api/v1/users" {
		t.Errorf("Log path = %v, want '/api/v1/users'", logEntry["path"])
	}

	if logEntry["method"] != "POST" {
		t.Errorf("Log method = %v, want 'POST'", logEntry["method"])
	}

	// Status should be a number
	status, ok := logEntry["status"].(float64)
	if !ok {
		t.Errorf("Log status should be a number, got %T: %v", logEntry["status"], logEntry["status"])
	} else if int(status) != 201 {
		t.Errorf("Log status = %v, want 201", int(status))
	}

	// xid should be a string
	if _, ok := logEntry["xid"].(string); !ok {
		t.Errorf("Log xid should be a string, got %T: %v", logEntry["xid"], logEntry["xid"])
	}

	// duration should be a string
	if _, ok := logEntry["duration"].(string); !ok {
		t.Errorf("Log duration should be a string, got %T: %v", logEntry["duration"], logEntry["duration"])
	}
}

func Test_handleLogMiddleware_Invoke_ContextLogger(t *testing.T) {
	var logBuffer bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logBuffer, nil))

	testHandler := &testLogHandler{err: nil}
	middleware := handleLogMiddleware{
		provider: func(ctx context.Context) *slog.Logger {
			return logger
		},
		decorated: testHandler,
	}

	input := &mockRequest{
		method: "GET",
		url:    &url.URL{Path: "/test"},
	}
	output := &Response{
		StatusCode: intPtr(200),
	}

	ctx := context.Background()
	err := middleware.Invoke(ctx, input, output)

	if err != nil {
		t.Errorf("Invoke() error = %v, want nil", err)
	}

	// Verify that logger was added to context
	ctxLogger := DefaultLoggerProvider(ctx)
	if ctxLogger == nil {
		t.Error("Expected logger in context but got nil")
	}
}

func Test_handleLogMiddleware_Invoke_NoPanicOnNilInputs(t *testing.T) {
	var logBuffer bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logBuffer, nil))

	testHandler := &testLogHandler{err: nil}
	middleware := handleLogMiddleware{
		provider: func(ctx context.Context) *slog.Logger {
			return logger
		},
		decorated: testHandler,
	}

	ctx := context.Background()

	// Test with nil input and nil output
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Invoke() panicked with nil inputs: %v", r)
		}
	}()

	err := middleware.Invoke(ctx, nil, nil)
	if err != nil {
		t.Errorf("Invoke() error = %v, want nil", err)
	}
}

func Test_handleLogMiddleware_Invoke_NoPanicOnNilLogger(t *testing.T) {
	testHandler := &testLogHandler{err: nil}
	middleware := handleLogMiddleware{
		provider: func(ctx context.Context) *slog.Logger {
			return nil
		},
		decorated: testHandler,
	}

	input := &mockRequest{
		method: "GET",
		url:    &url.URL{Path: "/test"},
	}
	output := &Response{
		StatusCode: intPtr(200),
	}

	ctx := context.Background()

	// Should not panic even with nil logger
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Invoke() panicked with nil logger: %v", r)
		}
	}()

	err := middleware.Invoke(ctx, input, output)
	if err != nil {
		t.Errorf("Invoke() error = %v, want nil", err)
	}
}

func TestNewHandleWithLoggerProviderMiddleware(t *testing.T) {
	provider := func(ctx context.Context) *slog.Logger {
		return slog.Default()
	}

	middlewareFunc := NewHandleWithLoggerProviderMiddleware(provider)
	if middlewareFunc == nil {
		t.Fatal("NewHandleWithLoggerProviderMiddleware() returned nil")
	}

	// Test that it returns a Handler when given a Handler
	testHandler := &testLogHandler{err: nil}
	wrappedHandler := middlewareFunc(testHandler)

	if wrappedHandler == nil {
		t.Fatal("Middleware function returned nil Handler")
	}

	// Verify it's the correct type
	if _, ok := wrappedHandler.(handleLogMiddleware); !ok {
		t.Errorf("Expected handleLogMiddleware type, got %T", wrappedHandler)
	}
}

func TestNewHandleWithSLogMiddleware(t *testing.T) {
	logger := slog.Default()

	middlewareFunc := NewHandleWithSLogMiddleware(logger)
	if middlewareFunc == nil {
		t.Fatal("NewHandleWithSLogMiddleware() returned nil")
	}

	// Test that it returns a Handler when given a Handler
	testHandler := &testLogHandler{err: nil}
	wrappedHandler := middlewareFunc(testHandler)

	if wrappedHandler == nil {
		t.Fatal("Middleware function returned nil Handler")
	}

	// Verify it's the correct type
	if _, ok := wrappedHandler.(handleLogMiddleware); !ok {
		t.Errorf("Expected handleLogMiddleware type, got %T", wrappedHandler)
	}

	// Test that the provider returns the correct logger
	hlm := wrappedHandler.(handleLogMiddleware)
	ctx := context.Background()
	returnedLogger := hlm.provider(ctx)

	if returnedLogger != logger {
		t.Errorf("Provider returned different logger, expected %p, got %p", logger, returnedLogger)
	}
}

func TestNewHandleWithSLogMiddleware_NilLogger(t *testing.T) {
	middlewareFunc := NewHandleWithSLogMiddleware(nil)
	if middlewareFunc == nil {
		t.Fatal("NewHandleWithSLogMiddleware() returned nil")
	}

	testHandler := &testLogHandler{err: nil}
	wrappedHandler := middlewareFunc(testHandler)

	// Should not panic
	ctx := context.Background()
	input := &mockRequest{
		method: "GET",
		url:    &url.URL{Path: "/test"},
	}
	output := &Response{
		StatusCode: intPtr(200),
	}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Invoke() panicked with nil logger: %v", r)
		}
	}()

	err := wrappedHandler.Invoke(ctx, input, output)
	if err != nil {
		t.Errorf("Invoke() error = %v, want nil", err)
	}
}

// Helper type for testing
type testLogHandler struct {
	err error
}

func (m *testLogHandler) Invoke(_ context.Context, _ Request, _ *Response) error {
	return m.err
}
