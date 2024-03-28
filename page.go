package paginator

import (
	"context"
	"maps"
	"math"
	"net/url"
)

type Query interface {
	Offset() int
	Limit() int
	Config() *Config
	FullPath() string
	PerPage() int
	CurrentPage() int
}

type query struct {
	config      *Config
	countFunc   QueryCountFunc
	resultFunc  QueryResultFunc
	fullPath    string
	perPage     int
	currentPage int
	total       int64
	values      url.Values
	result      PageResult
}

func (q *query) DoParse(ctx context.Context, obj Parsable, ignores ...string) (PageResult, error) {
	ignored := make(map[string]struct{})
	for i := range ignores {
		ignored[ignores[i]] = struct{}{}
	}

	var (
		pr  PageReady
		err error
	)
	switch v := obj.(type) {
	case *httpParse:
		query, err := v.parse(q.Config(), ignored)
		if err != nil {
			return PageResult{}, err
		}
		pr = query
		query.resultFunc = q.resultFunc
		query.countFunc = q.countFunc

	case *CustomParse:
		pr, err = v.Parse(q.Config(), ignored)
		if err != nil {
			return PageResult{}, err
		}
		pr.ResultHook(q.resultFunc)
		pr.CountHook(q.countFunc)
	default:
		return PageResult{}, ErrInvalidParserType
	}

	return pr.Do(ctx)
}

func (q *query) Offset() int {
	return (q.currentPage - q.config.StartOffset()) * q.perPage

}

func (q *query) Limit() int {
	return q.perPage
}

func (q *query) OffsetLimit() int {
	return q.Offset() + q.Limit()
}

func (q *query) Config() *Config {
	return q.config
}

func (q *query) FullPath() string {
	return q.fullPath
}

func (q *query) PerPage() int {
	return q.perPage
}

func (q *query) CurrentPage() int {
	return q.currentPage
}

func (q *query) Do(ctx context.Context) (PageResult, error) {
	var err error

	if err = q.calcTotal(ctx); err != nil {
		return PageResult{}, err
	}

	if q.total == 0 {
		return PageResult{}, nil
	}

	q.calcPath()

	var data any
	if q.resultFunc != nil {
		data, err = q.resultFunc(ctx, q)
		if err != nil {
			return PageResult{}, err
		}
	}
	q.result.Data = data
	q.result.Success = true
	return q.result, nil
}

func (q *query) calcTotal(ctx context.Context) error {
	if q.countFunc == nil {
		return nil
	}
	var err error
	q.total, err = q.countFunc(ctx, q)
	if err != nil {
		return err
	}

	if q.total == 0 {
		q.currentPage = 1
		return nil
	}

	q.result.Total = q.total
	q.result.PerPage = q.perPage

	lastPage := int(math.Ceil(float64(q.total) / float64(q.perPage)))
	q.result.LastPage = lastPage
	if q.currentPage <= 0 || q.currentPage > lastPage {
		q.currentPage = 1
	}
	if q.currentPage > 1 {
		q.result.PreviousPage = q.currentPage - 1
	}
	if lastPage > q.currentPage {
		q.result.NextPage = q.currentPage + 1
	}
	q.result.CurrentPage = q.currentPage
	return nil
}

func (q *query) calcPath() {
	q.result.Path = q.FullPath()

	if q.config.SkipPath {
		return
	}
	q.result.CurrentPath = q.currentURL()
	q.result.NextPageURL = q.nextURL()
	q.result.PrevPageURL = q.prevURL()
	q.result.LastPageURL = q.lastURL()
	q.result.FirstPageURL = q.firstURL()
}

func (q *query) CountHook(cf QueryCountFunc) PageReady {
	q.countFunc = cf
	return q
}

func (q *query) ResultHook(rf QueryResultFunc) PageReady {
	q.resultFunc = rf
	return q
}

func (q *query) nextURL() string {
	if q.result.LastPage > q.result.CurrentPage {
		enc := q.encodeValues()
		setPerPage(enc, q.config.GetKey(KeyPerPage), q.result.PerPage)
		setPage(enc, q.config.GetKey(KeyPage), q.result.CurrentPage+1)
		return q.FullPath() + "?" + enc.Encode()
	}
	return ""
}

func (q *query) prevURL() string {
	if q.result.CurrentPage > 1 {
		enc := q.encodeValues()
		setPerPage(enc, q.config.GetKey(KeyPerPage), q.result.PerPage)
		setPage(enc, q.config.GetKey(KeyPage), q.result.CurrentPage-1)
		return q.FullPath() + "?" + enc.Encode()
	}
	return ""
}

func (q *query) lastURL() string {
	if q.result.LastPage > 0 {
		enc := q.encodeValues()
		setPerPage(enc, q.config.GetKey(KeyPerPage), q.result.PerPage)
		setPage(enc, q.config.GetKey(KeyPage), q.result.LastPage)
		return q.FullPath() + "?" + enc.Encode()
	}
	return ""
}

func (q *query) firstURL() string {
	if q.result.Total > 0 {
		enc := q.encodeValues()
		setPerPage(enc, q.config.GetKey(KeyPerPage), q.result.PerPage)
		setPage(enc, q.config.GetKey(KeyPage), q.config.StartPage)
		return q.FullPath() + "?" + enc.Encode()
	}
	return ""
}

func (q *query) currentURL() string {
	enc := q.encodeValues()
	setPerPage(enc, q.config.GetKey(KeyPerPage), q.result.PerPage)
	setPage(enc, q.config.GetKey(KeyPage), q.result.CurrentPage)
	return q.FullPath() + "?" + enc.Encode()
}

func (q *query) encodeValues() url.Values {
	return maps.Clone(q.values)
}
