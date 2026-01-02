package gonethttp

import (
	"context"
	"errors"
	httpadpt "github.com/smart-libs/go-adapter/http/lib/pkg"
	"net/http"
	"net/http/httptest"
	"testing"
)

// mockHandler is a test implementation of httpadpt.Handler
type mockHandler struct {
	invokeFunc func(ctx context.Context, input httpadpt.Request, output *httpadpt.Response) error
}

func (m *mockHandler) Invoke(ctx context.Context, input httpadpt.Request, output *httpadpt.Response) error {
	if m.invokeFunc != nil {
		return m.invokeFunc(ctx, input, output)
	}
	return nil
}

func Test_buildPath(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		expectedResult string
	}{
		{
			name:           "method and path provided",
			method:         "GET",
			path:           "/api/users",
			expectedResult: "GET /api/users",
		},
		{
			name:           "empty method, path provided",
			method:         "",
			path:           "/api/users",
			expectedResult: "/api/users",
		},
		{
			name:           "POST method with path",
			method:         "POST",
			path:           "/api/users",
			expectedResult: "POST /api/users",
		},
		{
			name:           "PUT method with path",
			method:         "PUT",
			path:           "/api/users/123",
			expectedResult: "PUT /api/users/123",
		},
		{
			name:           "DELETE method with path",
			method:         "DELETE",
			path:           "/api/users/123",
			expectedResult: "DELETE /api/users/123",
		},
		{
			name:           "PATCH method with path",
			method:         "PATCH",
			path:           "/api/users/123",
			expectedResult: "PATCH /api/users/123",
		},
		{
			name:           "empty path with method",
			method:         "GET",
			path:           "",
			expectedResult: "GET ",
		},
		{
			name:           "empty method and empty path",
			method:         "",
			path:           "",
			expectedResult: "",
		},
		{
			name:           "path with query parameters",
			method:         "GET",
			path:           "/api/users?page=1",
			expectedResult: "GET /api/users?page=1",
		},
		{
			name:           "path with special characters",
			method:         "GET",
			path:           "/api/users/{id}",
			expectedResult: "GET /api/users/{id}",
		},
		{
			name:           "root path",
			method:         "GET",
			path:           "/",
			expectedResult: "GET /",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildPath(tt.method, tt.path)
			if result != tt.expectedResult {
				t.Errorf("buildPath(%q, %q) = %q, want %q", tt.method, tt.path, result, tt.expectedResult)
			}
		})
	}
}

func Test_buildHandler_Success(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     *int
		headers        map[string][]string
		body           []byte
		expectedStatus int
		expectedBody   string
		expectedHeader map[string][]string
	}{
		{
			name:           "success with status code",
			statusCode:     intPtr(201),
			headers:        nil,
			body:           []byte("created"),
			expectedStatus: 201,
			expectedBody:   "created",
			expectedHeader: nil,
		},
		{
			name:           "success with default status code",
			statusCode:     nil,
			headers:        nil,
			body:           []byte("ok"),
			expectedStatus: 200,
			expectedBody:   "ok",
			expectedHeader: nil,
		},
		{
			name:       "success with headers",
			statusCode: intPtr(200),
			headers: map[string][]string{
				"Content-Type": {"application/json"},
				"X-Custom":     {"value1", "value2"},
			},
			body:           []byte("{\"key\":\"value\"}"),
			expectedStatus: 200,
			expectedBody:   "{\"key\":\"value\"}",
			expectedHeader: map[string][]string{
				"Content-Type": {"application/json"},
				"X-Custom":     {"value1", "value2"},
			},
		},
		{
			name:           "success with empty body",
			statusCode:     intPtr(204),
			headers:        nil,
			body:           nil,
			expectedStatus: 204,
			expectedBody:   "",
			expectedHeader: nil,
		},
		{
			name:       "success with headers and no body",
			statusCode: intPtr(200),
			headers: map[string][]string{
				"Location": {"/api/users/123"},
			},
			body:           nil,
			expectedStatus: 200,
			expectedBody:   "",
			expectedHeader: map[string][]string{
				"Location": {"/api/users/123"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &mockHandler{
				invokeFunc: func(ctx context.Context, input httpadpt.Request, output *httpadpt.Response) error {
					output.StatusCode = tt.statusCode
					output.Header = tt.headers
					output.Body = tt.body
					return nil
				},
			}

			httpHandler := buildHandler(handler)
			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()

			httpHandler(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Status code = %d, want %d", w.Code, tt.expectedStatus)
			}

			if w.Body.String() != tt.expectedBody {
				t.Errorf("Body = %q, want %q", w.Body.String(), tt.expectedBody)
			}

			if tt.expectedHeader != nil {
				for k, expectedValues := range tt.expectedHeader {
					actualValues := w.Header()[k]
					if len(actualValues) != len(expectedValues) {
						t.Errorf("Header %q length = %d, want %d", k, len(actualValues), len(expectedValues))
					}
					for i, expectedValue := range expectedValues {
						if i < len(actualValues) && actualValues[i] != expectedValue {
							t.Errorf("Header %q[%d] = %q, want %q", k, i, actualValues[i], expectedValue)
						}
					}
				}
			}
		})
	}
}

func Test_buildHandler_Error(t *testing.T) {
	expectedError := errors.New("handler error")
	handler := &mockHandler{
		invokeFunc: func(ctx context.Context, input httpadpt.Request, output *httpadpt.Response) error {
			return expectedError
		},
	}

	httpHandler := buildHandler(handler)
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	httpHandler(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Status code = %d, want %d", w.Code, http.StatusInternalServerError)
	}

	if w.Body.String() != expectedError.Error()+"\n" {
		t.Errorf("Body = %q, want %q", w.Body.String(), expectedError.Error()+"\n")
	}
}

func Test_buildHandler_HeadersSetBeforeWriteHeader(t *testing.T) {
	// This test verifies that headers are set before WriteHeader is called
	// by checking that headers are actually present in the response
	handler := &mockHandler{
		invokeFunc: func(ctx context.Context, input httpadpt.Request, output *httpadpt.Response) error {
			output.StatusCode = intPtr(200)
			output.Header = map[string][]string{
				"Content-Type": {"application/json"},
			}
			output.Body = []byte("test")
			return nil
		},
	}

	httpHandler := buildHandler(handler)
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	httpHandler(w, req)

	// Verify header was set (if WriteHeader was called before setting headers, this would be empty)
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Content-Type header = %q, want %q", contentType, "application/json")
	}
}

func Test_buildAndAddHandles_Success(t *testing.T) {
	tests := []struct {
		name          string
		bindings      httpadpt.Bindings
		expectedPaths []string
		expectedError bool
	}{
		{
			name: "single binding with GET method",
			bindings: httpadpt.Bindings{
				{
					Condition: httpadpt.Condition{
						Path:    stringPtr("/api/users"),
						Methods: []string{"GET"},
					},
					Handler: &mockHandler{},
				},
			},
			expectedPaths: []string{"GET /api/users"},
			expectedError: false,
		},
		{
			name: "multiple bindings with different methods",
			bindings: httpadpt.Bindings{
				{
					Condition: httpadpt.Condition{
						Path:    stringPtr("/api/users"),
						Methods: []string{"GET", "POST"},
					},
					Handler: &mockHandler{},
				},
			},
			expectedPaths: []string{"GET /api/users", "POST /api/users"},
			expectedError: false,
		},
		{
			name: "binding with empty methods",
			bindings: httpadpt.Bindings{
				{
					Condition: httpadpt.Condition{
						Path:    stringPtr("/api/users"),
						Methods: []string{},
					},
					Handler: &mockHandler{},
				},
			},
			expectedPaths: []string{},
			expectedError: false,
		},
		{
			name: "multiple bindings",
			bindings: httpadpt.Bindings{
				{
					Condition: httpadpt.Condition{
						Path:    stringPtr("/api/users"),
						Methods: []string{"GET"},
					},
					Handler: &mockHandler{},
				},
				{
					Condition: httpadpt.Condition{
						Path:    stringPtr("/api/posts"),
						Methods: []string{"POST"},
					},
					Handler: &mockHandler{},
				},
			},
			expectedPaths: []string{"GET /api/users", "POST /api/posts"},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registeredPaths := make(map[string]bool)
			addHandle := func(path string, handler http.Handler) {
				registeredPaths[path] = true
			}

			err := buildAndAddHandles(addHandle, tt.bindings)

			if (err != nil) != tt.expectedError {
				t.Errorf("buildAndAddHandles() error = %v, want error = %v", err, tt.expectedError)
				return
			}

			if len(registeredPaths) != len(tt.expectedPaths) {
				t.Errorf("Registered paths count = %d, want %d", len(registeredPaths), len(tt.expectedPaths))
			}

			for _, expectedPath := range tt.expectedPaths {
				if !registeredPaths[expectedPath] {
					t.Errorf("Path %q not registered", expectedPath)
				}
			}
		})
	}
}

func Test_buildAndAddHandles_Error_NilPath(t *testing.T) {
	bindings := httpadpt.Bindings{
		{
			Condition: httpadpt.Condition{
				Path:    nil, // nil path should cause error
				Methods: []string{"GET"},
			},
			Handler: &mockHandler{},
		},
	}

	registeredPaths := make(map[string]bool)
	addHandle := func(path string, handler http.Handler) {
		registeredPaths[path] = true
	}

	err := buildAndAddHandles(addHandle, bindings)

	if err == nil {
		t.Error("buildAndAddHandles() expected error for nil path, got nil")
	}

	if len(registeredPaths) != 0 {
		t.Errorf("No paths should be registered when error occurs, got %d", len(registeredPaths))
	}
}

func Test_buildAndAddHandles_EmptyBindings(t *testing.T) {
	bindings := httpadpt.Bindings{}

	registeredPaths := make(map[string]bool)
	addHandle := func(path string, handler http.Handler) {
		registeredPaths[path] = true
	}

	err := buildAndAddHandles(addHandle, bindings)

	if err != nil {
		t.Errorf("buildAndAddHandles() error = %v, want nil", err)
	}

	if len(registeredPaths) != 0 {
		t.Errorf("No paths should be registered for empty bindings, got %d", len(registeredPaths))
	}
}

func Test_buildAndAddHandles_NilHandler(t *testing.T) {
	// Test with nil handler - should still register paths
	bindings := httpadpt.Bindings{
		{
			Condition: httpadpt.Condition{
				Path:    stringPtr("/api/test"),
				Methods: []string{"GET"},
			},
			Handler: nil, // nil handler
		},
	}

	registeredPaths := make(map[string]bool)
	var registeredHandler http.Handler
	addHandle := func(path string, handler http.Handler) {
		registeredPaths[path] = true
		registeredHandler = handler
	}

	err := buildAndAddHandles(addHandle, bindings)

	if err != nil {
		t.Errorf("buildAndAddHandles() error = %v, want nil", err)
	}

	if len(registeredPaths) != 1 {
		t.Errorf("Expected 1 path registered, got %d", len(registeredPaths))
	}

	if !registeredPaths["GET /api/test"] {
		t.Error("Expected path 'GET /api/test' to be registered")
	}

	// Handler should be created even if original was nil (buildHandler creates a wrapper)
	if registeredHandler == nil {
		t.Error("Expected handler to be created, got nil")
	}
}

func Test_buildAndAddHandles_MultipleMethodsSamePath(t *testing.T) {
	// Test multiple methods for the same path
	bindings := httpadpt.Bindings{
		{
			Condition: httpadpt.Condition{
				Path:    stringPtr("/api/users"),
				Methods: []string{"GET", "POST", "PUT", "DELETE"},
			},
			Handler: &mockHandler{},
		},
	}

	registeredPaths := make(map[string]bool)
	addHandle := func(path string, handler http.Handler) {
		registeredPaths[path] = true
	}

	err := buildAndAddHandles(addHandle, bindings)

	if err != nil {
		t.Errorf("buildAndAddHandles() error = %v, want nil", err)
	}

	expectedPaths := []string{"GET /api/users", "POST /api/users", "PUT /api/users", "DELETE /api/users"}
	if len(registeredPaths) != len(expectedPaths) {
		t.Errorf("Expected %d paths registered, got %d", len(expectedPaths), len(registeredPaths))
	}

	for _, expectedPath := range expectedPaths {
		if !registeredPaths[expectedPath] {
			t.Errorf("Expected path %q to be registered", expectedPath)
		}
	}
}

func Test_buildAndAddHandles_ErrorInMiddle(t *testing.T) {
	// Test that error stops processing and no paths are registered after error
	bindings := httpadpt.Bindings{
		{
			Condition: httpadpt.Condition{
				Path:    stringPtr("/api/valid"),
				Methods: []string{"GET"},
			},
			Handler: &mockHandler{},
		},
		{
			Condition: httpadpt.Condition{
				Path:    nil, // This will cause an error
				Methods: []string{"POST"},
			},
			Handler: &mockHandler{},
		},
		{
			Condition: httpadpt.Condition{
				Path:    stringPtr("/api/should-not-register"),
				Methods: []string{"PUT"},
			},
			Handler: &mockHandler{},
		},
	}

	registeredPaths := make(map[string]bool)
	addHandle := func(path string, handler http.Handler) {
		registeredPaths[path] = true
	}

	err := buildAndAddHandles(addHandle, bindings)

	if err == nil {
		t.Error("buildAndAddHandles() expected error for nil path, got nil")
	}

	// First binding should be registered before error
	if !registeredPaths["GET /api/valid"] {
		t.Error("Expected first valid path to be registered before error")
	}

	// Paths after error should not be registered
	if registeredPaths["PUT /api/should-not-register"] {
		t.Error("Expected path after error to not be registered")
	}
}

func Test_buildHandler_ContextPropagation(t *testing.T) {
	// Test that context is properly passed to handler
	var receivedCtx context.Context
	handler := &mockHandler{
		invokeFunc: func(ctx context.Context, input httpadpt.Request, output *httpadpt.Response) error {
			receivedCtx = ctx
			output.StatusCode = intPtr(200)
			return nil
		},
	}

	httpHandler := buildHandler(handler)
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	httpHandler(w, req)

	if receivedCtx == nil {
		t.Error("Expected context to be passed to handler, got nil")
	}

	if receivedCtx != req.Context() {
		t.Error("Expected context to be the request context")
	}
}

func Test_buildHandler_EmptyError(t *testing.T) {
	// Test with empty error message
	emptyError := errors.New("")
	handler := &mockHandler{
		invokeFunc: func(ctx context.Context, input httpadpt.Request, output *httpadpt.Response) error {
			return emptyError
		},
	}

	httpHandler := buildHandler(handler)
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	httpHandler(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Status code = %d, want %d", w.Code, http.StatusInternalServerError)
	}

	// Empty error should still result in newline
	if w.Body.String() != "\n" {
		t.Errorf("Body = %q, want %q", w.Body.String(), "\n")
	}
}

func Test_buildHandler_ResponseWithAllFields(t *testing.T) {
	// Test response with status code, headers, and body all set
	handler := &mockHandler{
		invokeFunc: func(ctx context.Context, input httpadpt.Request, output *httpadpt.Response) error {
			statusCode := 201
			output.StatusCode = &statusCode
			output.Header = map[string][]string{
				"Content-Type": {"application/json"},
				"Location":     {"/api/users/123"},
			}
			output.Body = []byte(`{"id":123,"name":"test"}`)
			return nil
		},
	}

	httpHandler := buildHandler(handler)
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	httpHandler(w, req)

	if w.Code != 201 {
		t.Errorf("Status code = %d, want %d", w.Code, 201)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Content-Type header = %q, want %q", w.Header().Get("Content-Type"), "application/json")
	}

	if w.Header().Get("Location") != "/api/users/123" {
		t.Errorf("Location header = %q, want %q", w.Header().Get("Location"), "/api/users/123")
	}

	expectedBody := `{"id":123,"name":"test"}`
	if w.Body.String() != expectedBody {
		t.Errorf("Body = %q, want %q", w.Body.String(), expectedBody)
	}
}

func Test_buildPath_AllHTTPMethods(t *testing.T) {
	// Test all common HTTP methods
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}
	path := "/api/resource"

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			result := buildPath(method, path)
			expected := method + " " + path
			if result != expected {
				t.Errorf("buildPath(%q, %q) = %q, want %q", method, path, result, expected)
			}
		})
	}
}

// Helper functions
func intPtr(i int) *int {
	return &i
}

func stringPtr(s string) *string {
	return &s
}
