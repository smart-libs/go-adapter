package tagbased

import (
	"fmt"
	sdkparam "github.com/smart-libs/go-adapter/sdk/lib/pkg/param"
	converter "github.com/smart-libs/go-crosscutting/converter/lib/pkg"
	"reflect"
)

type (
	defaultOutputParamSpecFactoryRegistry[Output any] struct {
		converters converter.Converters
		factories  []OutputParamSpecFactory[Output]
	}
)

func (o *defaultOutputParamSpecFactoryRegistry[Output]) AsList() []OutputParamSpecFactory[Output] {
	return o.factories
}

func (o *defaultOutputParamSpecFactoryRegistry[Output]) Add(factory OutputParamSpecFactory[Output]) OutputParamSpecFactoryRegistry[Output] {
	if factory == nil {
		return o
	}
	o.factories = append(o.factories, factory)
	return o
}

func (o *defaultOutputParamSpecFactoryRegistry[Output]) AddOption1(tagName, tagValue string, setter func(output Output, value any) error) OutputParamSpecFactoryRegistry[Output] {
	o.factories = append(o.factories, o.addOption1(tagName, tagValue, setter))
	return o
}

func (o *defaultOutputParamSpecFactoryRegistry[Output]) AddOption2(tagName string, setter func(output Output, tagValue string, value any) error) OutputParamSpecFactoryRegistry[Output] {
	o.factories = append(o.factories, o.addOption2(tagName, setter))
	return o
}

func (o *defaultOutputParamSpecFactoryRegistry[Output]) AddOption3(tagName string, setterFactory func(tagValue string) (func(Output, any) error, error)) OutputParamSpecFactoryRegistry[Output] {
	o.factories = append(o.factories, o.addOption3(tagName, setterFactory))
	return o
}

func (o *defaultOutputParamSpecFactoryRegistry[Output]) AddOption4(tagName string, factory func(tagValue string, field reflect.StructField) (getter func(output Output, value any) error, err error)) OutputParamSpecFactoryRegistry[Output] {
	o.factories = append(o.factories, o.addOption4(tagName, factory))
	return o
}

func (o *defaultOutputParamSpecFactoryRegistry[Output]) AddOption5(tagName string, factory func(tagValue string, field reflect.StructField) ([]sdkparam.Option, func(Output, any) error, error)) OutputParamSpecFactoryRegistry[Output] {
	o.factories = append(o.factories, o.addOption5(tagName, factory))
	return o
}

func (o *defaultOutputParamSpecFactoryRegistry[Output]) AddOption6(tagName string, factory func(tagValue string, field reflect.StructField, options []sdkparam.Option) (sdkparam.OutputParamSpec[Output], error)) OutputParamSpecFactoryRegistry[Output] {
	o.factories = append(o.factories, o.addOption6(tagName, factory))
	return o
}

func (o *defaultOutputParamSpecFactoryRegistry[Output]) AddOption7(tagName, tagValue string, factory func(field reflect.StructField, options []sdkparam.Option) (sdkparam.OutputParamSpec[Output], error)) OutputParamSpecFactoryRegistry[Output] {
	o.factories = append(o.factories, o.addOption7(tagName, tagValue, factory))
	return o
}

func (o *defaultOutputParamSpecFactoryRegistry[Output]) addOption1(tagName, expectedTagValue string, setter func(output Output, value any) error) OutputParamSpecFactoryFunc[Output] {
	return o.addOption4(tagName, func(tagValue string, _ reflect.StructField) (func(Output, any) error, error) {
		if tagValue == expectedTagValue {
			return setter, nil
		}
		return nil, nil
	})
}

func (o *defaultOutputParamSpecFactoryRegistry[Output]) addOption2(tagName string, setter func(output Output, tagValue string, value any) error) OutputParamSpecFactoryFunc[Output] {
	return o.addOption4(tagName, func(tagValue string, field reflect.StructField) (func(Output, any) error, error) {
		return func(output Output, a any) error {
			return setter(output, tagValue, a)
		}, nil
	})
}

func (o *defaultOutputParamSpecFactoryRegistry[Output]) addOption3(tagName string, setterFactory func(tagValue string) (func(Output, any) error, error)) OutputParamSpecFactoryFunc[Output] {
	return o.addOption4(tagName, func(tagValue string, field reflect.StructField) (func(Output, any) error, error) {
		return setterFactory(tagValue)
	})
}

func (o *defaultOutputParamSpecFactoryRegistry[Output]) addOption4(tagName string, factory func(tagValue string, field reflect.StructField) (func(Output, any) error, error)) OutputParamSpecFactoryFunc[Output] {
	return o.addOption5(tagName, func(tagValue string, field reflect.StructField) ([]sdkparam.Option, func(Output, any) error, error) {
		setter, err := factory(tagValue, field)
		if err != nil {
			return nil, nil, err
		}

		if setter == nil {
			return nil, nil, nil // try other factory
		}

		options, err := CreateParamSpecOptionsFromTag(field, o.converters)
		if err != nil {
			return nil, nil, err
		}

		return options, setter, nil
	})
}

func (o *defaultOutputParamSpecFactoryRegistry[Output]) addOption5(tagName string, factory func(tagValue string, field reflect.StructField) ([]sdkparam.Option, func(Output, any) error, error)) OutputParamSpecFactoryFunc[Output] {
	return func(field reflect.StructField) (sdkparam.OutputParamSpec[Output], error) {
		if tagValue, found := field.Tag.Lookup(tagName); found {
			options, setter, err := factory(tagValue, field)
			if err != nil {
				return nil, err
			}

			if setter == nil {
				return nil, nil // try other factory
			}

			specName := fmt.Sprintf("%s:%s", tagName, tagValue)
			spec := sdkparam.NewSpec(specName, options...)
			return sdkparam.NewOutputParamSpec(spec, setter), nil
		}

		return nil, nil
	}
}

func (o *defaultOutputParamSpecFactoryRegistry[Output]) addOption6(tagName string, factory func(tagValue string, field reflect.StructField, options []sdkparam.Option) (sdkparam.OutputParamSpec[Output], error)) OutputParamSpecFactoryFunc[Output] {
	return func(field reflect.StructField) (sdkparam.OutputParamSpec[Output], error) {
		if tagValue, found := field.Tag.Lookup(tagName); found {
			options, err := CreateParamSpecOptionsFromTag(field, o.converters)
			if err != nil {
				return nil, err
			}

			return factory(tagValue, field, options)
		}

		return nil, nil
	}
}

func (o *defaultOutputParamSpecFactoryRegistry[Output]) addOption7(tagName, desiredTagValue string, factory func(field reflect.StructField, options []sdkparam.Option) (sdkparam.OutputParamSpec[Output], error)) OutputParamSpecFactoryFunc[Output] {
	return o.addOption6(tagName, func(tagValue string, field reflect.StructField, options []sdkparam.Option) (sdkparam.OutputParamSpec[Output], error) {
		if tagValue == desiredTagValue {
			return factory(field, options)
		}
		return nil, nil
	})
}

func NewOutputParamSpecFactoryRegistry[Output any](converters converter.Converters) OutputParamSpecFactoryRegistry[Output] {
	return &defaultOutputParamSpecFactoryRegistry[Output]{converters: converters}
}
