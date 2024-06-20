package sdkusecasehandler

import (
	"github.com/smart-libs/go-adapter/interfaces/pkg/adapter"
	sdkhandler "github.com/smart-libs/go-adapter/sdk/lib/pkg/handler"
	sdkparam "github.com/smart-libs/go-adapter/sdk/lib/pkg/param"
)

type (
	builder[Input any, Output any] struct {
		adapter.UseCaseHandler
		sdkparam.InputSpecs[Input]
		sdkparam.OutputSpecs[Output]
	}
)

// NewBuilderForUseCase creates a Use Case Builder given the app.UseCase to be invoked.
func NewBuilderForUseCase[Input any, Output any](useCase adapter.UseCaseHandler) sdkhandler.InputSpecStep[Input, Output] {
	return &builder[Input, Output]{
		UseCaseHandler: useCase,
	}
}

func (b *builder[Input, Output]) Output() sdkhandler.OutputSpecStep[Input, Output] {
	return b
}

func (b *builder[Input, Output]) Build() sdkhandler.Handler[Input, Output] {
	return NewUseCaseHandler[Input, Output](b.UseCaseHandler, b.InputSpecs, b.OutputSpecs)
}

func (b *builder[Input, Output]) WithErrorParamSpec(spec sdkparam.OutputParamSpec[Output]) sdkhandler.OutputSpecStep[Input, Output] {
	b.OutputSpecs.AddErrorParamSpec(spec)
	return b
}

func (b *builder[Input, Output]) WithOutParamSpec(ref adapter.ParamRef, actionFunc sdkparam.OutputParamSpec[Output]) sdkhandler.OutputSpecStep[Input, Output] {
	b.OutputSpecs.AddParamSpec(ref, actionFunc)
	return b
}

func (b *builder[Input, Output]) WithInParamSpec(ref adapter.ParamRef, spec sdkparam.InputParamSpec[Input]) sdkhandler.InputSpecStep[Input, Output] {
	b.InputSpecs.AddParamSpec(ref, spec)
	return b
}
