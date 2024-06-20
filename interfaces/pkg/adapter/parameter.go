package adapter

type (
	// ParamRef is a reference to a parameter whose mapping should be known InputAccessor and/or OutputBuilder.
	ParamRef interface {
		GetID() string
	}

	StringParamRef string
)

// GetID identifies the parameter the UseCaseHandler needs from the input provided by the adapter.
func (s StringParamRef) GetID() string { return string(s) }
