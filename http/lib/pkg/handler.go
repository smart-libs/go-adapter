package httpadpt

import (
	"context"
	sdkhandler "github.com/smart-libs/go-adapter/sdk/lib/pkg/handler"
)

type (
	funcBasedHandler struct {
		invokeFunc func(ctx context.Context, input Request, output *Response) error
	}

	Handler = sdkhandler.Handler[Request, *Response]
)

func (f funcBasedHandler) Invoke(ctx context.Context, input Request, output *Response) error {
	if f.invokeFunc == nil {
		return nil // Return nil error if function is nil (defensive programming)
	}
	return f.invokeFunc(ctx, input, output)
}

func MakeHandler(invoker func(ctx context.Context, input Request, output *Response) error) Handler {
	if invoker == nil {
		// Return a handler that does nothing (no-op handler)
		return funcBasedHandler{
			invokeFunc: func(ctx context.Context, input Request, output *Response) error {
				return nil
			},
		}
	}
	return funcBasedHandler{invokeFunc: invoker}
}
