package tagbased

import (
	"fmt"
	sdkparam "github.com/smart-libs/go-adapter/sdk/lib/pkg/param"
	"reflect"
)

type (
	registryBasedInputParamSpecFactory[Input any] struct {
		registry InputParamSpecFactoryRegistry[Input]
	}
)

func (f registryBasedInputParamSpecFactory[Input]) CreateInputParamSpec(field reflect.StructField) (sdkparam.InputParamSpec[Input], error) {
	for _, factory := range f.registry.AsList() {
		spec, err := factory.CreateInputParamSpec(field)
		if err != nil {
			return nil, err
		}
		if spec != nil {
			return spec, nil
		}
	}
	return nil, ErrNoInputParamSpecCreatedForField{Field: field}
}

func NewsInputParamSpecFactory[Input any](registry InputParamSpecFactoryRegistry[Input]) InputParamSpecFactory[Input] {
	if registry == nil {
		panic(fmt.Errorf("InputParamSpecFactoryRegistry[Input] cannot be nil"))
	}
	return registryBasedInputParamSpecFactory[Input]{registry: registry}
}
