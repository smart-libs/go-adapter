package test

import (
	"context"
	"fmt"
	httpadpt "github.com/smart-libs/go-adapter/http/lib/pkg"
	serror "github.com/smart-libs/go-crosscutting/serror/lib/pkg"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type (
	testHandlerInput struct {
		Q1Value string `query:"q1"`
	}
	testHandlerOutput struct {
		ResultCode int `statuscode:""`
	}
)

func testHandlerType1(i testHandlerInput) (*testHandlerOutput, error) {
	if i.Q1Value != "" {
		if i.Q1Value == "error" {
			return nil, serror.IllegalConfigParamValue("q1", i.Q1Value)
		}
		return &testHandlerOutput{ResultCode: 201}, nil
	}
	return &testHandlerOutput{ResultCode: 204}, nil
}

func SuiteTest(t *testing.T, adapterFactory func(config httpadpt.Config) httpadpt.Adapter) {
	t.Run("Test Start/Stop", func(t *testing.T) {
		ctx := context.Background()
		adapter := adapterFactory(httpadpt.Config{
			Bindings: []httpadpt.Binding{
				httpadpt.NewBindingBuilderUsingPath("/v1/test").
					WithMethods(http.MethodGet).
					WithHandlerFunc(testHandlerType1),
			},
		})
		if assert.NoError(t, adapter.Start(ctx)) {
			assert.NoError(t, adapter.Stop(ctx))
		}
	})
	t.Run("Test GET /v1/test", func(t *testing.T) {
		port := 8080
		ctx := context.Background()
		adapter := adapterFactory(httpadpt.Config{
			Port: &port,
			Bindings: []httpadpt.Binding{
				httpadpt.NewBindingBuilderUsingPath("/v1/test").
					WithMethods(http.MethodGet).
					WithHandlerFunc(testHandlerType1),
			},
		})
		if assert.NoError(t, adapter.Start(ctx)) {
			resp, err := http.Get(fmt.Sprintf("http://localhost:%d/v1/test?q1=test", port))
			if assert.NoError(t, err) {
				if assert.Equal(t, 201, resp.StatusCode) {
					resp, err = http.Get(fmt.Sprintf("http://localhost:%d/v1/test", port))
					if assert.NoError(t, err) {
						assert.Equal(t, 204, resp.StatusCode)
					}
					resp, err = http.Get(fmt.Sprintf("http://localhost:%d/v1/test?q1=error", port))
					if assert.NoError(t, err) {
						assert.Equal(t, 400, resp.StatusCode)
					}
				}
			}
			assert.NoError(t, adapter.Stop(ctx))
		}
	})
}
