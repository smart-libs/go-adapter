package httpadpt

import "github.com/smart-libs/go-crosscutting/assertions/lib/pkg/check"

// HandleErrorHolder checks if the given error is non-nil, assigns it to errHolder if provided, and returns true if non-nil.
func HandleErrorHolder(errHolder *error, err error) bool {
	if check.IsNil(err) {
		return false
	}
	if errHolder != nil {
		*errHolder = err
	}
	return true
}
