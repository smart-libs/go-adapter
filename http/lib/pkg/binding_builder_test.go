package httpadpt

import (
	"testing"
)

func TestNewBindingBuilderUsingOtherCondition(t *testing.T) {
	other := "custom-condition"
	builder := NewBindingBuilderUsingOtherCondition(other)

	baseBuilder, ok := builder.(*BaseBuilder)
	if !ok {
		t.Fatalf("Expected *BaseBuilder, got %T", builder)
	}

	if baseBuilder.Condition.Other != other {
		t.Errorf("Expected Other = %q, got %q", other, baseBuilder.Condition.Other)
	}
}

func TestNewBindingBuilderUsingPath(t *testing.T) {
	path := "/api/users"
	builder := NewBindingBuilderUsingPath(path)

	baseBuilder, ok := builder.(*BaseBuilder)
	if !ok {
		t.Fatalf("Expected *BaseBuilder, got %T", builder)
	}

	if baseBuilder.Condition.Path == nil {
		t.Fatal("Expected Path to be set, got nil")
	}

	if *baseBuilder.Condition.Path != path {
		t.Errorf("Expected Path = %q, got %q", path, *baseBuilder.Condition.Path)
	}
}

func TestNewBindingBuilderUsingMethods(t *testing.T) {
	method1 := "GET"
	method2 := "POST"
	method3 := "PUT"

	builder := NewBindingBuilderUsingMethods(method1, method2, method3)

	baseBuilder, ok := builder.(*BaseBuilder)
	if !ok {
		t.Fatalf("Expected *BaseBuilder, got %T", builder)
	}

	if len(baseBuilder.Condition.Methods) != 3 {
		t.Fatalf("Expected 3 methods, got %d", len(baseBuilder.Condition.Methods))
	}

	// Verify the first method is first (this tests the bug fix)
	if baseBuilder.Condition.Methods[0] != method1 {
		t.Errorf("Expected first method = %q, got %q", method1, baseBuilder.Condition.Methods[0])
	}

	if baseBuilder.Condition.Methods[1] != method2 {
		t.Errorf("Expected second method = %q, got %q", method2, baseBuilder.Condition.Methods[1])
	}

	if baseBuilder.Condition.Methods[2] != method3 {
		t.Errorf("Expected third method = %q, got %q", method3, baseBuilder.Condition.Methods[2])
	}
}

func TestBaseBuilder_WithPath(t *testing.T) {
	builder := &BaseBuilder{}
	path := "/api/posts"

	result := builder.WithPath(path)

	if result != builder {
		t.Error("Expected WithPath to return the same builder")
	}

	if builder.Condition.Path == nil {
		t.Fatal("Expected Path to be set, got nil")
	}

	if *builder.Condition.Path != path {
		t.Errorf("Expected Path = %q, got %q", path, *builder.Condition.Path)
	}
}

func TestBaseBuilder_WithMethods(t *testing.T) {
	builder := &BaseBuilder{}
	method1 := "GET"
	method2 := "POST"
	method3 := "DELETE"

	result := builder.WithMethods(method1, method2, method3)

	if result != builder {
		t.Error("Expected WithMethods to return the same builder")
	}

	if len(builder.Condition.Methods) != 3 {
		t.Fatalf("Expected 3 methods, got %d", len(builder.Condition.Methods))
	}

	// Verify the first method is first (this tests the bug fix)
	if builder.Condition.Methods[0] != method1 {
		t.Errorf("Expected first method = %q, got %q", method1, builder.Condition.Methods[0])
	}

	if builder.Condition.Methods[1] != method2 {
		t.Errorf("Expected second method = %q, got %q", method2, builder.Condition.Methods[1])
	}

	if builder.Condition.Methods[2] != method3 {
		t.Errorf("Expected third method = %q, got %q", method3, builder.Condition.Methods[2])
	}
}

func TestBaseBuilder_WithMethods_SingleMethod(t *testing.T) {
	builder := &BaseBuilder{}
	method := "GET"

	builder.WithMethods(method)

	if len(builder.Condition.Methods) != 1 {
		t.Fatalf("Expected 1 method, got %d", len(builder.Condition.Methods))
	}

	if builder.Condition.Methods[0] != method {
		t.Errorf("Expected method = %q, got %q", method, builder.Condition.Methods[0])
	}
}

func TestBaseBuilder_WithMethods_NoVariadicArgs(t *testing.T) {
	builder := &BaseBuilder{}
	method := "POST"

	builder.WithMethods(method)

	if len(builder.Condition.Methods) != 1 {
		t.Fatalf("Expected 1 method, got %d", len(builder.Condition.Methods))
	}

	if builder.Condition.Methods[0] != method {
		t.Errorf("Expected method = %q, got %q", method, builder.Condition.Methods[0])
	}
}

func TestBaseBuilder_WithHandlerFunc(t *testing.T) {
	builder := &BaseBuilder{}
	// Handler function must have struct input and output with proper tags
	type handlerInput struct {
		Value string `query:"value"`
	}
	type handlerOutput struct {
		StatusCode int `statuscode:""`
	}
	handlerFunc := func(input handlerInput) (*handlerOutput, error) {
		return &handlerOutput{StatusCode: 200}, nil
	}

	binding := builder.WithHandlerFunc(handlerFunc)

	if binding.Handler == nil {
		t.Fatal("Expected Handler to be set, got nil")
	}

	if binding.Condition.Path != nil {
		t.Error("Expected Path to be nil when not set")
	}

	if len(binding.Condition.Methods) != 0 {
		t.Errorf("Expected empty Methods, got %d methods", len(binding.Condition.Methods))
	}
}

func TestBaseBuilder_Chaining(t *testing.T) {
	originalPath := "/api/users"
	builder := NewBindingBuilderUsingPath(originalPath)
	method1 := "GET"
	method2 := "POST"

	// Chain WithMethods
	result := builder.WithMethods(method1, method2)

	baseBuilder, ok := result.(*BaseBuilder)
	if !ok {
		t.Fatalf("Expected *BaseBuilder, got %T", result)
	}

	// Verify path is still set (WithMethods doesn't change the path)
	if baseBuilder.Condition.Path == nil {
		t.Fatal("Expected Path to still be set after WithMethods")
	}

	if *baseBuilder.Condition.Path != originalPath {
		t.Errorf("Expected Path = %q, got %q", originalPath, *baseBuilder.Condition.Path)
	}

	// Verify methods are set
	if len(baseBuilder.Condition.Methods) != 2 {
		t.Fatalf("Expected 2 methods, got %d", len(baseBuilder.Condition.Methods))
	}
}
