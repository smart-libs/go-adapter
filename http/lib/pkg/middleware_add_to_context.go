package httpadpt

import "context"

// AddToContextMiddleware creates a middleware that adds custom data to the context using the provided function.
func AddToContextMiddleware(addToContext func(ctx context.Context) context.Context) Middleware {
	return func(next Handler) Handler {
		return MakeHandler(func(ctx context.Context, input Request, output *Response) error {
			if addToContext == nil {
				// If addToContext is nil, use original context
				return next.Invoke(ctx, input, output)
			}
			newCtx := addToContext(ctx)
			return next.Invoke(newCtx, input, output)
		})
	}
}
