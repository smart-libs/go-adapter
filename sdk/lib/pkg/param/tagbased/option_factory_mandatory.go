package tagbased

import (
	"fmt"
	sdkparam "github.com/smart-libs/go-adapter/sdk/lib/pkg/param"
	converter "github.com/smart-libs/go-crosscutting/converter/lib/pkg"
	"reflect"
	"strings"
)

var (
	AssertOptionMap = map[string]func(field reflect.StructField, converters converter.Converters) (sdkparam.Option, error){
		"mandatory": func(_ reflect.StructField, _ converter.Converters) (sdkparam.Option, error) {
			return sdkparam.Mandatory(), nil
		},
		"notDefaultValue": func(_ reflect.StructField, _ converter.Converters) (sdkparam.Option, error) {
			return sdkparam.NotDefaultValueReflection(), nil
		},
		"notBlank": func(_ reflect.StructField, _ converter.Converters) (sdkparam.Option, error) {
			return sdkparam.NotBlankString(), nil
		},
		"notEmpty": func(_ reflect.StructField, _ converter.Converters) (sdkparam.Option, error) {
			return sdkparam.NotBlankString(), nil
		},
	}
)

func init() {
	OptionFactories = append(OptionFactories, createMandatoryOption)
}

func createMandatoryOption(field reflect.StructField, converters ...converter.Converters) (sdkparam.Option, error) {
	const tagName = "assert"
	if tagValue, found := field.Tag.Lookup(tagName); found {
		var assertionOptions []sdkparam.Option
		resolvedConverters := converter.ConvertersList(converters)
		assertions := strings.Split(tagValue, ",")
		for _, assertionName := range assertions {
			if optionFactory, found := AssertOptionMap[assertionName]; found {
				option, err := optionFactory(field, resolvedConverters)
				if err != nil {
					return nil, fmt.Errorf("failed to create assert named=[%s] in the tag %s", assertionName, tagName)
				}
				assertionOptions = append(assertionOptions, option)
			}
		}
		return sdkparam.AsSingleOptions(assertionOptions...), nil
	}
	return nil, nil
}
