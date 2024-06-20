package tagbased

import (
	sdkparam "github.com/smart-libs/go-adapter/sdk/lib/pkg/param"
	"reflect"
)

type (
	OutputParamSpecFactoryRegistry[Output any] interface {
		Add(factory OutputParamSpecFactory[Output]) OutputParamSpecFactoryRegistry[Output]

		// AddOption1 when the tag name and value are constants
		AddOption1(tagName, tagValue string, setter func(output Output, value any) error) OutputParamSpecFactoryRegistry[Output]

		// AddOption2 when the setter requires the tag value to set the output value
		AddOption2(tagName string, setter func(output Output, tagValue string, value any) error) OutputParamSpecFactoryRegistry[Output]

		// AddOption3 when the setter requires the tag value to set the output value, but it needs to convert tagValue first
		AddOption3(tagName string, setterFactory func(tagValue string) (func(Output, any) error, error)) OutputParamSpecFactoryRegistry[Output]

		// AddOption4 when the tag name is constant, but you need other tags
		AddOption4(tagName string, factory func(tagValue string, field reflect.StructField) (getter func(output Output, value any) error, err error)) OutputParamSpecFactoryRegistry[Output]

		// AddOption5 when the tag name is constant, but you want to return all options
		AddOption5(tagName string, factory func(tagValue string, field reflect.StructField) ([]sdkparam.Option, func(Output, any) error, error)) OutputParamSpecFactoryRegistry[Output]

		// AddOption6 when the tag name is constant, and you want to create the sdkparam.OutputParamSpec[Output]
		AddOption6(tagName string, factory func(tagValue string, field reflect.StructField, options []sdkparam.Option) (sdkparam.OutputParamSpec[Output], error)) OutputParamSpecFactoryRegistry[Output]

		// AddOption7 when the name and value are constant, and you want to create the sdkparam.OutputParamSpec[Output]
		AddOption7(tagName, tagValue string, factory func(field reflect.StructField, options []sdkparam.Option) (sdkparam.OutputParamSpec[Output], error)) OutputParamSpecFactoryRegistry[Output]

		// AsList returns the list of OptionFactory
		AsList() []OutputParamSpecFactory[Output]
	}
)
