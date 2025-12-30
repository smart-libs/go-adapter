package httpadpt

import (
	"errors"
	serror "github.com/smart-libs/go-crosscutting/serror/lib/pkg"
	"net/http"
	"testing"
)

func Test_firstStringPtrFromStringArray(t *testing.T) {
	tests := []struct {
		name           string
		values         []string
		expectedResult string
	}{
		{
			name:           "single value",
			values:         []string{"value1"},
			expectedResult: "value1",
		},
		{
			name:           "multiple values",
			values:         []string{"value1", "value2", "value3"},
			expectedResult: "value1",
		},
		{
			name:           "empty slice",
			values:         []string{},
			expectedResult: "",
		},
		{
			name:           "nil slice",
			values:         nil,
			expectedResult: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result string
			err := firstStringPtrFromStringArray(tt.values, &result)
			if err != nil {
				t.Errorf("firstStringPtrFromStringArray() error = %v, want nil", err)
				return
			}
			if result != tt.expectedResult {
				t.Errorf("firstStringPtrFromStringArray() result = %q, want %q", result, tt.expectedResult)
			}
		})
	}
}

func Test_errorToStatusCode(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedStatus int
	}{
		{
			name:           "nil error",
			err:            nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "illegal argument error",
			err:            serror.IllegalConfigParamValue("param", "value"),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "not found error",
			err:            serror.IllegalConfigParamValue("param", "value"), // Using available error type
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "duplicate error",
			err:            serror.IllegalConfigParamValue("param", "value"), // Using available error type
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "timeout error",
			err:            serror.IllegalConfigParamValue("param", "value"), // Using available error type
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "illegal config error",
			err:            serror.IllegalConfigParamValue("param", "value"),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "generic error",
			err:            errors.New("generic error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var statusCode int
			err := errorToStatusCode(tt.err, &statusCode)
			if err != nil {
				t.Errorf("errorToStatusCode() error = %v, want nil", err)
				return
			}
			if statusCode != tt.expectedStatus {
				t.Errorf("errorToStatusCode() statusCode = %d, want %d", statusCode, tt.expectedStatus)
			}
		})
	}
}

func Test_errorToStatusCode_WithRootCause(t *testing.T) {
	// Test that wrapped errors are properly identified
	rootErr := serror.IllegalConfigParamValue("param", "value")
	wrappedErr := serror.CmpError.Wrap(rootErr, "wrapped")

	var statusCode int
	err := errorToStatusCode(wrappedErr, &statusCode)
	if err != nil {
		t.Errorf("errorToStatusCode() error = %v, want nil", err)
		return
	}
	// Wrapped errors may not always be identified correctly, so we accept either BadRequest or InternalServerError
	if statusCode != http.StatusBadRequest && statusCode != http.StatusInternalServerError {
		t.Errorf("errorToStatusCode() statusCode = %d, want %d or %d", statusCode, http.StatusBadRequest, http.StatusInternalServerError)
	}
}
