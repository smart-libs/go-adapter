package cliadpt

import (
	"fmt"
	sdkparam "github.com/smart-libs/go-adapter/sdk/lib/pkg/param"
	"io"
	"os"
	"reflect"
)

const (
	printOutParam = "print"
	printStdout   = "stdout"
	printStderr   = "stderr"
)

func init() {
	getMask := func(field reflect.StructField) string {
		mask, found := field.Tag.Lookup("mask")
		if !found {
			mask = "%v"
		}
		return mask
	}

	getOutParamSpecFactoryRegistry().
		AddOption7(printOutParam, printStdout,
			func(field reflect.StructField, options []sdkparam.Option) (sdkparam.OutputParamSpec[*Output], error) {
				return NewPrintStdoutOutParamSpec(getMask(field), options...), nil
			},
		).
		AddOption7(printOutParam, printStderr,
			func(field reflect.StructField, options []sdkparam.Option) (sdkparam.OutputParamSpec[*Output], error) {
				return NewPrintStderrOutParamSpec(getMask(field), options...), nil
			},
		)
}

func newPrintOutParamSpec(file string, writer io.Writer, mask string, options ...sdkparam.Option) sdkparam.OutputParamSpec[*Output] {
	specName := fmt.Sprintf("%s:%s", printOutParam, file)
	return sdkparam.NewOutputParamSpec[*Output](sdkparam.NewSpec(specName, options...), printMessageTo(writer, mask))
}

func NewPrintStdoutOutParamSpec(mask string, options ...sdkparam.Option) sdkparam.OutputParamSpec[*Output] {
	return newPrintOutParamSpec(printStdout, os.Stdout, mask, options...)
}

func NewPrintStderrOutParamSpec(mask string, options ...sdkparam.Option) sdkparam.OutputParamSpec[*Output] {
	return newPrintOutParamSpec(printStderr, os.Stderr, mask, options...)
}

func printMessageTo(writer io.Writer, mask string) func(_ *Output, value any) error {
	return func(_ *Output, value any) error {
		_, err := fmt.Fprintf(writer, mask, value)
		return err
	}
}
