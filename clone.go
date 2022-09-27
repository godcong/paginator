package paginator

import (
	"context"
)

type wrapClone struct {
	query  Queryable
	parser Parser
}

func (m wrapClone) Count(ctx context.Context) (int64, error) {
	return m.query.Finder(m.parser).Count(ctx)
}

func (m wrapClone) Get(ctx context.Context, page PageQuery) (any, error) {
	return m.query.Finder(m.parser).Get(ctx, page)
}

func (m wrapClone) Clone() Finder {
	return m.query.Finder(m.parser)
}
