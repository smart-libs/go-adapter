package httpadpt

type (
	// Condition specifies the HTTP conditions to select the handler
	Condition struct {
		// Path can be a string, a string with path parameters surrounded by {}, or a regex
		Path *string
		// Methods identifies which HTTP verbs can be used with a given path
		Methods []string
		// Other can be used for implementation's conditions
		Other any
	}

	Binding struct {
		Condition
		Handler
	}

	// Bindings represents the bindings the HTTP handler should handle, the binding order in the list
	// can be used by the HTTP implementation to specify priority to evaluate the conditions to select the handler
	Bindings []Binding
)

func IsBindingValid(binding Binding) bool {
	return binding.Handler != nil
}
