package paginator

type Pager interface {
	Counter
	Finder
	Requester
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
