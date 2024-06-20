package sdkhandler

import "context"

type (
	// Handler is what the adapter invokes to delegate to SDK the app.UseCase invocation.
	// The adapter provides the Input and the Output objects that is what it has as input and what it expects as result
	// of the app.UseCase invocation. It is the responsibility of the Handler to create the app.InputAccessor and
	// the app.OutputBuilder to be provided to the app.UseCase.
	Handler[Input any, Output any] interface {
		Invoke(ctx context.Context, input Input, output Output) error
	}
)
