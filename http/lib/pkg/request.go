package httpadpt

import "github.com/smart-libs/go-crosscutting/assertions/lib/pkg"

type (
	QueryParams interface {
		// GetValue if the param name was specified, or it has a default value, then it returns the param value and true,
		// otherwise it returns nil and false.
		GetValue(flagName string) ([]string, bool)
	}

	Request interface {
		Query() QueryParams
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
