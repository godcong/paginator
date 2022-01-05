package paginator

// Pager ...
// @Description:
type Pager interface {
	Counter
	Finder
	Requester
}

// Page ...
// @Description:
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

type PrePager interface {
	Page() Page
	Offset() int
	To() int
	Limit() int
	Total() int64
	Current() int
	Iterator() Iterator
}

type prePage struct {
	page *Page
	from int
	to   int
	it   *iterator
}

func (p prePage) Page() Page {
	return *p.page
}

func (p prePage) Offset() int {
	return p.from
}

func (p prePage) Limit() int {
	return p.page.PerPage
}

func (p prePage) Total() int64 {
	return p.page.Total
}

func (p prePage) Current() int {
	return p.page.CurrentPage
}

func (p prePage) To() int {
	return p.to
}

func (p prePage) Iterator() Iterator {
	return p.it
}

var _ PrePager = (*prePage)(nil)
