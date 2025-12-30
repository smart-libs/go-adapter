package httpadpt

import (
	"context"
	"testing"
)

func TestIsBindingValid(t *testing.T) {
	tests := []struct {
		name     string
		binding  Binding
		expected bool
	}{
		{
			name:     "binding with nil handler",
			binding:  Binding{},
			expected: false,
		},
		{
			name: "binding with handler",
			binding: Binding{
				Handler: &mockHandler{},
			},
			expected: true,
		},
		{
			name: "binding with condition and handler",
			binding: Binding{
				Condition: Condition{
					Path:    stringPtr("/api/users"),
					Methods: []string{"GET"},
				},
				Handler: &mockHandler{},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsBindingValid(tt.binding)
			if result != tt.expected {
				t.Errorf("IsBindingValid() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// mockHandler is a minimal handler implementation for testing
type mockHandler struct{}

func (m *mockHandler) Invoke(_ context.Context, _ Request, output *Response) error {
	if output != nil {
		statusCode := 200
		output.StatusCode = &statusCode
	}
	return nil
}
