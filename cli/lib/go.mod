module github.com/smart-libs/go-adapter/cli/lib

go 1.22.1

require (
	github.com/smart-libs/go-adapter/interfaces v0.0.1
	github.com/smart-libs/go-adapter/sdk/lib v0.0.1
	github.com/smart-libs/go-crosscutting/converter/lib v0.0.1
	github.com/stretchr/testify v1.9.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/smart-libs/go-adapter/interfaces => ../../interfaces
	github.com/smart-libs/go-adapter/sdk/lib => ../../sdk/lib
	github.com/smart-libs/go-crosscutting/converter/lib => ../../../go-crosscutting/converter/lib
)
