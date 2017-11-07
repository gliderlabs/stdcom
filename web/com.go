package web

import (
	"net/http"

	"github.com/gliderlabs/com"
	"github.com/gliderlabs/stdcom/log"
)

// ComponentContextKey is a key for the Request context that returns the
// Component instance serving the request.
var ComponentContextKey = "component"

// Register will register an instance of Component with a com.Registry
func Register(registry *com.Registry) error {
	return registry.Register(&com.Object{Value: &Component{}})
}

func FromRequest(r *http.Request) *Component {
	return r.Context().Value(ComponentContextKey).(*Component)
}

// Component for web serving
type Component struct {
	Log        log.Logger         `com:"singleton"`
	Handlers   []Handler          `com:"extpoint"`
	Middleware MiddlewareInjector `com:"singleton"`

	http  *http.Server
	https *http.Server
	cr    *TLSCertReloader

	Config
}
