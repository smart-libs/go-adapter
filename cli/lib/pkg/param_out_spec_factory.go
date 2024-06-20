package cliadpt

import (
	"github.com/smart-libs/go-adapter/sdk/lib/pkg/param/tagbased"
)

var (
	outParamSpecFactoryRegistry tagbased.OutputParamSpecFactoryRegistry[*Output]
)

func getOutParamSpecFactoryRegistry() tagbased.OutputParamSpecFactoryRegistry[*Output] {
	if outParamSpecFactoryRegistry == nil {
		outParamSpecFactoryRegistry = tagbased.NewOutputParamSpecFactoryRegistry[*Output](Converters)
	}

	return outParamSpecFactoryRegistry
}

func createOutParamSpecFactory() tagbased.OutputParamSpecFactory[*Output] {
	return tagbased.NewOutputParamSpecFactory(outParamSpecFactoryRegistry)
}
