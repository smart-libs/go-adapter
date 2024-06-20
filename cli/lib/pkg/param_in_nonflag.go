package cliadpt

import (
	"fmt"
	sdkparam "github.com/smart-libs/go-adapter/sdk/lib/pkg/param"
	"strconv"
)

const (
	nonFlagTag = "non-flags"
)

func init() {
	getInputParamSpecFactoryRegistry().AddOption1(nonFlagTag, "*", getAllNonFlagInParamValue)
	getInputParamSpecFactoryRegistry().AddOption3(nonFlagTag, nonFlagInParamPosGetterFactory)
}

func nonFlagInParamPosGetterFactory(tagValue string) (func(input Input) (any, error), error) {
	pos, err := strconv.ParseInt(tagValue, 10, 16)
	if err != nil {
		return nil, err
	}
	return func(input Input) (any, error) {
		return getNonFlagInParamValue(input, int(pos))
	}, nil
}

func newNonFlagInParam(index any, options []sdkparam.Option, getter func(input Input) (any, error)) sdkparam.InputParamSpec[Input] {
	specName := fmt.Sprintf("%s[%v]", nonFlagTag, index)
	return newInputParamSpec(specName, getter, options...)
}

func newPosNonFlagInParam(pos int, options []sdkparam.Option) sdkparam.InputParamSpec[Input] {
	return newNonFlagInParam(pos, options, func(input Input) (any, error) {
		return getNonFlagInParamValue(input, pos)
	})
}

func newAllNonFlagInParam(options ...sdkparam.Option) sdkparam.InputParamSpec[Input] {
	return newNonFlagInParam("*", options, getAllNonFlagInParamValue)
}

func getAllNonFlagInParamValue(input Input) (any, error) {
	return input.FlagSet.Args(), nil
}

func getNonFlagInParamValue(input Input, pos int) (any, error) {
	args := input.FlagSet.Args()
	if pos >= len(args) {
		return nil, fmt.Errorf("index=[%d] out of bounds on input.FlagSet.Args()=[%v]", pos, args)
	}
	return args[pos], nil
}
