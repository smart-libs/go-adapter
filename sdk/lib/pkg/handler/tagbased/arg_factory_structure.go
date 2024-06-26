package tagbasedhandler

import (
	"context"
	"fmt"
	"github.com/smart-libs/go-adapter/interfaces/pkg/adapter"
	"github.com/smart-libs/go-adapter/sdk/lib/pkg/param/tagbased"
	"reflect"
)

type (
	// FieldSetter sets the structure field value that belongs to the structureInstance given and using the
	// given accessor that has the value to be set to the field.
	FieldSetter interface {
		set(ctx context.Context, accessor adapter.InputAccessor, structureInstance reflect.Value)
	}

	// recursiveFieldSetter is used when the field does not have the flag tag, but it is a structure that owns fields
	// that have the flag tag.
	recursiveFieldSetter struct {
		field             reflect.StructField
		paramRef          adapter.ParamRef
		fieldValueFactory argFactoryFunc
		isPointer         bool
	}
	// defaultFieldSetter is used when the field has the flag set and it is associated with a param spec.
	defaultFieldSetter struct {
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
		fieldSetters []FieldSetter

		// resolveReturnType is the function that knows whether the structure must be returned as a pointer or not.
		resolveReturnType func(structureInstance reflect.Value) reflect.Value
	}
)

func (f recursiveFieldSetter) set(ctx context.Context, accessor adapter.InputAccessor, structureInstance reflect.Value) {
	fieldValue := f.fieldValueFactory(ctx, accessor)
	field := structureInstance.Elem().FieldByIndex(f.field.Index)
	if f.isPointer {
		field.Set(fieldValue)
	} else {
		field.Set(fieldValue)
	}
}

func (f defaultFieldSetter) set(_ context.Context, accessor adapter.InputAccessor, structureInstance reflect.Value) {
	field := structureInstance.Elem().FieldByIndex(f.field.Index)
	err := accessor.CopyValue(f.paramRef, field.Addr().Interface())
	if err != nil {
		panic(err)
	}
}

// create is the function that creates and instantiate the structure
func (s structureArgFactory) create(ctx context.Context, factory adapter.InputAccessor) reflect.Value {
	structureInstance := reflect.New(s.structureType)
	for _, setter := range s.fieldSetters {
		setter.set(ctx, factory, structureInstance)
	}
	return s.resolveReturnType(structureInstance)
}

// tryFallback is invoked when the field does not have the flag tag. The fallback checks whether the field is another
// structure and go deeper trying to find flag tags in the inner fields recursively.
func tryFallback(structureType reflect.Type, field reflect.StructField, factory tagbased.AbstractInputSpecBuilder) FieldSetter {
	if field.Type.Kind() == reflect.Struct || (field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct) {
		structureFactory := createStructureArgFactory(field.Type, factory)
		if len(structureFactory.fieldSetters) > 0 {
			return recursiveFieldSetter{
				field:             field,
				paramRef:          adapter.StringParamRef(fmt.Sprintf("%s.%s", structureType.Name(), field.Name)),
				fieldValueFactory: structureFactory.create,
			}
		}
	}
	return nil
}

func createStructureArgFactory(givenType reflect.Type, factory tagbased.AbstractInputSpecBuilder) structureArgFactory {
	isPointer := false
	structureType := givenType
	if givenType.Kind() != reflect.Struct {
		if givenType.Kind() == reflect.Ptr && givenType.Elem().Kind() == reflect.Struct {
			isPointer = true
			structureType = givenType.Elem()
		} else {
			panic(fmt.Errorf("given type=[%s] is neither a structure nor a pointer to structure", givenType.String()))
		}
	}
	result := structureArgFactory{structureType: structureType}

	numOfFields := structureType.NumField()
	for i := 0; i < numOfFields; i++ {
		field := structureType.Field(i)
		ref := adapter.StringParamRef(fmt.Sprintf("%s.%s", structureType.Name(), field.Name))
		if err := factory.AddInputParamSpec(ref, field); err != nil {
			if _, ok := err.(tagbased.ErrNoInputParamSpecCreatedForField); ok {
				if fieldSetter := tryFallback(structureType, field, factory); fieldSetter != nil {
					result.fieldSetters = append(result.fieldSetters, fieldSetter)
				}
				continue // ignore the field without the flag tag
			}
			panic(err)
		}

		result.fieldSetters = append(result.fieldSetters, defaultFieldSetter{field: field, paramRef: ref})
	}

	result.resolveReturnType = func(structValue reflect.Value) reflect.Value {
		return structValue.Elem()
	}

	if isPointer {
		result.resolveReturnType = func(structValue reflect.Value) reflect.Value {
			return structValue
		}
	}

	return result
}

func createArgFactoryForStructure(structureType reflect.Type, factory tagbased.AbstractInputSpecBuilder) argFactoryFunc {
	return createStructureArgFactory(structureType, factory).create
}
