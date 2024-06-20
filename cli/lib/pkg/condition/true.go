package condition

import (
	cliadpt "github.com/smart-libs/go-adapter/cli/lib/pkg"
)

func True() func(cliadpt.Input) bool {
	return func(input cliadpt.Input) bool { return true }
}
