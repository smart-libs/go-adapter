package gonethttp

import (
	httpadpt "github.com/smart-libs/go-adapter/http/lib/pkg"
	"github.com/smart-libs/go-adapter/http/lib/test"
	"github.com/smart-libs/go-crosscutting/assertions/lib/pkg/must"
	"testing"
)

func Test_Suite(t *testing.T) {
	test.SuiteTest(t, func(config httpadpt.Config) httpadpt.Adapter {
		return must.SucceedWith1(NewAdapter(config))
	})
}
