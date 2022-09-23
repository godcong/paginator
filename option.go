package paginator

type s struct {
	perPage    int
	perPageKey string
	pageKey    string
}

type Option func(p *paginator)

// PerPageOption ...
// @Description: per page init on paginator
// @param int
// @return Option
func PerPageOption(perPage int) Option {
	return func(p *paginator) {
		p.perPage = perPage
	}
}

// PerPageKeyOption ...
// @Description: per page key init on paginator
// @param string
// @return Option
func PerPageKeyOption(key string) Option {
	return func(p *paginator) {
		p.perPageKey = key
	}
}

// PageKeyOption ...
// @Description: per key init on paginator
// @param string
// @return Option
func PageKeyOption(key string) Option {
	return func(p *paginator) {
		p.pageKey = key
	}
}
