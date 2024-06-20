package cliadpt

import (
	"github.com/smart-libs/go-adapter/interfaces/pkg/adapter"
	sdkhandler "github.com/smart-libs/go-adapter/sdk/lib/pkg/handler"
	tagbasedhandler "github.com/smart-libs/go-adapter/sdk/lib/pkg/handler/tagbased"
	sdkusecasehandler "github.com/smart-libs/go-adapter/sdk/lib/pkg/handler/usecase"
	sdkparam "github.com/smart-libs/go-adapter/sdk/lib/pkg/param"
)

type (
	UseCaseHandlerToStep interface {
		ToPanic(options ...sdkparam.Option) UseCaseHandlerOutputStep
		ToStderr(printfMask string, options ...sdkparam.Option) UseCaseHandlerOutputStep
		ToStdout(printfMask string, options ...sdkparam.Option) UseCaseHandlerOutputStep
		ToExitCode(options ...sdkparam.Option) UseCaseHandlerOutputStep
	}

	UseCaseHandlerOutputStep interface {
		BuildStep
		WithOutParamRef(ref adapter.ParamRef) UseCaseHandlerToStep
		WithOutParamError() UseCaseHandlerToStep
	}

	UseCaseHandlerFromStep interface {
		FromFlag(flagName string, options ...sdkparam.Option) UseCaseHandlerInputStep
		FromAllNonFlags(options ...sdkparam.Option) UseCaseHandlerInputStep
		FromNonFlagAtPos(pos int, options ...sdkparam.Option) UseCaseHandlerInputStep
		FromArgAtPos(pos int, options ...sdkparam.Option) UseCaseHandlerInputStep
		FromEnv(envName string, options ...sdkparam.Option) UseCaseHandlerInputStep
	}

	UseCaseHandlerInputStep interface {
		WithInParamRef(ref adapter.ParamRef) UseCaseHandlerFromStep
		AndResult() UseCaseHandlerOutputStep
	}

	HandlerStep interface {
		InvokeHandler(handler any) BuildStep
		InvokeUseCase(useCase adapter.UseCaseHandler) UseCaseHandlerInputStep
	}

	BuildStep interface {
		Build() Binding
	}

	builder struct {
		Condition
	}

	funcBindingBuilder struct {
		Condition
		innerBuilder sdkhandler.TagBasedInputSpecStep[Input, *Output]
	}

	useCaseBindingBuilder struct {
		Condition
		inBuilder    sdkhandler.InputSpecStep[Input, *Output]
		outBuilder   sdkhandler.OutputSpecStep[Input, *Output]
		lastParamRef adapter.ParamRef
	}
)

func (u *useCaseBindingBuilder) ToPanic(options ...sdkparam.Option) UseCaseHandlerOutputStep {
	if u.lastParamRef == nil {
		u.outBuilder.WithErrorParamSpec(NewPanicErrorOutParamSpec(options...))
	}
	return u
}

func (u *useCaseBindingBuilder) ToStderr(printfMask string, options ...sdkparam.Option) UseCaseHandlerOutputStep {
	if u.lastParamRef == nil {
		u.outBuilder.WithErrorParamSpec(NewPrintStderrOutParamSpec(printfMask, options...))
	} else {
		u.outBuilder.WithOutParamSpec(u.lastParamRef, NewPrintStderrOutParamSpec(printfMask, options...))
	}
	return u
}

func (u *useCaseBindingBuilder) ToStdout(printfMask string, options ...sdkparam.Option) UseCaseHandlerOutputStep {
	if u.lastParamRef == nil {
		u.outBuilder.WithErrorParamSpec(NewPrintStdoutOutParamSpec(printfMask, options...))
	} else {
		u.outBuilder.WithOutParamSpec(u.lastParamRef, NewPrintStdoutOutParamSpec(printfMask, options...))
	}
	return u
}

func (u *useCaseBindingBuilder) ToExitCode(options ...sdkparam.Option) UseCaseHandlerOutputStep {
	if u.lastParamRef == nil {
		u.outBuilder.WithErrorParamSpec(NewExitCodeOutParamSpec(options...))
	} else {
		u.outBuilder.WithOutParamSpec(u.lastParamRef, NewExitCodeOutParamSpec(options...))
	}
	return u
}

func (u *useCaseBindingBuilder) Build() Binding {
	return NewBinding(u.Condition, u.outBuilder.Build())
}

func (u *useCaseBindingBuilder) WithOutParamRef(ref adapter.ParamRef) UseCaseHandlerToStep {
	u.lastParamRef = ref
	return u
}

func (u *useCaseBindingBuilder) WithOutParamError() UseCaseHandlerToStep {
	u.lastParamRef = nil
	return u
}

func (u *useCaseBindingBuilder) FromFlag(flagName string, options ...sdkparam.Option) UseCaseHandlerInputStep {
	u.inBuilder.WithInParamSpec(u.lastParamRef, newFlagInParam(flagName, options...))
	return u
}

func (u *useCaseBindingBuilder) FromAllNonFlags(options ...sdkparam.Option) UseCaseHandlerInputStep {
	u.inBuilder.WithInParamSpec(u.lastParamRef, newAllNonFlagInParam(options...))
	return u
}

func (u *useCaseBindingBuilder) FromNonFlagAtPos(pos int, options ...sdkparam.Option) UseCaseHandlerInputStep {
	u.inBuilder.WithInParamSpec(u.lastParamRef, newPosNonFlagInParam(pos, options))
	return u
}

func (u *useCaseBindingBuilder) FromArgAtPos(pos int, options ...sdkparam.Option) UseCaseHandlerInputStep {
	u.inBuilder.WithInParamSpec(u.lastParamRef, newPosInParam(pos, options...))
	return u
}

func (u *useCaseBindingBuilder) FromEnv(envName string, options ...sdkparam.Option) UseCaseHandlerInputStep {
	u.inBuilder.WithInParamSpec(u.lastParamRef, newEnvInParam(envName, options...))
	return u
}

func (u *useCaseBindingBuilder) WithInParamRef(ref adapter.ParamRef) UseCaseHandlerFromStep {
	u.lastParamRef = ref
	return u
}

func (u *useCaseBindingBuilder) AndResult() UseCaseHandlerOutputStep {
	u.outBuilder = u.inBuilder.Output()
	return u
}

func (f funcBindingBuilder) Build() Binding {
	return NewBinding(
		f.Condition,
		f.innerBuilder.
			WithInTagBasedFactory(createInParamSpecFactory()).
			WithOutTagBasedFactory(createOutParamSpecFactory()).
			Build(),
	)
}

func (b builder) InvokeHandler(handler any) BuildStep {
	return funcBindingBuilder{
		Condition:    b.Condition,
		innerBuilder: tagbasedhandler.NewBuilderForFunc[Input, *Output](handler),
	}
}

func (b builder) InvokeUseCase(useCase adapter.UseCaseHandler) UseCaseHandlerInputStep {
	return &useCaseBindingBuilder{
		Condition: b.Condition,
		inBuilder: sdkusecasehandler.NewBuilderForUseCase[Input, *Output](useCase),
	}
}

func NewBindingBuilderWithCondition(cond Condition) HandlerStep {
	return builder{Condition: cond}
}
