package condition

import (
	cliadpt "github.com/smart-libs/go-adapter/cli/lib/pkg"
)

func HasFlag(flagName string) func(cliadpt.Input) bool {
	return func(input cliadpt.Input) bool {
		if input.FlagSet != nil {
			_, found := input.FlagSet.GetValue(flagName)
			return found
		}
		return false
	}
}
