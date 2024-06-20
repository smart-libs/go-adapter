package sdkparam

import (
	"cmp"
	"fmt"
	converter "github.com/smart-libs/go-crosscutting/converter/lib/pkg"
	converterdefault "github.com/smart-libs/go-crosscutting/converter/lib/pkg/default"
)

type (
	// InputParamSpec represents a parameter of the Input object
	InputParamSpec[Input any] interface {
		// Spec is the parameter specification
		Spec
		CopyValue(input Input, target any) error
		GetValue(input Input) (any, error)
	}

	defaultInputParam[Input any] struct {
		// Spec is the input specification that identifies the input name and options
		Spec

		// Converters are used by the CopyValue method
		converter.Converters

		// getValueFunc is the adapter method that knows how to get a value associated with the Spec given the Input
		getValueFunc func(input Input) (any, error)
	}
)

func (i defaultInputParam[Input]) CopyValue(input Input, target any) error {
	getResult, err := i.GetValue(input)
	if err != nil {
		return err
	}
	if isValueOfNil(getAsValueOf(getResult)) {
		return nil
	}
	fromToFunc := i.Converters.Convert
	return fromToFunc(getResult, target)
}

func (i defaultInputParam[Input]) GetValue(input Input) (any, error) {
	inputValue, err := i.getValueFunc(input)
	if err != nil {
		return nil, fmt.Errorf("failed to get param=[%s] value from input=[%T]", i.Spec.Name(), input)
	}
	return AsSingleOptions(i.Spec.Options()...)(i.Spec, inputValue)
}

// NewInputParamSpec is a helper factory without needing to provide the Spec
func NewInputParamSpec[Input any](specName string, options []Option, converters converter.Converters, getValueFunc func(input Input) (any, error)) InputParamSpec[Input] {
	spec := NewSpec(specName, options...)
	return NewInputParamSpecWithSpec(spec, converters, getValueFunc)
}

func NewInputParamSpecWithSpec[Input any](spec Spec, converters converter.Converters, getValueFunc func(input Input) (any, error)) InputParamSpec[Input] {
	return defaultInputParam[Input]{
		Spec:         spec,
		Converters:   cmp.Or(converters, converterdefault.Converters),
		getValueFunc: getValueFunc,
	}
}

var (
	_ InputParamSpec[any] = defaultInputParam[any]{}
)
