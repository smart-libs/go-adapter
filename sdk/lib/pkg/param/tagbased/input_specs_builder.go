package tagbased

import (
	"github.com/smart-libs/go-adapter/interfaces/pkg/adapter"
	sdkparam "github.com/smart-libs/go-adapter/sdk/lib/pkg/param"
	"reflect"
)

type (
	// AbstractInputSpecBuilder builds the InputParamSpec[Input] object without using generics
	AbstractInputSpecBuilder interface {
		AddInputParamSpec(ref adapter.ParamRef, field reflect.StructField) error
	}

	InputSpecBuilder[Input any] interface {
		AbstractInputSpecBuilder
		Build() sdkparam.InputSpecs[Input]
	}

	defaultInputSpecBuilder[Input any] struct {
		factory   InputParamSpecFactory[Input]
		inputSpec sdkparam.InputSpecs[Input]
	}
)

func NewInputSpecsBuilder[Input any](factory InputParamSpecFactory[Input]) InputSpecBuilder[Input] {
	return &defaultInputSpecBuilder[Input]{factory: factory}
}

func (d *defaultInputSpecBuilder[Input]) Build() sdkparam.InputSpecs[Input] {
	return d.inputSpec
}

func (d *defaultInputSpecBuilder[Input]) AddInputParamSpec(ref adapter.ParamRef, field reflect.StructField) (err error) {
	var inputParamSpec sdkparam.InputParamSpec[Input]
	if inputParamSpec, err = d.factory.CreateInputParamSpec(field); err != nil {
		return
	}
	d.inputSpec.AddParamSpec(ref, inputParamSpec)
	return nil
}
