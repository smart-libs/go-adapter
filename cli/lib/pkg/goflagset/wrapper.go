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

	dump = cliadpt.DumpVar
)

const (
	debugWrapper = "goflagset.wrapper"
)

func (f *wrapper) Usage() {
	_, _ = fmt.Fprintf(f.Output(), "Usage of %s:\n", os.Args[0])
	f.PrintDefaults()
}

// Wrap a copy of the flags to avoid side effect
func Wrap(fs flag.FlagSet) cliadpt.FlagSet {
	return &wrapper{&fs}
}

func (f *wrapper) GetValue(flagName string) (value any, foundVar bool) {
	if cliadpt.DebugEnabled {
		cliadpt.DebugDump(debugWrapper,
			dump{Name: "flagName", Value: flagName},
			dump{Name: "f.FlagSet", Value: f.FlagSet})
		defer func() {
			cliadpt.DebugDump(debugWrapper,
				dump{Name: "value", Value: value},
				dump{Name: "foundVar", Value: foundVar})
		}()
	}
	var flagFound *flag.Flag
	f.Visit(func(flagSet *flag.Flag) {
		if flagFound == nil && flagSet.Name == flagName {
			flagFound = flagSet
		}
	})
	foundVar = flagFound != nil
	if !foundVar {
		if flagFound = f.Lookup(flagName); flagFound != nil {
			if flagFound.DefValue != "" {
				foundVar = true
				value = flagFound.DefValue
			}
		}
	} else {
		value = flagFound.Value.String()
	}
	return
}

var _ cliadpt.FlagSet = &wrapper{}
