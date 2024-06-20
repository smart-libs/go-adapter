package adapter

type (
	// Output is an abstraction of the UseCaseHandler result
	Output any

	// OutputBuilder is provided by the input adapter to allow the UseCaseHandler to return a result
	OutputBuilder interface {
		WithParam(ref ParamRef, value any) OutputBuilder
		WithError(error) OutputBuilder
		Build() Output
	}
)
