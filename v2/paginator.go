package paginator

import (
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
)

const defaultPaginatorPerPage = 15
const defaultPageKey = "setPage"
const defaultPerPageKey = "per_page"

var DEBUG = false
var ErrArgumentRequest = fmt.Errorf("paginator: argument request is not a valid http.Request")

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
		fmt.Println("parse query", "per setPage", page.PerPage, "setPage", page.CurrentPage, "path", page.Path)
	}

	var conds []string
	if h, ok := t.(CustomHooker); ok {
		h.Hook()(t.Request())
	}

	var vv []string
	h, ok := t.(Hooker)
	if ok {
		for key, hook := range h.Hooks() {
			if vv, ok = v[key]; ok {
				if hook(vv) {
					conds = append(conds, key)
				}
			}
		}
	}

	it := &iterator{
		v:     v,
		conds: conds,
	}

	err := p.findTotal(page, it, t)
	if err != nil {
		return Page{}, err
	}
	from := (page.CurrentPage - 1) * page.PerPage
	to := from + page.PerPage

	found, err := t.Find(&prePage{
		page: page,
		from: from,
		to:   to,
		it:   it,
	})
	if err != nil {
		return Page{}, err
	}
	page.Data = found
	return *page, nil
}

func (p *paginator) findTotal(page *Page, it Iterator, t Turnable) error {
	count, err := t.Count(it)
	if err != nil {
		return err
	}
	if count == 0 {
		page.CurrentPage = 1
		return nil
	}

	page.Total = count
	page.LastPage = int(math.Ceil(float64(page.Total) / float64(page.PerPage)))
	if page.CurrentPage <= 0 || page.CurrentPage > page.LastPage {
		page.CurrentPage = 1
	}
	page.NextPageURL = p.next(page, it.Values())
	page.PrevPageURL = p.prev(page, it.Values())
	page.LastPageURL = p.last(page, it.Values())
	page.FirstPageURL = p.first(page, it.Values())
	return nil
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

func (p *paginator) initPage(r *http.Request, v url.Values) (page *Page) {
	return &Page{
		CurrentPage: p.getRequestCurrent(v),
		//LastPage:     0,
		PerPage: p.getRequestPerPage(v),
		//Data:         nil,
		//Total:        0,
		//FirstPageURL: "",
		//LastPageURL:  "",
		//NextPageURL:  "",
		//PrevPageURL:  "",
		Path: r.Host + r.URL.Path,
	}
}

func (p *paginator) next(page *Page, v url.Values) string {
	if page.LastPage > page.CurrentPage+1 {
		setPerPage(v, p.perPageKey, page.PerPage)
		setPage(v, p.pageKey, page.CurrentPage+1)
		return page.Path + "?" + v.Encode()
	}
	return ""
}

func (p *paginator) prev(page *Page, v url.Values) string {
	if page.CurrentPage-1 > 0 {
		setPerPage(v, p.perPageKey, page.PerPage)
		setPage(v, p.pageKey, page.CurrentPage-1)
		return page.Path + "?" + v.Encode()
	}
	return ""
}

func (p *paginator) last(page *Page, v url.Values) string {
	if page.LastPage > 0 {
		setPerPage(v, p.perPageKey, page.PerPage)
		setPage(v, p.pageKey, page.LastPage)
		return page.Path + "?" + v.Encode()
	}
	return ""
}
func (p *paginator) first(page *Page, v url.Values) string {
	if page.Total > 0 {
		setPerPage(v, p.perPageKey, page.PerPage)
		setPage(v, p.pageKey, 1)
		return page.Path + "?" + v.Encode()
	}
	return ""
}

func setPage(values url.Values, key string, i int) {
	values.Set(key, strconv.Itoa(i))
}

func setPerPage(values url.Values, key string, i int) {
	values.Set(key, strconv.Itoa(i))
}

func (p *paginator) setPage(page *Page, values url.Values) {
	setPage(values, p.pageKey, page.CurrentPage)
}

func (p *paginator) setPerPage(page *Page, values url.Values) {
	setPerPage(values, p.perPageKey, page.PerPage)
}

var _ Paginator = (*paginator)(nil)
