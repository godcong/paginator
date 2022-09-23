package paginator

type Parser[T any] interface {
	FindValue(key string, d string) string
	Get() T
}
