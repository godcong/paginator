package paginator

const (
	defaultPageKey          = "page"
	defaultStarPage         = 1
	defaultPerPageKey       = "per_page"
	defaultPaginatorPerPage = 50
	defaultDataKey          = "data"
	defaultFirstPageKey     = "first_page_url"
	defaultLastPageKey      = "last_page_url"
	defaultNextPageKey      = "next_page_url"
	defaultPrevPageKey      = "prev_page_url"
	defaultCurrentPageKey   = "current_page"
	defaultTotalKey         = "total"
	defaultPathKey          = "path"
)

type Option struct {
	staPage        int
	perPage        int
	perPageKey     string
	pageKey        string
	dataKey        string
	firstPageKey   string
	lastPageKey    string
	nextPageKey    string
	prevPageKey    string
	currentPageKey string
	totalKey       string
	pathKey        string
}

// SetPathKey
// @receiver *Option
// @param string
// @return *Option
func (o *Option) SetPathKey(pathKey string) *Option {
	o.pathKey = pathKey
	return o
}

// SetCurrentPageKey
// @receiver *Option
// @param string
// @return *Option
func (o *Option) SetCurrentPageKey(currentPageKey string) *Option {
	o.currentPageKey = currentPageKey
	return o
}

// SetTotalKey
// @receiver *Option
// @param string
// @return *Option
func (o *Option) SetTotalKey(totalKey string) *Option {
	o.totalKey = totalKey
	return o
}

// LastPageKey
// @receiver *Option
// @return string
func (o *Option) LastPageKey() string {
	return o.lastPageKey
}

// SetLastPageKey
// @receiver *Option
// @param string
// @return *Option
func (o *Option) SetLastPageKey(lastPageKey string) *Option {
	o.lastPageKey = lastPageKey
	return o
}

// NextPageKey
// @receiver *Option
// @return string
func (o *Option) NextPageKey() string {
	return o.nextPageKey
}

// SetNextPageKey
// @receiver *Option
// @param string
// @return *Option
func (o *Option) SetNextPageKey(nextPageKey string) *Option {
	o.nextPageKey = nextPageKey
	return o
}

// PrevPageKey
// @receiver *Option
// @return string
func (o *Option) PrevPageKey() string {
	return o.prevPageKey
}

// SetPrevPageKey
// @receiver *Option
// @param string
// @return *Option
func (o *Option) SetPrevPageKey(prevPageKey string) *Option {
	o.prevPageKey = prevPageKey
	return o
}

// SetDataKey
// @receiver *Option
// @param string
// @return *Option
func (o *Option) SetDataKey(dataKey string) *Option {
	o.dataKey = dataKey
	return o
}

// SetFirstPageKey
// @receiver *Option
// @param string
// @return *Option
func (o *Option) SetFirstPageKey(firstPageKey string) *Option {
	o.firstPageKey = firstPageKey
	return o
}

// DataKey
// @receiver *Option
// @return string
func (o *Option) DataKey() string {
	return o.dataKey
}

type OptionSet func(p *Option)

// StaPage
// @receiver *Option
// @return int
func (o *Option) StaPage() int {
	return o.staPage
}

// SetStaPage
// @receiver *Option
// @param int
func (o *Option) SetStaPage(staPage int) *Option {
	o.staPage = staPage
	return o
}

// PerPage
// @receiver *Option
// @return int
func (o *Option) PerPage() int {
	return o.perPage
}

// SetPerPage
// @receiver *Option
// @param int
func (o *Option) SetPerPage(perPage int) *Option {
	o.perPage = perPage
	return o
}

// PerPageKey
// @receiver *Option
// @return string
func (o *Option) PerPageKey() string {
	return o.perPageKey
}

// SetPerPageKey
// @receiver *Option
// @param string
func (o *Option) SetPerPageKey(perPageKey string) *Option {
	o.perPageKey = perPageKey
	return o
}

func (o *Option) PageKey() string {
	return o.pageKey
}

func (o *Option) SetPageKey(pageKey string) *Option {
	o.pageKey = pageKey
	return o
}

func (o *Option) FirstPageKey() string {
	return o.firstPageKey
}

func (o *Option) TotalKey() string {
	return o.totalKey
}

func (o *Option) CurrentPageKey() string {
	return o.currentPageKey
}

func (o *Option) PathKey() string {
	return o.pathKey
}

// PerPageOption ...
// @Description: per page init on paginator
// @param int
// @return OptionSet
func PerPageOption(perPage int) OptionSet {
	return func(p *Option) {
		p.perPage = perPage
	}
}

// DataKeyOption
// @param string
// @return OptionSet
func DataKeyOption(dataKey string) OptionSet {
	return func(p *Option) {
		p.dataKey = dataKey
	}
}

// PerPageKeyOption ...
// @Description: per page key init on paginator
// @param string
// @return OptionSet
func PerPageKeyOption(key string) OptionSet {
	return func(p *Option) {
		p.perPageKey = key
	}
}

// PageKeyOption ...
// @Description: per key init on paginator
// @param string
// @return OptionSet
func PageKeyOption(key string) OptionSet {
	return func(p *Option) {
		p.pageKey = key
	}
}

func defaultOption() *Option {
	return &Option{
		staPage:        defaultStarPage,
		perPage:        defaultPaginatorPerPage,
		perPageKey:     defaultPerPageKey,
		pageKey:        defaultPageKey,
		dataKey:        defaultDataKey,
		firstPageKey:   defaultFirstPageKey,
		lastPageKey:    defaultLastPageKey,
		nextPageKey:    defaultNextPageKey,
		prevPageKey:    defaultPrevPageKey,
		currentPageKey: defaultCurrentPageKey,
		totalKey:       defaultTotalKey,
		pathKey:        defaultPathKey,
	}
}
