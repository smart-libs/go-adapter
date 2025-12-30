package httpadpt

import (
	"testing"
)

func Test_getInputParamSpecFactoryRegistry(t *testing.T) {
	// Reset the global registry for testing
	inParamSpecFactoryRegistry = nil

	// First call should initialize
	registry1 := getInputParamSpecFactoryRegistry()
	if registry1 == nil {
		t.Fatal("Expected registry to be initialized, got nil")
	}

	// Second call should return the same instance
	registry2 := getInputParamSpecFactoryRegistry()
	if registry1 != registry2 {
		t.Error("Expected getInputParamSpecFactoryRegistry to return the same instance")
	}
}

func Test_createInParamSpecFactory(t *testing.T) {
	// Reset the global registry for testing
	inParamSpecFactoryRegistry = nil

	// This should not panic even if registry is nil (it should initialize it)
	factory := createInParamSpecFactory()
	if factory == nil {
		t.Fatal("Expected factory to be created, got nil")
	}

	// Verify registry was initialized
	if inParamSpecFactoryRegistry == nil {
		t.Error("Expected registry to be initialized after createInParamSpecFactory")
	}
}

func Test_createInParamSpecFactory_InitializesRegistry(t *testing.T) {
	// Reset the global registry
	inParamSpecFactoryRegistry = nil

	// Call createInParamSpecFactory which should initialize the registry
	_ = createInParamSpecFactory()

	// Verify registry is not nil
	if inParamSpecFactoryRegistry == nil {
		t.Error("Expected registry to be initialized, got nil")
	}
}
