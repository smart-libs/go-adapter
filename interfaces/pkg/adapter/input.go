package adapter

type (
	InputAccessor interface {
		GetValue(ref ParamRef) (any, error)
		CopyValue(ref ParamRef, target any) error
	}
)

// ValueAs is a helper function to retrieve the specified parameter value or panic
func ValueAs[T any](accessor InputAccessor, ref ParamRef) T {
	var t T
	if err := accessor.CopyValue(ref, &t); err != nil {
		panic(err)
	}
	return t
}
