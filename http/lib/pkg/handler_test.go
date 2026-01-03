package httpadpt

import (
	"context"
	"errors"
	"testing"
)

func Test_funcBasedHandler_Invoke(t *testing.T) {
	tests := []struct {
		name      string
		handler   funcBasedHandler
		expectErr bool
		errValue  error
	}{
		{
			name: "successful invocation",
			handler: funcBasedHandler{
				invokeFunc: func(ctx context.Context, input Request, output *Response) error {
					return nil
				},
			},
			expectErr: false,
		},
		{
			name: "handler returns error",
			handler: funcBasedHandler{
				invokeFunc: func(ctx context.Context, input Request, output *Response) error {
					return errors.New("handler error")
				},
			},
			expectErr: true,
			errValue:  errors.New("handler error"),
		},
		{
			name: "handler modifies output",
			handler: funcBasedHandler{
				invokeFunc: func(ctx context.Context, input Request, output *Response) error {
					statusCode := 201
					output.StatusCode = &statusCode
					output.Body = []byte("test body")
					return nil
				},
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			input := &mockRequest{method: "GET"}
			output := &Response{}

			err := tt.handler.Invoke(ctx, input, output)

			if (err != nil) != tt.expectErr {
				t.Errorf("Invoke() error = %v, expectErr = %v", err, tt.expectErr)
			}

			if tt.expectErr && err != nil && err.Error() != tt.errValue.Error() {
				t.Errorf("Invoke() error = %v, want %v", err, tt.errValue)
			}
		})
	}
}

func Test_funcBasedHandler_Invoke_NilFunction(t *testing.T) {
	// This test verifies the bug fix - nil function is handled gracefully
	handler := funcBasedHandler{
		invokeFunc: nil,
	}

	ctx := context.Background()
	input := &mockRequest{method: "GET"}
	output := &Response{}

	// Should not panic when invokeFunc is nil (bug fix)
	err := handler.Invoke(ctx, input, output)

	if err != nil {
		t.Errorf("Invoke() error = %v, want nil (nil function should return nil error)", err)
	}
}

func Test_funcBasedHandler_Invoke_WithNilInput(t *testing.T) {
	handler := funcBasedHandler{
		invokeFunc: func(ctx context.Context, input Request, output *Response) error {
			// Handler should be able to handle nil input
			if input == nil {
				statusCode := 400
				output.StatusCode = &statusCode
			}
			return nil
		},
	}

	ctx := context.Background()
	output := &Response{}

	err := handler.Invoke(ctx, nil, output)

	if err != nil {
		t.Errorf("Invoke() error = %v, want nil", err)
	}
}

func Test_funcBasedHandler_Invoke_WithNilOutput(t *testing.T) {
	handler := funcBasedHandler{
		invokeFunc: func(ctx context.Context, input Request, output *Response) error {
			// Handler should be able to handle nil output
			if output == nil {
				return errors.New("output is nil")
			}
			return nil
		},
	}

	ctx := context.Background()
	input := &mockRequest{method: "GET"}

	err := handler.Invoke(ctx, input, nil)

	if err == nil {
		t.Error("Invoke() should return error when output is nil and handler checks for it")
	}
}

func TestMakeHandler(t *testing.T) {
	invoker := func(ctx context.Context, input Request, output *Response) error {
		statusCode := 200
		output.StatusCode = &statusCode
		return nil
	}

	handler := MakeHandler(invoker)

	if handler == nil {
		t.Fatal("MakeHandler() returned nil")
	}

	// Verify it's the correct type
	if _, ok := handler.(funcBasedHandler); !ok {
		t.Errorf("Expected funcBasedHandler type, got %T", handler)
	}

	// Test that it works
	ctx := context.Background()
	input := &mockRequest{method: "GET"}
	output := &Response{}

	err := handler.Invoke(ctx, input, output)

	if err != nil {
		t.Errorf("Invoke() error = %v, want nil", err)
	}

	if output.StatusCode == nil || *output.StatusCode != 200 {
		t.Errorf("Expected StatusCode = 200, got %v", output.StatusCode)
	}
}

func TestMakeHandler_NilInvoker(t *testing.T) {
	// This test verifies the bug fix - MakeHandler handles nil invoker gracefully
	handler := MakeHandler(nil)

	if handler == nil {
		t.Fatal("MakeHandler() should not return nil even with nil invoker")
	}

	// Verify it's the correct type
	if _, ok := handler.(funcBasedHandler); !ok {
		t.Errorf("Expected funcBasedHandler type, got %T", handler)
	}

	// Should not panic when invoked (bug fix - returns no-op handler)
	ctx := context.Background()
	input := &mockRequest{method: "GET"}
	output := &Response{}

	err := handler.Invoke(ctx, input, output)

	if err != nil {
		t.Errorf("Invoke() error = %v, want nil (nil invoker should create no-op handler)", err)
	}
}

func TestMakeHandler_ErrorHandling(t *testing.T) {
	testErr := errors.New("test error")
	invoker := func(ctx context.Context, input Request, output *Response) error {
		return testErr
	}

	handler := MakeHandler(invoker)

	ctx := context.Background()
	input := &mockRequest{method: "GET"}
	output := &Response{}

	err := handler.Invoke(ctx, input, output)

	if err == nil {
		t.Error("Expected error from handler")
	}

	if err != testErr {
		t.Errorf("Invoke() error = %v, want %v", err, testErr)
	}
}

func TestMakeHandler_ContextPropagation(t *testing.T) {
	var receivedCtx context.Context
	invoker := func(ctx context.Context, input Request, output *Response) error {
		receivedCtx = ctx
		return nil
	}

	handler := MakeHandler(invoker)

	type ctxKey string
	ctx := context.WithValue(context.Background(), ctxKey("test-key"), "test-value")
	input := &mockRequest{method: "GET"}
	output := &Response{}

	err := handler.Invoke(ctx, input, output)

	if err != nil {
		t.Errorf("Invoke() error = %v, want nil", err)
	}

	if receivedCtx == nil {
		t.Error("Expected context to be passed to invoker")
	}

	if receivedCtx.Value(ctxKey("test-key")) != "test-value" {
		t.Error("Expected context values to be preserved")
	}
}

func TestMakeHandler_InputOutputModification(t *testing.T) {
	invoker := func(ctx context.Context, input Request, output *Response) error {
		// Modify output
		statusCode := 201
		output.StatusCode = &statusCode
		output.Body = []byte("created")
		output.Header = map[string][]string{"Location": {"/api/resource/123"}}

		return nil
	}

	handler := MakeHandler(invoker)

	ctx := context.Background()
	input := &mockRequest{method: "POST"}
	output := &Response{}

	err := handler.Invoke(ctx, input, output)

	if err != nil {
		t.Errorf("Invoke() error = %v, want nil", err)
	}

	if output.StatusCode == nil || *output.StatusCode != 201 {
		t.Errorf("Expected StatusCode = 201, got %v", output.StatusCode)
	}

	if string(output.Body) != "created" {
		t.Errorf("Expected Body = 'created', got %q", string(output.Body))
	}

	if output.Header == nil {
		t.Fatal("Expected Header to be set")
	}

	location, ok := output.Header["Location"]
	if !ok || len(location) == 0 || location[0] != "/api/resource/123" {
		t.Errorf("Expected Location header = '/api/resource/123', got %v", location)
	}
}
