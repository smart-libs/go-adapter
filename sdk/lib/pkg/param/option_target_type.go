package sdkparam

import (
	converter "github.com/smart-libs/go-crosscutting/converter/lib/pkg"
	converterdefault "github.com/smart-libs/go-crosscutting/converter/lib/pkg/default"
	"reflect"
)

func GivenTargetType(targetType reflect.Type, converters ...converter.Converters) Option {
	convertTo := converterdefault.Converters.ConvertToType
	if len(converters) > 0 {
		convertTo = converters[0].ConvertToType
	}
	return func(spec Spec, value any) (any, error) {
		return convertTo(value, targetType)
	}
}

func TargetType[T any](converters ...converter.Converters) Option {
	var t T
	targetType := reflect.TypeOf(t)
	return GivenTargetType(targetType, converters...)
}
