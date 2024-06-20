package tagbasedhandler

import (
	"context"
	"fmt"
	"github.com/smart-libs/go-adapter/sdk/lib/pkg/param/tagbased"
	"reflect"
)

var (
	contextType = reflect.TypeOf(new(context.Context)).Elem()
)

func createArgFactoriesForFunction(funcType reflect.Type, factory tagbased.AbstractInputSpecBuilder) []argFactoryFunc {
	numOfInputArgs := funcType.NumIn()
	var argFactories []argFactoryFunc

	for i := 0; i < numOfInputArgs; i++ {
		argType := funcType.In(i)
		if i == 0 {
			if argType.AssignableTo(contextType) {
				argFactories = append(argFactories, createContextArg)
				continue
			}
		}

		argStructType := assertIsStruct(argType)
		argFactories = append(argFactories, createArgFactoryForStructure(argStructType, factory))
	}
	return argFactories
}

func assertIsStruct(argType reflect.Type) reflect.Type {
	if argType == nil || argType.Kind() != reflect.Struct {
		panic(fmt.Errorf("handler argument must be a structure, not=[%s]", argType))
	}
	return argType
}
