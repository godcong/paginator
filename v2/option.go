package paginator

type Option func(p *paginator)

func PerPageOption(perPage int) Option {
	return func(p *paginator) {
		p.perPage = perPage
	}
}

func PerPageKeyOption(key string) Option {
	return func(p *paginator) {
		p.perPageKey = key
	}
}

func PageKeyOption(key string) Option {
	return func(p *paginator) {
		p.pageKey = key
	}
}
