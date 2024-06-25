package cliadpt

import (
	"context"
	"fmt"
	"os"
)

type (
	// MultiFlagSetAdapter is based on a map of Config organized per UseCaseName. It first try to identify the
	// UseCaseName using as input the arg[0] (binary name) or arg[1]. If it finds a Config associated with the
	// UseCaseName, then it runs using that config.
	MultiFlagSetAdapter struct {
		ConfigPerUseCaseName
	}
)

func NewMultiFlagSetAdapter(configPerUseCaseName ConfigPerUseCaseName) (Adapter, error) {
	if len(configPerUseCaseName) == 0 {
		return nil, NewInvalidConfigError(fmt.Errorf("ConfigPerUseCaseName is mandatory"))
	}

	return MultiFlagSetAdapter{ConfigPerUseCaseName: configPerUseCaseName}, nil
}

func (m MultiFlagSetAdapter) Run(ctx context.Context, args ...string) (exitCode int) {
	resolvedArgs := firstNotEmpty(args, os.Args)
	useCaseName := resolvedArgs[0]
	result, err := m.tryRun(ctx, m.find(useCaseName), resolvedArgs[1:])
	if err != nil && len(resolvedArgs) > 1 {
		useCaseName = resolvedArgs[1]
		result, err = m.tryRun(ctx, m.find(useCaseName), resolvedArgs[2:])
	}

	if err != nil {
		panic(ErrUseCaseNotFound{Args: resolvedArgs, UseCaseName: useCaseName})
	}

	return result
}

func (m MultiFlagSetAdapter) tryRun(ctx context.Context, config *Config, args []string) (int, error) {
	if config == nil {
		return 0, ErrNilConfig{}
	}

	adapter, err := NewSingleFlagSetAdapter(*config)
	if err != nil {
		panic(err)
	}

	return adapter.Run(ctx, args...), nil
}
