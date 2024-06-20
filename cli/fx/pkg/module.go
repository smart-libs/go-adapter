package cliadptfx

import (
	cliadpt "github.com/smart-libs/go-adapter/cli/lib/pkg"
	"go.uber.org/fx"
)

var MultiFlagSetModule = fx.Module("go-adapter/cli/lib",
	fx.Provide(
		cliadpt.NewMultiFlagSetAdapter,
	),
)

var SingleFlagSetModule = fx.Module("go-adapter/cli/lib",
	fx.Provide(
		cliadpt.NewSingleFlagSetAdapter,
	),
)
