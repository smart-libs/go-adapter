package gonethttp

import (
	"context"
	httpadpt "github.com/smart-libs/go-adapter/http/lib/pkg"
	"net/http"
)

func handleResponse(_ context.Context, w http.ResponseWriter, resp httpadpt.Response) {
	// Set headers before WriteHeader (headers must be set before WriteHeader)
	if len(resp.Header) > 0 {
		for k, v := range resp.Header {
			w.Header()[k] = v
		}
	}
	// Set status code
	if resp.StatusCode != nil {
		w.WriteHeader(*resp.StatusCode)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	// Write body
	if resp.Body != nil {
		if _, err := w.Write(resp.Body); err != nil {
			// Log error if possible, but response may already be committed
			// In production, consider using a logger here
			_ = err
		}
	}
}
