package tagbasedhandler

import "reflect"

func IsType[T any](v any) bool {
	if _, ok := v.(reflect.Type); ok {
		return v.(reflect.Type) == reflect.TypeFor[T]()
	}

	_, ok := v.(T)
	return ok
}
