package tagbased

import (
	"github.com/smart-libs/go-adapter/sdk/lib/pkg/param"
	converter "github.com/smart-libs/go-crosscutting/converter/lib/pkg"
	converterdefault "github.com/smart-libs/go-crosscutting/converter/lib/pkg/default"
	"reflect"
)

type (
	OptionFactory func(field reflect.StructField, converters ...converter.Converters) (sdkparam.Option, error)
)

var (
	// OptionFactories is initialized with Option factories driven by field tags
	OptionFactories []OptionFactory
)

func CreateParamSpecOptionsFromTag(field reflect.StructField, convertersList ...converter.Converters) ([]sdkparam.Option, error) {
	list := convertersList
	if len(list) == 0 {
		list = []converter.Converters{converterdefault.Converters}
	}
	converters := converter.NewConvertersList(list...)
	var options []sdkparam.Option
	for _, factory := range OptionFactories {
		option, err := factory(field, converters)
		if err != nil {
			return nil, err
		}
		if option != nil {
			options = append(options, option)
		}
	}

	return options, nil
}
