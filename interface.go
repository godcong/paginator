package paginator

// Counter is a counter for Query
type Counter interface {
	Count() (int64, error)
}

// Getter is a getter for Query
type Getter interface {
	Get(page PageQuery) (any, error)
}

// Finder is a Finder for Query
type Finder interface {
	Counter
	Getter
}

// Queryable is an interface for Query
type Queryable interface {
	Finder(Parser) Finder
}
