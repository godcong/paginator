package paginator

import (
	"net/http"
)

// HttpHooker ...
// @Description:
type HttpHooker[V any] func(*Parse, *http.Request) error

// CustomHooker ...
// @Description: return used keys
type CustomHooker interface {
	Hook() func(*http.Request) []string
}

// InitHooker ...
// @Description: call init before hook
type InitHooker interface {
	Initialize()
}
