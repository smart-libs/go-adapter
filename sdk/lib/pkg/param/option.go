package sdkparam

type (
	Option func(spec Spec, value any) (any, error)
)

func AsSingleOptions(options ...Option) Option {
	return func(spec Spec, inputValue any) (value any, err error) {
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
