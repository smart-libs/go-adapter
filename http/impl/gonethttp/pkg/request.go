package gonethttp

import (
	httpadpt "github.com/smart-libs/go-adapter/http/lib/pkg"
	"net/http"
	"net/url"
)

type (
	Request struct{ httpReq *http.Request }

	query  struct{ url *url.URL }
	header struct{ header http.Header }
	path   struct{ httpReq *http.Request }
)

func (p path) GetValue(pathParamName string) (string, bool) {
	if p.httpReq == nil {
		return "", false
	}
	value := p.httpReq.PathValue(pathParamName)
	// PathValue returns empty string if the parameter doesn't exist
	// However, we can't distinguish between "not found" and "empty value"
	// So we return true if the value is non-empty, false otherwise
	// Note: This means empty path parameter values will be treated as "not found"
	return value, value != ""
}

func (r Request) Header() httpadpt.HeaderParams {
	if r.httpReq == nil {
		return header{}
	}
	return header{header: r.httpReq.Header}
}

func (r Request) Path() httpadpt.PathParams {
	if r.httpReq == nil {
		return path{}
	}
	return path(r)
}

func (q header) GetValue(name string) ([]string, bool) {
	if len(q.header) == 0 {
		return nil, false
	}

	v := q.header.Values(name)
	return v, len(v) > 0
}

func (q query) GetValue(name string) ([]string, bool) {
	if q.url == nil {
		return nil, false
	}

	m := q.url.Query()
	if len(m) == 0 {
		return nil, false
	}
	v, found := m[name]
	return v, found
}

func (r Request) Query() httpadpt.QueryParams {
	if r.httpReq == nil {
		return query{}
	}
	return query{url: r.httpReq.URL}
}

func NewRequest(httpReq *http.Request) httpadpt.Request {
	return Request{httpReq: httpReq}
}
