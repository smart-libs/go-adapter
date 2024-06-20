package tagbasedhandler

import (
	"context"
	"github.com/smart-libs/go-adapter/interfaces/pkg/adapter"
	"reflect"
)

type (
	handler struct {
		targetFunc           reflect.Value
		argFactories         []argFactoryFunc
		outputBuilderActions []OutputBuilderActionFunc
	}
)

func (h handler) Invoke(ctx context.Context, accessor adapter.InputAccessor, builder adapter.OutputBuilder) adapter.Output {
	// Build the function input using the app.InputAccessor
	var args []reflect.Value
	for _, factory := range h.argFactories {
		args = append(args, factory(ctx, accessor))
	}

	// Execute the function that works like a handler given by the adapter user
	handlerFuncOutput := h.targetFunc.Call(args)

	// Build the app.Output using the function return and the app.OutputBuilder
	for i, outputBuilderAction := range h.outputBuilderActions {
		if outputBuilderAction == nil {
			continue
		}
		outputBuilderAction(builder, handlerFuncOutput[i])
	}

	return builder.Build()
}
