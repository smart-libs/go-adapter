package gonethttp

import (
	"context"
	"fmt"
	httpadpt "github.com/smart-libs/go-adapter/http/lib/pkg"
	"net/http"
)

type (
	DefaultAdapter struct {
		config     httpadpt.Config
		serveMux   *http.ServeMux
		server     *http.Server
		startError error
	}
)

func NewAdapter(config httpadpt.Config) (httpadpt.Adapter, error) {
	adapter := DefaultAdapter{config: config}
	port := 80
	if config.Port != nil {
		port = *config.Port
	}
	host := ""
	if config.Host != nil {
		host = *config.Host
	}

	adapter.serveMux = http.NewServeMux()
	adapter.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: adapter.serveMux,
	}

	if err := buildAndAddHandles(adapter.serveMux.Handle, config.Bindings); err != nil {
		return nil, err
	}
	return &adapter, nil
}

func (d *DefaultAdapter) Start(_ context.Context) error {
	if d.server == nil {
		return fmt.Errorf("server already shut down")
	}
	go func() { d.startError = d.server.ListenAndServe() }()
	return nil
}

func (d *DefaultAdapter) Stop(ctx context.Context) error {
	if d.server == nil {
		return fmt.Errorf("server not started")
	}
	if d.startError != nil {
		return d.startError
	}
	err := d.server.Shutdown(ctx)
	if err == nil {
		err = d.server.Close()
	}
	return err
}
