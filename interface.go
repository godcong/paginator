package paginator

import (
	"net/http"
)

type Counter interface {
	Count(pa *Parse[any]) (int64, error)
}

//type Finder interface {
//	Find(pager Pageable) (interface{}, error)
//}

type Requester interface {
	Request() *http.Request
}

type Turnable interface {
	Counter
	//Finder
	Requester
}

type Queryable[T any, Q any] interface {
	Count(pa *Parse[any]) (int64, error)
	Query() Q
	Find() (T, error)
	BeforeQuery(T) error
}
