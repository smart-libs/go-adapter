package cliadpt

type FlagSet interface {
	// Parsed returns true if the flag set was already parsed
	Parsed() bool
	// Parse the given arguments and sets the flags. arg[0] must not be the binary invoked
	Parse(arguments []string) error
	// GetValue if the flag name was specified or it has a default value, then it returns the flag value and true,
	// otherwise it returns nil and false.
	GetValue(flagName string) (any, bool)
	// Args are the list or values passed as argument that are not associated with any flag.
	Args() []string
	// Usage print the usage message to stderr
	Usage()
}
