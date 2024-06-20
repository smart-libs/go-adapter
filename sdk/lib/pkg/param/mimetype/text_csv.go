package mimetype

import (
	converter "github.com/smart-libs/go-crosscutting/converter/lib/pkg"
	"strings"
)

type (
	// FromTextCSV is used when you want to convert from a flattened string list separated by commas into an array of string
	FromTextCSV string
	// ToTextCSV is used when you want to convert an array of strings into a flattened string list separated by commas.
	ToTextCSV []string
)

func Load(registry converter.Registry) {
	converter.AddHandlerWithReturnNoError(registry, FromTextCSVToStringArray)
	converter.AddHandlerWithReturnNoError(registry, FromStringArrayToTextCSV)
}

func FromTextCSVToStringArray(str FromTextCSV) []string { return strings.Split(string(str), ",") }
func FromStringArrayToTextCSV(array ToTextCSV) string   { return strings.Join(array, ",") }
