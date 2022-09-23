package paginator

import (
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
)

const defaultPageKey = "page"
const defaultStarPage = 1
const defaultPerPageKey = "per_page"
const defaultPaginatorPerPage = 50

var DEBUG = false
var ErrArgumentRequest = fmt.Errorf("paginator: argument request is not a valid http.Request")

type Paginator interface {
	Parse(Parser[any]) (any, error)
}

//type T any

type paginator struct {
	staPage    int
	perPage    int
	perPageKey string
	pageKey    string
}

//func (p *paginator) getRequestPerPage(values url.Values) int {
//	perPageInt := defaultPaginatorPerPage
//	perPage := values.Get(p.perPageKey)
//	if perPage == "" {
//		return perPageInt
//	}
//	if DEBUG {
//		fmt.Println("per_page", perPage)
//	}
//	var err error
//	perPageInt, err = strconv.Atoi(perPage)
//	if err != nil {
//		return defaultPaginatorPerPage
//	}
//	if DEBUG {
//		fmt.Println("perPageInt", perPageInt)
//	}
//	return perPageInt
//}

//func (p *paginator) getRequestCurrent(values url.Values) int {
//	currentInt, err := strconv.Atoi(values.Get(p.pageKey))
//	if err != nil {
//		return 1
//	}
//	return currentInt
//}

func (p *paginator) PageKey() string {
	return p.pageKey
}

func (p *paginator) PerPageKey() string {
	return p.perPageKey
}

func (p *paginator) Parse(parser Parser[any]) (any, error) {
	page := p.initialize(parser)
	//var page pageQuery
	//parser, err := ps.Parse()
	//if err != nil {
	//	return page, err
	//}
	err := p.findTotal(page)
	if err != nil {
		return nil, err
	}

	//parser.
	//r := t.Request()
	//if r == nil {
	//	return Find{}, ErrArgumentRequest
	//}
	//v := r.URL.Query()
	//page := p.initPage(r, v)
	//if DEBUG {
	//	fmt.Println("parse query", "per setPage", page.PerPage, "setPage", page.CurrentPage, "path", page.Path)
	//}

	//var conds []string
	//if h, ok := t.(InitHooker); ok {
	//	h.Initialize()
	//}

	//if h, ok := t.(CustomHooker); ok {
	//	conds = h.Hook()(t.Request())
	//}

	//var vv []string
	//h, ok := t.(Hooker)
	//if ok {
	//	for key, hook := range h.Hooks() {
	//		if vv, ok = v[key]; ok {
	//			if hook(vv) {
	//				conds = append(conds, key)
	//			}
	//		}
	//	}
	//}

	//it := &iterator{
	//	v:     v,
	//	conds: conds,
	//}
	//if DEBUG {
	//	fmt.Println("page.per_page", page.PerPage)
	//}

	err = p.findTotal(parser)
	if err != nil {
		return page, err
	}
	from := (page.CurrentPage - 1) * parser.PerPage
	to := from + parser.PerPage
	q, err := p.q.Find()
	if err != nil {
		return page, err
	}
	//found, err := t.Find(&pageQuery{
	//	page: page,
	//	from: from,
	//	to:   to,
	//	it:   it,
	//})
	//if err != nil {
	//	return Find{}, err
	//}
	parser.Data = q
	return *page, nil
}

func (p *paginator) findTotal(pa *pageReady) error {
	count, err := p.counter.Count(pa)
	if err != nil {
		return err
	}
	if count == 0 {
		pa.CurrentPage = 1
		return nil
	}

	pa.Total = count
	pa.LastPage = int(math.Ceil(float64(pa.Total) / float64(pa.PerPage)))
	if pa.CurrentPage <= 0 || pa.CurrentPage > pa.LastPage {
		pa.CurrentPage = 1
	}
	//pa.NextPageURL = p.next(pa, it.Values())
	//pa.PrevPageURL = p.prev(pa, it.Values())
	//pa.LastPageURL = p.last(pa, it.Values())
	//pa.FirstPageURL = p.first(pa, it.Values())
	return nil
}

// New ...
// @Description: create paginator object for use anywhere
// @param ...Option
// @return Paginator
func New(ops ...Option) Paginator {
	p := &paginator{
		perPage:    defaultPaginatorPerPage,
		perPageKey: defaultPerPageKey,
		pageKey:    defaultPageKey,
		staPage:    defaultStarPage,
	}
	for i := range ops {
		ops[i](p)
	}
	return &paginator{
		q: q,
		s: s,
	}
}

func (p *paginator) PerPage() int {
	return p.perPage
}

func (p *paginator) SetPerPage(perPage int) {
	p.perPage = perPage
}

func (p *paginator) next(pa *Parse[any], v url.Values) string {
	if pa.LastPage > pa.CurrentPage+1 {
		setPerPage(v, p.perPageKey, pa.PerPage)
		setPage(v, p.pageKey, pa.CurrentPage+1)
		return pa.Path + "?" + v.Encode()
	}
	return ""
}

func (p *paginator) prev(pa *Parse[any], v url.Values) string {
	if pa.CurrentPage-1 > 0 {
		setPerPage(v, p.perPageKey, pa.PerPage)
		setPage(v, p.pageKey, pa.CurrentPage-1)
		return pa.Path + "?" + v.Encode()
	}
	return ""
}

func (p *paginator) last(pa *Parse[any], v url.Values) string {
	if pa.LastPage > 0 {
		setPerPage(v, p.perPageKey, pa.PerPage)
		setPage(v, p.pageKey, pa.LastPage)
		return pa.Path + "?" + v.Encode()
	}
	return ""
}
func (p *paginator) first(pa *Parse[any], v url.Values) string {
	if pa.Total > 0 {
		setPerPage(v, p.perPageKey, pa.PerPage)
		setPage(v, p.pageKey, 1)
		return pa.Path + "?" + v.Encode()
	}
	return ""
}

func setPage(values url.Values, key string, i int) {
	values.Set(key, strconv.Itoa(i))
}

func setPerPage(values url.Values, key string, i int) {
	values.Set(key, strconv.Itoa(i))
}

func (p *paginator) setPage(pa *Parse[any], values url.Values) {
	setPage(values, p.pageKey, pa.CurrentPage)
}

func (p *paginator) setPerPage(pa *Parse[any], values url.Values) {
	setPerPage(values, p.perPageKey, pa.PerPage)
}

func (p *paginator) initialize(parser Parser[any]) *pageReady {
	page := new(pageQuery)
	switch v := parser.(type) {
	case Parser[*http.Request]:
		r := v.Get()
		page.Path = r.URL.Scheme + r.Host + r.URL.Path
	}
	page.PerPage = stoi(parser.FindValue(p.perPageKey, ""), p.perPage)
	page.CurrentPage = stoi(parser.FindValue(p.pageKey, ""), p.staPage)
	return &pageReady{
		page:   page,
		parser: parser,
	}
}

func (p *paginator) Count(pa *pageReady) (int64, error) {

}

func stoi(s string, d int) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return d
	}
	return i
}

var _ Paginator = (*paginator)(nil)
