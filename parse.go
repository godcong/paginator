package paginator

type Values interface {
	Get(key string) string
	Set(key, value string)
	Add(key, value string)
	Del(key string)
	Has(key string) bool
	Encode() string
}

type Parser[T any] interface {
	FindValue(key string, d string) string
	FindArray(key string, d []string) []string
	FindNumber(key string, d float64) float64
	FindOthers() Values
	GetSource() T
}
