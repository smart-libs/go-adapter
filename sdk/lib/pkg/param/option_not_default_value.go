package sdkparam

import (
	"fmt"
)

// NotDefaultValue returns error if the given value is not of type T or its value is the default for the type T
func NotDefaultValue[T comparable]() Option {
	return func(spec Spec, value any) (any, error) {
		var defaultValueOfT T
		if instanceOfT, ok := value.(T); ok {
			if instanceOfT == defaultValueOfT {
				return nil, fmt.Errorf("%s=[%v] and it cannot be the default value=[%v]", spec.Name(), value, defaultValueOfT)
			}
		} else {
			return nil, fmt.Errorf("%s=[%v] is not [%T], it is [%T]", spec.Name(), value, defaultValueOfT, value)
		}

		return value, nil
	}
}

// NotDefaultValueReflection returns error if the given value is the default value for its type
func NotDefaultValueReflection() Option {
	return func(spec Spec, value any) (any, error) {
		if value == nil || getAsValueOf(value).IsZero() {
			return nil, fmt.Errorf("%s=[%v] and it cannot be the default value for [%T]", spec.Name(), value, value)
		}
		return value, nil
	}
}
