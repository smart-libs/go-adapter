package httpadpt

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"
)

func Test_defaultTypeBuilder(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: "nil",
		},
		{
			name:     "standard error",
			err:      errors.New("test error"),
			expected: "*errors.errorString",
		},
		{
			name:     "custom error type",
			err:      &customError{message: "custom"},
			expected: "*httpadpt.customError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := defaultTypeBuilder(tt.err)
			if result != tt.expected {
				t.Errorf("defaultTypeBuilder() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func Test_defaultTitleBuilder(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: "nil",
		},
		{
			name:     "non-nil error",
			err:      errors.New("test error"),
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := defaultTitleBuilder(tt.err)
			if result != tt.expected {
				t.Errorf("defaultTitleBuilder() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func Test_TypeBuilder(t *testing.T) {
	// Test that TypeBuilder variable points to defaultTypeBuilder
	if TypeBuilder == nil {
		t.Fatal("TypeBuilder should not be nil")
	}

	// Test with nil error
	result := TypeBuilder(nil)
	if result != "nil" {
		t.Errorf("TypeBuilder(nil) = %q, want %q", result, "nil")
	}

	// Test with non-nil error
	testErr := errors.New("test")
	result = TypeBuilder(testErr)
	expected := reflect.TypeOf(testErr).String()
	if result != expected {
		t.Errorf("TypeBuilder(testErr) = %q, want %q", result, expected)
	}
}

func Test_TitleBuilder(t *testing.T) {
	// Test that TitleBuilder variable points to defaultTitleBuilder
	if TitleBuilder == nil {
		t.Fatal("TitleBuilder should not be nil")
	}

	// Test with nil error
	result := TitleBuilder(nil)
	if result != "nil" {
		t.Errorf("TitleBuilder(nil) = %q, want %q", result, "nil")
	}

	// Test with non-nil error
	testErr := errors.New("test")
	result = TitleBuilder(testErr)
	if result != "" {
		t.Errorf("TitleBuilder(testErr) = %q, want %q", result, "")
	}
}

func TestProblemDetailFromError(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedType   string
		expectedTitle  string
		expectedDetail string
	}{
		{
			name:           "nil error",
			err:            nil,
			expectedType:   "nil",
			expectedTitle:  "nil",
			expectedDetail: "",
		},
		{
			name:           "standard error",
			err:            errors.New("test error message"),
			expectedType:   "*errors.errorString",
			expectedTitle:  "",
			expectedDetail: "test error message",
		},
		{
			name:           "custom error type",
			err:            &customError{message: "custom error"},
			expectedType:   "*httpadpt.customError",
			expectedTitle:  "",
			expectedDetail: "custom error",
		},
		{
			name:           "error with empty message",
			err:            errors.New(""),
			expectedType:   "*errors.errorString",
			expectedTitle:  "",
			expectedDetail: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProblemDetailFromError(tt.err)

			if result.Type != tt.expectedType {
				t.Errorf("ProblemDetailFromError().Type = %q, want %q", result.Type, tt.expectedType)
			}

			if result.Title != tt.expectedTitle {
				t.Errorf("ProblemDetailFromError().Title = %q, want %q", result.Title, tt.expectedTitle)
			}

			if result.Detail != tt.expectedDetail {
				t.Errorf("ProblemDetailFromError().Detail = %q, want %q", result.Detail, tt.expectedDetail)
			}

			// Verify other fields are empty/zero values
			if result.Instance != "" {
				t.Errorf("ProblemDetailFromError().Instance = %q, want empty", result.Instance)
			}

			if result.Status != "" {
				t.Errorf("ProblemDetailFromError().Status = %q, want empty", result.Status)
			}

			if result.Code != "" {
				t.Errorf("ProblemDetailFromError().Code = %q, want empty", result.Code)
			}

			if result.ID != "" {
				t.Errorf("ProblemDetailFromError().ID = %q, want empty", result.ID)
			}

			if result.AdditionalInfo != nil {
				t.Errorf("ProblemDetailFromError().AdditionalInfo = %v, want nil", result.AdditionalInfo)
			}
		})
	}
}

func TestProblemDetailFromError_NoPanic(t *testing.T) {
	// This test specifically verifies that the bug fix works - no panic on nil error
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("ProblemDetailFromError(nil) panicked: %v", r)
		}
	}()

	result := ProblemDetailFromError(nil)
	if result.Detail != "" {
		t.Errorf("ProblemDetailFromError(nil).Detail = %q, want empty string", result.Detail)
	}
}

func TestJSONProblemDetailFromError(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedType   string
		expectedTitle  string
		expectedDetail string
	}{
		{
			name:           "nil error",
			err:            nil,
			expectedType:   "nil",
			expectedTitle:  "nil",
			expectedDetail: "",
		},
		{
			name:           "standard error",
			err:            errors.New("test error message"),
			expectedType:   "*errors.errorString",
			expectedTitle:  "",
			expectedDetail: "test error message",
		},
		{
			name:           "custom error",
			err:            &customError{message: "custom error"},
			expectedType:   "*httpadpt.customError",
			expectedTitle:  "",
			expectedDetail: "custom error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBytes, err := JSONProblemDetailFromError(tt.err)
			if err != nil {
				t.Errorf("JSONProblemDetailFromError() error = %v, want nil", err)
				return
			}

			if jsonBytes == nil {
				t.Fatal("JSONProblemDetailFromError() returned nil bytes")
			}

			// Verify JSON is valid and can be unmarshaled
			var pd ProblemDetail
			if err := json.Unmarshal(jsonBytes, &pd); err != nil {
				t.Errorf("JSONProblemDetailFromError() returned invalid JSON: %v, bytes: %q", err, string(jsonBytes))
				return
			}

			if pd.Type != tt.expectedType {
				t.Errorf("Unmarshaled ProblemDetail.Type = %q, want %q", pd.Type, tt.expectedType)
			}

			if pd.Title != tt.expectedTitle {
				t.Errorf("Unmarshaled ProblemDetail.Title = %q, want %q", pd.Title, tt.expectedTitle)
			}

			if pd.Detail != tt.expectedDetail {
				t.Errorf("Unmarshaled ProblemDetail.Detail = %q, want %q", pd.Detail, tt.expectedDetail)
			}
		})
	}
}

func TestJSONProblemDetailFromError_NoPanic(t *testing.T) {
	// This test specifically verifies that the bug fix works - no panic on nil error
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("JSONProblemDetailFromError(nil) panicked: %v", r)
		}
	}()

	jsonBytes, err := JSONProblemDetailFromError(nil)
	if err != nil {
		t.Errorf("JSONProblemDetailFromError(nil) error = %v, want nil", err)
		return
	}

	if jsonBytes == nil {
		t.Fatal("JSONProblemDetailFromError(nil) returned nil bytes")
	}

	// Verify it's valid JSON
	var pd ProblemDetail
	if err := json.Unmarshal(jsonBytes, &pd); err != nil {
		t.Errorf("JSONProblemDetailFromError(nil) returned invalid JSON: %v", err)
	}
}

func TestProblemDetail_JSONMarshaling(t *testing.T) {
	// Test that omitempty works correctly
	pd := ProblemDetail{
		Type:   "test-type",
		Title:  "test-title",
		Detail: "test-detail",
	}

	jsonBytes, err := json.Marshal(pd)
	if err != nil {
		t.Errorf("json.Marshal(ProblemDetail) error = %v, want nil", err)
		return
	}

	jsonStr := string(jsonBytes)

	// Verify only set fields are in JSON
	if !strings.Contains(jsonStr, "\"type\":\"test-type\"") {
		t.Errorf("JSON should contain type field: %q", jsonStr)
	}
	if !strings.Contains(jsonStr, "\"title\":\"test-title\"") {
		t.Errorf("JSON should contain title field: %q", jsonStr)
	}
	if !strings.Contains(jsonStr, "\"detail\":\"test-detail\"") {
		t.Errorf("JSON should contain detail field: %q", jsonStr)
	}

	// Verify empty fields are omitted (omitempty behavior)
	if strings.Contains(jsonStr, "\"instance\"") {
		t.Errorf("JSON should not contain empty instance field: %q", jsonStr)
	}
	if strings.Contains(jsonStr, "\"status\"") {
		t.Errorf("JSON should not contain empty status field: %q", jsonStr)
	}
	if strings.Contains(jsonStr, "\"code\"") {
		t.Errorf("JSON should not contain empty code field: %q", jsonStr)
	}
	if strings.Contains(jsonStr, "\"id\"") {
		t.Errorf("JSON should not contain empty id field: %q", jsonStr)
	}
	if strings.Contains(jsonStr, "\"additional_info\"") {
		t.Errorf("JSON should not contain empty additional_info field: %q", jsonStr)
	}
}

// Helper types

type customError struct {
	message string
}

func (e *customError) Error() string {
	return e.message
}
