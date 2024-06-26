package tagbased

import (
	"fmt"
	"reflect"
)

type (
	// ErrNoInputParamSpecCreatedForField means the InputParamSpec factories did not find a way to create a InputParamSpec
	// for the given Field. This error was originally created because the field did not have the tag flag set. The
	// error was created to allow a fallback when the field is itself a structure that can have inner fields with the
	// tag flag set.
	ErrNoInputParamSpecCreatedForField struct {
		Field reflect.StructField
	}
)

func (e ErrNoInputParamSpecCreatedForField) Error() string {
	return fmt.Sprintf("no sdkparam.InputParamSpec[Input] instance created for Field=[%s] using tags=[%v]",
		e.Field.Name, e.Field.Tag)
}
