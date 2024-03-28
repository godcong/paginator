package paginator

type CustomValue map[string]any

type PageResult struct {
	Success      bool   `json:"success" xml:"success"`
	PerPage      int    `json:"per_page,omitempty" xml:"per_page"`
	Total        int64  `json:"total,omitempty" xml:"total"`
	CurrentPage  int    `json:"current_page,omitempty" xml:"current_page"`
	PreviousPage int    `json:"previous_page,omitempty" xml:"previous_page"`
	NextPage     int    `json:"next_page,omitempty" xml:"next_page"`
	LastPage     int    `json:"last_page,omitempty" xml:"last_page"`
	Path         string `json:"path,omitempty" xml:"path"`
	CurrentPath  string `json:"current_path,omitempty" xml:"current_path"`
	FirstPageURL string `json:"first_page_url,omitempty" xml:"first_page_url"`
	LastPageURL  string `json:"last_page_url,omitempty" xml:"last_page_url"`
	NextPageURL  string `json:"next_page_url,omitempty" xml:"next_page_url"`
	PrevPageURL  string `json:"prev_page_url,omitempty" xml:"prev_page_url"`
	Data         any    `json:"data,omitempty" xml:"data"`
	CustomValue  `json:"-" xml:"-"`
}
