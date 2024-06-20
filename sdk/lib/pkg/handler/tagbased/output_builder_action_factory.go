package tagbasedhandler

import (
	"github.com/smart-libs/go-adapter/interfaces/pkg/adapter"
	"reflect"
)

func creatOutputActionToInvokeOutputElement(action OutputBuilderActionFunc) OutputBuilderActionFunc {
	if action == nil {
		return nil
	}
	return func(builder adapter.OutputBuilder, output any) {
		action(builder, getAsValueOf(output).Elem())
	}
}

func createOutputActionToInvokeAllActions(actions []OutputBuilderActionFunc) OutputBuilderActionFunc {
	return func(builder adapter.OutputBuilder, output any) {
		for _, action := range actions {
			if action == nil {
				continue
			}
			action(builder, output)
		}
	}
}

func createOutputActionToCallAllElemOutputAction(elementAction OutputBuilderActionFunc) OutputBuilderActionFunc {
	return func(builder adapter.OutputBuilder, givenArray any) {
		array := getAsValueOf(givenArray)
		numOfElem := array.Len()
		for i := 0; i < numOfElem; i++ {
			elem := array.Index(i)
			elementAction(builder, elem)
		}
	}
}

func createOutputActionToHandleTheFieldValueOutputAction(ref adapter.ParamRef, field reflect.StructField) OutputBuilderActionFunc {
	return func(builder adapter.OutputBuilder, output any) {
		builder.WithParam(ref, getAsValueOf(output).FieldByName(field.Name))
	}
}

func getAsValueOf(v any) reflect.Value {
	if asValueOf, ok := v.(reflect.Value); ok {
		return asValueOf
	}
	return reflect.ValueOf(v)
}
