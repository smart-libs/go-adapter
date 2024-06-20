package cliadpt

import (
	sdkusecasehandler "github.com/smart-libs/go-adapter/sdk/lib/pkg/handler/usecase"
	sdkparam "github.com/smart-libs/go-adapter/sdk/lib/pkg/param"
	"os"
)

type (
	EnvGetter func(string) (value string, found bool)

	Input struct {
		FlagSet
		Args []string
		EnvGetter
	}

	InputFactory func(EnvGetter, FlagSet, []string) (Input, error)

	InputSpec = sdkparam.InputSpecs[Input]

	InputAccessor = sdkusecasehandler.InputAccessor[Input]
)

func NewInput(envGetter EnvGetter, flagSet FlagSet, args ...string) (Input, error) {
	if !flagSet.Parsed() {
		if err := flagSet.Parse(args); err != nil {
			return Input{}, err
		}
	}

	if envGetter == nil {
		envGetter = os.LookupEnv
	}

	return Input{
		Args:      args,
		FlagSet:   flagSet,
		EnvGetter: envGetter,
	}, nil
}
