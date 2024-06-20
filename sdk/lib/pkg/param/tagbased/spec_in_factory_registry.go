package tagbased

import (
	sdkparam "github.com/smart-libs/go-adapter/sdk/lib/pkg/param"
	"reflect"
)

type (
	InputParamSpecFactoryRegistry[Input any] interface {
		// Add adds an implementation of InputParamSpecFactory[Input], if factory returns a nil InputParamSpec[Input], then other factory will be invoked
		Add(InputParamSpecFactory[Input]) InputParamSpecFactoryRegistry[Input]

		// AddOption1 when the tag name and value are constants
		AddOption1(tagName, tagValue string, getter func(Input) (any, error)) InputParamSpecFactoryRegistry[Input]

		// AddOption2 when the getter requires the tag value to return the input value
		AddOption2(tagName string, getter func(Input, string) (any, error)) InputParamSpecFactoryRegistry[Input]

		// AddOption3 when the getter requires the tag value to return the input value but it needs to convert tagValue first
		AddOption3(tagName string, getterFactory func(tagValue string) (func(Input) (any, error), error)) InputParamSpecFactoryRegistry[Input]

		// AddOption4 when the tag name is constant but you need other tags
		AddOption4(tagName string, factory func(tagValue string, field reflect.StructField) (getter func(Input) (any, error), err error)) InputParamSpecFactoryRegistry[Input]

		// AddOption5 when the tag name is constant but you want to handle additional tags and options
		AddOption5(tagName string, factory func(tagValue string, field reflect.StructField) (options []sdkparam.Option, getter func(Input) (any, error), err error)) InputParamSpecFactoryRegistry[Input]

		// AddOption6 when the tag name is constant, and you want to create the sdkparam.inputParamSpec[input]
		AddOption6(tagName string, factory func(tagValue string, field reflect.StructField, options []sdkparam.Option) (sdkparam.InputParamSpec[Input], error)) InputParamSpecFactoryRegistry[Input]

		// AddOption7 when the name and value are constant, and you want to create the sdkparam.inputParamSpec[input]
		AddOption7(tagName, tagValue string, factory func(field reflect.StructField, options []sdkparam.Option) (sdkparam.InputParamSpec[Input], error)) InputParamSpecFactoryRegistry[Input]

		// AsList returns the internal list of InputParamSpecFactory registered
		AsList() []InputParamSpecFactory[Input]
	}
)
