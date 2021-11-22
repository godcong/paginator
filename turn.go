package paginator

import (
	"net/http"
)

type Counter interface {
	Count(it Iterator) (int64, error)
}

type Finder interface {
	Find(pager PrePager) (interface{}, error)
}

type Requester interface {
	Request() *http.Request
}

type Turnable interface {
	Counter
	Finder
	Requester
}
