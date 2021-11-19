package paginator

import (
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
)

const defaultPaginatorPerPage = 15
const defaultPageKey = "page"
const defaultPerPageKey = "per_page"

var DEBUG = false
var ErrArgumentRequest = fmt.Errorf("paginator: argument request is not a valid http.Request")

type Requester interface {
	Request() *http.Request
}

type Paginator interface {
	Parse(t Turnable) (Page, error)
}

type paginator struct {
	pageKey    string
	perPageKey string
	perPage    int
}

func (p *paginator) getRequestPerPage(values url.Values) int {
	perPageInt := defaultPaginatorPerPage
	if !values.Has(p.perPageKey) {
		return perPageInt
	}
	var err error
	perPageInt, err = strconv.Atoi(values.Get(p.perPageKey))
	if err != nil {
		return defaultPaginatorPerPage
	}
	return perPageInt
}

func (p *paginator) getRequestCurrent(values url.Values) int {
	currentInt, err := strconv.Atoi(values.Get(p.pageKey))
	if err != nil {
		return 1
	}
	return currentInt
}

func (p *paginator) PageKey() string {
	return p.pageKey
}

func (p *paginator) PerPageKey() string {
	return p.perPageKey
}

func (p *paginator) Parse(t Turnable) (Page, error) {
	r := t.Request()
	if r == nil {
		return Page{}, ErrArgumentRequest
	}
	v := r.URL.Query()
	page := p.initPage(r, v)
	if DEBUG {
		fmt.Println("parse query", "per page", page.PerPage, "page", page.CurrentPage, "path", page.Path)
	}
	var conds []string
	if cond, ok := t.(Conditioner); ok {
		conds = cond.Conditions()
	}
	it := &iterator{
		v:     v,
		conds: conds,
	}

	err := findTotal(page, it, t)
	if err != nil {
		return Page{}, err
	}
	//p.from = (p.CurrentPage - 1) * p.PerPage
	//p.to = p.from + p.PerPage
}

func findTotal(p *Page, it Iterator, t Turnable) error {
	count, err := t.Count(it)
	if err != nil {
		return err
	}
	if count == 0 {
		p.CurrentPage = 1
		return nil
	}

	p.Total = count
	p.LastPage = int(math.Ceil(float64(p.Total) / float64(p.PerPage)))
	if p.CurrentPage <= 0 || p.CurrentPage > p.LastPage {
		p.CurrentPage = 1
	}
	//todo(from,to)
	//p.from = (p.CurrentPage - 1) * p.PerPage
	//p.to = p.from + p.PerPage
	p.NextPageURL = p.next(it.Values())
	p.PrevPageURL = p.prev(it.Values())
	p.LastPageURL = p.last(it.Values())
	p.FirstPageURL = p.first(it.Values())
	return nil
}

func ParsePage(pager Pager, conditions ...string) (Page, error) {
	p := &paginator{
		counter:    pager,
		finder:     pager,
		perPage:    defaultPaginatorPerPage,
		page:       &Page{},
		conditions: conditions,
		values:     nil,
	}
	p.ParseRequest(pager.Request())
	return p.Find()
}

func New(opts ...Option) Paginator {
	p := &paginator{
		perPage:    defaultPaginatorPerPage,
		pageKey:    defaultPageKey,
		perPageKey: defaultPerPageKey,
	}
	for i := range opts {
		opts[i](p)
	}
	return p
}

func (p *paginator) PerPage() int {
	return p.perPage
}

func (p *paginator) SetPerPage(perPage int) {
	p.perPage = perPage
}

func (p *paginator) Offset() int {
	return p.from
}

func (p *paginator) Limit() int {
	return p.page.PerPage
}

func (p *paginator) Total() int64 {
	return p.page.Total
}

func (p *paginator) SetConditions(conditions ...string) {
	p.conditions = conditions
}

func (p *paginator) Conditions() (url.Values, bool) {
	values := url.Values{}
	if len(p.conditions) == 0 {
		return values, false
	}
	tmp := ""
	for _, condition := range p.conditions {
		tmp = p.values.Get(condition)
		if tmp != "" {
			values.Set(condition, tmp)
		}
	}
	return values, true
}

func (p *paginator) find() error {
	count, err := p.counter.Count(p)
	if DEBUG {
		fmt.Println("count", "error", err, "count", count)
	}
	if err != nil {
		return err
	}

	v, err := p.finder.Find(p)
	if err != nil {
		return err
	}
	p.page.Data = v

	return nil
}

func (p *paginator) initPage(r *http.Request, v url.Values) (page *Page) {
	return &Page{
		CurrentPage:  p.getRequestCurrent(v),
		LastPage:     0,
		PerPage:      p.getRequestPerPage(v),
		Data:         nil,
		Total:        0,
		FirstPageURL: "",
		LastPageURL:  "",
		NextPageURL:  "",
		PrevPageURL:  "",
		Path:         r.Host + r.URL.Path,
	}
}

func (p *paginator) getTotal(p2 *Page) interface{} {
	if count == 0 {
		page.CurrentPage = 1
		return
	}

	page.Total = count
	page.LastPage = int(math.Ceil(float64(page.Total) / float64(page.PerPage)))
	if page.CurrentPage <= 0 || page.CurrentPage > page.LastPage {
		page.CurrentPage = 1
	}
	p.from = (page.CurrentPage - 1) * page.PerPage
	p.to = p.from + page.PerPage
	v, _ := p.Conditions()
	page.NextPageURL = page.next(v)
	page.PrevPageURL = page.prev(v)
	page.LastPageURL = page.last(v)
	page.FirstPageURL = page.first(v)
	return
}

func (p *Page) next(v url.Values) string {
	if p.LastPage > p.CurrentPage+1 {
		p.perPage(v)
		page(v, p.CurrentPage+1)
		return p.Path + "?" + v.Encode()
	}
	return ""
}

func (p *Page) prev(v url.Values) string {
	if p.CurrentPage-1 > 0 {
		p.perPage(v)
		page(v, p.CurrentPage-1)
		return p.Path + "?" + v.Encode()
	}
	return ""
}

func (p *Page) last(v url.Values) string {
	if p.LastPage > 0 {
		p.perPage(v)
		page(v, p.LastPage)
		return p.Path + "?" + v.Encode()
	}
	return ""
}
func (p *Page) first(v url.Values) string {
	if p.Total > 0 {
		p.perPage(v)
		page(v, 1)
		return p.Path + "?" + v.Encode()
	}
	return ""
}

func page(values url.Values, i int) {
	values.Set("page", strconv.Itoa(i))
}

func perPage(values url.Values, i int) {
	values.Set("per_page", strconv.Itoa(i))
}

func (p *Page) page(values url.Values, i int) {
	values.Set("per_page", strconv.Itoa(i))
}

func (p *Page) perPage(values url.Values) {
	perPage(values, p.PerPage)
}

var _ Paginator = (*paginator)(nil)
