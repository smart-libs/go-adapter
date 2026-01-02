package httpadpt

import (
	"github.com/smart-libs/go-crosscutting/assertions/lib/pkg"
	"net/url"
)

type (
	QueryParams interface {
		// GetValue if the param name was sent, or it has a default value, then it returns the param value and true,
		// otherwise it returns nil and false.
		GetValue(queryParamName string) ([]string, bool)
	}

	HeaderParams interface {
		// GetValue if the param name was sent, or it has a default value, then it returns the param value and true,
		// otherwise it returns nil and false.
		GetValue(headerName string) ([]string, bool)
	}

	PathParams interface {
		// GetValue if the param name was sent, or it has a default value, then it returns the param value and true,
		// otherwise it returns nil and false.
		GetValue(pathParamName string) (string, bool)
	}

	Request interface {
		Query() QueryParams
		Header() HeaderParams
		Path() PathParams
		URL() *url.URL
		Method() string
	}
)

// IsRequestNil ensure the Request is not nil
func IsRequestNil(req Request, errHolder *error) bool {
	return HandleErrorHolder(errHolder, assertions.AnyIsNotNil(req))
}

// IsRequestQueryNil ensures that both the Request and Request.Query() are not nil
func IsRequestQueryNil(req Request, errHolder *error) bool {
	if !IsRequestNil(req, errHolder) {
		return false
	}
	return HandleErrorHolder(errHolder, assertions.AnyIsNotNil(req.Query()))
}

// IsRequestHeaderNil ensures that both the Request and Request.Header() are not nil
func IsRequestHeaderNil(req Request, errHolder *error) bool {
	if !IsRequestNil(req, errHolder) {
		return false
	}
	return HandleErrorHolder(errHolder, assertions.AnyIsNotNil(req.Header()))
}

// IsRequestPathNil ensures that both the Request and Request.Path() are not nil
func IsRequestPathNil(req Request, errHolder *error) bool {
	if !IsRequestNil(req, errHolder) {
		return false
	}
	return HandleErrorHolder(errHolder, assertions.AnyIsNotNil(req.Path()))
}
