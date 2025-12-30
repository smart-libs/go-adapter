package httpadpt

import (
	sdkparam "github.com/smart-libs/go-adapter/sdk/lib/pkg/param"
	"github.com/smart-libs/go-crosscutting/assertions/lib/pkg/check"
	serror "github.com/smart-libs/go-crosscutting/serror/lib/pkg"
)

type (
	OutErrorParamSpec struct{}
)

func (o OutErrorParamSpec) Name() string               { return "error" }
func (o OutErrorParamSpec) Options() []sdkparam.Option { return nil }

func (o OutErrorParamSpec) SetValue(output *Response, value any) error {
	if check.IsNil(value) {
		return nil // no error
	}
	if output == nil {
		return serror.CmpError.New("httpadpt.OutErrorParamSpec.SetValue: output is nil")
	}
	if err, ok := value.(error); ok {
		output.StatusCode = new(int)
		if convErr := errorToStatusCode(err, output.StatusCode); convErr != nil {
			return convErr
		}
		output.Body = []byte(err.Error())
	}
	return nil
}

func NewOutErrorParamSpec() sdkparam.OutputParamSpec[*Response] {
	return OutErrorParamSpec{}
}
