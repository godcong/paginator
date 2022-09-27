package paginator

import (
	"math"
)

// PageQuery ...
// @Description:
type PageQuery struct {
	CurrentPage  int
	LastPage     int
	PerPage      int
	Total        int64
	From         int
	To           int
	FirstPageURL string
	LastPageURL  string
	NextPageURL  string
	PrevPageURL  string
	Path         string
	Data         any
}

func (page *PageQuery) withTotal(total int64) *PageQuery {
	page.Total = total
	page.LastPage = int(math.Ceil(float64(page.Total) / float64(page.PerPage)))
	if page.CurrentPage <= 0 || page.CurrentPage > page.LastPage {
		page.CurrentPage = 1
	}
	return page
}

// Limit is an alias for Limit
// @receiver *PageQuery
// @return int
func (page *PageQuery) Limit() int {
	return page.PerPage
}

// Offset is an alias for From
// @receiver *PageQuery
// @return int
func (page *PageQuery) Offset() int {
	return page.From
}
