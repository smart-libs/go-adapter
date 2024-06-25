package sdk

import (
	"fmt"
	"strings"
)

// functions to help you to understand what is going on inside adapter SDK
type (
	DumpVar struct {
		Name   string
		Value  any
		Format string // Printf format
	}
	DumpTypedVar[T any] struct {
		Name   string
		Value  T
		Format string // Printf format
	}
)

var (
	DebugEnabled = false // change to true to see the debug info

	Debug func(string, ...any) = toStdout
)

func toStdout(format string, args ...any) {
	if DebugEnabled {
		fmt.Printf(format, args...)
	}
}

func (t DumpTypedVar[T]) toDumpVar() DumpVar {
	return DumpVar{
		Name:   t.Name,
		Value:  t.Value,
		Format: t.Format,
	}
}

func DebugDumpAndReturn[T any](ctxName string, v DumpTypedVar[T]) T {
	DebugDump(ctxName, v.toDumpVar())
	return v.Value
}

func DebugDumpAndReturn2[R1 any, R2 any](ctxName string, v1 DumpTypedVar[R1], v2 DumpTypedVar[R2]) (R1, R2) {
	DebugDump(ctxName, v1.toDumpVar(), v2.toDumpVar())
	return v1.Value, v2.Value
}

func DebugDump(ctxName string, vars ...DumpVar) {
	var (
		args   []any
		format strings.Builder
	)
	for _, dumpVar := range vars {
		format.WriteString(ctxName)
		format.WriteString(": ")
		format.WriteString(dumpVar.Name)
		format.WriteString("=[")
		if dumpVar.Format != "" {
			format.WriteString(dumpVar.Format)
		} else {
			format.WriteString("%v")
		}
		format.WriteString("]\n")
		args = append(args, dumpVar.Value)
	}
	Debug(format.String(), args...)
}
