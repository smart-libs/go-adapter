package sdkparam

import "github.com/smart-libs/go-adapter/interfaces/pkg/adapter"

type (
	// OutputSpecs has all the OutputParamSpec instances per ParamRef that must be used by the app.OutputBuilder
	// to create the app.Output.
	OutputSpecs[Output any] struct {
		// paramSpecs must be set by the real adapter with the actions that must be performed in order to
		// build the app.Output when each app.ParamRef is provided.
		paramSpecs map[adapter.ParamRef]OutputParamSpec[Output]

		// errorParamSpec identifies the OutputParamSpec that must be used when the app.OutputBuilder.WithError()
		// method is invoked.
		errorParamSpec OutputParamSpec[Output]
	}
)

func (o *OutputSpecs[Output]) AddParamSpec(ref adapter.ParamRef, spec OutputParamSpec[Output]) *OutputSpecs[Output] {
	if o.paramSpecs == nil {
		o.paramSpecs = make(map[adapter.ParamRef]OutputParamSpec[Output])
	}

	o.paramSpecs[ref] = spec
	return o
}

func (o OutputSpecs[Output]) GetErrorParamSpec() OutputParamSpec[Output] {
	return o.errorParamSpec
}

func (o OutputSpecs[Output]) GetParamSpec(ref adapter.ParamRef) OutputParamSpec[Output] {
	if o.paramSpecs != nil {
		return o.paramSpecs[ref]
	}
	return nil
}

func (o *OutputSpecs[Output]) AddErrorParamSpec(action OutputParamSpec[Output]) *OutputSpecs[Output] {
	o.errorParamSpec = action
	return o
}
