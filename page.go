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

type pageReady struct {
	parser Parser
	page   *PageQuery
	//enc    Values
}

func (page *PageQuery) withTotal(total int64) *PageQuery {
	page.Total = total
	page.LastPage = int(math.Ceil(float64(page.Total) / float64(page.PerPage)))
	if page.CurrentPage <= 0 || page.CurrentPage > page.LastPage {
		page.CurrentPage = 1
	}
	return page
}

func (page *PageQuery) values(op *Option) any {
	values := make(map[string]any)
	values[op.DataKey()] = page.Data
	values[op.PerPageKey()] = page.PerPage
	values[op.CurrentPageKey()] = page.CurrentPage
	values[op.TotalKey()] = page.Total
	values[op.PathKey()] = page.Path
	values[op.LastPageKey()] = page.LastPage
	values[op.FirstPageKey()] = page.FirstPageURL
	values[op.LastPageKey()] = page.LastPageURL
	values[op.NextPageKey()] = page.NextPageURL
	values[op.PrevPageKey()] = page.PrevPageURL
	return values
}
