package paginator

import (
	"fmt"
	"github.com/goextension/log/zap"
	"math"
	"net/http"
	"net/url"
	"strconv"

	"github.com/goextension/log"
)

const DefaultPaginatorPerPage = 10

var DEBUG = false

type Counter interface {
	Count(v *int64) error
}

type Finder interface {
	Find(p Paginator) (interface{}, error)
}

type Pager interface {
	Counter
	Finder
}

type Paginator interface {
	Find() (interface{}, error)
	ParseRequest(r *http.Request) Paginator
	Offset() int
	Limit() int
	Total() int64
}

type resultData struct {
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
	From         int         `json:"from"`
	To           int         `json:"to"`
}

type paginator struct {
	counter    Counter
	finder     Finder
	PerPage    int
	resultData *resultData
	values     url.Values
}

func init() {
	zap.InitZapSugar()
}

func New(counter Counter, finder Finder) Paginator {
	return &paginator{
		counter:    counter,
		finder:     finder,
		PerPage:    DefaultPaginatorPerPage,
		resultData: &resultData{},
		values:     nil,
	}
}

func (p *paginator) Find() (interface{}, error) {
	if err := p.find(); err != nil {
		return nil, err
	}
	return p.resultData, nil
}

func (p *paginator) ParseRequest(r *http.Request) Paginator {
	var err error

	p.values = r.URL.Query()
	perPage := p.values.Get("per_page")
	p.resultData.PerPage, err = strconv.Atoi(perPage)
	if err != nil {
		p.resultData.PerPage = p.PerPage
	}
	current := p.values.Get("page")
	p.resultData.CurrentPage, err = strconv.Atoi(current)
	if err != nil {
		p.resultData.CurrentPage = 1
	}

	p.resultData.Path = r.Host + r.URL.Path
	if DEBUG {
		fmt.Println("parse query", "raw", r.URL.RawQuery, "per page", perPage, "page", current, "path", p.resultData.Path)
	}
	return p
}

func (p *paginator) Offset() int {
	return p.resultData.From
}

func (p *paginator) Limit() int {
	return p.resultData.PerPage
}

func (p *paginator) Total() int64 {
	return p.resultData.Total
}

func (p *paginator) find() error {
	var count int64

	err := p.counter.Count(&count)
	if DEBUG {
		log.Infow("count", "error", err, "count", count)
	}
	if err != nil {
		return err
	}

	p.resultData.make(count)

	v, err := p.finder.Find(p)
	if err != nil {
		return err
	}
	p.resultData.Data = v

	return nil
}

func (p *resultData) make(count int64) {
	if count == 0 {
		p.CurrentPage = 1
		return
	}

	p.Total = count
	p.LastPage = int(math.Ceil(float64(p.Total) / float64(p.PerPage)))
	if p.CurrentPage <= 0 || p.CurrentPage > p.LastPage {
		p.CurrentPage = 1
	}
	p.From = (p.CurrentPage - 1) * p.PerPage
	p.To = p.From + p.PerPage
	p.NextPageURL = p.next()
	p.PrevPageURL = p.prev()
	p.LastPageURL = p.last()
	p.FirstPageURL = p.first()
	return
}

func (p *resultData) next() string {
	if p.LastPage > p.CurrentPage+1 {
		v := url.Values{}
		p.perPage(v)
		page(v, p.CurrentPage+1)
		return p.Path + "?" + v.Encode()
	}
	return ""
}

func (p *resultData) prev() string {
	if p.CurrentPage-1 > 0 {
		v := url.Values{}
		p.perPage(v)
		page(v, p.CurrentPage-1)
		return p.Path + "?" + v.Encode()
	}
	return ""
}

func (p *resultData) last() string {
	if p.LastPage > 0 {
		v := url.Values{}
		p.perPage(v)
		page(v, p.LastPage)
		return p.Path + "?" + v.Encode()
	}
	return ""
}
func (p *resultData) first() string {
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

func (p *resultData) page(values url.Values, i int) {
	values.Set("per_page", strconv.Itoa(i))
}

func (p *resultData) perPage(values url.Values) {
	perPage(values, p.PerPage)
}
