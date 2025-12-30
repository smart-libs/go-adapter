package httpadpt

import (
	"errors"
	serror "github.com/smart-libs/go-crosscutting/serror/lib/pkg"
	"net/http"
	"testing"
)

func TestOutErrorParamSpec_Name(t *testing.T) {
	spec := OutErrorParamSpec{}
	if spec.Name() != "error" {
		t.Errorf("Expected Name() = %q, got %q", "error", spec.Name())
	}
}

func TestOutErrorParamSpec_Options(t *testing.T) {
	spec := OutErrorParamSpec{}
	options := spec.Options()
	if len(options) != 0 {
		t.Errorf("Expected Options() to return nil or empty slice, got %v", options)
	}
}

func TestOutErrorParamSpec_SetValue_NilValue(t *testing.T) {
	spec := OutErrorParamSpec{}
	output := &Response{}

	err := spec.SetValue(output, nil)
	if err != nil {
		t.Errorf("Expected SetValue(nil) to return nil error, got %v", err)
	}

	if output.StatusCode != nil {
		t.Error("Expected StatusCode to remain nil when value is nil")
	}

	if output.Body != nil {
		t.Error("Expected Body to remain nil when value is nil")
	}
}

func TestOutErrorParamSpec_SetValue_NilOutput(t *testing.T) {
	spec := OutErrorParamSpec{}
	testErr := errors.New("test error")

	err := spec.SetValue(nil, testErr)
	if err == nil {
		t.Fatal("Expected SetValue(nil, error) to return error, got nil")
	}
}

func TestOutErrorParamSpec_SetValue_WithError(t *testing.T) {
	spec := OutErrorParamSpec{}
	output := &Response{}
	testErr := errors.New("test error message")

	err := spec.SetValue(output, testErr)
	if err != nil {
		t.Errorf("Expected SetValue to succeed, got error: %v", err)
		return
	}

	if output.StatusCode == nil {
		t.Fatal("Expected StatusCode to be set, got nil")
	}

	if *output.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected StatusCode = %d, got %d", http.StatusInternalServerError, *output.StatusCode)
	}

	if string(output.Body) != testErr.Error() {
		t.Errorf("Expected Body = %q, got %q", testErr.Error(), string(output.Body))
	}
}

func TestOutErrorParamSpec_SetValue_WithIllegalArgumentError(t *testing.T) {
	spec := OutErrorParamSpec{}
	output := &Response{}
	testErr := serror.IllegalConfigParamValue("param", "value")

	err := spec.SetValue(output, testErr)
	if err != nil {
		t.Errorf("Expected SetValue to succeed, got error: %v", err)
		return
	}

	if output.StatusCode == nil {
		t.Fatal("Expected StatusCode to be set, got nil")
	}

	if *output.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected StatusCode = %d, got %d", http.StatusBadRequest, *output.StatusCode)
	}
}

func TestOutErrorParamSpec_SetValue_WithNotFoundError(t *testing.T) {
	spec := OutErrorParamSpec{}
	output := &Response{}
	testErr := serror.IllegalConfigParamValue("param", "value")

	err := spec.SetValue(output, testErr)
	if err != nil {
		t.Errorf("Expected SetValue to succeed, got error: %v", err)
		return
	}

	if output.StatusCode == nil {
		t.Fatal("Expected StatusCode to be set, got nil")
	}

	if *output.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected StatusCode = %d, got %d", http.StatusBadRequest, *output.StatusCode)
	}
}

func TestOutErrorParamSpec_SetValue_NonErrorType(t *testing.T) {
	spec := OutErrorParamSpec{}
	output := &Response{}

	// SetValue with non-error type should not set anything
	err := spec.SetValue(output, "not an error")
	if err != nil {
		t.Errorf("Expected SetValue to return nil error, got %v", err)
	}

	if output.StatusCode != nil {
		t.Error("Expected StatusCode to remain nil when value is not an error")
	}

	if output.Body != nil {
		t.Error("Expected Body to remain nil when value is not an error")
	}
}

func TestNewOutErrorParamSpec(t *testing.T) {
	spec := NewOutErrorParamSpec()
	if spec == nil {
		t.Fatal("Expected NewOutErrorParamSpec to return non-nil spec")
	}

	// Verify it's the correct type
	if _, ok := spec.(OutErrorParamSpec); !ok {
		t.Errorf("Expected OutErrorParamSpec type, got %T", spec)
	}
}
