package tagbased

import (
	"fmt"
	sdkparam "github.com/smart-libs/go-adapter/sdk/lib/pkg/param"
	"reflect"
)

type (
	registryBasedOutputParamSpecFactory[Output any] struct {
		registry OutputParamSpecFactoryRegistry[Output]
	}
)

func (r registryBasedOutputParamSpecFactory[Output]) CreateOutputParamSpec(field reflect.StructField) (sdkparam.OutputParamSpec[Output], error) {
	for _, factory := range r.registry.AsList() {
		spec, err := factory.CreateOutputParamSpec(field)
		if err != nil {
			return nil, err
		}
		if spec != nil {
			return spec, nil
		}
	}
	return nil, fmt.Errorf("no sdkparam.OutputParamSpec[Output] instance created for Field=[%s] using tags=[%v]",
		field.Name, field.Tag)
}

func NewOutputParamSpecFactory[Output any](registry OutputParamSpecFactoryRegistry[Output]) OutputParamSpecFactory[Output] {
	return registryBasedOutputParamSpecFactory[Output]{
		registry: registry,
	}
}
