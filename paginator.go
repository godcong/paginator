// Package paginator is a paginator for data pagination
package paginator

import (
	"fmt"
	"net/http"
	"strconv"
)

// Paginator is a paginator for data pagination
type Paginator interface {
	SetDefaultQuery(queryable Queryable) Paginator
	Apply(opts ...OptionSet) Paginator
	Parse(Parser) (any, error)
	ParseWithQuery(Parser, Queryable) (any, error)
}

type paginator struct {
	query Queryable
	op    *Option
}

// New ...
// @Description: create paginator object for use anywhere
// @param ...OptionSet
// @return Paginator
func New(ops ...OptionSet) Paginator {
	op := defaultOption()
	for i := range ops {
		ops[i](op)
	}
	return &paginator{
		op: op,
	}
}

func (p *paginator) SetDefaultQuery(queryable Queryable) Paginator {
	p.query = queryable
	return p
}

func (p *paginator) Apply(opts ...OptionSet) Paginator {
	for i := range opts {
		opts[i](p.op)
	}
	return p
}

func (p *paginator) ParseWithQuery(parser Parser, finder Queryable) (any, error) {
	return p.parse(parser, finder)
}

func (p *paginator) Parse(parser Parser) (any, error) {
	return p.parse(parser, p.query)
}

func (p *paginator) parse(parser Parser, query Queryable) (any, error) {
	pr := p.initialize(parser)
	finder := p.getFinder(parser, query)

	count, err := finder.Count(parser.Context())
	if err != nil {
		return nil, err
	}
	if !p.total(pr, count) {
		return pr.values(pr.page, p.op), nil
	}

	pr.page.From = (pr.page.CurrentPage - p.op.StartOffset()) * pr.page.PerPage
	pr.page.To = pr.page.From + pr.page.PerPage

	v, err := finder.Clone().Get(parser.Context(), *pr.page)
	if err != nil {
		return nil, err
	}
	pr.page.Data = v
	fmt.Printf("%+v", pr.page)
	return pr.values(pr.page, p.op), err
}

func (p *paginator) total(pr *pageReady, count int64) bool {
	if count == 0 {
		pr.page.CurrentPage = 1
		return false
	}
	pr.page.withTotal(count)
	pr.page.NextPageURL = p.nextURL(pr)
	pr.page.PrevPageURL = p.prevURL(pr)
	pr.page.LastPageURL = p.lastURL(pr)
	pr.page.FirstPageURL = p.firstURL(pr)
	return true
}

func (p *paginator) nextURL(pa *pageReady) string {
	enc := pa.parser.FindOthers()
	if pa.page.LastPage > pa.page.CurrentPage {
		setPerPage(enc, p.op.PerPageKey(), p.op.PerPage())
		setPage(enc, p.op.PageKey(), pa.page.CurrentPage+1)
		return pa.page.Path + "?" + enc.Encode()
	}
	return ""
}

func (p *paginator) prevURL(pa *pageReady) string {
	enc := pa.parser.FindOthers()
	if pa.page.CurrentPage > 1 {
		setPerPage(enc, p.op.PerPageKey(), p.op.PerPage())
		setPage(enc, p.op.PageKey(), pa.page.CurrentPage-1)
		return pa.page.Path + "?" + enc.Encode()
	}
	return ""
}

func (p *paginator) lastURL(pa *pageReady) string {
	enc := pa.parser.FindOthers()
	if pa.page.LastPage > 0 {
		setPerPage(enc, p.op.PerPageKey(), p.op.PerPage())
		setPage(enc, p.op.PageKey(), pa.page.LastPage)
		return pa.page.Path + "?" + enc.Encode()
	}
	return ""
}
func (p *paginator) firstURL(pa *pageReady) string {
	enc := pa.parser.FindOthers()
	if pa.page.Total > 0 {
		setPerPage(enc, p.op.PerPageKey(), p.op.PerPage())
		setPage(enc, p.op.PageKey(), 1)
		return pa.page.Path + "?" + enc.Encode()
	}
	return ""
}

func setPage(values Values, key string, i int) {
	values.Set(key, strconv.Itoa(i))
}

func setPerPage(values Values, key string, i int) {
	values.Set(key, strconv.Itoa(i))
}

func (p *paginator) initialize(parser Parser) *pageReady {
	page := new(PageQuery)

	src := parser.GetSource()
	switch v := src.(type) {
	case *http.Request:
		page.Path = requestPath(v)
	}
	page.PerPage = stoi(parser.FindValue(p.op.perPageKey, ""), p.op.perPage)
	page.CurrentPage = stoi(parser.FindValue(p.op.pageKey, ""), p.op.startPage)
	values := orderedValues
	if !p.op.order {
		values = unorderedValues
	}
	return &pageReady{
		values: values,
		page:   page,
		parser: parser,
	}
}

func (p *paginator) getFinder(parser Parser, query Queryable) CloneFinder {
	f := query.Finder(parser)
	if c, ok := f.(CloneFinder); ok {
		return c
	}
	return wrapClone{
		query:  query,
		parser: parser,
	}
}

func stoi(s string, d int) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return d
	}
	return i
}

var _ Paginator = (*paginator)(nil)
