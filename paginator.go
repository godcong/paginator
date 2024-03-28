// Package paginator is a paginator for data pagination
package paginator

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
)

// Paginator is a paginator for data pagination
type Paginator interface {
	Parse(Parsable, ...string) (PageReady, error)
	Ready() PageReady
	ToJSON(PageResult) ([]byte, error)
}

type paginator struct {
	cfg *Config
}

var ErrInvalidParserType = errors.New("invalid parser type")

// New ...
// @Description: create paginator object for use anywhere
// @param ...SettingOption
// @return Paginator
func New(opts ...*Option) Paginator {
	cfg := NewConfig(opts...)
	return &paginator{
		cfg: cfg,
	}
}

func (p *paginator) Parse(obj Parsable, ignores ...string) (PageReady, error) {
	ignored := make(map[string]struct{})
	for i := range ignores {
		ignored[ignores[i]] = struct{}{}
	}

	var (
		pr  PageReady
		err error
	)
	switch v := obj.(type) {
	case *httpParse:
		pr, err = v.parse(p.cfg.Clone(), ignored)
	case *CustomParse:
		pr, err = v.Parse(p.cfg.Clone(), ignored)
	default:
		return nil, ErrInvalidParserType
	}
	return pr, err
}

func (p *paginator) Ready() PageReady {
	var pr query
	pr.config = p.cfg
	return &pr
}

func (p *paginator) ToJSON(result PageResult) ([]byte, error) {
	v := p.cfg.Output(&result)
	return json.Marshal(v)
}

func setPage(values url.Values, key string, i int) {
	values.Set(key, strconv.Itoa(i))
}

func setPerPage(values url.Values, key string, i int) {
	values.Set(key, strconv.Itoa(i))
}

func stoi(s string, d int) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return d
	}
	return i
}

var _ Paginator = (*paginator)(nil)
