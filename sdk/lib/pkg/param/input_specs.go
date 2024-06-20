package sdkparam

import (
	"github.com/smart-libs/go-adapter/interfaces/pkg/adapter"
)

type (
	// InputSpecs has all the InputParamSpec per app.ParamRef that can be used by the app.InputAccessor to return the
	// input values required by the sdkusecasehandler.UseCaseHandler or by the sdkfunchandler.FuncHandler.
	InputSpecs[Input any] struct {
		paramSpecs map[adapter.ParamRef]InputParamSpec[Input]
	}
)

func (i *InputSpecs[Input]) AddParamSpec(ref adapter.ParamRef, spec InputParamSpec[Input]) (oldParam InputParamSpec[Input]) {
	if i.paramSpecs == nil {
		i.paramSpecs = make(map[adapter.ParamRef]InputParamSpec[Input])
	}

	oldParam, i.paramSpecs[ref] = i.paramSpecs[ref], spec
	return oldParam
}

func (i *InputSpecs[Input]) GetParamSpec(ref adapter.ParamRef) InputParamSpec[Input] {
	if i.paramSpecs != nil {
		return i.paramSpecs[ref]
	}
	return nil
}
