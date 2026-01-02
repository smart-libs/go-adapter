package httpadpt

import (
	converter "github.com/smart-libs/go-crosscutting/converter/lib/pkg"
	serror "github.com/smart-libs/go-crosscutting/serror/lib/pkg"
)

const (
	TagBody = "body"
)

func init() {
	getOutParamSpecFactoryRegistry().AddOption1(TagBody, "", setBodyBytes)
}

func setBodyBytes(output *Response, value any) error {
	const fName = "httpadpt.setBodyBytes"
	if err := IsResponseNil(output); err != nil {
		return serror.CmpError.Wrap(err, "%s: failed to convert value=[%v]", fName, value)
	}

	bodyBytes, err := converter.To[[]byte](Converters, value)
	if err != nil {
		return serror.CmpError.Wrap(err, "%s: failed to convert value=[%v]", fName, value)
	}
	output.Body = bodyBytes
	return nil
}
