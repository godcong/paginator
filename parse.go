package paginator

// Values is a collection of values
type Values interface {
	Get(key string) string
	Set(key, value string)
	Add(key, value string)
	Del(key string)
	Has(key string) bool
	Encode() string
}

// Parser is a parser for Paginator
type Parser interface {
	FindValue(key string, d string) string
	FindArray(key string, d []string) []string
	FindNumber(key string, d float64) float64
	FindOthers() Values
	GetSource() any
}
