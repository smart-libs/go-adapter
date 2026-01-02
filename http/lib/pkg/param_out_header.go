package httpadpt

import (
	converter "github.com/smart-libs/go-crosscutting/converter/lib/pkg"
	serror "github.com/smart-libs/go-crosscutting/serror/lib/pkg"
)

const (
	TagResponseHeader = "header"
)

func init() {
	getOutParamSpecFactoryRegistry().AddOption2(TagResponseHeader, setResponseHeader)
}

func setResponseHeader(output *Response, headerName string, value any) error {
	const fName = "httpadpt.setResponseHeader"
	if err := IsResponseNil(output); err != nil {
		return serror.CmpError.Wrap(err, "%s: failed to convert value=[%v]", fName, value)
	}

	responseReaderValue, err := converter.To[[]string](Converters, value)
	if err != nil {
		return serror.CmpError.Wrap(err, "%s: failed to convert value=[%v]", fName, value)
	}
	if output.Header == nil {
		output.Header = make(map[string][]string)
	}
	output.Header[headerName] = responseReaderValue
	return nil
}
