package paginator

import (
	"maps"
)

const (
	defaultStartIndex = 0
	defaultStartPage  = 1
	defaultPerPage    = 50
)

var NULL struct{}

// Key is the key for the paginator
// ENUM(success,page, per_page, last_page,next_page,previous_page, data,current_path, first_page_url, last_page_url, next_page_url, prev_page_url, current_page, total, path,key_max)
type Key int

// Value is a Value of paginator
// ENUM(start_index, start_page, per_page)
type Value int

type Config struct {
	Keys            map[Key]string
	Ignored         map[string]struct{}
	OnQueryKeys     map[string]struct{}
	NoSuccess       bool
	SkipPath        bool
	StartIndex      int
	StartPage       int
	PerPage         int
	UsePathPaginate bool
}

var defaultOnQueryKeys = map[string]struct{}{
	KeyPage.String():    {},
	KeyPerPage.String(): {},
}

// StartOffset ...
// @receiver *Config
// @return int
func (c *Config) StartOffset() int {
	return c.StartPage - c.StartIndex
}

func (c *Config) Customize() bool {
	return len(c.Keys) > 0
}

func (c *Config) GetKey(key Key) string {
	if v, ok := c.Keys[key]; ok {
		return v
	}
	return key.String()
}

func (c *Config) SetValue(key Value, value int) {
	switch key {
	case ValueStartIndex:
		c.StartIndex = value
	case ValueStartPage:
		c.StartPage = value
	case ValuePerPage:
		c.PerPage = value
	}
}

func (c *Config) IsIgnored(key Key) bool {
	if _, ok := c.Ignored[key.String()]; ok {
		return true
	}
	return false
}

func (c *Config) Output(result *PageResult) any {
	if c.UsePathPaginate {
		result.CurrentPage = 0
		result.PreviousPage = 0
		result.NextPage = 0
		result.LastPage = 0
		result.Path = ""
	}
	if c.Customize() {
		return customOutput(c, result)
	}
	return result

}

func (c *Config) parse(opts ...*Option) {
	var (
		kk Key
		kv string
		vk Value
		vv int
	)
	for i := range opts {
		if opts[i].usePathPaginate {
			c.UsePathPaginate = true
		}

		if opts[i].noSuccess {
			c.NoSuccess = true
		}
		for vk, vv = range opts[i].valSet {
			c.SetValue(vk, vv)
		}

		if !opts[i].customize {
			continue
		}

		for kk, kv = range opts[i].keySet {
			// if _, ok := c.OnQueryKeys[kk.String()]; ok {
			//     delete(c.OnQueryKeys, kk.String())
			//     c.OnQueryKeys[kv] = NULL
			// }
			c.Keys[kk] = kv
		}

		if opts[i].ignoreKeys != nil {
			for _, k := range opts[i].ignoreKeys {
				c.Ignored[k.String()] = NULL
			}
		}
	}
}

func NewConfig(opts ...*Option) *Config {
	c := new(Config)
	c.Keys = make(map[Key]string)
	c.Ignored = make(map[string]struct{})
	c.OnQueryKeys = maps.Clone(defaultOnQueryKeys)
	c.NoSuccess = false
	c.SkipPath = false
	c.StartIndex = defaultStartIndex
	c.StartPage = defaultStartPage
	c.PerPage = defaultPerPage

	c.parse(opts...)

	urlCount := 0
	if c.IsIgnored(KeyPrevPageUrl) {
		urlCount++
	}
	if c.IsIgnored(KeyNextPageUrl) {
		urlCount++
	}
	if c.IsIgnored(KeyFirstPageUrl) {
		urlCount++
	}
	if c.IsIgnored(KeyLastPageUrl) {
		urlCount++
	}
	if urlCount == 4 {
		c.SkipPath = true
	}

	if c.SkipPath {
		c.Ignored[c.GetKey(KeyPrevPageUrl)] = NULL
		c.Ignored[c.GetKey(KeyNextPageUrl)] = NULL
		c.Ignored[c.GetKey(KeyFirstPageUrl)] = NULL
		c.Ignored[c.GetKey(KeyLastPageUrl)] = NULL
	}

	if !c.SkipPath && c.UsePathPaginate {
		c.Ignored[c.GetKey(KeyCurrentPage)] = NULL
		c.Ignored[c.GetKey(KeyPreviousPage)] = NULL
		c.Ignored[c.GetKey(KeyLastPage)] = NULL
		c.Ignored[c.GetKey(KeyNextPage)] = NULL
	}

	if c.NoSuccess {
		c.Ignored[c.GetKey(KeySuccess)] = NULL
	}

	return c
}

func (c *Config) Clone() *Config {
	cfg := new(Config)
	cfg.Keys = maps.Clone(cfg.Keys)
	cfg.OnQueryKeys = maps.Clone(cfg.OnQueryKeys)
	cfg.SkipPath = c.SkipPath
	cfg.StartIndex = c.StartIndex
	cfg.StartPage = c.StartPage
	cfg.PerPage = c.PerPage
	return cfg
}
