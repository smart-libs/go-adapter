package tagbasedhandler

import (
	"context"
	"github.com/smart-libs/go-adapter/interfaces/pkg/adapter"
	"reflect"
)

// Before invoking the handler function it is needed to prepare the input arguments to be passed to the function call.
// In the GO reflection package, the function invocation (reflect.Value.Call(args)) requires an array of reflect.Value
// instances each one with one the input arguments.
// The argFactoryFunc data type is created during the function input argument analysis for each argument needed to create
// the array of arguments.
type (
	// argFactoryFunc is the function associated with an input argument that will use the app.InputAccessor object
	// to retrieve values needed to create the input argument.
	// Notice that the argFactoryFunc does not know the input object from which the value will be retrieved, for this
	// it uses the app.InputAccessor. The app.InputAccessor object keeps argFactoryFunc away from using Go generics.
	argFactoryFunc func(ctx context.Context, accessor adapter.InputAccessor) reflect.Value
)
