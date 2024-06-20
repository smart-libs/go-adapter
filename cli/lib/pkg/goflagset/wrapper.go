package goflagset

import (
	"flag"
	"fmt"
	cliadpt "github.com/smart-libs/go-adapter/cli/lib/pkg"
	"os"
)

type (
	wrapper struct {
		*flag.FlagSet
	}
)

func (f *wrapper) Usage() {
	_, _ = fmt.Fprintf(f.Output(), "Usage of %s:\n", os.Args[0])
	f.PrintDefaults()
}

// Wrap a copy of the flags to avoid side effect
func Wrap(fs flag.FlagSet) cliadpt.FlagSet {
	return &wrapper{&fs}
}

func (f *wrapper) GetValue(flagName string) (any, bool) {
	var found *flag.Flag
	f.Visit(func(flagSet *flag.Flag) {
		if found == nil && flagSet.Name == flagName {
			found = flagSet
		}
	})
	if found == nil {
		if found = f.Lookup(flagName); found != nil {
			if found.DefValue != "" {
				return getFlagInParamValueDebug(f.FlagSet, flagName, found.DefValue, true)
			}
		}
		return getFlagInParamValueDebug(f.FlagSet, flagName, nil, false)
	}
	return getFlagInParamValueDebug(f.FlagSet, flagName, found.Value.String(), true)
}

func getFlagInParamValueDebug(flagSet *flag.FlagSet, flagName string, value any, found bool) (any, bool) {
	if cliadpt.DebugEnabled {
		fmt.Printf("goflagset.wrapper: FlagSet[%v], flagName=[%s], return (value=[%v], found=[%v])\n", flagSet, flagName, value, found)
	}
	return value, found
}

var _ cliadpt.FlagSet = &wrapper{}
