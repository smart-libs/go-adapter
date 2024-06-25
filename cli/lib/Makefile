GO=go
TEST_FOLDERS=./pkg/... ./test/...

#run-tests: @ run the tests
run-tests:
	@echo "****[VET]**********************************************************"
	@$(GO) vet $(TEST_FOLDERS)
	@mkdir -p ./coverage
	@APP_ENV=dev $(GO) test -v $(TEST_FOLDERS) -coverprofile coverage/coverage.fmt fmt
	@APP_ENV=dev $(GO) tool cover -html=coverage/coverage.fmt -o coverage/coverage.html

#coverage: @ generate the coverage report
coverage:
	$(GO) tool cover -html=coverage/fmtcoverage.html

#lint-install: @ install the lint tool
lint-install:
	@echo "****[Installing lint]**********************************************"
	@$(GO) install honnef.co/go/tools/cmd/staticcheck@latest

#lint: @ run lint to check code style
lint:
	@echo "****[Running lint]*************************************************"
	@staticcheck ./...

#test: @ run the tests and the lint tool
test: run-tests lint
