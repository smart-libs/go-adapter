package sdkusecasehandler

import (
	"context"
	"fmt"
	"github.com/smart-libs/go-adapter/interfaces/pkg/adapter"
	sdkhandler "github.com/smart-libs/go-adapter/sdk/lib/pkg/handler"
	sdkparam "github.com/smart-libs/go-adapter/sdk/lib/pkg/param"
)

type (
	// useCaseHandler is the use case that is used to invoke a specific app.UseCase instance
	useCaseHandler[Input any, Output any] struct {
		adapter.UseCaseHandler
		sdkparam.InputSpecs[Input]
		sdkparam.OutputSpecs[Output]
	}
)

func (u useCaseHandler[Input, Output]) Invoke(ctx context.Context, input Input, output Output) error {
	accessor := NewInputAccessor[Input](input, u.InputSpecs)
	outBuilder := NewOutputBuilder[Output](output, u.OutputSpecs)
	result, ok := u.UseCaseHandler.Invoke(ctx, accessor, outBuilder).(Output)
	if ok {
		return nil
	}

	return fmt.Errorf("UseCase=[%T] did not return [%T], but [%T]=[%v]", u.UseCaseHandler, output, result, result)
}

// NewUseCaseHandler creates a new app.UseCase handler
func NewUseCaseHandler[Input any, Output any](
	useCase adapter.UseCaseHandler,
	inSpecs sdkparam.InputSpecs[Input],
	outSpecs sdkparam.OutputSpecs[Output],
) sdkhandler.Handler[Input, Output] {
	return useCaseHandler[Input, Output]{
		UseCaseHandler: useCase,
		InputSpecs:     inSpecs,
		OutputSpecs:    outSpecs,
	}
}
