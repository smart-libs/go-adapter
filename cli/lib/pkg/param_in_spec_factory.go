package cliadpt

import (
	sdkparam "github.com/smart-libs/go-adapter/sdk/lib/pkg/param"
	"github.com/smart-libs/go-adapter/sdk/lib/pkg/param/tagbased"
)

var (
	inParamSpecFactoryRegistry tagbased.InputParamSpecFactoryRegistry[Input]
)

func getInputParamSpecFactoryRegistry() tagbased.InputParamSpecFactoryRegistry[Input] {
	if inParamSpecFactoryRegistry == nil {
		inParamSpecFactoryRegistry = tagbased.NewInputParamSpecFactoryRegistry[Input](Converters)
	}

	return inParamSpecFactoryRegistry
}

func createInParamSpecFactory() tagbased.InputParamSpecFactory[Input] {
	return tagbased.NewsInputParamSpecFactory(inParamSpecFactoryRegistry)
}

// newInputParamSpec adds the CLI converters
func newInputParamSpec(specName string, getter func(input Input) (any, error), options ...sdkparam.Option) sdkparam.InputParamSpec[Input] {
	return sdkparam.NewInputParamSpec[Input](specName, options, Converters, getter)
}
