package httpadpt

import (
	"testing"
)

func TestIsResponseNil(t *testing.T) {
	tests := []struct {
		name     string
		resp     *Response
		expected bool
	}{
		{
			name:     "nil response",
			resp:     nil,
			expected: true,
		},
		{
			name:     "non-nil response",
			resp:     &Response{},
			expected: false,
		},
		{
			name: "response with status code",
			resp: &Response{
				StatusCode: intPtr(200),
			},
			expected: false,
		},
		{
			name: "response with body",
			resp: &Response{
				Body: []byte("test"),
			},
			expected: false,
		},
		{
			name: "response with headers",
			resp: &Response{
				Header: map[string][]string{
					"Content-Type": {"application/json"},
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := IsResponseNil(tt.resp)
			if (err != nil) != tt.expected {
				t.Errorf("IsResponseNil() error = %v, expected error = %v", err, tt.expected)
			}
		})
	}
}
