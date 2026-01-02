package gonethttp

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRequest(t *testing.T) {
	tests := []struct {
		name    string
		httpReq *http.Request
		wantNil bool
	}{
		{
			name:    "valid request",
			httpReq: httptest.NewRequest("GET", "/test", nil),
			wantNil: false,
		},
		{
			name:    "nil request",
			httpReq: nil,
			wantNil: false, // Should return a Request struct, not nil
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := NewRequest(tt.httpReq)
			if req == nil {
				t.Error("NewRequest() returned nil")
			}
		})
	}
}

func TestRequest_Query(t *testing.T) {
	tests := []struct {
		name           string
		httpReq        *http.Request
		queryParamName string
		expectedValue  []string
		expectedFound  bool
	}{
		{
			name:           "query param exists",
			httpReq:        httptest.NewRequest("GET", "/test?name=value", nil),
			queryParamName: "name",
			expectedValue:  []string{"value"},
			expectedFound:  true,
		},
		{
			name:           "query param with multiple values",
			httpReq:        httptest.NewRequest("GET", "/test?tags=go&tags=test", nil),
			queryParamName: "tags",
			expectedValue:  []string{"go", "test"},
			expectedFound:  true,
		},
		{
			name:           "query param not found",
			httpReq:        httptest.NewRequest("GET", "/test?name=value", nil),
			queryParamName: "missing",
			expectedValue:  nil,
			expectedFound:  false,
		},
		{
			name:           "no query params",
			httpReq:        httptest.NewRequest("GET", "/test", nil),
			queryParamName: "name",
			expectedValue:  nil,
			expectedFound:  false,
		},
		{
			name:           "empty query param value",
			httpReq:        httptest.NewRequest("GET", "/test?name=", nil),
			queryParamName: "name",
			expectedValue:  []string{""},
			expectedFound:  true,
		},
		{
			name:           "multiple query params",
			httpReq:        httptest.NewRequest("GET", "/test?name=value&age=25", nil),
			queryParamName: "age",
			expectedValue:  []string{"25"},
			expectedFound:  true,
		},
		{
			name:           "nil request",
			httpReq:        nil,
			queryParamName: "name",
			expectedValue:  nil,
			expectedFound:  false,
		},
		{
			name:           "query param with special characters",
			httpReq:        httptest.NewRequest("GET", "/test?q=hello%20world", nil),
			queryParamName: "q",
			expectedValue:  []string{"hello world"},
			expectedFound:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := NewRequest(tt.httpReq)
			queryParams := req.Query()
			value, found := queryParams.GetValue(tt.queryParamName)

			if found != tt.expectedFound {
				t.Errorf("GetValue() found = %v, want %v", found, tt.expectedFound)
			}

			if len(value) != len(tt.expectedValue) {
				t.Errorf("GetValue() value length = %d, want %d", len(value), len(tt.expectedValue))
			}

			for i, v := range value {
				if i < len(tt.expectedValue) && v != tt.expectedValue[i] {
					t.Errorf("GetValue() value[%d] = %q, want %q", i, v, tt.expectedValue[i])
				}
			}
		})
	}
}

func TestRequest_Header(t *testing.T) {
	tests := []struct {
		name          string
		httpReq       *http.Request
		headerName    string
		expectedValue []string
		expectedFound bool
	}{
		{
			name:          "header exists",
			httpReq:       createRequestWithHeader("Content-Type", "application/json"),
			headerName:    "Content-Type",
			expectedValue: []string{"application/json"},
			expectedFound: true,
		},
		{
			name:          "header with multiple values",
			httpReq:       createRequestWithHeaders("Accept", []string{"application/json", "application/xml"}),
			headerName:    "Accept",
			expectedValue: []string{"application/json", "application/xml"},
			expectedFound: true,
		},
		{
			name:          "header not found",
			httpReq:       createRequestWithHeader("Content-Type", "application/json"),
			headerName:    "X-Custom",
			expectedValue: nil,
			expectedFound: false,
		},
		{
			name:          "no headers",
			httpReq:       httptest.NewRequest("GET", "/test", nil),
			headerName:    "Content-Type",
			expectedValue: nil,
			expectedFound: false,
		},
		{
			name:          "empty header value",
			httpReq:       createRequestWithHeader("X-Empty", ""),
			headerName:    "X-Empty",
			expectedValue: []string{""},
			expectedFound: true,
		},
		{
			name:          "nil request",
			httpReq:       nil,
			headerName:    "Content-Type",
			expectedValue: nil,
			expectedFound: false,
		},
		{
			name:          "case insensitive header lookup",
			httpReq:       createRequestWithHeader("Content-Type", "application/json"),
			headerName:    "content-type",
			expectedValue: []string{"application/json"},
			expectedFound: true,
		},
		{
			name:          "Authorization header",
			httpReq:       createRequestWithHeader("Authorization", "Bearer token123"),
			headerName:    "Authorization",
			expectedValue: []string{"Bearer token123"},
			expectedFound: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := NewRequest(tt.httpReq)
			headerParams := req.Header()
			value, found := headerParams.GetValue(tt.headerName)

			if found != tt.expectedFound {
				t.Errorf("GetValue() found = %v, want %v", found, tt.expectedFound)
			}

			if len(value) != len(tt.expectedValue) {
				t.Errorf("GetValue() value length = %d, want %d", len(value), len(tt.expectedValue))
			}

			for i, v := range value {
				if i < len(tt.expectedValue) && v != tt.expectedValue[i] {
					t.Errorf("GetValue() value[%d] = %q, want %q", i, v, tt.expectedValue[i])
				}
			}
		})
	}
}

func TestRequest_Path(t *testing.T) {
	tests := []struct {
		name          string
		httpReq       *http.Request
		pathParamName string
		expectedValue string
		expectedFound bool
	}{
		{
			name:          "path param exists",
			httpReq:       createRequestWithPathValue("id", "123"),
			pathParamName: "id",
			expectedValue: "123",
			expectedFound: true,
		},
		{
			name:          "path param not found",
			httpReq:       createRequestWithPathValue("id", "123"),
			pathParamName: "missing",
			expectedValue: "",
			expectedFound: false,
		},
		{
			name:          "no path params",
			httpReq:       httptest.NewRequest("GET", "/test", nil),
			pathParamName: "id",
			expectedValue: "",
			expectedFound: false,
		},
		{
			name:          "empty path param value (treated as not found)",
			httpReq:       createRequestWithPathValue("id", ""),
			pathParamName: "id",
			expectedValue: "",
			expectedFound: false, // Empty values are treated as not found due to PathValue limitation
		},
		{
			name:          "nil request",
			httpReq:       nil,
			pathParamName: "id",
			expectedValue: "",
			expectedFound: false,
		},
		{
			name:          "path param with special characters",
			httpReq:       createRequestWithPathValue("slug", "hello-world-123"),
			pathParamName: "slug",
			expectedValue: "hello-world-123",
			expectedFound: true,
		},
		{
			name:          "numeric path param",
			httpReq:       createRequestWithPathValue("id", "456"),
			pathParamName: "id",
			expectedValue: "456",
			expectedFound: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := NewRequest(tt.httpReq)
			pathParams := req.Path()
			value, found := pathParams.GetValue(tt.pathParamName)

			if found != tt.expectedFound {
				t.Errorf("GetValue() found = %v, want %v", found, tt.expectedFound)
			}

			if value != tt.expectedValue {
				t.Errorf("GetValue() value = %q, want %q", value, tt.expectedValue)
			}
		})
	}
}

func TestRequest_AllMethods(t *testing.T) {
	// Test that all methods work together
	httpReq := httptest.NewRequest("GET", "/test?name=value", nil)
	httpReq.Header.Set("Content-Type", "application/json")
	// Note: PathValue requires a mux pattern, so we'll test it separately

	req := NewRequest(httpReq)

	// Test Query
	queryParams := req.Query()
	queryValue, queryFound := queryParams.GetValue("name")
	if !queryFound || len(queryValue) == 0 || queryValue[0] != "value" {
		t.Errorf("Query() GetValue() = %v, %v, want [value], true", queryValue, queryFound)
	}

	// Test Header
	headerParams := req.Header()
	headerValue, headerFound := headerParams.GetValue("Content-Type")
	if !headerFound || len(headerValue) == 0 || headerValue[0] != "application/json" {
		t.Errorf("Header() GetValue() = %v, %v, want [application/json], true", headerValue, headerFound)
	}

	// Test Path (will be empty for standard request)
	pathParams := req.Path()
	pathValue, pathFound := pathParams.GetValue("id")
	if pathFound {
		t.Errorf("Path() GetValue() found = true, want false for standard request")
	}
	if pathValue != "" {
		t.Errorf("Path() GetValue() value = %q, want empty string", pathValue)
	}
}

func TestQuery_GetValue_EmptyMap(t *testing.T) {
	// Test query with empty map (no query params)
	httpReq := httptest.NewRequest("GET", "/test", nil)
	req := NewRequest(httpReq)
	queryParams := req.Query()

	value, found := queryParams.GetValue("any")
	if found {
		t.Errorf("GetValue() found = true, want false for empty query")
	}
	if value != nil {
		t.Errorf("GetValue() value = %v, want nil", value)
	}
}

func TestHeader_GetValue_EmptyMap(t *testing.T) {
	// Test header with empty map (no headers)
	httpReq := httptest.NewRequest("GET", "/test", nil)
	req := NewRequest(httpReq)
	headerParams := req.Header()

	value, found := headerParams.GetValue("any")
	if found {
		t.Errorf("GetValue() found = true, want false for empty headers")
	}
	if value != nil {
		t.Errorf("GetValue() value = %v, want nil", value)
	}
}

func TestHeader_GetValue_NilHeader(t *testing.T) {
	// Test header with nil header (when Request.httpReq is nil)
	req := NewRequest(nil)
	headerParams := req.Header()

	value, found := headerParams.GetValue("any")
	if found {
		t.Errorf("GetValue() found = true, want false for nil header")
	}
	if value != nil {
		t.Errorf("GetValue() value = %v, want nil", value)
	}
}

func TestQuery_GetValue_NilURL(t *testing.T) {
	// Test query with nil URL (when Request.httpReq is nil)
	req := NewRequest(nil)
	queryParams := req.Query()

	value, found := queryParams.GetValue("any")
	if found {
		t.Errorf("GetValue() found = true, want false for nil URL")
	}
	if value != nil {
		t.Errorf("GetValue() value = %v, want nil", value)
	}
}

func TestPath_GetValue_NilRequest(t *testing.T) {
	// Test path with nil request
	req := NewRequest(nil)
	pathParams := req.Path()

	value, found := pathParams.GetValue("any")
	if found {
		t.Errorf("GetValue() found = true, want false for nil request")
	}
	if value != "" {
		t.Errorf("GetValue() value = %q, want empty string", value)
	}
}

// Helper functions

func createRequestWithHeader(name, value string) *http.Request {
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set(name, value)
	return req
}

func createRequestWithHeaders(name string, values []string) *http.Request {
	req := httptest.NewRequest("GET", "/test", nil)
	for _, v := range values {
		req.Header.Add(name, v)
	}
	return req
}

func createRequestWithPathValue(name, value string) *http.Request {
	req := httptest.NewRequest("GET", "/test", nil)
	// PathValue requires a mux pattern to be set up
	// For testing, we'll use a workaround by setting the path value directly
	// This is a limitation of testing PathValue without a full mux setup
	// In real usage, PathValue is populated by the mux router
	req.SetPathValue(name, value)
	return req
}
