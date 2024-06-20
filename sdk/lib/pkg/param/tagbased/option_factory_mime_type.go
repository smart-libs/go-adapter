package tagbased

import (
	"encoding/json"
	"fmt"
	sdkparam "github.com/smart-libs/go-adapter/sdk/lib/pkg/param"
	"github.com/smart-libs/go-adapter/sdk/lib/pkg/param/mimetype"
	converter "github.com/smart-libs/go-crosscutting/converter/lib/pkg"
	"reflect"
)

// The mime-type option can be used as input or output. The decision is based on the input data type. If you don't want to
// depend on the input type, then use i-format or o-format.

var (
	MimeTypeOptionMap = map[string]func(field reflect.StructField, converters converter.Converters) (sdkparam.Option, error){
		"text/csv": func(field reflect.StructField, converters converter.Converters) (sdkparam.Option, error) {
			return optionTextCSV, nil
		},
		"application/json": func(field reflect.StructField, converters converter.Converters) (sdkparam.Option, error) {
			return optionApplicationJSON, nil
		},
	}
)

func init() {
	OptionFactories = append(OptionFactories, createMimeTypeOption)
}

func createMimeTypeOption(field reflect.StructField, converters ...converter.Converters) (sdkparam.Option, error) {
	const tagName = "mime-type"
	resolvedConverters := converter.ConvertersList(converters)
	if tagValue, found := field.Tag.Lookup(tagName); found {
		if optionFactory, found := MimeTypeOptionMap[tagValue]; found {
			return optionFactory(field, resolvedConverters)
		}
		return nil, fmt.Errorf("unknow value=[%s] using with tag %s", tagValue, tagName)
	}
	return nil, nil
}

func optionTextCSV(spec sdkparam.Spec, inputValue any) (any, error) {
	switch value := inputValue.(type) {
	// input is a string, then output will be []string
	case string:
		return mimetype.FromTextCSVToStringArray(mimetype.FromTextCSV(value)), nil
	// input is a FromTextCSV, then output will be []string
	case mimetype.FromTextCSV:
		return mimetype.FromTextCSVToStringArray(value), nil
	// input is a []stringV, then output will be string
	case []string:
		return mimetype.FromStringArrayToTextCSV(value), nil
	// input is a ToTextCSV, then output will be string
	case mimetype.ToTextCSV:
		return mimetype.FromStringArrayToTextCSV(value), nil
	default:
		return nil, fmt.Errorf("%s: input-type[%T], value=[%v] cannot be used as mime-type=text/csv",
			spec.Name(), inputValue, inputValue)
	}
}

func optionApplicationJSON(_ sdkparam.Spec, inputValue any) (any, error) {
	switch value := inputValue.(type) {
	// input is a string, then output will be json.RawMessage to allow a converter to parse to the right data type given (unmarshalling)
	case string:
		return json.RawMessage(value), nil
	// input is a RawMessage, then output will be json.RawMessage to allow a converter to parse to the right data type given (unmarshalling)
	case json.RawMessage:
		return value, nil
	// input is a []byte, then output will be json.RawMessage to allow a converter to parse to the right data type given (unmarshalling)
	case []byte:
		return json.RawMessage(value), nil
	default:
		// will convert to json.RawMessage assuming the caller wants to marshalling to JSON
		return json.Marshal(value)
	}
}
