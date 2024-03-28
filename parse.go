package paginator

type ParserFunc = func(cfg *Config, ignored map[string]struct{}) (PageReady, error)

type parser interface {
	ParserFunc | *httpParse | *CustomParse
}

// Parsable is a parser for Paginator
type Parsable interface {
	Type() string
	Parse(cfg *Config, ignored map[string]struct{}) (PageReady, error)
}
