package cliadpt

import (
	converter "github.com/smart-libs/go-crosscutting/converter/lib/pkg"
	converterdefault "github.com/smart-libs/go-crosscutting/converter/lib/pkg/default"
)

var (
	// ConverterRegistry it is the default converters registry for the CLI adapter
	ConverterRegistry = converterdefault.NewRegistry()

	// Converters is a list of converter.Converters that will be used by the CLI Adapter. This implementations tries
	// first the CLI adapter conversions and, if no one succeeded, it tries to use the default converter.Converters.
	// The main idea it to used first converter functions specialized for the CLI Adapter.
	Converters = converter.NewConvertersList(
		converterdefault.NewConverters(ConverterRegistry), // This is the local converters for the CLI Adapter
		converterdefault.Converters,                       // default as fallback
	)
)
