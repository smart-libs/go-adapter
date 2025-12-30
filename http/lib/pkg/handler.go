package httpadpt

import sdkhandler "github.com/smart-libs/go-adapter/sdk/lib/pkg/handler"

type Handler = sdkhandler.Handler[Request, *Response]
