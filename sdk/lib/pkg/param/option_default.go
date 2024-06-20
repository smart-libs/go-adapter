package sdkparam

func Default(defaultValue any, isNilList ...func(any) bool) Option {
	isNil := func(v any) bool {
		valueOf := getAsValueOf(v)
		return isValueOfNil(valueOf)
	}
	if len(isNilList) > 0 {
		isNil = isNilList[0]
	}
	return func(spec Spec, v any) (any, error) {
		if isNil(v) {
			return defaultValue, nil
		}
		return v, nil
	}
}
