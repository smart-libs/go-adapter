package httpadpt

import (
	"encoding/json"
	"reflect"

	"github.com/smart-libs/go-crosscutting/assertions/lib/pkg/check"
)

type (
	//ProblemDetail to inform error
	ProblemDetail struct {
		//These are specified for https://datatracker.ietf.org/doc/html/rfc7807
		Type           string                 `json:"type,omitempty"`
		Instance       string                 `json:"instance,omitempty"`
		AdditionalInfo map[string]interface{} `json:"additional_info,omitempty"`
		//These fields are specific for https://jsonapi.org/format/#errors
		Status string `json:"status,omitempty"`
		Code   string `json:"code,omitempty"`
		ID     string `json:"id,omitempty"`
		//These fields are common to RFC and JSON API
		Title  string `json:"title,omitempty"`
		Detail string `json:"detail,omitempty"`
	}
)

var (
	TypeBuilder  = defaultTypeBuilder
	TitleBuilder = defaultTitleBuilder
)

func defaultTypeBuilder(err error) string {
	if check.IsNil(err) {
		return "nil"
	}
	return reflect.TypeOf(err).String()
}

func defaultTitleBuilder(err error) string {
	if check.IsNil(err) {
		return "nil"
	}
	return ""
}

func ProblemDetailFromError(err error) ProblemDetail {
	detail := ""
	if !check.IsNil(err) {
		detail = err.Error()
	}
	return ProblemDetail{
		Type:   TypeBuilder(err),
		Title:  TitleBuilder(err),
		Detail: detail,
	}
}

func JSONProblemDetailFromError(err error) ([]byte, error) {
	pd := ProblemDetailFromError(err)
	return json.Marshal(pd)
}
