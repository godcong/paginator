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

type Counter interface {
	Count(v *int64) error
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
	counter Counter
	finder  Finder
	from    int
	to      int
	perPage int
	page    *Page
	values  url.Values
}

func ParsePage(pager Pager) (Page, error) {
	p := &paginator{
		counter: pager,
		finder:  pager,
		perPage: DefaultPaginatorPerPage,
		page:    &Page{},
		values:  nil,
	}
	p.ParseRequest(pager.Request())
	return p.Find()
}

func New(counter Counter, finder Finder) Paginator {
	return &paginator{
		counter: counter,
		finder:  finder,
		perPage: DefaultPaginatorPerPage,
		page:    &Page{},
		values:  nil,
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

func (p *paginator) find() error {
	var count int64

	err := p.counter.Count(&count)
	if DEBUG {
		fmt.Println("count", "error", err, "count", count)
	}
	if err != nil {
		return err
	}

	p.from, p.to = p.page.make(count)

	v, err := p.finder.Find(p)
	if err != nil {
		return err
	}
	p.page.Data = v

	return nil
}

func (p *Page) make(count int64) (from, to int) {
	if count == 0 {
		p.CurrentPage = 1
		return
	}

	p.Total = count
	p.LastPage = int(math.Ceil(float64(p.Total) / float64(p.PerPage)))
	if p.CurrentPage <= 0 || p.CurrentPage > p.LastPage {
		p.CurrentPage = 1
	}
	from = (p.CurrentPage - 1) * p.PerPage
	to = from + p.PerPage
	p.NextPageURL = p.next()
	p.PrevPageURL = p.prev()
	p.LastPageURL = p.last()
	p.FirstPageURL = p.first()
	return
}

func (p *Page) next() string {
	if p.LastPage > p.CurrentPage+1 {
		v := url.Values{}
		p.perPage(v)
		page(v, p.CurrentPage+1)
		return p.Path + "?" + v.Encode()
	}
	return ""
}

func (p *Page) prev() string {
	if p.CurrentPage-1 > 0 {
		v := url.Values{}
		p.perPage(v)
		page(v, p.CurrentPage-1)
		return p.Path + "?" + v.Encode()
	}
	return ""
}

func (p *Page) last() string {
	if p.LastPage > 0 {
		v := url.Values{}
		p.perPage(v)
		page(v, p.LastPage)
		return p.Path + "?" + v.Encode()
	}
	return ""
}
func (p *Page) first() string {
	if p.Total > 0 {
		v := url.Values{}
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
