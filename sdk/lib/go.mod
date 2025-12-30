module github.com/smart-libs/go-adapter/sdk/lib

go 1.22.1

require (
	github.com/smart-libs/go-adapter/interfaces v0.0.1
	github.com/smart-libs/go-crosscutting/assertions/lib v0.0.6
	github.com/smart-libs/go-crosscutting/converter/lib v0.0.2
)

replace github.com/smart-libs/go-adapter/interfaces => ../../interfaces
