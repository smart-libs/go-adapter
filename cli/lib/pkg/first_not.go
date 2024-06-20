package cliadpt

import "reflect"

func firstNotNil[T any](vs ...T) T {
	var defaultValue T
	for _, v := range vs {
		if !reflect.ValueOf(v).IsZero() {
			return v
		}
	}
	return defaultValue
}

func firstNotEmpty[T any](vs ...[]T) []T {
	for _, v := range vs {
		if len(v) > 0 {
			return v
		}
	}
	return nil
}
