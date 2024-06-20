package cliadpt

import (
	"fmt"
	sdkparam "github.com/smart-libs/go-adapter/sdk/lib/pkg/param"
)

const (
	envTag = "env"
)

func init() {
	getInputParamSpecFactoryRegistry().AddOption2(envTag, getEnvInParamValue)
}

func newEnvInParam(envName string, options ...sdkparam.Option) sdkparam.InputParamSpec[Input] {
	specName := fmt.Sprintf("%s[%s]", envTag, envName)
	getter := func(input Input) (any, error) { return getEnvInParamValue(input, envName) }
	return newInputParamSpec(specName, getter, options...)
}

func getEnvInParamValue(input Input, envName string) (any, error) {
	if input.EnvGetter == nil {
		return nil, nil
	}
	if value, found := input.EnvGetter(envName); found {
		return value, nil
	}
	return nil, nil
}
