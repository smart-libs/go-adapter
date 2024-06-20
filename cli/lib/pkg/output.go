package cliadpt

import (
	sdkparam "github.com/smart-libs/go-adapter/sdk/lib/pkg/param"
)

type (
	Output struct {
		ExitActionFunc func() int
	}

	OutputSpec = sdkparam.OutputSpecs[*Output]
)

func NewOutput() *Output {
	return &Output{ExitActionFunc: func() int { return 0 }}
}
