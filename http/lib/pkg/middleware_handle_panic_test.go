package httpadpt

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func Test_handlePanicMiddleware_Invoke_NoPanic(t *testing.T) {
	mockHandler := &testPanicHandler{shouldPanic: false, err: nil}
	middleware := handlePanicMiddleware{decorated: mockHandler}

	ctx := context.Background()
	input := &mockRequest{method: "GET"}
	output := &Response{StatusCode: intPtr(200)}

	err := middleware.Invoke(ctx, input, output)

	if err != nil {
		t.Errorf("Invoke() error = %v, want nil", err)
	}

	// Verify output was not modified
	if output.StatusCode == nil || *output.StatusCode != 200 {
		t.Errorf("Output StatusCode should remain 200, got %v", output.StatusCode)
	}
}

func Test_handlePanicMiddleware_Invoke_PanicWithError(t *testing.T) {
	panicErr := errors.New("test panic error")
	mockHandler := &testPanicHandler{shouldPanic: true, panicValue: panicErr}
	middleware := handlePanicMiddleware{decorated: mockHandler}

	ctx := context.Background()
	input := &mockRequest{method: "GET"}
	output := &Response{}

	err := middleware.Invoke(ctx, input, output)

	// Should not return error (panic was recovered)
	if err != nil {
		t.Errorf("Invoke() error = %v, want nil", err)
	}

	// Verify panic was handled
	if output.StatusCode == nil {
		t.Fatal("Expected StatusCode to be set after panic")
	}

	if *output.StatusCode != 500 {
		t.Errorf("Expected StatusCode = 500, got %d", *output.StatusCode)
	}

	// Verify Content-Type header
	if output.Header == nil {
		t.Fatal("Expected Header to be set after panic")
	}

	contentType, ok := output.Header["Content-Type"]
	if !ok || len(contentType) == 0 {
		t.Fatal("Expected Content-Type header to be set")
	}

	if contentType[0] != ContentTypeProblemDetail {
		t.Errorf("Expected Content-Type = %q, got %q", ContentTypeProblemDetail, contentType[0])
	}

	// Verify body contains problem detail
	if len(output.Body) == 0 {
		t.Fatal("Expected Body to be set after panic")
	}

	// Parse JSON body
	var pd ProblemDetail
	if err := json.Unmarshal(output.Body, &pd); err != nil {
		t.Fatalf("Failed to parse problem detail JSON: %v, body: %s", err, string(output.Body))
	}

	// Verify problem detail contains error information
	if pd.Detail == "" {
		t.Error("Expected ProblemDetail.Detail to contain error message")
	}

	if !strings.Contains(pd.Detail, "test panic error") {
		t.Errorf("Expected ProblemDetail.Detail to contain 'test panic error', got %q", pd.Detail)
	}
}

func Test_handlePanicMiddleware_Invoke_PanicWithString(t *testing.T) {
	panicValue := "string panic"
	mockHandler := &testPanicHandler{shouldPanic: true, panicValue: panicValue}
	middleware := handlePanicMiddleware{decorated: mockHandler}

	ctx := context.Background()
	input := &mockRequest{method: "GET"}
	output := &Response{}

	err := middleware.Invoke(ctx, input, output)

	// Should not return error (panic was recovered)
	if err != nil {
		t.Errorf("Invoke() error = %v, want nil", err)
	}

	// Verify panic was handled
	if output.StatusCode == nil || *output.StatusCode != 500 {
		t.Errorf("Expected StatusCode = 500, got %v", output.StatusCode)
	}

	// Verify body contains problem detail
	if len(output.Body) == 0 {
		t.Fatal("Expected Body to be set after panic")
	}

	// Parse JSON body
	var pd ProblemDetail
	if err := json.Unmarshal(output.Body, &pd); err != nil {
		t.Fatalf("Failed to parse problem detail JSON: %v", err)
	}

	// Verify problem detail contains panic value
	if pd.Detail == "" {
		t.Error("Expected ProblemDetail.Detail to contain panic value")
	}

	if !strings.Contains(pd.Detail, "string panic") {
		t.Errorf("Expected ProblemDetail.Detail to contain 'string panic', got %q", pd.Detail)
	}
}

func Test_handlePanicMiddleware_Invoke_PanicWithInt(t *testing.T) {
	panicValue := 42
	mockHandler := &testPanicHandler{shouldPanic: true, panicValue: panicValue}
	middleware := handlePanicMiddleware{decorated: mockHandler}

	ctx := context.Background()
	input := &mockRequest{method: "GET"}
	output := &Response{}

	err := middleware.Invoke(ctx, input, output)

	// Should not return error (panic was recovered)
	if err != nil {
		t.Errorf("Invoke() error = %v, want nil", err)
	}

	// Verify panic was handled
	if output.StatusCode == nil || *output.StatusCode != 500 {
		t.Errorf("Expected StatusCode = 500, got %v", output.StatusCode)
	}

	// Verify body contains problem detail
	if len(output.Body) == 0 {
		t.Fatal("Expected Body to be set after panic")
	}

	// Parse JSON body
	var pd ProblemDetail
	if err := json.Unmarshal(output.Body, &pd); err != nil {
		t.Fatalf("Failed to parse problem detail JSON: %v", err)
	}

	// Verify problem detail contains panic value
	if pd.Detail == "" {
		t.Error("Expected ProblemDetail.Detail to contain panic value")
	}

	if !strings.Contains(pd.Detail, "42") {
		t.Errorf("Expected ProblemDetail.Detail to contain '42', got %q", pd.Detail)
	}
}

func Test_handlePanicMiddleware_Invoke_PanicWithNilOutput(t *testing.T) {
	// This test verifies the bug fix - no panic when output is nil
	panicErr := errors.New("test panic")
	mockHandler := &testPanicHandler{shouldPanic: true, panicValue: panicErr}
	middleware := handlePanicMiddleware{decorated: mockHandler}

	ctx := context.Background()
	input := &mockRequest{method: "GET"}

	// Test with nil output - should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Invoke() panicked with nil output: %v", r)
		}
	}()

	err := middleware.Invoke(ctx, input, nil)

	// Should not return error (panic was recovered, even though output is nil)
	if err != nil {
		t.Errorf("Invoke() error = %v, want nil", err)
	}
}

func Test_handlePanicMiddleware_Invoke_PanicWithNilOutput_NoPanic(t *testing.T) {
	// This test specifically verifies the bug fix works
	panicValue := "panic with nil output"
	mockHandler := &testPanicHandler{shouldPanic: true, panicValue: panicValue}
	middleware := handlePanicMiddleware{decorated: mockHandler}

	ctx := context.Background()
	input := &mockRequest{method: "GET"}

	// Should not panic even when output is nil
	panicked := false
	func() {
		defer func() {
			if r := recover(); r != nil {
				panicked = true
			}
		}()
		_ = middleware.Invoke(ctx, input, nil)
	}()

	if panicked {
		t.Error("Invoke() should not panic when output is nil and handler panics")
	}
}

func Test_handlePanicMiddleware_Invoke_PanicPreservesExistingOutput(t *testing.T) {
	panicErr := errors.New("test panic")
	mockHandler := &testPanicHandler{shouldPanic: true, panicValue: panicErr}
	middleware := handlePanicMiddleware{decorated: mockHandler}

	ctx := context.Background()
	input := &mockRequest{method: "GET"}
	output := &Response{
		StatusCode: intPtr(200),
		Header:     map[string][]string{"X-Custom": {"value"}},
		Body:       []byte("original body"),
	}

	err := middleware.Invoke(ctx, input, output)

	if err != nil {
		t.Errorf("Invoke() error = %v, want nil", err)
	}

	// Verify StatusCode was changed to 500
	if output.StatusCode == nil || *output.StatusCode != 500 {
		t.Errorf("Expected StatusCode = 500, got %v", output.StatusCode)
	}

	// Verify Header was replaced (not merged)
	if output.Header == nil {
		t.Fatal("Expected Header to be set")
	}

	// Should have Content-Type
	if _, ok := output.Header["Content-Type"]; !ok {
		t.Error("Expected Content-Type header to be set")
	}

	// Original header should be gone (header was replaced, not merged)
	if _, ok := output.Header["X-Custom"]; ok {
		t.Error("Expected original headers to be replaced, not merged")
	}

	// Verify Body was replaced
	if len(output.Body) == 0 {
		t.Fatal("Expected Body to be set")
	}

	// Body should be problem detail JSON, not original body
	var pd ProblemDetail
	if err := json.Unmarshal(output.Body, &pd); err != nil {
		t.Errorf("Expected Body to be problem detail JSON, got: %s", string(output.Body))
	}
}

func Test_handlePanicMiddleware_Invoke_MultiplePanics(t *testing.T) {
	// Test that middleware can handle multiple panics in sequence
	mockHandler1 := &testPanicHandler{shouldPanic: true, panicValue: "panic 1"}
	middleware1 := handlePanicMiddleware{decorated: mockHandler1}

	mockHandler2 := &testPanicHandler{shouldPanic: true, panicValue: "panic 2"}
	middleware2 := handlePanicMiddleware{decorated: mockHandler2}

	ctx := context.Background()
	input := &mockRequest{method: "GET"}

	// First panic
	output1 := &Response{}
	err1 := middleware1.Invoke(ctx, input, output1)
	if err1 != nil {
		t.Errorf("First Invoke() error = %v, want nil", err1)
	}
	if output1.StatusCode == nil || *output1.StatusCode != 500 {
		t.Error("First panic should set StatusCode to 500")
	}

	// Second panic
	output2 := &Response{}
	err2 := middleware2.Invoke(ctx, input, output2)
	if err2 != nil {
		t.Errorf("Second Invoke() error = %v, want nil", err2)
	}
	if output2.StatusCode == nil || *output2.StatusCode != 500 {
		t.Error("Second panic should set StatusCode to 500")
	}
}

func Test_handlePanicMiddleware_Invoke_ErrorReturnedFromHandler(t *testing.T) {
	// Test that errors returned from handler are not treated as panics
	handlerErr := errors.New("handler error")
	mockHandler := &testPanicHandler{shouldPanic: false, err: handlerErr}
	middleware := handlePanicMiddleware{decorated: mockHandler}

	ctx := context.Background()
	input := &mockRequest{method: "GET"}
	output := &Response{StatusCode: intPtr(200)}

	err := middleware.Invoke(ctx, input, output)

	// Error should be returned, not converted to panic response
	if err == nil {
		t.Error("Expected error to be returned from handler")
	}

	if err != handlerErr {
		t.Errorf("Expected error = %v, got %v", handlerErr, err)
	}

	// Output should not be modified (no panic occurred)
	if output.StatusCode == nil || *output.StatusCode != 200 {
		t.Errorf("Expected StatusCode to remain 200, got %v", output.StatusCode)
	}

	if output.Header != nil {
		t.Error("Expected Header to remain nil (no panic occurred)")
	}

	if len(output.Body) > 0 {
		t.Error("Expected Body to remain empty (no panic occurred)")
	}
}

func TestHandlePanic(t *testing.T) {
	mockHandler := &testPanicHandler{shouldPanic: false, err: nil}
	wrappedHandler := HandlePanic(mockHandler)

	if wrappedHandler == nil {
		t.Fatal("HandlePanic() returned nil")
	}

	// Verify it's the correct type
	if _, ok := wrappedHandler.(handlePanicMiddleware); !ok {
		t.Errorf("Expected handlePanicMiddleware type, got %T", wrappedHandler)
	}

	// Test that it works
	ctx := context.Background()
	input := &mockRequest{method: "GET"}
	output := &Response{StatusCode: intPtr(200)}

	err := wrappedHandler.Invoke(ctx, input, output)
	if err != nil {
		t.Errorf("Invoke() error = %v, want nil", err)
	}
}

func TestHandlePanic_NilHandler(t *testing.T) {
	// Test with nil handler - should not panic during wrapping
	wrappedHandler := HandlePanic(nil)

	if wrappedHandler == nil {
		t.Fatal("HandlePanic() should not return nil even with nil handler")
	}

	// Should panic when invoked (because decorated handler is nil)
	// But the panic should be recovered and handled
	ctx := context.Background()
	input := &mockRequest{method: "GET"}
	output := &Response{}

	// The panic should be recovered by the middleware
	err := wrappedHandler.Invoke(ctx, input, output)

	// Should not return error (panic was recovered)
	if err != nil {
		t.Errorf("Invoke() error = %v, want nil (panic should be recovered)", err)
	}

	// Verify the panic was handled correctly
	if output.StatusCode == nil || *output.StatusCode != 500 {
		t.Errorf("Expected panic to set StatusCode = 500, got %v", output.StatusCode)
	}

	if output.Header == nil {
		t.Error("Expected Header to be set after panic")
	}

	if len(output.Body) == 0 {
		t.Error("Expected Body to be set after panic")
	}
}

func Test_handlePanicMiddleware_Invoke_ProblemDetailFormat(t *testing.T) {
	panicErr := errors.New("test error message")
	mockHandler := &testPanicHandler{shouldPanic: true, panicValue: panicErr}
	middleware := handlePanicMiddleware{decorated: mockHandler}

	ctx := context.Background()
	input := &mockRequest{method: "POST"}
	output := &Response{}

	err := middleware.Invoke(ctx, input, output)

	if err != nil {
		t.Errorf("Invoke() error = %v, want nil", err)
	}

	// Verify body is valid JSON
	if len(output.Body) == 0 {
		t.Fatal("Expected Body to be set")
	}

	var pd ProblemDetail
	if err := json.Unmarshal(output.Body, &pd); err != nil {
		t.Fatalf("Failed to parse problem detail JSON: %v", err)
	}

	// Verify problem detail structure
	if pd.Type == "" {
		t.Error("Expected ProblemDetail.Type to be set")
	}

	if pd.Detail == "" {
		t.Error("Expected ProblemDetail.Detail to be set")
	}

	if !strings.Contains(pd.Detail, "test error message") {
		t.Errorf("Expected ProblemDetail.Detail to contain error message, got %q", pd.Detail)
	}
}

// Helper type for testing
type testPanicHandler struct {
	shouldPanic bool
	panicValue  interface{}
	err         error
}

func (h *testPanicHandler) Invoke(_ context.Context, _ Request, _ *Response) error {
	if h.shouldPanic {
		panic(h.panicValue)
	}
	return h.err
}
