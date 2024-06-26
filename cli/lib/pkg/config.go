package cliadpt

type (
	Config struct {
		// OsArgsUseDisabled enable or disable the use of os.Args when no args are provided to the SingleFlagSetAdapter.Run method
		OsArgsUseDisabled bool
		// ArgumentsValidationDisabled enable or disable the check to find any flag in the non-flags Array
		ArgumentsValidationDisabled bool
		// Bindings are rules used to identify which function/use case to be invoked based on the FlagSet and the Args
		Bindings
		// EnvGetter is the function used to get environment variables
		EnvGetter
		// FlagSet is the set of arguments accepted by the CLI adapter
		FlagSet
	}

	ConfigPerUseCaseName map[UseCaseName]Config

	UseCaseName = string
)

func (c ConfigPerUseCaseName) find(useCaseName UseCaseName) *Config {
	if config, found := c[useCaseName]; found {
		return &config
	}
	return nil
}
