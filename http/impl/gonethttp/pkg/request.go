package gonethttp

import (
	httpadpt "github.com/smart-libs/go-adapter/http/lib/pkg"
	"net/http"
)

type (
	Request struct{ httpReq *http.Request }

	query struct{ httpReq *http.Request }
)

func (q query) GetValue(name string) ([]string, bool) {
	if q.httpReq == nil || q.httpReq.URL == nil {
		return nil, false
	}

	m := q.httpReq.URL.Query()
	if len(m) == 0 {
		return nil, false
	}
	v, found := m[name]
	return v, found
}

func (r Request) Query() httpadpt.QueryParams {
	return query(r)
}

func NewRequest(httpReq *http.Request) httpadpt.Request {
	return Request{httpReq: httpReq}
}
