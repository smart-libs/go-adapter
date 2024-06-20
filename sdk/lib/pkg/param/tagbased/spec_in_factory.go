package tagbased

import (
	"github.com/smart-libs/go-adapter/sdk/lib/pkg/param"
	"reflect"
)

type (
	InputParamSpecFactory[Input any] interface {
		CreateInputParamSpec(reflect.StructField) (sdkparam.InputParamSpec[Input], error)
	}

	InputParamSpecFactoryFunc[Input any] func(reflect.StructField) (sdkparam.InputParamSpec[Input], error)
)

func (f InputParamSpecFactoryFunc[Input]) CreateInputParamSpec(field reflect.StructField) (sdkparam.InputParamSpec[Input], error) {
	return f(field)
}
