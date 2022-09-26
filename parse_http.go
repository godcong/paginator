package paginator

import (
	"net/http"
	"net/url"
	"strconv"
)

type httpParser struct {
	r      *http.Request
	query  url.Values
	others url.Values
}

func (p *httpParser) GetEncoder() Values {
	return p.r.URL.Query()
}

func (p *httpParser) GetSource() any {
	return any(p.r)
}

func (p *httpParser) FindValue(key string, d string) string {
	if p.query.Has(key) {
		p.others.Del(key)
		return p.query.Get(key)
	}
	return d
}

func (p *httpParser) FindArray(key string, d []string) []string {
	if p.query.Has(key) {
		p.others.Del(key)
		return p.query[key]
	}
	return d
}

func (p *httpParser) FindNumber(key string, d float64) float64 {
	if p.query.Has(key) {
		p.others.Del(key)
		if n, err := strconv.ParseFloat(p.query.Get(key), 10); err == nil {
			return n
		}
	}
	return d
}

func (p *httpParser) FindOthers() Values {
	return p.others
}

func NewHttpParser(r *http.Request) Parser {
	return &httpParser{
		r:      r,
		query:  r.URL.Query(),
		others: r.URL.Query(),
	}
}

func requestPath(r *http.Request) string {
	scheme := "http"
	if r.URL.Scheme != "" {
		scheme = r.URL.Scheme
	}
	return scheme + "://" + r.Host + r.URL.Path
}
