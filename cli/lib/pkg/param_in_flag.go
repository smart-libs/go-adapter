package cliadpt

import (
	"fmt"
	sdkparam "github.com/smart-libs/go-adapter/sdk/lib/pkg/param"
)

const (
	flagTag = "flag"
)

func init() {
	getInputParamSpecFactoryRegistry().AddOption2(flagTag, getFlagInParamValue)
}

func newFlagInParam(flagName string, options ...sdkparam.Option) sdkparam.InputParamSpec[Input] {
	specName := fmt.Sprintf("arg[%s:%s]", flagTag, flagName)
	getter := func(input Input) (any, error) { return getFlagInParamValue(input, flagName) }
	return newInputParamSpec(specName, getter, options...)
}

func getFlagInParamValue(input Input, flagName string) (any, error) {
	value, _ := input.FlagSet.GetValue(flagName)
	return value, nil
}
