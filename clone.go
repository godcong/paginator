package paginator

type wrapClone struct {
	query  Queryable
	parser Parser
}

func (m wrapClone) Count() (int64, error) {
	return m.query.Finder(m.parser).Count()
}

func (m wrapClone) Get(page PageQuery) (any, error) {
	return m.query.Finder(m.parser).Get(page)
}

func (m wrapClone) Clone() Finder {
	return m.query.Finder(m.parser)
}
