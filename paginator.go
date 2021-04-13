package paginator

import (
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
)

const DefaultPaginatorPerPage = 10

var DEBUG = false

type Conditioner interface {
	Conditions() (url.Values, bool)
}

type Counter interface {
	Count(cond Conditioner) (int64, error)
}

type Finder interface {
	Find(p Paginator) (interface{}, error)
}

type Requester interface {
	Request() *http.Request
}

type Pager interface {
	Counter
	Finder
	Requester
}

type Paginator interface {
	Find() (Page, error)
	SetPerPage(int)
	PerPage() int
	ParseRequest(r *http.Request) Paginator
	Conditions() (url.Values, bool)
	Offset() int
	Limit() int
	Total() int64
}

type Page struct {
	CurrentPage  int         `json:"current_page"`
	LastPage     int         `json:"last_page"`
	PerPage      int         `json:"per_page"`
	Data         interface{} `json:"data"`
	Total        int64       `json:"total"`
	FirstPageURL string      `json:"first_page_url"`
	LastPageURL  string      `json:"last_page_url"`
	NextPageURL  string      `json:"next_page_url"`
	PrevPageURL  string      `json:"prev_page_url"`
	Path         string      `json:"path"`
}

type paginator struct {
	counter    Counter
	finder     Finder
	from       int
	to         int
	conditions []string
	perPage    int
	page       *Page
	values     url.Values
}

func ParsePage(pager Pager, conditions ...string) (Page, error) {
	p := &paginator{
		counter:    pager,
		finder:     pager,
		perPage:    DefaultPaginatorPerPage,
		page:       &Page{},
		conditions: conditions,
		values:     nil,
	}
	p.ParseRequest(pager.Request())
	return p.Find()
}

func New(counter Counter, finder Finder) Paginator {
	return &paginator{
		counter:    counter,
		finder:     finder,
		perPage:    DefaultPaginatorPerPage,
		page:       &Page{},
		conditions: []string{},
		values:     nil,
	}
}

func (p *paginator) PerPage() int {
	return p.perPage
}

func (p *paginator) SetPerPage(perPage int) {
	p.perPage = perPage
}

func (p *paginator) Find() (Page, error) {
	if err := p.find(); err != nil {
		return Page{}, err
	}
	return *p.page, nil
}

func (p *paginator) ParseRequest(r *http.Request) Paginator {
	var err error

	p.values = r.URL.Query()
	perPage := p.values.Get("per_page")
	p.page.PerPage, err = strconv.Atoi(perPage)
	if err != nil {
		p.page.PerPage = p.PerPage()
	}
	current := p.values.Get("page")
	p.page.CurrentPage, err = strconv.Atoi(current)
	if err != nil {
		p.page.CurrentPage = 1
	}

	p.page.Path = r.Host + r.URL.Path
	if DEBUG {
		fmt.Println("parse query", "raw", r.URL.RawQuery, "per page", perPage, "page", current, "path", p.page.Path)
	}
	return p
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

	p.page = p.makePage(count)

	v, err := p.finder.Find(p)
	if err != nil {
		return err
	}
	p.page.Data = v

	return nil
}

func (p *paginator) makePage(count int64) (page *Page) {
	page = p.page
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
