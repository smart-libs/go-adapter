package tagbased

import (
	"github.com/smart-libs/go-adapter/interfaces/pkg/adapter"
	sdkparam "github.com/smart-libs/go-adapter/sdk/lib/pkg/param"
	"reflect"
)

type (
	AbstractOutputSpecBuilder interface {
		AddOutputParamSpec(ref adapter.ParamRef, field reflect.StructField) error
	}

	OutputSpecBuilder[Output any] interface {
		AbstractOutputSpecBuilder
		Build() sdkparam.OutputSpecs[Output]
	}

	defaultOutputSpecBuilder[Output any] struct {
		factory    OutputParamSpecFactory[Output]
		outputSpec sdkparam.OutputSpecs[Output]
	}
)

func NewOutputSpecsBuilder[Output any](factory OutputParamSpecFactory[Output]) OutputSpecBuilder[Output] {
	return &defaultOutputSpecBuilder[Output]{factory: factory}
}

func (d *defaultOutputSpecBuilder[Output]) Build() sdkparam.OutputSpecs[Output] {
	return d.outputSpec
}

func (d *defaultOutputSpecBuilder[Output]) AddOutputParamSpec(ref adapter.ParamRef, field reflect.StructField) (err error) {
	var inputParamSpec sdkparam.OutputParamSpec[Output]
	if inputParamSpec, err = d.factory.CreateOutputParamSpec(field); err != nil {
		return
	}
	d.outputSpec.AddParamSpec(ref, inputParamSpec)
	return nil
}
