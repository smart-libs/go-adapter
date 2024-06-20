package sdkparam

type (
	// OutputParamSpec represents a parameter of the Output object
	OutputParamSpec[Output any] interface {
		// Spec is the parameter specification
		Spec
		// SetValue calls the specialized code to set a value into the given Output instance
		SetValue(output Output, value any) error
	}

	defaultOutputParam[Output any] struct {
		Spec
		setValueFunc func(output Output, value any) error
	}
)

func NewOutputParamSpec[Output any](spec Spec, setValueFunc func(output Output, value any) error) OutputParamSpec[Output] {
	return defaultOutputParam[Output]{Spec: spec, setValueFunc: setValueFunc}
}

func (o defaultOutputParam[Output]) SetValue(output Output, value any) error {
	applyOptions := AsSingleOptions(o.Options()...)
	finalValue, err := applyOptions(o, value)
	if err != nil {
		return err
	}
	return o.setValueFunc(output, getAsValueOf(finalValue).Interface())
}

var (
	_ OutputParamSpec[any] = defaultOutputParam[any]{}
)
