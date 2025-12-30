package httpadpt

import (
	"github.com/smart-libs/go-adapter/sdk/lib/pkg/param/tagbased"
)

var (
	outParamSpecFactoryRegistry tagbased.OutputParamSpecFactoryRegistry[*Response]
)

func getOutParamSpecFactoryRegistry() tagbased.OutputParamSpecFactoryRegistry[*Response] {
	if outParamSpecFactoryRegistry == nil {
		outParamSpecFactoryRegistry = tagbased.NewOutputParamSpecFactoryRegistry[*Response](Converters)
	}

	return outParamSpecFactoryRegistry
}

func createOutParamSpecFactory() tagbased.OutputParamSpecFactory[*Response] {
	return tagbased.NewOutputParamSpecFactory(getOutParamSpecFactoryRegistry())
}
