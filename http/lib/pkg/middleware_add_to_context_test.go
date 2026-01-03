package httpadpt

import (
	"context"
	"errors"
	"testing"
)

func TestAddToContextMiddleware(t *testing.T) {
	type ctxKey string
	key := ctxKey("test-key")
	value := "test-value"

	addToContext := func(ctx context.Context) context.Context {
		return context.WithValue(ctx, key, value)
	}

	middleware := AddToContextMiddleware(addToContext)

	if middleware == nil {
		t.Fatal("AddToContextMiddleware() returned nil")
	}

	// Create a test handler that verifies the context was modified
	testHandler := &testContextHandler{
		verifyKey:     key,
		expectedValue: value,
	}

	wrappedHandler := middleware(testHandler)

	if wrappedHandler == nil {
		t.Fatal("Middleware returned nil Handler")
	}

	ctx := context.Background()
	input := &mockRequest{method: "GET"}
	output := &Response{}

	err := wrappedHandler.Invoke(ctx, input, output)

	if err != nil {
		t.Errorf("Invoke() error = %v, want nil", err)
	}

	// Verify the handler received the modified context
	if !testHandler.contextModified {
		t.Error("Expected context to be modified, but it wasn't")
	}
}

func TestAddToContextMiddleware_ContextPropagation(t *testing.T) {
	type ctxKey1 string
	type ctxKey2 string
	key1 := ctxKey1("key1")
	key2 := ctxKey2("key2")
	value1 := "value1"
	value2 := "value2"

	// First middleware adds key1
	middleware1 := AddToContextMiddleware(func(ctx context.Context) context.Context {
		return context.WithValue(ctx, key1, value1)
	})

	// Second middleware adds key2
	middleware2 := AddToContextMiddleware(func(ctx context.Context) context.Context {
		return context.WithValue(ctx, key2, value2)
	})

	testHandler := &testContextHandler{
		verifyKey:      key1,
		expectedValue:  value1,
		verifyKey2:     key2,
		expectedValue2: value2,
	}

	// Chain middlewares
	wrappedHandler := middleware1(middleware2(testHandler))

	ctx := context.Background()
	input := &mockRequest{method: "GET"}
	output := &Response{}

	err := wrappedHandler.Invoke(ctx, input, output)

	if err != nil {
		t.Errorf("Invoke() error = %v, want nil", err)
	}

	// Verify both context values are present
	if !testHandler.contextModified {
		t.Error("Expected context to be modified with key1")
	}

	if !testHandler.contextModified2 {
		t.Error("Expected context to be modified with key2")
	}
}

func TestAddToContextMiddleware_NilAddToContextFunction(t *testing.T) {
	// This test verifies the bug fix - nil addToContext function is handled gracefully
	middleware := AddToContextMiddleware(nil)

	if middleware == nil {
		t.Fatal("AddToContextMiddleware() should not return nil even with nil function")
	}

	testHandler := &testContextHandler{}
	wrappedHandler := middleware(testHandler)

	ctx := context.Background()
	input := &mockRequest{method: "GET"}
	output := &Response{}

	// Should not panic when addToContext is nil (bug fix)
	err := wrappedHandler.Invoke(ctx, input, output)

	if err != nil {
		t.Errorf("Invoke() error = %v, want nil (nil function should use original context)", err)
	}
}

func TestAddToContextMiddleware_ReturnsOriginalContext(t *testing.T) {
	// Test that if addToContext returns the same context, it still works
	middleware := AddToContextMiddleware(func(ctx context.Context) context.Context {
		return ctx // Return original context unchanged
	})

	testHandler := &testContextHandler{}
	wrappedHandler := middleware(testHandler)

	ctx := context.Background()
	input := &mockRequest{method: "GET"}
	output := &Response{}

	err := wrappedHandler.Invoke(ctx, input, output)

	if err != nil {
		t.Errorf("Invoke() error = %v, want nil", err)
	}
}

func TestAddToContextMiddleware_ErrorPropagation(t *testing.T) {
	testErr := errors.New("handler error")
	testHandler := &testContextHandler{err: testErr}

	middleware := AddToContextMiddleware(func(ctx context.Context) context.Context {
		type ctxKey string
		return context.WithValue(ctx, ctxKey("key"), "value")
	})

	wrappedHandler := middleware(testHandler)

	ctx := context.Background()
	input := &mockRequest{method: "GET"}
	output := &Response{}

	err := wrappedHandler.Invoke(ctx, input, output)

	if err == nil {
		t.Error("Expected error from handler")
	}

	if err != testErr {
		t.Errorf("Invoke() error = %v, want %v", err, testErr)
	}
}

func TestAddToContextMiddleware_WithNilInput(t *testing.T) {
	type ctxKey string
	middleware := AddToContextMiddleware(func(ctx context.Context) context.Context {
		return context.WithValue(ctx, ctxKey("key"), "value")
	})

	testHandler := &testContextHandler{}
	wrappedHandler := middleware(testHandler)

	ctx := context.Background()
	output := &Response{}

	err := wrappedHandler.Invoke(ctx, nil, output)

	if err != nil {
		t.Errorf("Invoke() error = %v, want nil", err)
	}
}

func TestAddToContextMiddleware_WithNilOutput(t *testing.T) {
	type ctxKey string
	middleware := AddToContextMiddleware(func(ctx context.Context) context.Context {
		return context.WithValue(ctx, ctxKey("key"), "value")
	})

	testHandler := &testContextHandler{}
	wrappedHandler := middleware(testHandler)

	ctx := context.Background()
	input := &mockRequest{method: "GET"}

	err := wrappedHandler.Invoke(ctx, input, nil)

	if err != nil {
		t.Errorf("Invoke() error = %v, want nil", err)
	}
}

func TestAddToContextMiddleware_MultipleContextValues(t *testing.T) {
	type ctxKey1 string
	type ctxKey2 string
	type ctxKey3 string

	middleware := AddToContextMiddleware(func(ctx context.Context) context.Context {
		ctx = context.WithValue(ctx, ctxKey1("key1"), "value1")
		ctx = context.WithValue(ctx, ctxKey2("key2"), "value2")
		ctx = context.WithValue(ctx, ctxKey3("key3"), "value3")
		return ctx
	})

	testHandler := &testContextHandler{
		verifyKey:      ctxKey1("key1"),
		expectedValue:  "value1",
		verifyKey2:     ctxKey2("key2"),
		expectedValue2: "value2",
		verifyKey3:     ctxKey3("key3"),
		expectedValue3: "value3",
	}

	wrappedHandler := middleware(testHandler)

	ctx := context.Background()
	input := &mockRequest{method: "GET"}
	output := &Response{}

	err := wrappedHandler.Invoke(ctx, input, output)

	if err != nil {
		t.Errorf("Invoke() error = %v, want nil", err)
	}

	if !testHandler.contextModified {
		t.Error("Expected context to contain key1")
	}

	if !testHandler.contextModified2 {
		t.Error("Expected context to contain key2")
	}

	if !testHandler.contextModified3 {
		t.Error("Expected context to contain key3")
	}
}

func TestAddToContextMiddleware_ContextReplacement(t *testing.T) {
	type ctxKey string
	originalKey := ctxKey("original")
	replacementKey := ctxKey("replacement")

	// Create context with original value
	ctx := context.WithValue(context.Background(), originalKey, "original-value")

	middleware := AddToContextMiddleware(func(ctx context.Context) context.Context {
		// Replace the context entirely (though this is unusual, it should work)
		return context.WithValue(context.Background(), replacementKey, "replacement-value")
	})

	testHandler := &testContextHandler{
		verifyKey:     replacementKey,
		expectedValue: "replacement-value",
	}

	wrappedHandler := middleware(testHandler)

	input := &mockRequest{method: "GET"}
	output := &Response{}

	err := wrappedHandler.Invoke(ctx, input, output)

	if err != nil {
		t.Errorf("Invoke() error = %v, want nil", err)
	}

	// Replacement context value should be present
	if !testHandler.contextModified {
		t.Error("Expected replacement context to contain replacement-key")
	}
}

// Helper type for testing
type testContextHandler struct {
	err              error
	verifyKey        interface{}
	expectedValue    interface{}
	contextModified  bool
	verifyKey2       interface{}
	expectedValue2   interface{}
	contextModified2 bool
	verifyKey3       interface{}
	expectedValue3   interface{}
	contextModified3 bool
}

func (h *testContextHandler) Invoke(ctx context.Context, _ Request, _ *Response) error {
	if h.verifyKey != nil {
		if ctx.Value(h.verifyKey) == h.expectedValue {
			h.contextModified = true
		}
	}

	if h.verifyKey2 != nil {
		if ctx.Value(h.verifyKey2) == h.expectedValue2 {
			h.contextModified2 = true
		}
	}

	if h.verifyKey3 != nil {
		if ctx.Value(h.verifyKey3) == h.expectedValue3 {
			h.contextModified3 = true
		}
	}

	return h.err
}
