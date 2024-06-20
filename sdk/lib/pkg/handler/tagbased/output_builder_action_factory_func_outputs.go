package tagbasedhandler

import (
	"fmt"
	"github.com/smart-libs/go-adapter/interfaces/pkg/adapter"
	"github.com/smart-libs/go-adapter/sdk/lib/pkg/param/tagbased"
	"reflect"
)

type (
	// AbstractOutputSpecBuilder builds the OutputParamSpec[Input] object without using generics

	// OutputBuilderActionsFactory is a generic factory to be reused by any adapter ro create the OutputSpecs object.
	// The adapter identifies which GO lang tag name belongs to it and should be used to invoke OutputActionFactory.
	// For instance, the CLI adapter uses the tag name 'cli' to say to OutputBuilderActionsFactory that whenever the tag 'cli'
	// is found in a structure field, then OutputBuilderActionsFactory function must be invoked so that the CLI adapter
	// will return the sdkparam.OutputBuilderActionFunc that must be executed and this struct field value is retrieved.
	OutputBuilderActionsFactory struct {
		tagbased.AbstractOutputSpecBuilder
		// OnTypeNotSupported is used when the output type is not supported by this factory
		OnTypeNotSupported func(outputType reflect.Type) (OutputBuilderActionFunc, error)
	}
)

func (m OutputBuilderActionsFactory) Create(funcType reflect.Type) ([]OutputBuilderActionFunc, error) {
	var actions []OutputBuilderActionFunc

	numOfOut := funcType.NumOut()
	for i := 0; i < numOfOut; i++ {
		outputType := funcType.Out(i)
		action, err := m.createOutputActionForAny(outputType)
		if err != nil {
			return nil, err
		}

		actions = append(actions, action)
	}
	return actions, nil
}

func (m OutputBuilderActionsFactory) createOutputActionForArray(arrayOutputType reflect.Type) (OutputBuilderActionFunc, error) {
	elem := arrayOutputType.Elem()
	elemAction, err := m.createOutputActionForAny(elem)
	if err != nil {
		return nil, err
	}
	return createOutputActionToCallAllElemOutputAction(elemAction), nil
}

func (m OutputBuilderActionsFactory) createOutputActionForStruct(structOutputType reflect.Type) (OutputBuilderActionFunc, error) {
	var result []OutputBuilderActionFunc
	numOfFields := structOutputType.NumField()
	for i := 0; i < numOfFields; i++ {
		field := structOutputType.Field(i)
		ref := adapter.StringParamRef(fmt.Sprintf("%s.%s", structOutputType.Name(), field.Name))
		if err := m.AbstractOutputSpecBuilder.AddOutputParamSpec(ref, field); err != nil {
			return nil, err
		}

		result = append(result, createOutputActionToHandleTheFieldValueOutputAction(ref, field))
	}
	return createOutputActionToInvokeAllActions(result), nil
}

func (m OutputBuilderActionsFactory) createOutputActionForAny(anyOutputType reflect.Type) (OutputBuilderActionFunc, error) {
	switch {
	case anyOutputType.Kind() == reflect.Ptr:
		action, err := m.createOutputActionForAny(anyOutputType.Elem())
		if err != nil {
			return nil, err
		}
		return creatOutputActionToInvokeOutputElement(action), nil

	case anyOutputType.Kind() == reflect.Struct:
		action, err := m.createOutputActionForStruct(anyOutputType)
		if err != nil {
			return nil, err
		}
		return action, nil

	case anyOutputType.Kind() == reflect.Slice || anyOutputType.Kind() == reflect.Array:
		action, err := m.createOutputActionForArray(anyOutputType)
		if err != nil {
			return nil, err
		}
		return action, nil

	case IsType[error](anyOutputType):
		return func(builder adapter.OutputBuilder, outputParamValue any) {
			asValueOf := getAsValueOf(outputParamValue)
			if asValueOf.IsNil() {
				builder.WithError(nil)
			} else {
				builder.WithError(asValueOf.Interface().(error))
			}
		}, nil

	case m.OnTypeNotSupported != nil:
		return m.OnTypeNotSupported(anyOutputType)

	default:
		return nil, fmt.Errorf("no rule to invoke OutputBuilder for handler output with type=[%s]", anyOutputType)
	}
}
