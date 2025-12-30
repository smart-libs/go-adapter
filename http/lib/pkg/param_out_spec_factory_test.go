package httpadpt

import (
	"testing"
)

func Test_getOutParamSpecFactoryRegistry(t *testing.T) {
	// Reset the global registry for testing
	outParamSpecFactoryRegistry = nil

	// First call should initialize
	registry1 := getOutParamSpecFactoryRegistry()
	if registry1 == nil {
		t.Fatal("Expected registry to be initialized, got nil")
	}

	// Second call should return the same instance
	registry2 := getOutParamSpecFactoryRegistry()
	if registry1 != registry2 {
		t.Error("Expected getOutParamSpecFactoryRegistry to return the same instance")
	}
}

func Test_createOutParamSpecFactory(t *testing.T) {
	// Reset the global registry for testing
	outParamSpecFactoryRegistry = nil

	// This should not panic even if registry is nil (it should initialize it)
	factory := createOutParamSpecFactory()
	if factory == nil {
		t.Fatal("Expected factory to be created, got nil")
	}

	// Verify registry was initialized
	if outParamSpecFactoryRegistry == nil {
		t.Error("Expected registry to be initialized after createOutParamSpecFactory")
	}
}

func Test_createOutParamSpecFactory_InitializesRegistry(t *testing.T) {
	// Reset the global registry
	outParamSpecFactoryRegistry = nil

	// Call createOutParamSpecFactory which should initialize the registry
	_ = createOutParamSpecFactory()

	// Verify registry is not nil
	if outParamSpecFactoryRegistry == nil {
		t.Error("Expected registry to be initialized, got nil")
	}
}
