package paginator

type Counter interface {
	Count() (int64, error)
}

type Getter interface {
	Get() (any, error)
}

type Finder interface {
	Counter
	Getter
}

type Queryable interface {
	Finder(Parser[any]) Finder
}
