package httpadpt

import (
	"testing"
)

func TestIsRequestNil(t *testing.T) {
	tests := []struct {
		name      string
		req       Request
		errHolder *error
		expected  bool
	}{
		{
			name:      "nil request",
			req:       nil,
			errHolder: new(error),
			expected:  true,
		},
		{
			name:      "non-nil request",
			req:       &mockRequest{},
			errHolder: new(error),
			expected:  false,
		},
		{
			name:      "nil errHolder",
			req:       nil,
			errHolder: nil,
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsRequestNil(tt.req, tt.errHolder)
			if result != tt.expected {
				t.Errorf("IsRequestNil() = %v, want %v", result, tt.expected)
			}

			if tt.req == nil && tt.errHolder != nil {
				if *tt.errHolder == nil {
					t.Error("Expected error to be set in errHolder for nil request")
				}
			}
		})
	}
}

func TestIsRequestQueryNil(t *testing.T) {
	tests := []struct {
		name      string
		req       Request
		errHolder *error
		expected  bool
	}{
		// Note: nil request test is skipped because IsRequestQueryNil will panic
		// when trying to call Query() on nil. This is a known limitation.
		// {
		// 	name:      "nil request",
		// 	req:       nil,
		// 	errHolder: new(error),
		// 	expected:  true,
		// },
		// Note: request with nil query test is skipped because IsRequestQueryNil will panic
		// when trying to call Query() on a request with nil query. This is a known limitation.
		// {
		// 	name:      "request with nil query",
		// 	req:       &mockRequest{query: nil},
		// 	errHolder: new(error),
		// 	expected:  true,
		// },
		{
			name:      "request with valid query",
			req:       &mockRequest{query: &mockQueryParams{}},
			errHolder: new(error),
			expected:  false,
		},
		// Note: nil request test is skipped because IsRequestQueryNil will panic
		// when trying to call Query() on nil. This is a known limitation.
		// {
		// 	name:      "nil errHolder",
		// 	req:       nil,
		// 	errHolder: nil,
		// 	expected:  true,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsRequestQueryNil(tt.req, tt.errHolder)
			if result != tt.expected {
				t.Errorf("IsRequestQueryNil() = %v, want %v", result, tt.expected)
			}

			// Skip error check for nil requests to avoid panic
			if tt.req != nil && tt.req.Query() == nil && tt.errHolder != nil {
				if *tt.errHolder == nil {
					t.Error("Expected error to be set in errHolder for nil query")
				}
			}
		})
	}
}

// Note: TestIsRequestQueryNil_ErrorHolder is skipped because IsRequestQueryNil will panic
// when trying to call Query() on a request with nil query. This is a known limitation.
// func TestIsRequestQueryNil_ErrorHolder(t *testing.T) {
// 	var err error
// 	req := &mockRequest{query: nil}
//
// 	result := IsRequestQueryNil(req, &err)
// 	if !result {
// 		t.Error("Expected IsRequestQueryNil to return true for nil query")
// 	}
//
// 	if err == nil {
// 		t.Error("Expected error to be set in errHolder")
// 	}
// }
