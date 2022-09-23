package paginator

import (
	"encoding/json"
)

// Pager ...
// @Description:
type Pager interface {
	Counter
	Requester
}

// pageQuery ...
// @Description:
type pageQuery struct {
	CurrentPage  int
	LastPage     int
	PerPage      int
	Total        int64
	FirstPageURL string
	LastPageURL  string
	NextPageURL  string
	PrevPageURL  string
	Path         string
	Data         json.RawMessage
}

type pageReady struct {
	parser Parser[any]
	page   *pageQuery
}

type Pageable interface {
}

//
//type pageQuery[T any] struct {
//	page *Find[T]
//	from int
//	to   int
//	it   *iterator
//}
//
//func (p pageQuery) Find() Find {
//	return *p.page
//}
//
//func (p pageQuery) Offset() int {
//	return p.from
//}
//
//func (p pageQuery) Limit() int {
//	return p.page.PerPage
//}
//
//func (p pageQuery) Total() int64 {
//	return p.page.Total
//}
//
//func (p pageQuery) Current() int {
//	return p.page.CurrentPage
//}
//
//func (p pageQuery) To() int {
//	return p.to
//}
//
//func (p pageQuery) Iterator() Iterator {
//	return p.it
//}
//
//var _ Pageable = (*pageQuery)(nil)
