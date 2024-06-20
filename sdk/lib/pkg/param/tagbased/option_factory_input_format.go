package tagbased

import (
	"fmt"
	sdkparam "github.com/smart-libs/go-adapter/sdk/lib/pkg/param"
	converter "github.com/smart-libs/go-crosscutting/converter/lib/pkg"
	convertertypes "github.com/smart-libs/go-crosscutting/converter/lib/pkg/types"
	"reflect"
)

var (
	// InputFormatOptionMap should have all types in the ConvertTo as pointer because the parameter can be optional and so be nil
	InputFormatOptionMap = map[string]func(field reflect.StructField, converters converter.Converters) (sdkparam.Option, error){
		"RFC3339Year": func(_ reflect.StructField, converters converter.Converters) (sdkparam.Option, error) {
			type alias = convertertypes.RFC3339Year
			return sdkparam.ConvertTo[*convertertypes.StringTime[alias]](converters), nil
		},
		"RFC3339YearMonth": func(_ reflect.StructField, converters converter.Converters) (sdkparam.Option, error) {
			type alias = convertertypes.RFC3339YearMonth
			return sdkparam.ConvertTo[*convertertypes.StringTime[alias]](converters), nil
		},
		"RFC3339Date": func(_ reflect.StructField, converters converter.Converters) (sdkparam.Option, error) {
			type alias = convertertypes.RFC3339Date
			return sdkparam.ConvertTo[*convertertypes.StringTime[alias]](converters), nil
		},
		"RFC3339WithMin": func(_ reflect.StructField, converters converter.Converters) (sdkparam.Option, error) {
			type alias = convertertypes.RFC3339WithMin
			return sdkparam.ConvertTo[*convertertypes.StringTime[alias]](converters), nil
		},
		"RFC3339": func(_ reflect.StructField, converters converter.Converters) (sdkparam.Option, error) {
			type alias = convertertypes.RFC3339
			return sdkparam.ConvertTo[*convertertypes.StringTime[alias]](converters), nil
		},
		"RFC3339WithFraction": func(_ reflect.StructField, converters converter.Converters) (sdkparam.Option, error) {
			type alias = convertertypes.RFC3339WithFraction
			return sdkparam.ConvertTo[*convertertypes.StringTime[alias]](converters), nil
		},
		"DateCompressed": func(_ reflect.StructField, converters converter.Converters) (sdkparam.Option, error) {
			type alias = convertertypes.DateCompressed
			return sdkparam.ConvertTo[*convertertypes.StringTime[alias]](converters), nil
		},
		"DateFree": func(_ reflect.StructField, converters converter.Converters) (sdkparam.Option, error) {
			type alias = convertertypes.DateFree
			return sdkparam.ConvertTo[*convertertypes.StringTime[alias]](converters), nil
		},
		"UNIX-time": func(_ reflect.StructField, converters converter.Converters) (sdkparam.Option, error) {
			type alias = convertertypes.UNIX
			return sdkparam.ConvertTo[*convertertypes.StringUnixTime[alias]](converters), nil
		},
		"UNIXMilli-time": func(_ reflect.StructField, converters converter.Converters) (sdkparam.Option, error) {
			type alias = convertertypes.UNIXMilli
			return sdkparam.ConvertTo[*convertertypes.StringUnixTime[alias]](converters), nil
		},
		"UNIXMicro-time": func(_ reflect.StructField, converters converter.Converters) (sdkparam.Option, error) {
			type alias = convertertypes.UNIXMicro
			return sdkparam.ConvertTo[*convertertypes.StringUnixTime[alias]](converters), nil
		},
		"UNIXNano-time": func(_ reflect.StructField, converters converter.Converters) (sdkparam.Option, error) {
			type alias = convertertypes.UNIXNano
			return sdkparam.ConvertTo[*convertertypes.StringUnixTime[alias]](converters), nil
		},
	}
)

func AddInputFormatAlias(alias, current string) {
	if currentFunc := InputFormatOptionMap[current]; currentFunc != nil {
		InputFormatOptionMap[alias] = currentFunc
		return
	}
	panic(fmt.Errorf("no i-format=[%s] found to create alias=[%s]", current, alias))
}

func init() {
	OptionFactories = append(OptionFactories, createInputFormatOption)
	AddInputFormatAlias("year", "RFC3339Year")
	AddInputFormatAlias("yyyy-mm-dd", "RFC3339Date")
	AddInputFormatAlias("yyyymmdd", "DateCompressed")
}

func createInputFormatOption(field reflect.StructField, converters ...converter.Converters) (sdkparam.Option, error) {
	const tagName = "i-format"
	resolvedConverters := converter.ConvertersList(converters)
	if tagValue, found := field.Tag.Lookup(tagName); found {
		if optionFactory, found := InputFormatOptionMap[tagValue]; found {
			return optionFactory(field, resolvedConverters)
		}
		return nil, fmt.Errorf("unknow value=[%s] using with tag %s", tagValue, tagName)
	}
	return nil, nil
}
