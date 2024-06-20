package tagbased

import (
	"fmt"
	sdkparam "github.com/smart-libs/go-adapter/sdk/lib/pkg/param"
	converter "github.com/smart-libs/go-crosscutting/converter/lib/pkg"
	"reflect"
)

type (
	defaultInputParamSpecFactoryRegistry[Input any] struct {
		converters converter.Converters
		factories  []InputParamSpecFactory[Input]
	}
)

func (d *defaultInputParamSpecFactoryRegistry[Input]) Add(factory InputParamSpecFactory[Input]) InputParamSpecFactoryRegistry[Input] {
	if factory == nil {
		return d
	}
	d.factories = append(d.factories, factory)
	return d
}

func (d *defaultInputParamSpecFactoryRegistry[Input]) AsList() []InputParamSpecFactory[Input] {
	return d.factories
}

func (d *defaultInputParamSpecFactoryRegistry[Input]) AddOption1(tagName, tagValue string, getter func(Input) (any, error)) InputParamSpecFactoryRegistry[Input] {
	d.factories = append(d.factories, d.addOption1(tagName, tagValue, getter))
	return d
}

func (d *defaultInputParamSpecFactoryRegistry[Input]) AddOption2(tagName string, getter func(input Input, tagValue string) (any, error)) InputParamSpecFactoryRegistry[Input] {
	d.factories = append(d.factories, d.addOption2(tagName, getter))
	return d
}

func (d *defaultInputParamSpecFactoryRegistry[Input]) AddOption3(tagName string, getterFactory func(tagValue string) (func(Input) (any, error), error)) InputParamSpecFactoryRegistry[Input] {
	d.factories = append(d.factories, d.addOption3(tagName, getterFactory))
	return d
}

func (d *defaultInputParamSpecFactoryRegistry[Input]) AddOption4(tagName string, factory func(tagValue string, field reflect.StructField) (func(Input) (any, error), error)) InputParamSpecFactoryRegistry[Input] {
	d.factories = append(d.factories, d.addOption4(tagName, factory))
	return d
}

func (d *defaultInputParamSpecFactoryRegistry[Input]) AddOption5(tagName string, factory func(tagValue string, field reflect.StructField) (options []sdkparam.Option, getter func(Input) (any, error), err error)) InputParamSpecFactoryRegistry[Input] {
	d.factories = append(d.factories, d.addOption5(tagName, factory))
	return d
}

func (d *defaultInputParamSpecFactoryRegistry[Input]) AddOption6(tagName string, factory func(tagValue string, field reflect.StructField, options []sdkparam.Option) (sdkparam.InputParamSpec[Input], error)) InputParamSpecFactoryRegistry[Input] {
	d.factories = append(d.factories, d.addOption6(tagName, factory))
	return d
}

func (d *defaultInputParamSpecFactoryRegistry[Input]) AddOption7(tagName, tagValue string, factory func(field reflect.StructField, options []sdkparam.Option) (sdkparam.InputParamSpec[Input], error)) InputParamSpecFactoryRegistry[Input] {
	d.factories = append(d.factories, d.addOption7(tagName, tagValue, factory))
	return d
}

func (d *defaultInputParamSpecFactoryRegistry[Input]) addOption1(tagName, expectedTagValue string, getter func(Input) (any, error)) InputParamSpecFactory[Input] {
	return d.addOption4(tagName, func(tagValue string, _ reflect.StructField) (func(Input) (any, error), error) {
		if tagValue == expectedTagValue {
			return getter, nil
		}
		return nil, nil
	})
}

func (d *defaultInputParamSpecFactoryRegistry[Input]) addOption2(tagName string, getter func(input Input, tagValue string) (any, error)) InputParamSpecFactory[Input] {
	return d.addOption4(tagName, func(tagValue string, _ reflect.StructField) (func(Input) (any, error), error) {
		return func(input Input) (any, error) {
			return getter(input, tagValue)
		}, nil
	})
}

func (d *defaultInputParamSpecFactoryRegistry[Input]) addOption3(tagName string, getterFactory func(tagValue string) (func(Input) (any, error), error)) InputParamSpecFactory[Input] {
	return d.addOption4(tagName, func(tagValue string, _ reflect.StructField) (func(Input) (any, error), error) {
		getter, err := getterFactory(tagValue)
		if err != nil {
			return nil, err
		}
		return getter, nil
	})
}

func (d *defaultInputParamSpecFactoryRegistry[Input]) addOption4(tagName string, factory func(tagValue string, field reflect.StructField) (func(Input) (any, error), error)) InputParamSpecFactoryFunc[Input] {
	return d.addOption5(tagName, func(tagValue string, field reflect.StructField) ([]sdkparam.Option, func(Input) (any, error), error) {
		getter, err := factory(tagValue, field)
		if err != nil {
			return nil, nil, err
		}

		if getter == nil {
			return nil, nil, nil // try other factory
		}

		options, err := CreateParamSpecOptionsFromTag(field, d.converters)
		if err != nil {
			return nil, nil, err
		}

		return options, getter, nil
	})
}

func (d *defaultInputParamSpecFactoryRegistry[Input]) addOption5(tagName string, factory func(tagValue string, field reflect.StructField) ([]sdkparam.Option, func(Input) (any, error), error)) InputParamSpecFactoryFunc[Input] {
	return func(field reflect.StructField) (sdkparam.InputParamSpec[Input], error) {
		if tagValue, found := field.Tag.Lookup(tagName); found {
			options, getter, err := factory(tagValue, field)
			if err != nil {
				return nil, err
			}

			if getter == nil {
				return nil, nil // try other factory
			}

			specName := fmt.Sprintf("%s:%s", tagName, tagValue)
			spec := sdkparam.NewSpec(specName, options...)
			return sdkparam.NewInputParamSpecWithSpec(spec, d.converters, getter), nil
		}

		return nil, nil
	}
}

func (d *defaultInputParamSpecFactoryRegistry[Input]) addOption6(tagName string, factory func(tagValue string, field reflect.StructField, options []sdkparam.Option) (sdkparam.InputParamSpec[Input], error)) InputParamSpecFactoryFunc[Input] {
	return func(field reflect.StructField) (sdkparam.InputParamSpec[Input], error) {
		if tagValue, found := field.Tag.Lookup(tagName); found {
			options, err := CreateParamSpecOptionsFromTag(field, d.converters)
			if err != nil {
				return nil, err
			}
			return factory(tagValue, field, options)
		}

		return nil, nil
	}
}

func (d *defaultInputParamSpecFactoryRegistry[Input]) addOption7(tagName, desiredTagValue string, factory func(field reflect.StructField, options []sdkparam.Option) (sdkparam.InputParamSpec[Input], error)) InputParamSpecFactoryFunc[Input] {
	return d.addOption6(tagName, func(tagValue string, field reflect.StructField, options []sdkparam.Option) (sdkparam.InputParamSpec[Input], error) {
		if tagValue == desiredTagValue {
			return factory(field, options)
		}
		return nil, nil
	})
}

func NewInputParamSpecFactoryRegistry[Input any](converters converter.Converters) InputParamSpecFactoryRegistry[Input] {
	return &defaultInputParamSpecFactoryRegistry[Input]{converters: converters}
}
