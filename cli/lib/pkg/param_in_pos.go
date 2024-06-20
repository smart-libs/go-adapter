package cliadpt

import (
	"fmt"
	sdkparam "github.com/smart-libs/go-adapter/sdk/lib/pkg/param"
	"strconv"
)

const (
	posTag = "pos"
)

func init() {
	getInputParamSpecFactoryRegistry().AddOption2(posTag, func(input Input, tagValue string) (any, error) {
		pos, err := strconv.ParseInt(tagValue, 10, 16)
		if err != nil {
			return nil, err
		}
		return getPosInParamValue(input, int(pos))
	})
}

func newPosInParam(pos int, options ...sdkparam.Option) sdkparam.InputParamSpec[Input] {
	specName := fmt.Sprintf("%s[%d]", envTag, pos)
	getter := func(input Input) (any, error) { return getPosInParamValue(input, pos) }
	return newInputParamSpec(specName, getter, options...)
}

func getPosInParamValue(input Input, pos int) (any, error) {
	if input.Args == nil || pos >= len(input.Args) {
		return nil, nil
	}

	return input.Args[pos], nil
}
