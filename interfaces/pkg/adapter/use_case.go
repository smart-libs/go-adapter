package adapter

import "context"

type (
	// UseCaseHandler is the abstraction layer between the external world (the adapter) and the app world.
	// The UseCaseHandler uses the InputAccessor to retrieve the input from the external world provided by an input
	// adapter that knows the protocol used to retrieve the parameter. With the input received it invokes an application
	// use case and with its result uses the OutputBuilder, provided by the input adapter, to build the Output to be
	// returned, if needed.
	// Any panic generated by the UseCaseHandler must be considered a failure.
	// Each input adapter is responsible for providing a mapping strategy to retrieve and return parameters. This
	// mapping strategy will be used by the accessor to identify what shall be retrieved, and what the OutputBuilder
	// shall do with the value received to build the target Output.
	// Each input adapter is responsible for providing a binding strategy to identify which UseCaseHandler to invoke.
	UseCaseHandler interface {
		Invoke(ctx context.Context, accessor InputAccessor, builder OutputBuilder) Output
	}
)