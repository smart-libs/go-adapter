package tagbased

import (
	"github.com/smart-libs/go-adapter/sdk/lib/pkg/param"
	"reflect"
)

type (
	OutputParamSpecFactory[Output any] interface {
		CreateOutputParamSpec(reflect.StructField) (sdkparam.OutputParamSpec[Output], error)
	}

	OutputParamSpecFactoryFunc[Output any] func(reflect.StructField) (sdkparam.OutputParamSpec[Output], error)
)

func (f OutputParamSpecFactoryFunc[Output]) CreateOutputParamSpec(field reflect.StructField) (sdkparam.OutputParamSpec[Output], error) {
	return f(field)
}

var _ OutputParamSpecFactory[any] = OutputParamSpecFactoryFunc[any](nil)
