package tagbasedhandler

import (
	"context"
	"fmt"
	"github.com/smart-libs/go-adapter/interfaces/pkg/adapter"
	"github.com/smart-libs/go-adapter/sdk/lib/pkg/param/tagbased"
	"reflect"
)

type (
	fieldSetter struct {
		field    reflect.StructField
		paramRef adapter.ParamRef
	}

	// structureArgFactory creates the argFactoryFunc that creates and instantiate a structure setting its fields with
	// input values to be used by the handler. This instantiated structure will be an argument in the array of
	// arguments that must be passed to invoke the handler function.
	structureArgFactory struct {
		// structureType is the type of the structure to be instantiated
		structureType reflect.Type

		// fieldSetters has all the functions needed to set each field value of the structure specified by structureType
		fieldSetters []fieldSetter

		// resolveReturnType is the function that knows whether the structure must be returned as a pointer or not.
		resolveReturnType func(structureInstance reflect.Value) reflect.Value
	}
)

func (f fieldSetter) set(accessor adapter.InputAccessor, structureInstance reflect.Value) {
	field := structureInstance.Elem().FieldByIndex(f.field.Index)
	err := accessor.CopyValue(f.paramRef, field.Addr().Interface())
	if err != nil {
		panic(err)
	}
}

// create is the function that creates and instantiate the structure
func (s structureArgFactory) create(_ context.Context, factory adapter.InputAccessor) reflect.Value {
	structureInstance := reflect.New(s.structureType)
	for _, setter := range s.fieldSetters {
		setter.set(factory, structureInstance)
	}
	return s.resolveReturnType(structureInstance)
}

func createArgFactoryForStructure(structureType reflect.Type, factory tagbased.AbstractInputSpecBuilder) argFactoryFunc {
	result := structureArgFactory{structureType: structureType}

	numOfFields := structureType.NumField()
	for i := 0; i < numOfFields; i++ {
		field := structureType.Field(i)
		ref := adapter.StringParamRef(fmt.Sprintf("%s.%s", structureType.Name(), field.Name))
		if err := factory.AddInputParamSpec(ref, field); err != nil {
			panic(err)
		}

		result.fieldSetters = append(result.fieldSetters, fieldSetter{field: field, paramRef: ref})
	}

	result.resolveReturnType = func(structValue reflect.Value) reflect.Value {
		return structValue.Elem()
	}

	if structureType.Kind() == reflect.Ptr {
		result.resolveReturnType = func(structValue reflect.Value) reflect.Value {
			return structValue
		}
	}

	return result.create
}
