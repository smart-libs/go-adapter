package condition

import (
	"fmt"
	cliadpt "github.com/smart-libs/go-adapter/cli/lib/pkg"
)

func FlagValueIs(flagName, value string) func(cliadpt.Input) bool {
	return func(input cliadpt.Input) bool {
		if input.FlagSet != nil {
			if flagValue, found := input.FlagSet.GetValue(flagName); found {
				if str, ok := flagValue.(string); ok {
					return str == value
				}
				if flagValue != nil {
					return fmt.Sprintf("%v", flagValue) == value
				}
			}
		}
		return false
	}
}
