package sdkparam

import sdk "github.com/smart-libs/go-adapter/sdk/lib/pkg"

type (
	Option func(spec Spec, value any) (any, error)
)

func AsSingleOptions(options ...Option) Option {
	return func(spec Spec, inputValue any) (value any, err error) {
		if sdk.DebugEnabled {
			defer func() {
				sdk.DebugDump("sdkparam.AsSingleOptions",
					sdk.DumpVar{Name: "spec", Value: spec},
					sdk.DumpVar{Name: "inputValue", Value: inputValue},
					sdk.DumpVar{Name: "value", Value: value},
					sdk.DumpVar{Name: "err", Value: err})
			}()
		}
		value = inputValue
		for _, option := range options {
			value, err = option(spec, value)
			if err != nil {
				break
			}
		}
		return
	}
}
