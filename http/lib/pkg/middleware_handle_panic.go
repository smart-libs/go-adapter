package httpadpt

import (
	"context"
	"fmt"
	serror "github.com/smart-libs/go-crosscutting/serror/lib/pkg"
)

type (
	handlePanicMiddleware struct {
		decorated Handler
	}
)

const (
	ContentTypeProblemDetail = "application/problem+json"
)

func (h handlePanicMiddleware) Invoke(ctx context.Context, input Request, output *Response) error {
	defer func() {
		if panicArg := recover(); panicArg != nil {
			if output == nil {
				// Cannot set response fields if output is nil
				// Log or handle this case appropriately
				return
			}
			statusCode := 500
			output.StatusCode = &statusCode
			output.Header = map[string][]string{"Content-Type": {ContentTypeProblemDetail}}
			err, ok := panicArg.(error)
			if !ok {
				err = serror.WrapAsInternalError(fmt.Errorf("%v", panicArg))
			}
			output.Body, _ = JSONProblemDetailFromError(err)
		}
	}()

	return h.decorated.Invoke(ctx, input, output)
}

func HandlePanic(handler Handler) Handler {
	return handlePanicMiddleware{decorated: handler}
}
