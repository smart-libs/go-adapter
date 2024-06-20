module github.com/smart-libs/go-adapter/sdk/lib

go 1.22.1

require (
	github.com/smart-libs/go-adapter/interfaces v0.0.1
	github.com/smart-libs/go-crosscutting/converter/lib v0.0.1
)

replace (
	github.com/smart-libs/go-adapter/interfaces => ../../interfaces
	github.com/smart-libs/go-crosscutting/converter/lib => ../../../go-crosscutting/converter/lib
)
