package gonethttp

import (
	"fmt"
	httpadpt "github.com/smart-libs/go-adapter/http/lib/pkg"
	serror "github.com/smart-libs/go-crosscutting/serror/lib/pkg"
	"net/http"
)

func buildPath(method, path string) string {
	if method != "" {
		return method + " " + path
	}
	return path
}

func buildHandler(handler httpadpt.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		req := NewRequest(r)
		resp := httpadpt.Response{}
		err := handler.Invoke(ctx, req, &resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		handleResponse(ctx, w, resp)
	}
}

func buildAndAddHandles(addHandle func(path string, handler http.Handler), bindings httpadpt.Bindings) error {
	fName := "httpadpt.buildAndAddHandles"
	for i, binding := range bindings {
		if len(binding.Condition.Methods) > 0 {
			if binding.Condition.Path == nil {
				return serror.IllegalConfigParamValue(
					fmt.Sprintf("%s.Config.Bindings[%d].Condition.Path", fName, i), binding.Condition.Path)
			}
			for _, method := range binding.Condition.Methods {
				addHandle(buildPath(method, *binding.Condition.Path), buildHandler(binding.Handler))
			}
		}
	}
	return nil
}
