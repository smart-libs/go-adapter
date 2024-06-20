package sdkparam

import (
	"fmt"
	"strings"
)

func stringAssertion(checks ...func(Spec, string) (any, error)) Option {
	return func(spec Spec, value any) (any, error) {
		check := func(spec Spec, str string) (any, error) {
			for _, checkString := range checks {
				if _, err := checkString(spec, str); err != nil {
					return nil, err
				}
			}
			return str, nil
		}

		switch str := value.(type) {
		case string:
			return check(spec, str)
		case *string:
			if str != nil {
				return check(spec, *str)
			}
			return nil, fmt.Errorf("%s=[nil] and it cannot be blank", spec.Name())
		}
		return nil, fmt.Errorf("%s=[%v] is not string, it is [%T]", spec.Name(), value, value)
	}
}

func notBlankAssertion(spec Spec, str string) (any, error) {
	if len(str) > 0 && strings.TrimSpace(str) == "" {
		return nil, fmt.Errorf("%s=[%v] and it cannot be blank", spec.Name(), str)
	}
	return str, nil
}

func notEmptyAssertion(spec Spec, str string) (any, error) {
	if str == "" {
		return nil, fmt.Errorf("%s=[%v] and it cannot be empty", spec.Name(), str)
	}
	return str, nil
}

func NotBlankString() Option        { return stringAssertion(notBlankAssertion) }
func NotEmptyString() Option        { return stringAssertion(notEmptyAssertion) }
func NotBlankOrEmptyString() Option { return stringAssertion(notBlankAssertion, notEmptyAssertion) }
