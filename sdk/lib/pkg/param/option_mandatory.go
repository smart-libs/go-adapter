package sdkparam

import "fmt"

func Mandatory(isNilList ...func(any) bool) Option {
	isNil := func(v any) bool {
		valueOf := getAsValueOf(v)
		return isValueOfNil(valueOf)
	}
	if len(isNilList) > 0 {
		isNil = isNilList[0]
	}
	return func(spec Spec, v any) (any, error) {
		if isNil(v) {
			return nil, fmt.Errorf("param=[%s] is mandatory, given value=[%v], type=[%T]", spec.Name(), v, v)
		}
		return v, nil
	}
}
