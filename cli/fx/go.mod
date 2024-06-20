module github.com/smart-libs/go-adapter/cli/fx

go 1.22.1

require (
	github.com/smart-libs/go-adapter/cli/lib v0.0.0-00010101000000-000000000000
	go.uber.org/fx v1.22.0
)

require (
	github.com/smart-libs/go-adapter/interfaces v0.0.1 // indirect
	github.com/smart-libs/go-adapter/sdk/lib v0.0.1 // indirect
)

require (
	github.com/smart-libs/go-crosscutting/converter/lib v0.0.1 // indirect
	go.uber.org/dig v1.17.1 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.26.0 // indirect
	golang.org/x/sys v0.0.0-20220412211240-33da011f77ad // indirect
)

replace (
	github.com/smart-libs/go-adapter/cli/lib => ../../cli/lib
	github.com/smart-libs/go-adapter/interfaces => ../../interfaces
	github.com/smart-libs/go-adapter/sdk/lib => ../../sdk/lib
	github.com/smart-libs/go-crosscutting/converter/lib => ../../../go-crosscutting/converter/lib
)
