package paginator

import (
	"net/http"
	"net/url"
)

type httpParser struct {
	r     *http.Request
	query url.Values
	clone url.Values
}

func (p *httpParser) Get() *http.Request {
	return p.r
}

func (p *httpParser) FindValue(key string, d string) string {
	if p.query.Has(key) {
		p.clone.Del(key)
		return p.query.Get(key)
	}
	return d
}

func (p *httpParser) Other() url.Values {
	return p.clone
}

func NewHttpParser(r *http.Request) Parser[*http.Request] {
	return &httpParser{
		r:     r,
		query: r.URL.Query(),
		clone: r.URL.Query(),
	}
}
