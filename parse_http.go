package paginator

import (
	"net/http"
	"net/url"
)

type httpParse struct {
	r     *http.Request
	query url.Values
}

func (p *httpParse) Parse(cfg *Config, ignored map[string]struct{}) (PageReady, error) {
	return p.parse(cfg, ignored)
}

func (p *httpParse) Type() string {
	return "http"
}

func (p *httpParse) parse(cfg *Config, ignored map[string]struct{}) (*query, error) {
	var pr query
	pr.config = cfg
	for s := range ignored {
		delete(p.query, s)
	}

	if !cfg.SkipPath {
		pr.fullPath = fullRequestPath(p.r)
	}
	p.query = p.r.URL.Query()
	var key string
	key = cfg.GetKey(KeyPerPage)
	pr.perPage = stoi(p.query.Get(key), cfg.PerPage)
	delete(p.query, key)

	key = cfg.GetKey(KeyPage)
	pr.currentPage = stoi(p.query.Get(key), cfg.StartPage)
	delete(p.query, key)

	pr.values = p.query
	return &pr, nil
}

func getValueOrDefault(v url.Values, key string, dv int) int {
	return stoi(v.Get(key), dv)
}

// FromHTTP creates a new http parser
// @param *http.Request
// @return Parsable
func FromHTTP(r *http.Request) Parsable {
	return &httpParse{
		r: r,
	}
}

func fullRequestPath(r *http.Request) string {
	scheme := "http"
	if r.URL.Scheme != "" {
		scheme = r.URL.Scheme
	}
	return scheme + "://" + r.Host + r.URL.Path
}
