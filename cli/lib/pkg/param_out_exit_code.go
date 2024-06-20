package cliadpt

import (
	sdkparam "github.com/smart-libs/go-adapter/sdk/lib/pkg/param"
	converter "github.com/smart-libs/go-crosscutting/converter/lib/pkg"
)

const exitOutParam = "exit"

func init() {
	getOutParamSpecFactoryRegistry().AddOption1("status", "", setExitActionFunc)
}

func NewExitCodeOutParamSpec(options ...sdkparam.Option) sdkparam.OutputParamSpec[*Output] {
	return sdkparam.NewOutputParamSpec[*Output](sdkparam.NewSpec(exitOutParam, options...), setExitActionFunc)
}

func setExitActionFunc(output *Output, value any) error {
	exitCode, err := converter.To[int](Converters, value)
	if err != nil {
		return err
	}
	output.ExitActionFunc = func() int {
		return exitCode
	}
	return nil
}
