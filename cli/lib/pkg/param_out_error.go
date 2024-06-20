package cliadpt

import (
	sdkparam "github.com/smart-libs/go-adapter/sdk/lib/pkg/param"
	"reflect"
)

func init() {
	getOutParamSpecFactoryRegistry().AddOption7("error", "panic",
		func(_ reflect.StructField, options []sdkparam.Option) (sdkparam.OutputParamSpec[*Output], error) {
			return NewPanicErrorOutParamSpec(options...), nil
		},
	)
}

const panicOutParam = "error:panic"

func NewPanicErrorOutParamSpec(options ...sdkparam.Option) sdkparam.OutputParamSpec[*Output] {
	return sdkparam.NewOutputParamSpec[*Output](sdkparam.NewSpec(panicOutParam, options...), setPanicExitActionFunc)
}

func setPanicExitActionFunc(output *Output, value any) error {
	output.ExitActionFunc = func() int {
		if value != nil {
			panic(value)
		}
		return 0
	}
	return nil
}
