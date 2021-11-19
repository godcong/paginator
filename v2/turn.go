package paginator

type Turnable interface {
	Counter
	Finder
	Requester
}

type Counter interface {
	Count(it Iterator) (int64, error)
}

type Finder interface {
	Find(p Paginator) (interface{}, error)
}
