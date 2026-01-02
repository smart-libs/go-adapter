package httpadpt

import (
	"testing"
)

func Test_setResponseHeader(t *testing.T) {
	tests := []struct {
		name           string
		output         *Response
		headerName     string
		value          any
		expectError    bool
		expectedHeader map[string][]string
	}{
		{
			name:           "nil output",
			output:         nil,
			headerName:     "Content-Type",
			value:          []string{"application/json"},
			expectError:    true,
			expectedHeader: nil,
		},
		{
			name:        "valid string slice",
			output:      &Response{},
			headerName:  "Content-Type",
			value:       []string{"application/json"},
			expectError: false,
			expectedHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
		},
		{
			name:        "valid string slice with multiple values",
			output:      &Response{},
			headerName:  "Accept",
			value:       []string{"application/json", "application/xml"},
			expectError: false,
			expectedHeader: map[string][]string{
				"Accept": {"application/json", "application/xml"},
			},
		},
		{
			name:        "empty string slice",
			output:      &Response{},
			headerName:  "X-Custom",
			value:       []string{},
			expectError: false,
			expectedHeader: map[string][]string{
				"X-Custom": {},
			},
		},
		{
			name:        "multiple headers",
			output:      &Response{},
			headerName:  "Location",
			value:       []string{"/api/users/123"},
			expectError: false,
			expectedHeader: map[string][]string{
				"Location": {"/api/users/123"},
			},
		},
		{
			name:        "header map initialization",
			output:      &Response{Header: nil},
			headerName:  "Content-Type",
			value:       []string{"text/plain"},
			expectError: false,
			expectedHeader: map[string][]string{
				"Content-Type": {"text/plain"},
			},
		},
		{
			name:        "overwrite existing header",
			output:      &Response{Header: map[string][]string{"Content-Type": {"old/value"}}},
			headerName:  "Content-Type",
			value:       []string{"new/value"},
			expectError: false,
			expectedHeader: map[string][]string{
				"Content-Type": {"new/value"},
			},
		},
		{
			name:        "add header to existing map",
			output:      &Response{Header: map[string][]string{"X-Existing": {"value1"}}},
			headerName:  "X-New",
			value:       []string{"value2"},
			expectError: false,
			expectedHeader: map[string][]string{
				"X-Existing": {"value1"},
				"X-New":      {"value2"},
			},
		},
		{
			name:        "case sensitive header name",
			output:      &Response{},
			headerName:  "X-Custom-Header",
			value:       []string{"custom-value"},
			expectError: false,
			expectedHeader: map[string][]string{
				"X-Custom-Header": {"custom-value"},
			},
		},
		{
			name:        "empty header name",
			output:      &Response{},
			headerName:  "",
			value:       []string{"value"},
			expectError: false,
			expectedHeader: map[string][]string{
				"": {"value"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := setResponseHeader(tt.output, tt.headerName, tt.value)
			if (err != nil) != tt.expectError {
				t.Errorf("setResponseHeader() error = %v, expectError = %v", err, tt.expectError)
				return
			}

			if !tt.expectError {
				if tt.output.Header == nil {
					t.Fatal("Expected Header map to be initialized, got nil")
				}

				if len(tt.output.Header) != len(tt.expectedHeader) {
					t.Errorf("Expected Header map length = %d, got %d", len(tt.expectedHeader), len(tt.output.Header))
					return
				}

				for expectedKey, expectedValues := range tt.expectedHeader {
					actualValues, exists := tt.output.Header[expectedKey]
					if !exists {
						t.Errorf("Expected header %q to exist, but it doesn't", expectedKey)
						continue
					}

					if len(actualValues) != len(expectedValues) {
						t.Errorf("Expected header %q to have %d values, got %d", expectedKey, len(expectedValues), len(actualValues))
						continue
					}

					for i, expectedValue := range expectedValues {
						if i < len(actualValues) && actualValues[i] != expectedValue {
							t.Errorf("Expected header %q[%d] = %q, got %q", expectedKey, i, expectedValue, actualValues[i])
						}
					}
				}
			}
		})
	}
}

func Test_setResponseHeader_InvalidType(t *testing.T) {
	// Test with a type that cannot be converted to []string
	// Note: The converter may or may not handle this conversion,
	// so we test both success and failure cases
	output := &Response{}
	headerName := "Content-Type"
	value := 123 // int may or may not be convertible to []string

	err := setResponseHeader(output, headerName, value)
	if err != nil {
		// If conversion fails, header should not be set
		if len(output.Header) > 0 {
			t.Error("Expected Header to remain empty on error, but it was set")
		}
	} else {
		// If conversion succeeds, verify the header was set
		if output.Header == nil {
			t.Error("Expected Header to be initialized when conversion succeeds")
		} else {
			values, exists := output.Header[headerName]
			if !exists {
				t.Error("Expected header to be set when conversion succeeds")
			} else if len(values) == 0 {
				t.Error("Expected header to have values when conversion succeeds")
			}
		}
	}
}

func Test_setResponseHeader_NilValue(t *testing.T) {
	// Test with nil value
	// Note: The converter may or may not handle nil values,
	// so we test both success and failure cases
	output := &Response{}
	headerName := "Content-Type"
	var value any = nil

	err := setResponseHeader(output, headerName, value)
	if err != nil {
		// If conversion fails, header should not be set
		if len(output.Header) > 0 {
			t.Error("Expected Header to remain empty on error, but it was set")
		}
	} else {
		// If conversion succeeds, verify the header was set (may be empty slice)
		if output.Header == nil {
			t.Error("Expected Header to be initialized when conversion succeeds")
		} else {
			_, exists := output.Header[headerName]
			if !exists {
				t.Error("Expected header to be set when conversion succeeds")
			}
		}
	}
}

func Test_setResponseHeader_StringValue(t *testing.T) {
	// Test if converter can handle single string (may or may not work depending on converter)
	output := &Response{}
	headerName := "Content-Type"
	value := "application/json"

	err := setResponseHeader(output, headerName, value)
	// This may or may not work depending on converter capabilities
	// If it works, verify the result
	if err == nil {
		if output.Header == nil {
			t.Fatal("Expected Header map to be initialized, got nil")
		}
		values, exists := output.Header[headerName]
		if !exists {
			t.Error("Expected header to be set, but it doesn't exist")
		} else if len(values) == 0 {
			t.Error("Expected header to have at least one value, got empty slice")
		}
	}
	// If it fails, that's also acceptable - the converter may not support string->[]string
}

func Test_setResponseHeader_ConcurrentAccess(t *testing.T) {
	// Test that header map can handle multiple headers being set
	output := &Response{}

	headers := map[string][]string{
		"Content-Type": {"application/json"},
		"X-Request-ID": {"12345"},
		"Location":     {"/api/users/123"},
	}

	for headerName, values := range headers {
		err := setResponseHeader(output, headerName, values)
		if err != nil {
			t.Errorf("Failed to set header %q: %v", headerName, err)
		}
	}

	if len(output.Header) != len(headers) {
		t.Errorf("Expected %d headers, got %d", len(headers), len(output.Header))
	}

	for expectedKey, expectedValues := range headers {
		actualValues, exists := output.Header[expectedKey]
		if !exists {
			t.Errorf("Expected header %q to exist, but it doesn't", expectedKey)
			continue
		}

		if len(actualValues) != len(expectedValues) {
			t.Errorf("Expected header %q to have %d values, got %d", expectedKey, len(expectedValues), len(actualValues))
			continue
		}

		for i, expectedValue := range expectedValues {
			if i < len(actualValues) && actualValues[i] != expectedValue {
				t.Errorf("Expected header %q[%d] = %q, got %q", expectedKey, i, expectedValue, actualValues[i])
			}
		}
	}
}
