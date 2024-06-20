package tagbasedhandler

import (
	"context"
	"github.com/smart-libs/go-adapter/interfaces/pkg/adapter"
	"reflect"
)

// createContextArg creates is the argFactoryFunc instance that returns the context to be passed to the handler
func createContextArg(ctx context.Context, _ adapter.InputAccessor) reflect.Value {
	return reflect.ValueOf(ctx)
}
