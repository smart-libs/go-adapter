package httpadpt

import (
	"github.com/smart-libs/go-adapter/sdk/lib/pkg/param/tagbased"
)

var (
	inParamSpecFactoryRegistry tagbased.InputParamSpecFactoryRegistry[Request]
)

func getInputParamSpecFactoryRegistry() tagbased.InputParamSpecFactoryRegistry[Request] {
	if inParamSpecFactoryRegistry == nil {
		inParamSpecFactoryRegistry = tagbased.NewInputParamSpecFactoryRegistry[Request](Converters)
	}

	return inParamSpecFactoryRegistry
}

func createInParamSpecFactory() tagbased.InputParamSpecFactory[Request] {
	return tagbased.NewsInputParamSpecFactory(getInputParamSpecFactoryRegistry())
}
