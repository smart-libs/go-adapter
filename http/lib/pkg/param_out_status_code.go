package httpadpt

import (
	converter "github.com/smart-libs/go-crosscutting/converter/lib/pkg"
	serror "github.com/smart-libs/go-crosscutting/serror/lib/pkg"
)

const (
	TagStatusCode = "statuscode"
)

func init() {
	getOutParamSpecFactoryRegistry().AddOption1(TagStatusCode, "", setStatusCode)
}

func setStatusCode(output *Response, value any) error {
	const fName = "httpadpt.setStatusCode"
	if err := IsResponseNil(output); err != nil {
		return serror.CmpError.Wrap(err, "%s: failed to convert value=[%v]", fName, value)
	}

	statusCode, err := converter.To[int](Converters, value)
	if err != nil {
		return serror.CmpError.Wrap(err, "%s: failed to convert value=[%v]", fName, value)
	}
	output.StatusCode = &statusCode
	return nil
}
