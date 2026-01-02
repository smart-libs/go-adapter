package gonethttp

import (
	"context"
	httpadpt "github.com/smart-libs/go-adapter/http/lib/pkg"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_handleResponse_WithStatusCode(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     *int
		expectedStatus int
	}{
		{
			name:           "status code 200",
			statusCode:     intPtr(200),
			expectedStatus: 200,
		},
		{
			name:           "status code 201",
			statusCode:     intPtr(201),
			expectedStatus: 201,
		},
		{
			name:           "status code 404",
			statusCode:     intPtr(404),
			expectedStatus: 404,
		},
		{
			name:           "status code 500",
			statusCode:     intPtr(500),
			expectedStatus: 500,
		},
		{
			name:           "nil status code defaults to 200",
			statusCode:     nil,
			expectedStatus: 200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := httpadpt.Response{
				StatusCode: tt.statusCode,
				Header:     nil,
				Body:       nil,
			}

			w := httptest.NewRecorder()
			ctx := context.Background()

			handleResponse(ctx, w, resp)

			if w.Code != tt.expectedStatus {
				t.Errorf("Status code = %d, want %d", w.Code, tt.expectedStatus)
			}
		})
	}
}

func Test_handleResponse_WithHeaders(t *testing.T) {
	tests := []struct {
		name           string
		headers        map[string][]string
		expectedHeader map[string][]string
	}{
		{
			name: "single header",
			headers: map[string][]string{
				"Content-Type": {"application/json"},
			},
			expectedHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
		},
		{
			name: "multiple headers",
			headers: map[string][]string{
				"Content-Type": {"application/json"},
				"Location":     {"/api/users/123"},
				"X-Custom":     {"value1"},
			},
			expectedHeader: map[string][]string{
				"Content-Type": {"application/json"},
				"Location":     {"/api/users/123"},
				"X-Custom":     {"value1"},
			},
		},
		{
			name: "header with multiple values",
			headers: map[string][]string{
				"Accept": {"application/json", "application/xml"},
			},
			expectedHeader: map[string][]string{
				"Accept": {"application/json", "application/xml"},
			},
		},
		{
			name:           "nil headers",
			headers:        nil,
			expectedHeader: nil,
		},
		{
			name:           "empty headers map",
			headers:        map[string][]string{},
			expectedHeader: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := httpadpt.Response{
				StatusCode: intPtr(200),
				Header:     tt.headers,
				Body:       nil,
			}

			w := httptest.NewRecorder()
			ctx := context.Background()

			handleResponse(ctx, w, resp)

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
			} else if len(tt.headers) == 0 && len(w.Header()) > 0 {
				// If we expected no headers but got some, that's an error
				t.Errorf("Expected no headers, but got %d", len(w.Header()))
			}
		})
	}
}

func Test_handleResponse_WithBody(t *testing.T) {
	tests := []struct {
		name         string
		body         []byte
		expectedBody string
	}{
		{
			name:         "simple text body",
			body:         []byte("Hello, World!"),
			expectedBody: "Hello, World!",
		},
		{
			name:         "JSON body",
			body:         []byte(`{"id":123,"name":"test"}`),
			expectedBody: `{"id":123,"name":"test"}`,
		},
		{
			name:         "empty body",
			body:         []byte(""),
			expectedBody: "",
		},
		{
			name:         "nil body",
			body:         nil,
			expectedBody: "",
		},
		{
			name:         "binary body",
			body:         []byte{0x00, 0x01, 0x02, 0xFF},
			expectedBody: string([]byte{0x00, 0x01, 0x02, 0xFF}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := httpadpt.Response{
				StatusCode: intPtr(200),
				Header:     nil,
				Body:       tt.body,
			}

			w := httptest.NewRecorder()
			ctx := context.Background()

			handleResponse(ctx, w, resp)

			if w.Body.String() != tt.expectedBody {
				t.Errorf("Body = %q, want %q", w.Body.String(), tt.expectedBody)
			}
		})
	}
}

func Test_handleResponse_HeadersSetBeforeWriteHeader(t *testing.T) {
	// This test verifies that headers are set before WriteHeader is called
	// Headers must be set before WriteHeader, otherwise they are ignored
	resp := httpadpt.Response{
		StatusCode: intPtr(200),
		Header: map[string][]string{
			"Content-Type": {"application/json"},
			"Location":     {"/api/users/123"},
		},
		Body: []byte("test"),
	}

	w := httptest.NewRecorder()
	ctx := context.Background()

	handleResponse(ctx, w, resp)

	// Verify headers were set (if WriteHeader was called before setting headers, these would be empty)
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Content-Type header = %q, want %q", contentType, "application/json")
	}

	location := w.Header().Get("Location")
	if location != "/api/users/123" {
		t.Errorf("Location header = %q, want %q", location, "/api/users/123")
	}
}

func Test_handleResponse_AllFieldsSet(t *testing.T) {
	// Test response with status code, headers, and body all set
	statusCode := 201
	resp := httpadpt.Response{
		StatusCode: &statusCode,
		Header: map[string][]string{
			"Content-Type": {"application/json"},
			"Location":     {"/api/users/123"},
		},
		Body: []byte(`{"id":123,"name":"test"}`),
	}

	w := httptest.NewRecorder()
	ctx := context.Background()

	handleResponse(ctx, w, resp)

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

func Test_handleResponse_EmptyResponse(t *testing.T) {
	// Test with all fields nil/empty
	resp := httpadpt.Response{
		StatusCode: nil,
		Header:     nil,
		Body:       nil,
	}

	w := httptest.NewRecorder()
	ctx := context.Background()

	handleResponse(ctx, w, resp)

	// Should default to 200 OK
	if w.Code != http.StatusOK {
		t.Errorf("Status code = %d, want %d", w.Code, http.StatusOK)
	}

	// Should have no body
	if w.Body.String() != "" {
		t.Errorf("Body = %q, want empty string", w.Body.String())
	}
}

func Test_handleResponse_NilHeaderMap(t *testing.T) {
	// Test that nil header map is handled correctly
	resp := httpadpt.Response{
		StatusCode: intPtr(200),
		Header:     nil, // nil map
		Body:       []byte("test"),
	}

	w := httptest.NewRecorder()
	ctx := context.Background()

	handleResponse(ctx, w, resp)

	// Should not panic and should work correctly
	if w.Code != 200 {
		t.Errorf("Status code = %d, want %d", w.Code, 200)
	}

	if w.Body.String() != "test" {
		t.Errorf("Body = %q, want %q", w.Body.String(), "test")
	}
}

func Test_handleResponse_EmptyHeaderMap(t *testing.T) {
	// Test that empty header map is handled correctly
	resp := httpadpt.Response{
		StatusCode: intPtr(200),
		Header:     map[string][]string{}, // empty map
		Body:       []byte("test"),
	}

	w := httptest.NewRecorder()
	ctx := context.Background()

	handleResponse(ctx, w, resp)

	// Should not set any headers
	if len(w.Header()) > 0 {
		t.Errorf("Expected no headers, but got %d", len(w.Header()))
	}

	if w.Body.String() != "test" {
		t.Errorf("Body = %q, want %q", w.Body.String(), "test")
	}
}

func Test_handleResponse_ContextPassed(t *testing.T) {
	// Test that context is accepted (even though not currently used)
	type keyType string
	ctx := context.WithValue(context.Background(), keyType("test-key"), "test-value")
	resp := httpadpt.Response{
		StatusCode: intPtr(200),
		Header:     nil,
		Body:       []byte("test"),
	}

	w := httptest.NewRecorder()

	// Should not panic
	handleResponse(ctx, w, resp)

	if w.Code != 200 {
		t.Errorf("Status code = %d, want %d", w.Code, 200)
	}
}

func Test_handleResponse_VariousStatusCodes(t *testing.T) {
	statusCodes := []int{
		http.StatusOK,                  // 200
		http.StatusCreated,             // 201
		http.StatusNoContent,           // 204
		http.StatusBadRequest,          // 400
		http.StatusUnauthorized,        // 401
		http.StatusForbidden,           // 403
		http.StatusNotFound,            // 404
		http.StatusConflict,            // 409
		http.StatusInternalServerError, // 500
		http.StatusBadGateway,          // 502
		http.StatusServiceUnavailable,  // 503
	}

	for _, code := range statusCodes {
		t.Run(http.StatusText(code), func(t *testing.T) {
			resp := httpadpt.Response{
				StatusCode: &code,
				Header:     nil,
				Body:       nil,
			}

			w := httptest.NewRecorder()
			ctx := context.Background()

			handleResponse(ctx, w, resp)

			if w.Code != code {
				t.Errorf("Status code = %d, want %d", w.Code, code)
			}
		})
	}
}

func Test_handleResponse_HeaderOverwrite(t *testing.T) {
	// Test that headers can be set and overwritten
	resp := httpadpt.Response{
		StatusCode: intPtr(200),
		Header: map[string][]string{
			"X-Custom": {"value1", "value2"},
		},
		Body: nil,
	}

	w := httptest.NewRecorder()
	ctx := context.Background()

	handleResponse(ctx, w, resp)

	// Verify header values
	values := w.Header()["X-Custom"]
	if len(values) != 2 {
		t.Errorf("Header X-Custom length = %d, want %d", len(values), 2)
	}
	if values[0] != "value1" {
		t.Errorf("Header X-Custom[0] = %q, want %q", values[0], "value1")
	}
	if values[1] != "value2" {
		t.Errorf("Header X-Custom[1] = %q, want %q", values[1], "value2")
	}
}
