package httpadpt

import (
	converter "github.com/smart-libs/go-crosscutting/converter/lib/pkg"
	converterdefault "github.com/smart-libs/go-crosscutting/converter/lib/pkg/default"
	serror "github.com/smart-libs/go-crosscutting/serror/lib/pkg"
	"net/http"
)

var (
	// ConverterRegistry it is the default converters registry for the HTTP adapter
	ConverterRegistry = converterdefault.NewRegistry()

	// Converters is a list of converter.Converters that will be used by the HTTP Adapter. This implementations tries
	// first the HTTP adapter conversions and, if no one succeeded, it tries to use the default converter.Converters.
	// The main idea it to used first converter functions specialized for the HTTP Adapter.
	Converters = converter.NewConvertersList(
		converterdefault.NewConverters(ConverterRegistry), // This is the local converters for the HTTP Adapter
		converterdefault.Converters,                       // default as fallback
	)
)

func init() {
	converter.AddHandler[error, int](ConverterRegistry, errorToStatusCode)
	converter.AddHandler[[]string, string](ConverterRegistry, firstStringPtrFromStringArray)
}

// firstStringPtrFromStringArray is used to return the first element of []string. In the HTTP protocol
// Query and Header parameter values are by default an array of string, but most of the time it has only one
// element, so this conversion function returns only the first or an empty string.
func firstStringPtrFromStringArray(values []string, first *string) error {
	if len(values) > 0 {
		*first = values[0]
	} else {
		*first = ""
	}
	return nil
}

func errorToStatusCode(err error, to *int) error {
	if err == nil {
		*to = http.StatusOK
		return nil
	}

	found := serror.IdentifyRootCause(
		err,
		func(err error) { *to = http.StatusInternalServerError }, // fallback
		serror.CallbackCondition{
			Condition: serror.IsIllegalArgumentError,
			Callback:  func(err error) { *to = http.StatusBadRequest },
		},
		serror.CallbackCondition{
			Condition: serror.IsNotFoundError,
			Callback:  func(err error) { *to = http.StatusNotFound },
		},
		serror.CallbackCondition{
			Condition: serror.IsDuplicateError,
			Callback:  func(err error) { *to = http.StatusConflict },
		},
		serror.CallbackCondition{
			Condition: serror.IsTimeoutError,
			Callback:  func(err error) { *to = http.StatusGatewayTimeout },
		},
		serror.CallbackCondition{
			Condition: serror.IsIllegalConfigError,
			Callback:  func(err error) { *to = http.StatusBadRequest },
		},
	)
	if !found {
		*to = http.StatusInternalServerError
	}
	return nil
}
