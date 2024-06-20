package sdkparam

import "reflect"

func in[T comparable](value T, listElem1 T, otherElems ...T) bool {
	if value == listElem1 {
		return true
	}
	for _, elem := range otherElems {
		if value == elem {
			return true
		}
	}
	return false
}

var ValueOfCreatedWithNil = reflect.Value{}

func isValueOfNil(valueOf reflect.Value) bool {
	return valueOf == ValueOfCreatedWithNil ||
		(in(valueOf.Kind(), reflect.Ptr, reflect.Map, reflect.Array, reflect.Slice) && valueOf.IsNil())
}

func getAsValueOf(v any) reflect.Value {
	if asValueOf, ok := v.(reflect.Value); ok {
		return asValueOf
	}
	return reflect.ValueOf(v)
}
