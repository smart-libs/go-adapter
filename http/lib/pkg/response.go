package httpadpt

import assertions "github.com/smart-libs/go-crosscutting/assertions/lib/pkg"

type (
	ParamName = string

	Response struct {
		StatusCode *int
		Body       []byte
		Header     map[ParamName][]string
	}
)

func IsResponseNil(resp *Response) error {
	if resp == nil {
		return assertions.WrapAsIllegalArgumentValue("Response", resp)
	}

	return nil
}
