package paginator

import (
	"net/http"
)

// Hooker ...
// @Description:
type Hooker interface {
	Hooks() map[string]func([]string) bool
}

// CustomHooker ...
// @Description: return used keys
type CustomHooker interface {
	Hook() func(*http.Request) []string
}
