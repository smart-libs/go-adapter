package tagbased

import (
	sdkparam "github.com/smart-libs/go-adapter/sdk/lib/pkg/param"
	converter "github.com/smart-libs/go-crosscutting/converter/lib/pkg"
	"reflect"
)

func init() {
	OptionFactories = append(OptionFactories, createDefaultOption)
}

func createDefaultOption(field reflect.StructField, converters ...converter.Converters) (sdkparam.Option, error) {
	if tagValue, found := field.Tag.Lookup("default"); found {
		defaultValue, err := converter.ConvertersList(converters).ConvertToType(tagValue, field.Type)
		if err != nil {
			return nil, err
		}
		return sdkparam.Default(defaultValue), nil
	}
	return nil, nil
}
