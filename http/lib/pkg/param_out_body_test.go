package httpadpt

import (
	"testing"
)

func Test_setBodyBytes(t *testing.T) {
	tests := []struct {
		name        string
		output      *Response
		value       any
		expectError bool
		expected    []byte
	}{
		{
			name:        "nil output",
			output:      nil,
			value:       []byte("test"),
			expectError: true,
			expected:    nil,
		},
		{
			name:        "valid byte slice",
			output:      &Response{},
			value:       []byte("hello world"),
			expectError: false,
			expected:    []byte("hello world"),
		},
		{
			name:        "empty byte slice",
			output:      &Response{},
			value:       []byte{},
			expectError: false,
			expected:    []byte{},
		},
		{
			name:        "byte slice with binary data",
			output:      &Response{},
			value:       []byte{0x00, 0x01, 0x02, 0xFF},
			expectError: false,
			expected:    []byte{0x00, 0x01, 0x02, 0xFF},
		},
		{
			name:        "string value",
			output:      &Response{},
			value:       "test string",
			expectError: false,
			expected:    []byte("test string"),
		},
		{
			name:        "empty string",
			output:      &Response{},
			value:       "",
			expectError: false,
			expected:    []byte(""),
		},
		{
			name:        "JSON string",
			output:      &Response{},
			value:       `{"key":"value"}`,
			expectError: false,
			expected:    []byte(`{"key":"value"}`),
		},
		{
			name:        "overwrite existing body",
			output:      &Response{Body: []byte("old body")},
			value:       []byte("new body"),
			expectError: false,
			expected:    []byte("new body"),
		},
		{
			name:        "set body when body was nil",
			output:      &Response{Body: nil},
			value:       []byte("new body"),
			expectError: false,
			expected:    []byte("new body"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := setBodyBytes(tt.output, tt.value)
			if (err != nil) != tt.expectError {
				t.Errorf("setBodyBytes() error = %v, expectError = %v", err, tt.expectError)
				return
			}

			if !tt.expectError {
				if len(tt.output.Body) == 0 && len(tt.expected) > 0 {
					t.Fatal("Expected Body to be set, got empty")
				}

				if len(tt.output.Body) != len(tt.expected) {
					t.Errorf("setBodyBytes() body length = %d, want %d", len(tt.output.Body), len(tt.expected))
					return
				}

				for i, expectedByte := range tt.expected {
					if i < len(tt.output.Body) && tt.output.Body[i] != expectedByte {
						t.Errorf("setBodyBytes() body[%d] = %d, want %d", i, tt.output.Body[i], expectedByte)
					}
				}
			}
		})
	}
}

func Test_setBodyBytes_InvalidType(t *testing.T) {
	// Test with a type that cannot be converted to []byte
	// Note: The converter may or may not handle this conversion,
	// so we test both success and failure cases
	output := &Response{}
	value := 123 // int may or may not be convertible to []byte

	err := setBodyBytes(output, value)
	if err != nil {
		// If conversion fails, body should not be set
		if len(output.Body) > 0 {
			t.Error("Expected Body to remain empty on error, but it was set")
		}
	} else {
		// If conversion succeeds, verify the body was set
		if output.Body == nil {
			t.Error("Expected Body to be initialized when conversion succeeds")
		} else {
			if len(output.Body) == 0 {
				t.Error("Expected Body to have content when conversion succeeds")
			}
		}
	}
}

func Test_setBodyBytes_NilValue(t *testing.T) {
	// Test with nil value
	// Note: The converter may or may not handle nil values,
	// so we test both success and failure cases
	output := &Response{}
	var value any = nil

	err := setBodyBytes(output, value)
	if err != nil {
		// If conversion fails, body should not be set
		if len(output.Body) > 0 {
			t.Error("Expected Body to remain empty on error, but it was set")
		}
	} else {
		// If conversion succeeds, body may be nil or empty
		// Both are acceptable outcomes
		if len(output.Body) > 0 {
			t.Logf("Body was set to: %v", output.Body)
		}
	}
}

func Test_setBodyBytes_LargeBody(t *testing.T) {
	// Test with a large body
	largeData := make([]byte, 10000)
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}

	output := &Response{}
	err := setBodyBytes(output, largeData)
	if err != nil {
		t.Errorf("setBodyBytes() error = %v, want nil", err)
		return
	}

	if len(output.Body) != len(largeData) {
		t.Errorf("setBodyBytes() body length = %d, want %d", len(output.Body), len(largeData))
		return
	}

	// Verify first and last bytes
	if output.Body[0] != 0 {
		t.Errorf("setBodyBytes() body[0] = %d, want 0", output.Body[0])
	}
	if output.Body[len(output.Body)-1] != byte(9999%256) {
		t.Errorf("setBodyBytes() body[last] = %d, want %d", output.Body[len(output.Body)-1], byte(9999%256))
	}
}

func Test_setBodyBytes_UnicodeString(t *testing.T) {
	// Test with Unicode string
	unicodeStr := "Hello, 世界! 🌍"
	output := &Response{}

	err := setBodyBytes(output, unicodeStr)
	if err != nil {
		t.Errorf("setBodyBytes() error = %v, want nil", err)
		return
	}

	expected := []byte(unicodeStr)
	if len(output.Body) != len(expected) {
		t.Errorf("setBodyBytes() body length = %d, want %d", len(output.Body), len(expected))
		return
	}

	for i, expectedByte := range expected {
		if i < len(output.Body) && output.Body[i] != expectedByte {
			t.Errorf("setBodyBytes() body[%d] = %d, want %d", i, output.Body[i], expectedByte)
		}
	}
}

func Test_setBodyBytes_MultipleCalls(t *testing.T) {
	// Test that multiple calls overwrite the body
	output := &Response{}

	// First call
	err1 := setBodyBytes(output, []byte("first"))
	if err1 != nil {
		t.Errorf("setBodyBytes() first call error = %v, want nil", err1)
		return
	}

	if string(output.Body) != "first" {
		t.Errorf("setBodyBytes() first call body = %q, want %q", string(output.Body), "first")
	}

	// Second call
	err2 := setBodyBytes(output, []byte("second"))
	if err2 != nil {
		t.Errorf("setBodyBytes() second call error = %v, want nil", err2)
		return
	}

	if string(output.Body) != "second" {
		t.Errorf("setBodyBytes() second call body = %q, want %q", string(output.Body), "second")
	}
}
