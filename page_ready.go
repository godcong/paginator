package paginator

import (
	"context"
)

type QueryCountFunc func(ctx context.Context, query Query) (int64, error)
type QueryResultFunc func(ctx context.Context, query Query) (any, error)

type PageReady interface {
	CountHook(fn QueryCountFunc) PageReady
	ResultHook(fn QueryResultFunc) PageReady
	DoParse(ctx context.Context, obj Parsable, ignores ...string) (PageResult, error)
	Do(ctx context.Context) (PageResult, error)
}
