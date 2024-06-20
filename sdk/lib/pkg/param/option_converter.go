package sdkparam

import (
	"fmt"
	converter "github.com/smart-libs/go-crosscutting/converter/lib/pkg"
)

func IfNotNil(options ...Option) Option {
	return func(spec Spec, value any) (any, error) {
		if value == nil {
			return nil, nil
		}
		return AsSingleOptions(options...)(spec, value)
	}
}

func ConverterTyped[F any, T any](fromToFunc converter.FromToTypedFunc[F, T]) Option {
	converterFunc := converter.MakeFromToFunc(fromToFunc)
	return func(spec Spec, v any) (any, error) {
		var target T
		err := converterFunc(v, &target)
		if err != nil {
			return nil, fmt.Errorf("%s: Converter option failed: %w", spec.Name(), err)
		}
		return target, nil
	}
}

func ConvertTo[T any](registry converter.Converters) Option {
	return func(spec Spec, v any) (any, error) {
		return converter.To[T](registry, v)
	}
}

//func ConvertTo[F any, T any](toTypedFunc func(F) (T, error)) Option {
//	converterFunc := converter.MakeToFunc(toTypedFunc)
//	return func(spec Spec, v any) (any, error) {
//		target, err := converterFunc(v)
//		if err != nil {
//			return nil, fmt.Errorf("%s: Converter option failed: %w", spec.Name(), err)
//		}
//		return target, nil
//	}
//}
