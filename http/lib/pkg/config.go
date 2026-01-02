package httpadpt

type (
	Config struct {
		// Bindings are rules used to identify which function/use case to be invoked based on the FlagSet and the Args
		Bindings

		// Middlewares that shall be created by the implementation
		Middlewares
		Host *string
		Port *int

		// Other is used to provide additional implementation specific configuration
		Other any
	}
)
