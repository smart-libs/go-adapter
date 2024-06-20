package tagbasedhandler

import (
	"fmt"
	sdkhandler "github.com/smart-libs/go-adapter/sdk/lib/pkg/handler"
	sdkusecasehandler "github.com/smart-libs/go-adapter/sdk/lib/pkg/handler/usecase"
	sdkparam "github.com/smart-libs/go-adapter/sdk/lib/pkg/param"
	"github.com/smart-libs/go-adapter/sdk/lib/pkg/param/tagbased"
	"reflect"
)

type (
	builder[Input any, Output any] struct {
		targetFunc           reflect.Value
		argFactories         []argFactoryFunc
		outputBuilderActions []OutputBuilderActionFunc
		outputSpecs          sdkparam.OutputSpecs[Output]
		inputSpecs           sdkparam.InputSpecs[Input]
	}
)

// NewBuilderForFunc creates a handler based on the given function
func NewBuilderForFunc[Input any, Output any](handlerFunc any) sdkhandler.TagBasedInputSpecStep[Input, Output] {
	_ = assertIsFunc(handlerFunc)
	return &builder[Input, Output]{
		targetFunc: reflect.ValueOf(handlerFunc),
	}
}

func (b *builder[Input, Output]) Build() sdkhandler.Handler[Input, Output] {
	useCase := handler{
		targetFunc:           b.targetFunc,
		argFactories:         b.argFactories,
		outputBuilderActions: b.outputBuilderActions,
	}
	return sdkusecasehandler.NewUseCaseHandler[Input, Output](useCase, b.inputSpecs, b.outputSpecs)
}

func (b *builder[Input, Output]) WithOutTagBasedFactory(factory tagbased.OutputParamSpecFactory[Output]) sdkhandler.TagBasedBuildStep[Input, Output] {
	var err error
	outSpecBuilder := tagbased.NewOutputSpecsBuilder[Output](factory)
	outputBuilderActionsFactory := OutputBuilderActionsFactory{AbstractOutputSpecBuilder: outSpecBuilder}
	b.outputBuilderActions, err = outputBuilderActionsFactory.Create(b.targetFunc.Type())
	if err != nil {
		panic(err)
	}
	b.outputSpecs = outSpecBuilder.Build()
	return b
}

func (b *builder[Input, Output]) WithInTagBasedFactory(factory tagbased.InputParamSpecFactory[Input]) sdkhandler.TagBasedOutputSpecStep[Input, Output] {
	inSpecBuilder := tagbased.NewInputSpecsBuilder[Input](factory)
	b.argFactories = createArgFactoriesForFunction(b.targetFunc.Type(), inSpecBuilder)
	b.inputSpecs = inSpecBuilder.Build()
	return b
}

func assertIsFunc(input any) reflect.Type {
	inputType := reflect.TypeOf(input)
	if inputType == nil || inputType.Kind() != reflect.Func {
		panic(fmt.Errorf("handler must be a function, not=[%T]", input))
	}
	return inputType
}
