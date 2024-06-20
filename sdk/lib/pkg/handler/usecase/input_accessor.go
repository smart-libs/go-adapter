package sdkusecasehandler

import (
	"fmt"
	"github.com/smart-libs/go-adapter/interfaces/pkg/adapter"
	sdkparam "github.com/smart-libs/go-adapter/sdk/lib/pkg/param"
)

type (
	InputAccessor[Input any] struct {
		Input Input
		sdkparam.InputSpecs[Input]
	}
)

func NewInputAccessor[Input any](input Input, spec sdkparam.InputSpecs[Input]) adapter.InputAccessor {
	return InputAccessor[Input]{Input: input, InputSpecs: spec}
}

func (i InputAccessor[Input]) GetValue(ref adapter.ParamRef) (any, error) {
	spec, err := i.getParam(ref)
	if err != nil {
		return nil, err
	}
	return spec.GetValue(i.Input)
}

func (i InputAccessor[Input]) CopyValue(ref adapter.ParamRef, target any) error {
	spec, err := i.getParam(ref)
	if err != nil {
		return err
	}
	return spec.CopyValue(i.Input, target)
}

func (i InputAccessor[Input]) getParam(ref adapter.ParamRef) (sdkparam.InputParamSpec[Input], error) {
	if spec := i.InputSpecs.GetParamSpec(ref); spec != nil {
		return spec, nil
	}
	return nil, fmt.Errorf("InputAccessor: no ParamRef [%s] found", ref.GetID())
}
