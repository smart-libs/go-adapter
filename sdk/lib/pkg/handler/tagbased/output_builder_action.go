package tagbasedhandler

import "github.com/smart-libs/go-adapter/interfaces/pkg/adapter"

type (
	// OutputBuilderActionFunc represents the action that must be performed when `OutputBuilder.WithParam()` is invoked with a
	// specific `ParamRef` or when `OutputBuilder.WithError()` is invoked. These actions build the target app.Output instance.
	// This action is created by `MappedOutputSpecFactory.OutputActionFactory` when processing an output type.
	// The app.Output is what the adapter must build to return the app.UseCase result to the caller. It can be an
	// HTTP response or an SQS message for instance.
	// The outputParamValue argument is the output parameter value set by app.UseCase using the app.OutputBuilder to
	// build an output result.
	OutputBuilderActionFunc func(builder adapter.OutputBuilder, outputParamValue any)
)
