package cliadpt

import sdkhandler "github.com/smart-libs/go-adapter/sdk/lib/pkg/handler"

type (
	Condition = func(Input) bool
	Handler   = sdkhandler.Handler[Input, *Output]

	Binding interface {
		EvaluateCondition(Input) bool
		Handler
	}

	Bindings []Binding

	baseBinding struct {
		Condition func(Input) bool
		Handler
	}
)

func NewBinding(condition Condition, handler Handler) Binding {
	return baseBinding{Condition: condition, Handler: handler}
}

func (b baseBinding) EvaluateCondition(input Input) bool { return b.Condition(input) }
