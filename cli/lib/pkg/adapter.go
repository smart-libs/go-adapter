package cliadpt

import "context"

type (
	Adapter interface {
		// Run evaluates the binding conditions selecting the first condition that returns true and executes the
		// handler associated with this condition.
		//Run(ctx context.Context, input Input) (exitCode int)
		Run(ctx context.Context, args ...string) (exitCode int)
	}
)

var (
	DebugEnabled = false
)
