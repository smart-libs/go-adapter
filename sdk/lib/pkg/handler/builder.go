package sdkhandler

import (
	"github.com/smart-libs/go-adapter/interfaces/pkg/adapter"
	sdkparam "github.com/smart-libs/go-adapter/sdk/lib/pkg/param"
	"github.com/smart-libs/go-adapter/sdk/lib/pkg/param/tagbased"
)

type (
	OutputSpecStep[Input any, Output any] interface {
		WithOutParamSpec(ref adapter.ParamRef, spec sdkparam.OutputParamSpec[Output]) OutputSpecStep[Input, Output]
		WithErrorParamSpec(spec sdkparam.OutputParamSpec[Output]) OutputSpecStep[Input, Output]
		Build() Handler[Input, Output]
	}

	InputSpecStep[Input any, Output any] interface {
		WithInParamSpec(ref adapter.ParamRef, actionFunc sdkparam.InputParamSpec[Input]) InputSpecStep[Input, Output]
		Output() OutputSpecStep[Input, Output]
	}

	TagBasedWithErrorBuildStep[Input any, Output any] interface {
		WithOutErrorParamSpec(spec sdkparam.OutputParamSpec[Output]) TagBasedBuildStep[Input, Output]
		TagBasedBuildStep[Input, Output]
	}

	TagBasedBuildStep[Input any, Output any] interface {
		Build() Handler[Input, Output]
	}

	TagBasedOutputSpecStep[Input any, Output any] interface {
		WithOutTagBasedFactory(factory tagbased.OutputParamSpecFactory[Output]) TagBasedWithErrorBuildStep[Input, Output]
	}

	TagBasedInputSpecStep[Input any, Output any] interface {
		WithInTagBasedFactory(factory tagbased.InputParamSpecFactory[Input]) TagBasedOutputSpecStep[Input, Output]
	}
)
