package sdkusecasehandler

import (
	"github.com/smart-libs/go-adapter/interfaces/pkg/adapter"
	"github.com/smart-libs/go-adapter/sdk/lib/pkg/param"
)

type (
	// outputBuilder does not know the target app.Output that needs to be created by the adapter. Because of this,
	// its constructor requires an app.Output object as input that will be passed to all OutputBuilderActionFunc instances
	// executed
	outputBuilder[Output any] struct {
		// OutputSpecs with all actions to be performed per app.ParamRef
		sdkparam.OutputSpecs[Output]

		// Output is the target object that must be built by the OutputBuilderActionFunc executed
		Output Output
	}
)

func NewOutputBuilder[Output any](output Output, spec sdkparam.OutputSpecs[Output]) adapter.OutputBuilder {
	return outputBuilder[Output]{OutputSpecs: spec, Output: output}
}

func (o outputBuilder[Output]) WithParam(ref adapter.ParamRef, value any) adapter.OutputBuilder {
	if outputParamSpec := o.GetParamSpec(ref); outputParamSpec != nil {
		if err := outputParamSpec.SetValue(o.Output, value); err != nil {
			panic(err)
		}
	}
	return o
}

func (o outputBuilder[Output]) WithError(err error) adapter.OutputBuilder {
	if outputParamSpec := o.GetErrorParamSpec(); outputParamSpec != nil {
		if err2 := outputParamSpec.SetValue(o.Output, err); err2 != nil {
			panic(err2)
		}
	}
	return o
}

func (o outputBuilder[Output]) Build() adapter.Output { return o.Output }
