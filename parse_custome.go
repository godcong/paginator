package paginator

type CustomParse struct {
	ParseHook ParserFunc
}

func (p CustomParse) Parse(cfg *Config, ignored map[string]struct{}) (PageReady, error) {
	return p.ParseHook(cfg, ignored)
}

func (p CustomParse) Type() string {
	return "custom"
}

// func (p CustomParse) parse(cfg *cfg, ignored map[string]struct{}) (*queryT, error) {
//     pr, err := p.ParseHook(cfg, ignored)
//     if err != nil {
//         return nil, err
//     }
//     q, ok := pr.(*queryT)
//     if !ok {
//         return nil, errors.New("invalid return type")
//     }
//     return q, nil
// }

func NewCustomParse(h ParserFunc) Parsable {
	return &CustomParse{
		ParseHook: h,
	}
}
