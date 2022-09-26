package paginator

import (
	"fmt"
	"net/http"
	"strconv"
)

var DEBUG = false
var ErrArgumentRequest = fmt.Errorf("paginator: argument request is not a valid http.Request")

type Paginator interface {
	SetDefaultQuery(queryable Queryable) Paginator
	Apply(opts ...OptionSet) Paginator
	Parse(Parser[any]) (any, error)
	ParseWithQuery(Parser[any], Queryable) (any, error)
}

//type T any

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

func (p *paginator) ParseWithQuery(parser Parser[any], finder Queryable) (any, error) {
	return p.parse(parser, finder)
}

func (p *paginator) Parse(parser Parser[any]) (any, error) {
	return p.parse(parser, p.query)
}

// Count
// @receiver *paginator
// @param Parser[any]
// @param Queryable
// @return int64
// @return error
func (p *paginator) Count(pa Parser[any], cf Queryable) (int64, error) {
	return cf.Finder(pa).Count()
}

// Get
// @receiver *paginator
// @param Parser[any]
// @param Queryable
// @return any
// @return error
func (p *paginator) Get(pa Parser[any], cf Queryable) (any, error) {
	return cf.Finder(pa).Get()
}

func (p *paginator) parse(parser Parser[any], cf Queryable) (any, error) {
	pr := p.initialize(parser)
	count, err := p.Count(parser, cf)
	if err != nil {
		return nil, err
	}
	if !p.total(pr, count) {
		return pr.page.values(p.op), nil
	}

	pr.page.From = (pr.page.CurrentPage - 1) * pr.page.PerPage
	pr.page.To = pr.page.From + pr.page.PerPage
	v, err := p.Get(parser, cf)
	if err != nil {
		return nil, err
	}
	pr.page.Data = v
	return pr.page.values(p.op), err
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
	if pa.page.LastPage > pa.page.CurrentPage+1 {
		setPerPage(enc, p.op.PerPageKey(), p.op.PerPage())
		setPage(enc, p.op.PageKey(), pa.page.CurrentPage+1)
		return pa.page.Path + "?" + pa.enc.Encode()
	}
	return ""
}

func (p *paginator) prevURL(pa *pageReady) string {
	enc := pa.parser.FindOthers()
	if pa.page.CurrentPage-1 > 0 {
		setPerPage(enc, p.op.PerPageKey(), p.op.PerPage())
		setPage(enc, p.op.PageKey(), pa.page.CurrentPage-1)
		return pa.page.Path + "?" + pa.enc.Encode()
	}
	return ""
}

func (p *paginator) lastURL(pa *pageReady) string {
	enc := pa.parser.FindOthers()
	if pa.page.LastPage > 0 {
		setPerPage(enc, p.op.PerPageKey(), p.op.PerPage())
		setPage(enc, p.op.PageKey(), pa.page.LastPage)
		return pa.page.Path + "?" + pa.enc.Encode()
	}
	return ""
}
func (p *paginator) firstURL(pa *pageReady) string {
	enc := pa.parser.FindOthers()
	if pa.page.Total > 0 {
		setPerPage(enc, p.op.PerPageKey(), p.op.PerPage())
		setPage(enc, p.op.PageKey(), 1)
		return pa.page.Path + "?" + pa.enc.Encode()
	}
	return ""
}

func setPage(values Values, key string, i int) {
	values.Set(key, strconv.Itoa(i))
}

func setPerPage(values Values, key string, i int) {
	values.Set(key, strconv.Itoa(i))
}

func (p *paginator) initialize(parser Parser[any]) *pageReady {
	page := new(pageQuery)
	switch v := parser.(type) {
	case Parser[*http.Request]:
		r := v.GetSource()
		page.Path = r.URL.Scheme + r.Host + r.URL.Path
	}
	page.PerPage = stoi(parser.FindValue(p.op.perPageKey, ""), p.op.perPage)
	page.CurrentPage = stoi(parser.FindValue(p.op.pageKey, ""), p.op.staPage)
	return &pageReady{
		page:   page,
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
