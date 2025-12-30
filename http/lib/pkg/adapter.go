package httpadpt

import "context"

type (
	Adapter interface {
		Start(ctx context.Context) error
		Stop(ctx context.Context) error
	}
)
