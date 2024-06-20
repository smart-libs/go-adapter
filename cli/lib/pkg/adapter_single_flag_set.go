package cliadpt

import (
	"context"
	"fmt"
	"os"
	"strings"
)

type (
	// SingleFlagSetAdapter is based on a single Config object that has a single FlagSet
	SingleFlagSetAdapter struct {
		Config
	}
)

func NewSingleFlagSetAdapter(config Config) (Adapter, error) {
	if config.Bindings == nil {
		return nil, NewInvalidConfigError(fmt.Errorf("config.Bindings is mandatory"))
	}
	if config.FlagSet == nil {
		return nil, NewInvalidConfigError(fmt.Errorf("config.FlagSet is mandatory"))
	}
	return SingleFlagSetAdapter{
		Config: Config{
			ArgumentsValidationDisabled: config.ArgumentsValidationDisabled,
			Bindings:                    config.Bindings,
			EnvGetter:                   firstNotNil(config.EnvGetter, os.LookupEnv),
			FlagSet:                     config.FlagSet,
		},
	}, nil
}

func (s SingleFlagSetAdapter) Run(ctx context.Context, args ...string) (exitCode int) {
	input, err := NewInput(s.Config.EnvGetter, s.Config.FlagSet, firstNotEmpty(args, os.Args[1:])...)
	if err != nil {
		panic(err)
	}
	return s.doRun(ctx, s.parse(input))
}

func (s SingleFlagSetAdapter) doRun(ctx context.Context, input Input) int {
	var (
		useCaseFound bool
	)

	output := NewOutput()
	for _, binding := range s.Bindings {
		if binding.EvaluateCondition(input) {
			if err := binding.Invoke(ctx, input, output); err != nil {
				panic(err)
			}
			useCaseFound = true
			break
		}
	}

	if !useCaseFound {
		panic(ErrUseCaseNotFound{Args: input.Args})
	}

	return output.ExitActionFunc()
}

func (s SingleFlagSetAdapter) isDashDashWasUsed() bool {
	for _, arg := range os.Args {
		if arg == "--" {
			return true
		}
	}
	return false
}

func (s SingleFlagSetAdapter) isThereAnyFlagInTheNonFlagsArray(input Input) error {
	for _, arg := range input.FlagSet.Args() {
		if strings.HasPrefix(arg, "-") {
			input.FlagSet.Usage()
			return fmt.Errorf("command line argument [%s] not expected as non-flag argument", arg)
		}
	}
	return nil
}

func (s SingleFlagSetAdapter) parse(input Input) Input {
	if !input.FlagSet.Parsed() {
		if err := input.FlagSet.Parse(input.Args); err != nil {
			panic(err)
		}
	}

	if s.Config.ArgumentsValidationDisabled {
		return input
	}

	// If -- is used to separate flags from non-flag arguments, then no additional validation is performed
	if s.isDashDashWasUsed() {
		return input
	}

	if err := s.isThereAnyFlagInTheNonFlagsArray(input); err != nil {
		panic(err)
	}

	return input
}
