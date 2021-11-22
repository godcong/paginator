package paginator

import (
	"net/url"
)

type Iterator interface {
	Range(fn func(key string, val []string) bool)
	Values() url.Values
	Conditions() []string
}

type iterator struct {
	v     url.Values
	conds []string
}

func (i iterator) Values() url.Values {
	return i.v
}

func (i iterator) Conditions() []string {
	return i.conds
}

func (i iterator) Range(fn func(key string, val []string) bool) {
	for _, cond := range i.conds {
		if !fn(cond, i.v[cond]) {
			return
		}
	}
}

var _ Iterator = (*iterator)(nil)
