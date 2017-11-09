package web

import (
	"net/http"

	"github.com/gliderlabs/com/objects"
	"github.com/gliderlabs/stdcom/log"
)

// ComponentContextKey is a key for the Request context that returns the
// Component instance serving the request.
var ComponentContextKey = "component"

// Register will register an instance of Component with a com.Registry
func Register(reg *objects.Registry) error {
	return reg.Register(&objects.Object{Value: &Component{}})
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
