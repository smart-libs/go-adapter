package test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	httpadpt "github.com/smart-libs/go-adapter/http/lib/pkg"
	serror "github.com/smart-libs/go-crosscutting/serror/lib/pkg"
	"github.com/stretchr/testify/assert"
	"io"
	"log/slog"
	"net/http"
	"testing"
)

type (
	testHandlerInput struct {
		Q1Value string `query:"q1"`
		H1Value string `header:"h1"`
		P1Value string `path:"p1"`
	}
	testHandlerOutput struct {
		ResultCode  int    `statuscode:""`
		ContentType string `header:"Content-Type"`
		Body        string `body:""`
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

func testHandlerType2(i testHandlerInput) (*testHandlerOutput, error) {
	if i.H1Value != "" {
		if i.H1Value == "error" {
			return nil, serror.IllegalConfigParamValue("h1", i.Q1Value)
		}
		return &testHandlerOutput{ResultCode: 201}, nil
	}
	return &testHandlerOutput{ResultCode: 204}, nil
}

func testHandlerType3(i testHandlerInput) (*testHandlerOutput, error) {
	if i.P1Value != "10" {
		if i.P1Value == "error" {
			return nil, serror.IllegalConfigParamValue("p1", i.Q1Value)
		}
		return &testHandlerOutput{ResultCode: 201}, nil
	}
	return &testHandlerOutput{ResultCode: 204}, nil
}

func testHandlerType4() (*testHandlerOutput, error) {
	return &testHandlerOutput{
		ResultCode:  200,
		ContentType: "application/json",
		Body:        `{"test":"test"}`,
	}, nil
}

func testHandlerForPanicWithAnyArg() (*testHandlerOutput, error) {
	panic("test panic")
}

func testHandlerForPanicWithError() (*testHandlerOutput, error) {
	panic(errors.New("test panic with error"))
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
	t.Run("Test GET /v1/test using query param", func(t *testing.T) {
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
	t.Run("Test GET /v1/test using header param", func(t *testing.T) {
		port := 8080
		ctx := context.Background()
		adapter := adapterFactory(httpadpt.Config{
			Port: &port,
			Bindings: []httpadpt.Binding{
				httpadpt.NewBindingBuilderUsingPath("/v1/test").
					WithMethods(http.MethodGet).
					WithHandlerFunc(testHandlerType2),
			},
		})
		if assert.NoError(t, adapter.Start(ctx)) {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%d/v1/test", port), nil)
			if assert.NoError(t, err) {
				// NO HEADER SET
				resp, respErr := http.DefaultClient.Do(req)
				if assert.NoError(t, respErr) {
					if assert.Equal(t, 204, resp.StatusCode) {
						// HEADER SET
						req.Header["h1"] = []string{"test"}
						resp, respErr = http.DefaultClient.Do(req)
						if assert.NoError(t, respErr) {
							assert.Equal(t, 201, resp.StatusCode)
						}

						req.Header["h1"] = []string{"error"}
						resp, respErr = http.DefaultClient.Do(req)
						if assert.NoError(t, respErr) {
							assert.Equal(t, 400, resp.StatusCode)
						}
					}
				}
			}
			assert.NoError(t, adapter.Stop(ctx))
		}
	})
	t.Run("Test GET /v1/test using path param", func(t *testing.T) {
		port := 8080
		ctx := context.Background()
		adapter := adapterFactory(httpadpt.Config{
			Port: &port,
			Bindings: []httpadpt.Binding{
				httpadpt.NewBindingBuilderUsingPath("/v1/test/{p1}").
					WithMethods(http.MethodGet).
					WithHandlerFunc(testHandlerType3),
			},
		})
		if assert.NoError(t, adapter.Start(ctx)) {
			// path = 10
			resp, respErr := http.Get(fmt.Sprintf("http://localhost:%d/v1/test/10", port))
			if assert.NoError(t, respErr) {
				if assert.Equal(t, 204, resp.StatusCode) {
					// path != 10
					resp, respErr = http.Get(fmt.Sprintf("http://localhost:%d/v1/test/11", port))
					if assert.NoError(t, respErr) {
						assert.Equal(t, 201, resp.StatusCode)
					}

					resp, respErr = http.Get(fmt.Sprintf("http://localhost:%d/v1/test/error", port))
					if assert.NoError(t, respErr) {
						assert.Equal(t, 400, resp.StatusCode)
					}
				}
			}
			assert.NoError(t, adapter.Stop(ctx))
		}
	})
	t.Run("Test GET /v1/test to return body", func(t *testing.T) {
		port := 8080
		ctx := context.Background()
		adapter := adapterFactory(httpadpt.Config{
			Port: &port,
			Bindings: []httpadpt.Binding{
				httpadpt.NewBindingBuilderUsingPath("/v1/test").
					WithMethods(http.MethodGet).
					WithHandlerFunc(testHandlerType4),
			},
		})
		if assert.NoError(t, adapter.Start(ctx)) {
			// path = 10
			resp, respErr := http.Get(fmt.Sprintf("http://localhost:%d/v1/test", port))
			if assert.NoError(t, respErr) {
				if assert.Equal(t, 200, resp.StatusCode) {
					assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
					all, err := io.ReadAll(resp.Body)
					if assert.NoError(t, err) {
						assert.Equal(t, `{"test":"test"}`, string(all))
					}
				}
			}
			assert.NoError(t, adapter.Stop(ctx))
		}
	})
	t.Run("Test GET /v1/test that panic with any", func(t *testing.T) {
		port := 8080
		ctx := context.Background()
		adapter := adapterFactory(httpadpt.Config{
			Port: &port,
			Bindings: []httpadpt.Binding{
				httpadpt.NewBindingBuilderUsingPath("/v1/test").
					WithMethods(http.MethodGet).
					WithHandlerFunc(testHandlerForPanicWithAnyArg),
			},
			Middlewares: httpadpt.Middlewares{
				httpadpt.HandlePanic,
			},
		})
		if assert.NoError(t, adapter.Start(ctx)) {
			// path = 10
			resp, respErr := http.Get(fmt.Sprintf("http://localhost:%d/v1/test", port))
			if assert.NoError(t, respErr) {
				if assert.Equal(t, 500, resp.StatusCode) {
					assert.Equal(t, httpadpt.ContentTypeProblemDetail, resp.Header.Get("Content-Type"))
					all, err := io.ReadAll(resp.Body)
					if assert.NoError(t, err) {
						assert.Equal(t, `{"type":"*errorx.Error","detail":"common.internal_error: test panic"}`, string(all))
					}
				}
			}
			assert.NoError(t, adapter.Stop(ctx))
		}
	})
	t.Run("Test GET /v1/test that panic with error", func(t *testing.T) {
		port := 8080
		ctx := context.Background()
		adapter := adapterFactory(httpadpt.Config{
			Port: &port,
			Bindings: []httpadpt.Binding{
				httpadpt.NewBindingBuilderUsingPath("/v1/test").
					WithMethods(http.MethodGet).
					WithHandlerFunc(testHandlerForPanicWithError),
			},
			Middlewares: httpadpt.Middlewares{
				httpadpt.HandlePanic,
			},
		})
		if assert.NoError(t, adapter.Start(ctx)) {
			// path = 10
			resp, respErr := http.Get(fmt.Sprintf("http://localhost:%d/v1/test", port))
			if assert.NoError(t, respErr) {
				if assert.Equal(t, 500, resp.StatusCode) {
					assert.Equal(t, httpadpt.ContentTypeProblemDetail, resp.Header.Get("Content-Type"))
					all, err := io.ReadAll(resp.Body)
					if assert.NoError(t, err) {
						assert.Equal(t, `{"type":"*errors.errorString","detail":"test panic with error"}`, string(all))
					}
				}
			}
			assert.NoError(t, adapter.Stop(ctx))
		}
	})
	t.Run("Test GET /v1/test that logs HTTP request", func(t *testing.T) {
		var buf bytes.Buffer
		logger := slog.New(slog.NewJSONHandler(&buf, nil))
		ctx := context.Background()
		port := 8080
		adapter := adapterFactory(httpadpt.Config{
			Port: &port,
			Bindings: []httpadpt.Binding{
				httpadpt.NewBindingBuilderUsingPath("/v1/test").
					WithMethods(http.MethodGet).
					WithHandlerFunc(testHandlerType1),
			},
			Middlewares: httpadpt.Middlewares{
				httpadpt.HandlePanic,
				httpadpt.NewHandleWithSLogMiddleware(logger),
			},
		})
		if assert.NoError(t, adapter.Start(ctx)) {
			// path = 10
			resp, respErr := http.Get(fmt.Sprintf("http://localhost:%d/v1/test", port))
			if assert.NoError(t, respErr) {
				if assert.Equal(t, 204, resp.StatusCode) {
					assert.Contains(t, buf.String(), `"path":"/v1/test"`)
					assert.Contains(t, buf.String(), `"xid":"`)
					assert.Contains(t, buf.String(), `"duration":"`)
					assert.Contains(t, buf.String(), `"method":"GET"`)
				}
			}
			assert.NoError(t, adapter.Stop(ctx))
		}
	})
}
