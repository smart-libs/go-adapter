package httpadpt

import (
	tagbasedhandler "github.com/smart-libs/go-adapter/sdk/lib/pkg/handler/tagbased"
)

type (
	StandardConditionBuildingStep interface {
		PathConditionBuildingStep
		MethodsConditionBuildingStep
	}

	HandlerBuildingStep interface {
		WithHandlerFunc(handler any) Binding
	}

	PathConditionBuildingStep interface {
		// WithPath accepts an HTTP Path where path parameters shall be surrounded by {}
		WithPath(string) HandlerBuildingStep
	}

	MethodsConditionBuildingStep interface {
		WithMethods(string, ...string) HandlerBuildingStep
	}

	BaseBuilder struct {
		Binding
	}
)

func NewBindingBuilderUsingOtherCondition(other any) StandardConditionBuildingStep {
	return &BaseBuilder{Binding{Condition: Condition{Other: other}}}
}

func NewBindingBuilderUsingPath(path string) MethodsConditionBuildingStep {
	result := &BaseBuilder{}
	result.WithPath(path)
	return result
}

func NewBindingBuilderUsingMethods(method string, methods ...string) MethodsConditionBuildingStep {
	result := &BaseBuilder{}
	result.WithMethods(method, methods...)
	return result
}

func (b *BaseBuilder) WithPath(path string) HandlerBuildingStep {
	b.Condition.Path = &path
	return b
}

func (b *BaseBuilder) WithMethods(method string, methods ...string) HandlerBuildingStep {
	b.Condition.Methods = append([]string{method}, methods...)
	return b
}

func (b *BaseBuilder) WithHandlerFunc(handler any) Binding {
	b.Handler = tagbasedhandler.NewBuilderForFunc[Request, *Response](handler).
		WithInTagBasedFactory(createInParamSpecFactory()).
		WithOutTagBasedFactory(createOutParamSpecFactory()).
		WithOutErrorParamSpec(NewOutErrorParamSpec()).
		Build()
	return b.Binding
}
