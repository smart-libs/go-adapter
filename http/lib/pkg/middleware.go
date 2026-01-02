package httpadpt

type (
	Middleware = func(next Handler) Handler

	Middlewares = []Middleware
)

func WrapHandlerWithMiddlewares(handler Handler, middlewares Middlewares) Handler {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}
